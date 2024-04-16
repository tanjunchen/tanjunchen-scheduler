package example

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

const ExampledName = "example"

var _ framework.ReservePlugin = &ExamplePlugin{}
var _ framework.PreBindPlugin = &ExamplePlugin{}
var _ framework.PreFilterPlugin = &ExamplePlugin{}
var _ framework.FilterPlugin = &ExamplePlugin{}

type ExamplePlugin struct{}

// NewExampleSchedPlugin initializes a new plugin and returns it.
func NewExampleSchedPlugin(_ runtime.Object, _ framework.Handle) (framework.Plugin, error) {
	return &ExamplePlugin{}, nil
}

func (e *ExamplePlugin) PreFilter(ctx context.Context, state *framework.CycleState, pod *v1.Pod) (*framework.PreFilterResult, *framework.Status) {
	klog.InfoS("tanjunchen-scheduler PreFilter", "pod_name", pod.Name)
	return nil, framework.NewStatus(framework.Success, "")
}

func (e *ExamplePlugin) Filter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	cpu := nodeInfo.Allocatable.MilliCPU
	memory := nodeInfo.Allocatable.Memory
	klog.InfoS("tanjunchen-scheduler Filter", "pod_name", pod.Name, "current node", nodeInfo.Node().Name, "cpu", cpu, "memory", memory)
	return framework.NewStatus(framework.Success, "")
}

func (e *ExamplePlugin) PreFilterExtensions() framework.PreFilterExtensions {
	return nil
}

func (e *ExamplePlugin) PreBind(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) *framework.Status {
	if pod == nil {
		return framework.NewStatus(framework.Error, "pod cannot be nil")
	}
	klog.InfoS("tanjunchen-scheduler PreBind", "pod_name", pod.Name, "current node", nodeName)
	return nil
}

func (e *ExamplePlugin) Reserve(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) *framework.Status {
	if pod == nil {
		return framework.NewStatus(framework.Error, "pod cannot be nil")
	}
	klog.InfoS("tanjunchen-scheduler Reserve", "pod_name", pod.Name, "current node", nodeName)
	return nil
}

func (e *ExamplePlugin) Unreserve(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) {
	klog.InfoS("tanjunchen-scheduler Unreserve", "pod_name", pod.Name, "current node", nodeName)
}

func (e *ExamplePlugin) Name() string {
	return ExampledName
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
