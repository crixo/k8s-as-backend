package k8sasbackend

import (
	"context"
	"fmt"

	k8sasbackendv1alpha1 "github.com/crixo/k8s-as-backend/operator/pkg/apis/k8sasbackend/v1alpha1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type ResourceFactory interface {
	createEmpty() runtime.Object
	ensure(r *ReconcileK8sAsBackend, request reconcile.Request, i *k8sasbackendv1alpha1.K8sAsBackend) (*reconcile.Result, error)
	create(name string, i *k8sasbackendv1alpha1.K8sAsBackend) runtime.Object
	getNames() []string
}

//type resourceFactory func(name string, instance *k8sasbackendv1alpha1.K8sAsBackend) runtime.Object

func (r *ReconcileK8sAsBackend) ensureResource(request reconcile.Request,
	instance *k8sasbackendv1alpha1.K8sAsBackend,
	resourceNamespacedName types.NamespacedName,
	resourceFactory ResourceFactory,
) (*reconcile.Result, error) {

	found := resourceFactory.createEmpty()

	name := resourceNamespacedName.Name
	ns := resourceNamespacedName.Namespace
	kind := fmt.Sprintf("%T", found) //metaAccessor.Kind(found)

	// See if deployment already exists and create if it doesn't

	err := r.client.Get(context.TODO(), resourceNamespacedName, found)
	//log.Error(err, fmt.Sprintf("[check]Failed to get %s", kind))
	if err != nil && errors.IsNotFound(err) {

		resource := resourceFactory.create(name, instance)
		metaObject, _ := meta.Accessor(resource)

		if metaObject.GetNamespace() != "" {
			// not tracking cluster-wide resources (eg. crd, csr)
			controllerutil.SetControllerReference(instance, metaObject, r.scheme)
		}

		// Create the resource
		log.Info(fmt.Sprintf("Creating a new %s", kind),
			fmt.Sprintf("%s.Namespace", kind), ns,
			fmt.Sprintf("%s.Name", kind), name)
		err = r.client.Create(context.TODO(), resource)

		if err != nil {
			// resource creation failed
			log.Error(err, fmt.Sprintf("Failed to create new %s", kind), fmt.Sprintf("%s.Namespace", kind), ns, fmt.Sprintf("%s.Name", kind), name)
			return &reconcile.Result{}, err
		} else {
			// resource creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the resource not existing
		log.Error(err, fmt.Sprintf("Failed to get %s", kind))
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func createMeta(name string, namespace string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
	}
}

func createNamespacedName(name string, namespace string) types.NamespacedName {
	return types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}
}

// type resourceAccessor struct{}

// func (resourceAccessor) Namespace(obj runtime.Object) (string, error) {
// 	accessor, err := Accessor(obj)
// 	if err != nil {
// 		return "", err
// 	}
// 	return accessor.GetNamespace(), nil
// }

// func (resourceAccessor) SetNamespace(obj runtime.Object, namespace string) error {
// 	accessor, err := Accessor(obj)
// 	if err != nil {
// 		return err
// 	}
// 	accessor.SetNamespace(namespace)
// 	return nil
// }

// func (resourceAccessor) Name(obj runtime.Object) (string, error) {
// 	accessor, err := Accessor(obj)
// 	if err != nil {
// 		return "", err
// 	}
// 	return accessor.GetName(), nil
// }

// func Accessor(obj interface{}) (metav1.Object, error) {
// 	switch t := obj.(type) {
// 	case metav1.Object:
// 		return t, nil
// 	case metav1.ObjectMetaAccessor:
// 		if m := t.GetObjectMeta(); m != nil {
// 			return m, nil
// 		}
// 		return nil, errNotObject
// 	default:
// 		return nil, errNotObject
// 	}
// }
