package dynamic

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/framework"

	config "github.com/tanjunchen/tanjunchen-scheduler/apis/config"
)

const DynamicName = "dynamic"

var _ framework.FilterPlugin = &DynamicPlugin{}

type DynamicPlugin struct {
	handle      framework.Handle
	NodeCache   Cache
	DynamicArgs *config.DynamicArgs
}

// NewDynamicPlugin initializes a new plugin and returns it.
func NewDynamicPlugin(plArgs runtime.Object, handle framework.Handle) (framework.Plugin, error) {
	args, ok := plArgs.(*config.DynamicArgs)
	if !ok {
		return nil, fmt.Errorf("want args to be of type DynamicArgs, got %T", args)
	}
	cfg := handle.KubeConfig()
	nc, err := NewNodeCache(cfg)
	if err != nil {
		return nil, err
	}

	return &DynamicPlugin{
		DynamicArgs: args,
		handle:      handle,
		NodeCache:   nc,
	}, nil
}

func (dp *DynamicPlugin) Filter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	node := nodeInfo.Node()
	if node == nil {
		return framework.NewStatus(framework.Error, "node not found")
	}

	nodeName := node.Name
	nodesStat := dp.NodeCache.GetNodeInfo(nodeName, pod)

	fmt.Printf("nodesStat: %+v\n", nodesStat)

	if nodesStat.RealCPURate > dp.DynamicArgs.ToleranceCPURate {
		return framework.NewStatus(framework.Unschedulable, fmt.Sprintf("Real cpu rate > %v", dp.DynamicArgs.ToleranceCPURate))
	}

	if nodesStat.RealMemoryRate > dp.DynamicArgs.ToleranceMemoryRate {
		return framework.NewStatus(framework.Unschedulable, fmt.Sprintf("Real memory rate > %v", dp.DynamicArgs.ToleranceMemoryRate))
	}

	return framework.NewStatus(framework.Success, "")
}

func (dp *DynamicPlugin) Name() string {
	return DynamicName
}
