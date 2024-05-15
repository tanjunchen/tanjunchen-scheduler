package example

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"

	"github.com/tanjunchen/tanjunchen-scheduler/pkg/names"
)

var _ framework.FilterPlugin = &ExamplePlugin{}

type ExamplePlugin struct{}

// NewExampleSchedPlugin initializes a new plugin and returns it.
func NewExamplePlugin(_ runtime.Object, _ framework.Handle) (framework.Plugin, error) {
	return &ExamplePlugin{}, nil
}

func (e *ExamplePlugin) Filter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	cpu := nodeInfo.Allocatable.MilliCPU
	memory := nodeInfo.Allocatable.Memory
	klog.InfoS("tanjunchen-scheduler Filter", "pod_name", pod.Name, "current node", nodeInfo.Node().Name, "cpu", cpu, "memory", memory)
	return framework.NewStatus(framework.Success, "")
}

func (e *ExamplePlugin) Name() string {
	return names.ExampleName
}

type exampleStateData struct {
	data string
}

func (s *exampleStateData) Clone() framework.StateData {
	copy := &exampleStateData{
		data: s.data,
	}
	return copy
}
