package authz

import (
	k8sasbackendv1alpha1 "github.com/crixo/k8s-as-backend/operator/pkg/apis/k8sasbackend/v1alpha1"
	"github.com/crixo/k8s-as-backend/operator/pkg/controller/k8sasbackend/common"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

//TODO: make secondary resource names unique
var (
	BaseName string      = "todo-crd"
	log      logr.Logger = common.Log
)

type Authz struct {
	Manager manager.Manager
}

func NewAuthz(mgr manager.Manager) *Authz {
	return &Authz{
		Manager: mgr,
	}
}

func (a Authz) GetWatchedResources() []runtime.Object {
	return []runtime.Object{
		&corev1.ServiceAccount{},
		&rbac.RoleBinding{},
	}
}

func (a *Authz) Reconcile(i *k8sasbackendv1alpha1.K8sAsBackend) (*reconcile.Result, error) {

	res, err := a.ensureAccount(i)
	if common.PrepareComponentResult(res, err) {
		return res, err
	}

	res, err = a.ensureRoleBinding(i)
	if common.PrepareComponentResult(res, err) {
		return res, err
	}

	return nil, nil
}
