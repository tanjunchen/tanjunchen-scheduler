package nodefilter

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

const Name = "nodeFilter"

// 定义这个插件的结构体
type NodeFilter struct{}

// 实现 Name 方法
func (nodeFilter *NodeFilter) Name() string {
	return Name
}

// 实现 Filter 方法
func (nodeFilter *NodeFilter) Filter(ctx context.Context, _ *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	cpu := nodeInfo.Allocatable.MilliCPU
	memory := nodeInfo.Allocatable.Memory
	klog.InfoS("tanjunchen-scheduler nodeFilter filter", "pod_name", pod.Name, "current node", nodeInfo.Node().Name, "cpu", cpu, "memory", memory)
	return nil
}

// 编写 New 函数
func New(_ runtime.Object, _ framework.Handle) (framework.Plugin, error) {
	return &NodeFilter{}, nil
}
