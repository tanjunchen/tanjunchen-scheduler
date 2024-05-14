package scheme

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	kubeschedulerscheme "k8s.io/kubernetes/pkg/scheduler/apis/config/scheme"

	config "github.com/tanjunchen/tanjunchen-scheduler/apis/config"
	v1 "github.com/tanjunchen/tanjunchen-scheduler/apis/config/v1"
	v1beta2 "github.com/tanjunchen/tanjunchen-scheduler/apis/config/v1beta2"
	v1beta3 "github.com/tanjunchen/tanjunchen-scheduler/apis/config/v1beta3"
)

var (
	// Re-use the in-tree Scheme.
	Scheme = kubeschedulerscheme.Scheme

	// Codecs provides access to encoding and decoding for the scheme.
	Codecs = serializer.NewCodecFactory(Scheme, serializer.EnableStrict)
)

func init() {
	AddToScheme(Scheme)
}

// AddToScheme builds the kubescheduler scheme using all known versions of the kubescheduler api.
func AddToScheme(scheme *runtime.Scheme) {
	utilruntime.Must(config.AddToScheme(scheme))
	utilruntime.Must(v1.AddToScheme(scheme))
	utilruntime.Must(v1beta2.AddToScheme(scheme))
	utilruntime.Must(v1beta3.AddToScheme(scheme))
}
