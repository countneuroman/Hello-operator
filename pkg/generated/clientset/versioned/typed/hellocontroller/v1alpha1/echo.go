// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	context "context"

	hellocontrollerv1alpha1 "github.com/countneuroman/hello-operator/pkg/apis/hellocontroller/v1alpha1"
	scheme "github.com/countneuroman/hello-operator/pkg/generated/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	gentype "k8s.io/client-go/gentype"
)

// EchosGetter has a method to return a EchoInterface.
// A group's client should implement this interface.
type EchosGetter interface {
	Echos(namespace string) EchoInterface
}

// EchoInterface has methods to work with Echo resources.
type EchoInterface interface {
	Create(ctx context.Context, echo *hellocontrollerv1alpha1.Echo, opts v1.CreateOptions) (*hellocontrollerv1alpha1.Echo, error)
	Update(ctx context.Context, echo *hellocontrollerv1alpha1.Echo, opts v1.UpdateOptions) (*hellocontrollerv1alpha1.Echo, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*hellocontrollerv1alpha1.Echo, error)
	List(ctx context.Context, opts v1.ListOptions) (*hellocontrollerv1alpha1.EchoList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *hellocontrollerv1alpha1.Echo, err error)
	EchoExpansion
}

// echos implements EchoInterface
type echos struct {
	*gentype.ClientWithList[*hellocontrollerv1alpha1.Echo, *hellocontrollerv1alpha1.EchoList]
}

// newEchos returns a Echos
func newEchos(c *HelloV1alpha1Client, namespace string) *echos {
	return &echos{
		gentype.NewClientWithList[*hellocontrollerv1alpha1.Echo, *hellocontrollerv1alpha1.EchoList](
			"echos",
			c.RESTClient(),
			scheme.ParameterCodec,
			namespace,
			func() *hellocontrollerv1alpha1.Echo { return &hellocontrollerv1alpha1.Echo{} },
			func() *hellocontrollerv1alpha1.EchoList { return &hellocontrollerv1alpha1.EchoList{} },
		),
	}
}
