package kubernetes

import (
	"context"

	"github.com/rancher/opni/pkg/ident"
	"github.com/rancher/opni/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type kubernetesProvider struct {
	KubernetesIdentOptions
	clientset *kubernetes.Clientset
}

type KubernetesIdentOptions struct {
	restConfig *rest.Config
}

type KubernetesIdentOption func(*KubernetesIdentOptions)

func (o *KubernetesIdentOptions) apply(opts ...KubernetesIdentOption) {
	for _, op := range opts {
		op(o)
	}
}

func WithRestConfig(rc *rest.Config) KubernetesIdentOption {
	return func(o *KubernetesIdentOptions) {
		o.restConfig = rc
	}
}

func NewKubernetesProvider(opts ...KubernetesIdentOption) ident.Provider {
	options := KubernetesIdentOptions{}
	options.apply(opts...)
	if options.restConfig == nil {
		options.restConfig = util.Must(rest.InClusterConfig())
	}
	cs := kubernetes.NewForConfigOrDie(options.restConfig)
	return &kubernetesProvider{
		clientset: cs,
	}
}

func (p *kubernetesProvider) UniqueIdentifier(ctx context.Context) (string, error) {
	ns, err := p.clientset.CoreV1().
		Namespaces().
		Get(ctx, "kube-system", metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return string(ns.GetUID()), nil
}

func init() {
	util.Must(ident.RegisterProvider("kubernetes", func() ident.Provider {
		return NewKubernetesProvider()
	}))
}
