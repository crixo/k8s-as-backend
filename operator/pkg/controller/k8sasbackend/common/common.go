package common

import (
	"context"
	"fmt"
	"os"
	"unicode/utf8"

	k8sasbackendv1alpha1 "github.com/crixo/k8s-as-backend/operator/pkg/apis/k8sasbackend/v1alpha1"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"k8s.io/client-go/rest"
)

var (
	Log                  = logf.Log.WithName("controller_k8sasbackend")
	AppState             = &State{}
	InformerImage string = "crixo/k8s-as-backend-informer"
)

//type UniqueSecondaryResourceNameFactory func(crName string) string
func CreateUniqueSecondaryResourceName(i *k8sasbackendv1alpha1.K8sAsBackend, baseName string) string {
	return fmt.Sprintf("%s-%s", i.Name, baseName)
}

func CreateMatchingLabels(i *k8sasbackendv1alpha1.K8sAsBackend, baseName string) map[string]string {
	return map[string]string{
		"app": CreateUniqueSecondaryResourceName(i, baseName),
	}
}

func CreateImageName(i *k8sasbackendv1alpha1.K8sAsBackend, imageName string) string {
	return fmt.Sprintf("%s:%s", imageName, i.Spec.ProductVersion)
}

func GetKind(obj runtime.Object) string {
	return fmt.Sprintf("%T", obj)
}

func CreateMeta(name string, namespace string) metav1.ObjectMeta {
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

func FileNotExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return true
	}
	return false
}

func CreateSecret(nsName types.NamespacedName, data map[string][]byte) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: CreateMeta(nsName.Name, nsName.Namespace),
		Type:       "Opaque",
		Data:       data,
	}
}

type State struct {
	ClientConfig *rest.Config
}

type ResourceUtils struct {
	Scheme          *runtime.Scheme
	Client          client.Client
	ResourceFactory func(namespacedName types.NamespacedName, i *k8sasbackendv1alpha1.K8sAsBackend) runtime.Object
}

func EnsureResource(found runtime.Object,
	resourceNamespacedName types.NamespacedName,
	i *k8sasbackendv1alpha1.K8sAsBackend,
	resUtils *ResourceUtils,
) error {
	// if resourceNamespacedName.Namespace != "" && i.Namespace != resourceNamespacedName.Namespace {
	// 	panic(fmt.Sprintf(
	// 		"resource ns %s must match with primary resource ns %s",
	// 		resourceNamespacedName.Namespace, i.Namespace))
	// }
	nsn := resourceNamespacedName
	err := resUtils.Client.Get(context.TODO(), nsn, found)
	Log.Info(fmt.Sprintf("Getting %T", found))
	if err != nil && errors.IsNotFound(err) {
		resource := resUtils.ResourceFactory(nsn, i)
		// object must be the same instance used for the Client.Create
		// this call must occur earlier then Client.Create
		metaObj := resource.(metav1.Object)
		// handle cluster-wide resource such as ValidatingWebhookConfiguration and CRD
		if i != nil {
			if err := controllerutil.SetControllerReference(i, metaObj, resUtils.Scheme); err != nil {
				panic(err)
			}
		}
		//}
		err = resUtils.Client.Create(context.TODO(), resource)
		if err != nil {
			// Resource Creation failed
			Log.Error(err, fmt.Sprintf("Failed to create new %T", found), "Namespace", nsn.Namespace, "Name", nsn.Name)
			return err
		} else {
			// Resource Creation was successful
			Log.Info(fmt.Sprintf("Creation succeded for new %T", found))
			return nil
		}
	} else if err != nil {
		// Error that isn't due to the resource not existing
		//log.Error(err, fmt.Sprintf("Failed to get %s", common.GetKind(found)))
		Log.Error(err, fmt.Sprintf("Failed to get %T", found))
		return err
	}

	return nil
}

func UpdateResource(client client.Client, found runtime.Object) (*reconcile.Result, error) {
	metaAccessor := meta.NewAccessor()
	kind, err := metaAccessor.Kind(found)
	ns, err := metaAccessor.Namespace(found)
	name, err := metaAccessor.Name(found)
	if err != nil {
		panic(err)
	}
	err = client.Update(context.TODO(), found)
	if err != nil {
		Log.Error(err, fmt.Sprintf("Failed to update existing %s", kind), "Namespace", ns, "Name", name)
		return &reconcile.Result{}, err
	}
	// Spec updated - return and requeue
	return &reconcile.Result{Requeue: true}, nil
}

// Component must return nil if all good
func PrepareComponentResult(res *reconcile.Result, err error) (mustReturn bool) {
	if err != nil {
		return true
	} else if res != nil {
		return true
	}
	return false
}

func TrimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

func GetEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

//type resourceFactory func(name string, instance *k8sasbackendv1alpha1.K8sAsBackend) runtime.Object

// func (r *ReconcileK8sAsBackend) ensureResource(request reconcile.Request,
// 	instance *k8sasbackendv1alpha1.K8sAsBackend,
// 	resourceNamespacedName types.NamespacedName,
// 	resourceFactory ResourceFactory,
// ) (*reconcile.Result, error) {

// 	found := resourceFactory.createEmpty()

// 	name := resourceNamespacedName.Name
// 	ns := resourceNamespacedName.Namespace
// 	kind := fmt.Sprintf("%T", found) //metaAccessor.Kind(found)

// 	// See if deployment already exists and create if it doesn't

// 	err := r.client.Get(context.TODO(), resourceNamespacedName, found)
// 	//log.Error(err, fmt.Sprintf("[check]Failed to get %s", kind))
// 	if err != nil && errors.IsNotFound(err) {

// 		resource := resourceFactory.create(name, instance)
// 		metaObject, _ := meta.Accessor(resource)

// 		if metaObject.GetNamespace() != "" {
// 			// not tracking cluster-wide resources (eg. crd, csr)
// 			controllerutil.SetControllerReference(instance, metaObject, r.scheme)
// 		}

// 		// Create the resource
// 		log.Info(fmt.Sprintf("Creating a new %s", kind),
// 			fmt.Sprintf("%s.Namespace", kind), ns,
// 			fmt.Sprintf("%s.Name", kind), name)
// 		err = r.client.Create(context.TODO(), resource)

// 		if err != nil {
// 			// resource creation failed
// 			log.Error(err, fmt.Sprintf("Failed to create new %s", kind), fmt.Sprintf("%s.Namespace", kind), ns, fmt.Sprintf("%s.Name", kind), name)
// 			return &reconcile.Result{}, err
// 		} else {
// 			// resource creation was successful
// 			return nil, nil
// 		}
// 	} else if err != nil {
// 		// Error that isn't due to the resource not existing
// 		log.Error(err, fmt.Sprintf("Failed to get %s", kind))
// 		return &reconcile.Result{}, err
// 	}

// 	return nil, nil
// }

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
