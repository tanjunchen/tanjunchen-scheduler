package dynamic

import (
	"math"

	"k8s.io/apimachinery/pkg/api/resource"
)

type NodeInfo struct {
	NodeName    string
	Labels      map[string]string
	Annotations map[string]string

	RealCPURate          float64
	RequestCPURate       float64
	RemainAllocatableCPU resource.Quantity

	RealMemoryRate          float64
	RequestMemoryRate       float64
	RemainAllocatableMemory resource.Quantity
}

type NodeInfos []NodeInfo

// 1.RealMemoryRate
// 2.RealCPURate
// 3.RequestMemoryRate
// 4.RequestCPURate
func (s NodeInfos) Less(i, j int) bool {
	e := 5.0
	if math.Abs(float64(s[i].RealMemoryRate-s[j].RealMemoryRate)) < e {
		if math.Abs(float64(s[i].RealCPURate-s[j].RealCPURate)) < e {
			if math.Abs(float64(s[i].RequestMemoryRate-s[j].RealMemoryRate)) < e {
				return s[i].RequestCPURate > s[j].RequestCPURate
			}
			return s[i].RequestMemoryRate > s[j].RequestMemoryRate
		}
		return s[i].RealCPURate > s[j].RealCPURate
	}
	return s[i].RealMemoryRate > s[j].RealMemoryRate
}

// Swap s
func (s NodeInfos) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Len l
func (s NodeInfos) Len() int {
	return len(s)
}
