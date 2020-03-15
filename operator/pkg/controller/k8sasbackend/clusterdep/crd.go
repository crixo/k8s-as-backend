package clusterdep

import (
	k8sasbackendv1alpha1 "github.com/crixo/k8s-as-backend/operator/pkg/apis/k8sasbackend/v1alpha1"
	"github.com/crixo/k8s-as-backend/operator/pkg/controller/k8sasbackend/common"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (cd ClusterDependencies) ensureTodosCrd() (*reconcile.Result, error) {
	resUtils := &common.ResourceUtils{
		Scheme:          cd.Scheme,
		Client:          cd.Client,
		ResourceFactory: createTodosCrd,
	}
	nsn := types.NamespacedName{Name: TodosCrdName, Namespace: ""}
	found := &apiextensionsv1beta1.CustomResourceDefinition{}
	return nil, common.EnsureResource(found, nsn, nil, resUtils)
}

func createTodosCrd(nsName types.NamespacedName, i *k8sasbackendv1alpha1.K8sAsBackend) runtime.Object {
	return &apiextensionsv1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: nsName.Name,
		},
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
