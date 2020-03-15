package clusterdep

import (
	"github.com/crixo/k8s-as-backend/operator/pkg/controller/k8sasbackend/common"
	"github.com/go-logr/logr"
	arv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	k8sscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	ValidationWebhookConfigurationName string      = "admission-webhook-example-validation-webhook"
	TodosCrdName                       string      = "todos.k8sasbackend.com"
	log                                logr.Logger = common.Log
	//caBundle []byte      = common.AppState.ClientConfig.CAData
)

// Users/cristiano/Coding/golang/k8s-as-backend/operator/certs/server-cert.pem
// Users/cristiano/Coding/golang/k8s-as-backend/operator/certs/server-key.pem
type ClusterDependencies struct {
	Client client.Client
	Scheme *runtime.Scheme
}

func NewClusterDependencies(cl client.Client, s *runtime.Scheme) *ClusterDependencies {
	apiextensionsv1beta1.AddToScheme(k8sscheme.Scheme)
	return &ClusterDependencies{
		Client: cl,
		Scheme: s,
	}
}

func (cd ClusterDependencies) GetWatchedResources() []runtime.Object {
	return []runtime.Object{
		&arv1beta1.ValidatingWebhookConfiguration{},
	}
}

func (cd *ClusterDependencies) Reconcile() (*reconcile.Result, error) {

	res, err := cd.ensureValidationWebhookConfiguration()
	if common.PrepareComponentResult(res, err) {
		return res, err
	}

	res, err = cd.ensureTodosCrd()
	if common.PrepareComponentResult(res, err) {
		return res, err
	}

	return nil, nil
}
