package dynamic

import (
	"container/list"
	"context"
	"fmt"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	coreinformerv1 "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
)

var (
	reSyncPeriod     = time.Minute
	nodePodIndexName = "node_pod"
)

// Cache cache
type Cache interface {
	GetNodeInfos(nodeNames []string, pod *corev1.Pod) NodeInfos
	GetNodeInfo(nodeName string, pod *corev1.Pod) NodeInfo
	Init() error
	Close()
}

// NodeCache
type NodeCache struct {
	metricsClient metricsclientset.Interface
	clientSet     kubernetes.Interface
	nodeInformer  coreinformerv1.NodeInformer
	podInformer   cache.SharedIndexInformer
	stopCh        chan struct{}
	nodeMetrics   map[string]*list.List
	sync.RWMutex
}

// NewNodeCache new node cache
func NewNodeCache(kc *rest.Config) (*NodeCache, error) {
	metricsClient, err := metricsclientset.NewForConfig(kc)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(kc)
	if err != nil {
		return nil, err
	}

	nc := NodeCache{
		clientSet:     client,
		stopCh:        make(chan struct{}),
		metricsClient: metricsClient,
		nodeMetrics:   make(map[string]*list.List),
	}
	nc.Init()
	return &nc, nil
}

// Init init node cache
func (nc *NodeCache) Init() error {
	informerFactory := informers.NewSharedInformerFactory(nc.clientSet, reSyncPeriod)

	nodeInformer := informerFactory.Core().V1().Nodes()
	go nodeInformer.Informer().Run(nc.stopCh)

	if !cache.WaitForCacheSync(nc.stopCh, nodeInformer.Informer().HasSynced) {
		return fmt.Errorf("wait for node cache sync error")
	}
	klog.Info("sync node cache successful.")

	podInformer := coreinformerv1.NewPodInformer(
		nc.clientSet, metav1.NamespaceAll, 0,
		cache.Indexers{
			nodePodIndexName: func(obj interface{}) ([]string, error) {
				pod, ok := obj.(*corev1.Pod)
				if !ok {
					return nil, fmt.Errorf("type error")
				}
				return []string{pod.Spec.NodeName}, nil
			},
		},
	)

	go podInformer.Run(nc.stopCh)

	if !cache.WaitForCacheSync(nc.stopCh, podInformer.HasSynced) {
		return fmt.Errorf("wait for all pod cache sync error")
	}
	klog.Info("sync all pod cache successful.")

	nc.nodeInformer = nodeInformer
	nc.podInformer = podInformer

	for i := 0; i < 5; i++ {
		if nc.scrapeNodeMetrics() {
			klog.Info("init all node metrics successful.")
			break
		}
		klog.Warningf("try init all node metrics failed, on %d times", i)
	}

	go func() {
		wait.Until(func() { nc.scrapeNodeMetrics() }, time.Minute, nc.stopCh)
	}()

	return nil
}

func (nc *NodeCache) calcNodeRequestResourceTotal(nodeName string,
	pod *corev1.Pod) (totalCPU resource.Quantity, totalMemory resource.Quantity) {
	objs, err := nc.podInformer.GetIndexer().ByIndex(nodePodIndexName, nodeName)
	if err != nil {
		return
	}

	klog.V(3).Infof("get %d pod for node %v", len(objs), nodeName)

	objs = append(objs, pod)

	for _, item := range objs {
		p, ok := item.(*corev1.Pod)
		if !ok {
			klog.Errorf("kind is not *corev1.Pod")
			continue
		}

		// TODO initContainer ?
		for _, c := range p.Spec.Containers {
			if c.Resources.Requests.Cpu() != nil {
				totalCPU.Add(*c.Resources.Requests.Cpu())
			}

			if c.Resources.Requests.Memory() != nil {
				totalMemory.Add(*c.Resources.Requests.Memory())
			}
		}
	}

	return
}

