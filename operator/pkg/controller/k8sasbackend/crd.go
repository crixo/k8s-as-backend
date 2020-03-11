package k8sasbackend

import (
	k8sasbackendv1alpha1 "github.com/crixo/k8s-as-backend/operator/pkg/apis/k8sasbackend/v1alpha1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type CrdFactory struct{}

func (f CrdFactory) AddToScheme() {
	s := k8sscheme.Scheme
	apiextensionsv1beta1.AddToScheme(s)
}

func (f CrdFactory) createEmpty() runtime.Object {
	return &apiextensionsv1beta1.CustomResourceDefinition{}
}

func (f CrdFactory) getNames() []string {
	return []string{"todos.k8sasbackend.com"}
}

func (f CrdFactory) ensure(r *ReconcileK8sAsBackend,
	request reconcile.Request,
	i *k8sasbackendv1alpha1.K8sAsBackend) (*reconcile.Result, error) {
	return r.ensureResource(request, i, createNamespacedName(f.getNames()[0], ""), f)
}

func (f CrdFactory) create(name string, i *k8sasbackendv1alpha1.K8sAsBackend) runtime.Object {
	return &apiextensionsv1beta1.CustomResourceDefinition{
		ObjectMeta: createMeta(name, ""),
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   "k8sasbackend.com",
			Version: "v1",
			Scope:   apiextensionsv1beta1.NamespaceScoped,
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Plural: "todos",
				Kind:   "Todo",
			},
		},
	}
}
