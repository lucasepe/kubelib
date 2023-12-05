package kubelib

import (
	"fmt"
	"net/http"

	gruntime "runtime"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type CreateRESTClientOption func(*rest.Config)

func APIPath(apiPath string) CreateRESTClientOption {
	return func(rc *rest.Config) {
		rc.APIPath = apiPath
	}
}

func GroupVersion(gv schema.GroupVersion) CreateRESTClientOption {
	return func(rc *rest.Config) {
		rc.GroupVersion = &gv
	}
}

func Verbose(verbose bool) CreateRESTClientOption {
	return func(rc *rest.Config) {
		if !verbose {
			return
		}

		rc.Wrap(func(rt http.RoundTripper) http.RoundTripper {
			return &wireTracer{RoundTripper: rt}
		})
	}
}

func UserAgent(ua string) CreateRESTClientOption {
	return func(rc *rest.Config) {
		rc.UserAgent = ua
	}
}

func CreateRESTClient(rc *rest.Config, opts ...CreateRESTClientOption) (*rest.RESTClient, error) {
	config := *rc
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	for _, fn := range opts {
		fn(&config)
	}

	if len(config.APIPath) == 0 {
		config.APIPath = "/apis"
	}

	if len(config.UserAgent) == 0 {
		config.UserAgent = fmt.Sprintf("kubelib (%s/%s)", gruntime.GOOS, gruntime.GOARCH)
	}

	return rest.RESTClientFor(&config)
}