func (nc *NodeCache) scrapeNodeMetrics() bool {
	klog.V(3).Infof("start to scrape node metrics...")

	nodes, err := nc.nodeInformer.Lister().List(labels.Everything())
	if err != nil {
		klog.Errorf("list node from node informer err: %v", err)
		return false
	}

	for _, n := range nodes {
		metrics, err := nc.metricsClient.MetricsV1beta1().NodeMetricses().Get(context.TODO(), n.Name, metav1.GetOptions{})
		if err != nil {
			klog.Warningf("get node: %v metrics err: %v", n.Name, err)
			continue
		}

		klog.V(3).Infof("get %v node metrics success", n.Name)

		nc.Lock()
		if _, ok := nc.nodeMetrics[n.Name]; !ok {
			nc.nodeMetrics[n.Name] = list.New()
		}
		nc.nodeMetrics[n.Name].PushBack(metrics)
		if nc.nodeMetrics[n.Name].Len() > 15 {
			nc.nodeMetrics[n.Name].Remove(nc.nodeMetrics[n.Name].Front())
		}
		nc.Unlock()

		time.Sleep(time.Millisecond * 500)
	}

	nc.RLock()
	res := len(nc.nodeMetrics) == len(nodes)
	nc.RUnlock()

	return res
}

func (nc *NodeCache) getNodeMetrics(nodeName string) *metricsv1beta1.NodeMetrics {
	nc.RLock()
	defer nc.RUnlock()

	if nc.nodeMetrics[nodeName] == nil {
		return nil
	}

	// TODO consider better algorithm
	e := nc.nodeMetrics[nodeName].Back()
	if e == nil {
		return nil
	}

	return e.Value.(*metricsv1beta1.NodeMetrics)
}

// GetNodeInfos get nodes cpu state
func (nc *NodeCache) GetNodeInfos(nodeNames []string, pod *corev1.Pod) NodeInfos {
	infos := make([]NodeInfo, len(nodeNames))

	wg := sync.WaitGroup{}
	for i, name := range nodeNames {
		wg.Add(1)
		go func(i int, name string) {
			defer wg.Done()
			infos[i] = nc.GetNodeInfo(name, pod)
		}(i, name)
	}

	wg.Wait()
	return infos
}

// GetNodeInfo get single node cpu state
func (nc *NodeCache) GetNodeInfo(nodeName string, pod *corev1.Pod) NodeInfo {
	info := NodeInfo{
		NodeName:          nodeName,
		RealCPURate:       100,
		RequestCPURate:    100,
		RealMemoryRate:    100,
		RequestMemoryRate: 100,
	}

	node, err := nc.nodeInformer.Lister().Get(nodeName)
	if err != nil {
		klog.Errorf("ge%v err: %v", nodeName, err)
		return info
	}

	info.Labels = node.Labels
	info.Annotations = node.Annotations

	metrics := nc.getNodeMetrics(nodeName)
	if metrics == nil {
		return info
	}

	totalRequestCPU, totalRequestMemory := nc.calcNodeRequestResourceTotal(nodeName, pod)
	if use, ok := metrics.Usage[corev1.ResourceCPU]; ok {
		lastTotalCPU := *node.Status.Allocatable.Cpu()
		lastTotalCPU.Sub(totalRequestCPU)
		info.RemainAllocatableCPU = lastTotalCPU
		info.RequestCPURate = 100 * (float64(totalRequestCPU.MilliValue()) / float64(node.Status.Allocatable.Cpu().MilliValue()))
		info.RealCPURate = 100 * (float64(use.MilliValue()) / float64(node.Status.Capacity.Cpu().MilliValue()))
	}

	if use, ok := metrics.Usage[corev1.ResourceMemory]; ok {
		lastTotalMemory := *node.Status.Allocatable.Memory()
		lastTotalMemory.Sub(totalRequestMemory)
		info.RemainAllocatableMemory = lastTotalMemory
		info.RequestMemoryRate = 100 * (float64(totalRequestMemory.MilliValue()) / float64(node.Status.Allocatable.Memory().MilliValue()))
		info.RealMemoryRate = 100 * (float64(use.MilliValue()) / float64(node.Status.Capacity.Memory().MilliValue()))
	}
	return info
}

// Close close every thing
func (nc *NodeCache) Close() {
	klog.Infof("close node cache")
	close(nc.stopCh)
}
