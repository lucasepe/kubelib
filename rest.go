package kubelib

import (
	"fmt"
	"net/http"

	gruntime "runtime"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type CreateRESTClientOptions struct {
	APIPath   string
	UserAgent string
	Verbose   bool
}

func CreateRESTClient(rc *rest.Config, gv schema.GroupVersion, opts CreateRESTClientOptions) (*rest.RESTClient, error) {
	config := *rc
	config.GroupVersion = &gv
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	config.APIPath = opts.APIPath
	if len(config.APIPath) == 0 {
		config.APIPath = "/apis"
	}

	config.UserAgent = opts.UserAgent
	if len(config.UserAgent) == 0 {
		config.UserAgent = fmt.Sprintf("kubelib (%s/%s)", gruntime.GOOS, gruntime.GOARCH)
	}

	if opts.Verbose {
		config.WrapTransport = func(rt http.RoundTripper) http.RoundTripper {
			return &wireTracer{RoundTripper: rt}
		}
	}

	return rest.RESTClientFor(&config)
}
