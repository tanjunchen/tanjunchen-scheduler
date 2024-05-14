package v1beta3

var (
	DefaultToleranceCPURate    float64 = 80
	DefaultToleranceMemoryRate float64 = 80
)

func SetDefaults_DynamicArgs(obj *DynamicArgs) {
	if obj.ToleranceCPURate == 0 {
		obj.ToleranceCPURate = DefaultToleranceCPURate
	}
	if obj.ToleranceMemoryRate == 0 {
		obj.ToleranceMemoryRate = DefaultToleranceMemoryRate
	}
}
