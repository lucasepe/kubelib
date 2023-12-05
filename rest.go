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

// APIPath set a sub-path that points to
// an API root (optional, default '/apis')
func APIPath(apiPath string) CreateRESTClientOption {
	return func(rc *rest.Config) {
		rc.APIPath = apiPath
	}
}

// GroupVersion is the API version to talk to (required)
func GroupVersion(gv schema.GroupVersion) CreateRESTClientOption {
	return func(rc *rest.Config) {
		rc.GroupVersion = &gv
	}
}

// Verbose if true adds a http wire dump
// transport middleware (optional)
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

// UserAgent specifies the caller of this
// request (optional)
func UserAgent(ua string) CreateRESTClientOption {
	return func(rc *rest.Config) {
		rc.UserAgent = ua
	}
}

// CreateRESTClient returns a rest.RESTClient that satisfies
// the requested options on a rest.Config
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
