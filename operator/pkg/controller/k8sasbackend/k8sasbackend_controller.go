package k8sasbackend

import (
	"context"
	"fmt"

	k8sasbackendv1alpha1 "github.com/crixo/k8s-as-backend/operator/pkg/apis/k8sasbackend/v1alpha1"
	common "github.com/crixo/k8s-as-backend/operator/pkg/controller/k8sasbackend/common"
	webhookserver "github.com/crixo/k8s-as-backend/operator/pkg/controller/k8sasbackend/webhookserver"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	certv1beta1 "k8s.io/client-go/kubernetes/typed/certificates/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var (
	webhookServer *webhookserver.WebhookServer
	log           logr.Logger = common.Log
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new K8sAsBackend Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {

	reconciler := &ReconcileK8sAsBackend{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		//certClient: certClient.CertificateSigningRequests(),
	}

	certCl, _ := certv1beta1.NewForConfig(mgr.GetConfig())
	webhookServer = &webhookserver.WebhookServer{
		CerFilePath: "/Users/cristiano/Coding/golang/k8s-as-backend/operator/certs/server-cert.pem",
		KeyFilePath: "/Users/cristiano/Coding/golang/k8s-as-backend/operator/certs/server-key.pem",
		Client:      reconciler.Client,
		Scheme:      reconciler.Scheme,
		CertClient:  certCl.CertificateSigningRequests(),
	}

	return reconciler
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("k8sasbackend-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource K8sAsBackend
	err = c.Watch(&source.Kind{Type: &k8sasbackendv1alpha1.K8sAsBackend{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	watchedObjects := []runtime.Object{}
	watchedObjects = append(watchedObjects,
		webhookServer.GetWatchedResources()...)
	log.Info("Watching", "watchedObjects", len(watchedObjects))
	watchedObjects = removeDuplicates(watchedObjects)
	log.Info("Watching after remove duplicates", "watchedObjects", len(watchedObjects))
	for _, obj := range watchedObjects {
		log.Info("Watching", "resource", fmt.Sprintf("%T", obj))
		err = c.Watch(&source.Kind{Type: obj}, &handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &k8sasbackendv1alpha1.K8sAsBackend{},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// blank assignment to verify that ReconcileK8sAsBackend implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileK8sAsBackend{}

// ReconcileK8sAsBackend reconciles a K8sAsBackend object
type ReconcileK8sAsBackend struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	Client client.Client
	Scheme *runtime.Scheme
	//certClient certv1beta1.CertificateSigningRequestInterface
}

// Reconcile reads that state of the cluster for a K8sAsBackend object and makes changes based on the state read
// and what is in the K8sAsBackend.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileK8sAsBackend) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling K8sAsBackend")

	// Fetch the K8sAsBackend instance
	instance := &k8sasbackendv1alpha1.K8sAsBackend{}
	err := r.Client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	result, err := webhookServer.Reconcile(instance)
	if result != nil {
		return *result, err
	}

	// nsName := types.NamespacedName{
	// 	Name:      "admission-webhook-example-certs",
	// 	Namespace: instance.Namespace,
	// }
	// // sec := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{
	// // 	Name:      nsName.Name,
	// // 	Namespace: nsName.Namespace,
	// // }}
	// // if err := controllerutil.SetControllerReference(instance, sec, r.Scheme); err != nil {
	// // 	panic(err)
	// // }
	// cert, _ := ioutil.ReadFile(webhookServer.CerFilePath)
	// key, _ := ioutil.ReadFile(webhookServer.KeyFilePath)
	// sec := &corev1.Secret{
	// 	ObjectMeta: common.CreateMeta(nsName.Name, nsName.Namespace),
	// 	Type:       "Opaque",

	// 	Data: map[string][]byte{
	// 		"key.pem":  cert,
	// 		"cert.pem": key,
	// 	},
	// }
	// if err := controllerutil.SetControllerReference(instance, sec, r.Scheme); err != nil {
	// 	panic(err)
	// }

	// secret := &corev1.Secret{
	// 	ObjectMeta: metav1.ObjectMeta{
	// 		Name:      "my-secret",
	// 		Namespace: instance.Namespace,
	// 	},
	// 	Type: "Opaque",

	// 	StringData: map[string]string{
	// 		"key.pem":  "cert",
	// 		"cert.pem": "key",
	// 	},
	// }
	// // if err := controllerutil.SetControllerReference(instance, secret, r.Scheme); err != nil {
	// // 	log.Error(err, "secret SetControllerReference")
	// // 	return reconcile.Result{}, err
	// // }
	// controllerutil.SetControllerReference(instance, secret, r.Scheme)
	// founds := &corev1.Secret{}
	// err = r.Client.Get(context.TODO(), types.NamespacedName{Name: secret.Name, Namespace: secret.Namespace}, founds)
	// if err != nil && errors.IsNotFound(err) {
	// 	reqLogger.Info("Creating a new Secret", "Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)
	// 	err = r.Client.Create(context.TODO(), secret)

	// 	//
	// 	if err != nil {
	// 		return reconcile.Result{}, err
	// 	}
	// }

	// // // Define a new Pod object
	// pod := newPodForCR(instance)
	// // Set K8sAsBackend instance as the owner and controller
	// if err := controllerutil.SetControllerReference(instance, pod, r.Scheme); err != nil {
	// 	return reconcile.Result{}, err
	// }
	// // Check if this Pod already exists
	// found := &corev1.Pod{}
	// err = r.Client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	// if err != nil && errors.IsNotFound(err) {
	// 	reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
	// 	err = r.Client.Create(context.TODO(), pod)
	// 	//controllerutil.SetControllerReference(instance, pod, r.Scheme)
	// 	if err != nil {
	// 		return reconcile.Result{}, err
	// 	}
	// }

	// == Finish ==========
	// Everything went fine, don't requeue
	return reconcile.Result{}, nil
}

func removeDuplicates(elements []runtime.Object) []runtime.Object {
	// Use map to record duplicates as we find them.
	encountered := map[runtime.Object]bool{}
	result := []runtime.Object{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

// // Define a new Pod object
// pod := newPodForCR(instance)

// // Set K8sAsBackend instance as the owner and controller
// if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
// 	return reconcile.Result{}, err
// }

// // Check if this Pod already exists
// found := &corev1.Pod{}
// err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
// if err != nil && errors.IsNotFound(err) {
// 	reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
// 	err = r.client.Create(context.TODO(), pod)
// 	if err != nil {
// 		return reconcile.Result{}, err
// 	}

// 	// Pod created successfully - don't requeue
// 	return reconcile.Result{}, nil
// } else if err != nil {
// 	return reconcile.Result{}, err
// }

// // Pod already exists - don't requeue
// reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)

//newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *k8sasbackendv1alpha1.K8sAsBackend) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}
