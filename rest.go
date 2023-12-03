package kubelib

import (
	"fmt"
	"net/http"

	gruntime "runtime"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type RESTClientOptions struct {
	GroupVersion schema.GroupVersion
	APIPath      string
	Verbose      bool
}

func RESTClientForAPI(restConfig *rest.Config, opts RESTClientOptions) (*rest.RESTClient, error) {
	config := *restConfig
	config.GroupVersion = &opts.GroupVersion
	config.APIPath = opts.APIPath
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = fmt.Sprintf("kubelib (%s/%s)", gruntime.GOOS, gruntime.GOARCH)

	if opts.Verbose {
		config.WrapTransport = func(rt http.RoundTripper) http.RoundTripper {
			return &wireTracer{RoundTripper: rt}
		}
	}

	return rest.RESTClientFor(&config)
}
