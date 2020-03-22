package todoapp

import (
	"fmt"

	k8sasbackendv1alpha1 "github.com/crixo/k8s-as-backend/operator/pkg/apis/k8sasbackend/v1alpha1"
	"github.com/crixo/k8s-as-backend/operator/pkg/controller/k8sasbackend/common"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

//TODO: make secondary resource names unique
var (
	todoAppUrlSegmentIdentifier string      = "todo-app"
	BaseName                    string      = "todo-app"
	todoAppImage                string      = "crixo/k8s-as-backend-todo-app"
	kubectlImage                string      = "bitnami/kubectl:1.16"
	SvcPort                     int         = 5000
	todoPort                    int         = 80
	kubectlApiPort              int32       = 8080
	log                         logr.Logger = common.Log
	//caBundle []byte      = common.AppState.ClientConfig.CAData
)

// Users/cristiano/Coding/golang/k8s-as-backend/operator/certs/server-cert.pem
// Users/cristiano/Coding/golang/k8s-as-backend/operator/certs/server-key.pem
type TodoApp struct {
	Manager manager.Manager
}

func NewTodoApp(mgr manager.Manager) *TodoApp {
	return &TodoApp{
		Manager: mgr,
	}
}

func (t TodoApp) GetWatchedResources() []runtime.Object {
	return []runtime.Object{
		&appsv1.Deployment{},
		&corev1.Service{},
		&extv1beta1.Ingress{},
	}
}

func (t *TodoApp) Reconcile(i *k8sasbackendv1alpha1.K8sAsBackend) (*reconcile.Result, error) {

	res, err := t.ensureDeployment(i)
	if common.PrepareComponentResult(res, err) {
		return res, err
	}

	res, err = t.ensureService(i)
	if common.PrepareComponentResult(res, err) {
		return res, err
	}

	res, err = t.ensureIngress(i)
	if common.PrepareComponentResult(res, err) {
		return res, err
	}

	return nil, nil
}

func getAppBaseUrl(i *k8sasbackendv1alpha1.K8sAsBackend) string {
	return fmt.Sprintf("/%s/%s/%s", i.Namespace, i.Name, todoAppUrlSegmentIdentifier)
}
