package k8sasbackend

import (
	"context"

	k8sasbackendv1alpha1 "github.com/crixo/k8s-as-backend/operator/pkg/apis/k8sasbackend/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	certv1beta1 "k8s.io/client-go/kubernetes/typed/certificates/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_k8sasbackend")

var ownedResources = []ResourceFactory{}

// func init() {
// 	log.Info("init controller")
// 	registerCrd()
// }

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
	certClient, _ := certv1beta1.NewForConfig(mgr.GetConfig())
	return &ReconcileK8sAsBackend{
		client:     mgr.GetClient(),
		scheme:     mgr.GetScheme(),
		certClient: certClient.CertificateSigningRequests(),
	}
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

	crdFactory := &CrdFactory{}
	crdFactory.AddToScheme() // TODO try init method
	ownedResources = []ResourceFactory{
		&RoleFactory{},
		crdFactory,
		&AccountFactory{},
		&CertFactory{},
		&SecretFactory{},
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// TODO CRD
	for _, r := range ownedResources {
		err = c.Watch(&source.Kind{Type: r.createEmpty()}, &handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &k8sasbackendv1alpha1.K8sAsBackend{},
		})
		if err != nil {
			return err
		}
	}

	//apiVersion: apps/v1
	// kind: Deployment

	//apiVersion: v1
	//kind: Service

	//apiVersion: extensions/v1beta1
	//kind: Ingress

	//apiVersion: v1
	//kind: ServiceAccount

	//kind: ClusterRole
	//apiVersion: rbac.authorization.k8s.io/v1beta1
	// err = c.Watch(&source.Kind{Type: &admissionregistrationv1beta1.ValidatingWebhookConfiguration{}}, &handler.EnqueueRequestForOwner{
	// 	IsController: true,
	// 	OwnerType:    &k8sasbackendv1alpha1.K8sAsBackend{},
	// })
	// if err != nil {
	// 	return err
	// }

	//kind: ClusterRoleBinding
	//apiVersion: rbac.authorization.k8s.io/v1beta1

	// apiVersion: admissionregistration.k8s.io/v1beta1
	//kind: ValidatingWebhookConfiguration
	// err = c.Watch(&source.Kind{Type: &admissionregistrationv1beta1.ValidatingWebhookConfiguration{}}, &handler.EnqueueRequestForOwner{
	// 	IsController: true,
	// 	OwnerType:    &k8sasbackendv1alpha1.K8sAsBackend{},
	// })
	// if err != nil {
	// 	return err
	// }

	// // apiVersion: admissionregistration.k8s.io/v1
	// //kind: ValidatingWebhookConfiguration
	// err = c.Watch(&source.Kind{Type: &admissionregistrationv1.ValidatingWebhookConfiguration{}}, &handler.EnqueueRequestForOwner{
	// 	IsController: true,
	// 	OwnerType:    &k8sasbackendv1alpha1.K8sAsBackend{},
	// })
	// if err != nil {
	// 	return err
	// }

	return nil
}

// blank assignment to verify that ReconcileK8sAsBackend implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileK8sAsBackend{}

// ReconcileK8sAsBackend reconciles a K8sAsBackend object
type ReconcileK8sAsBackend struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client     client.Client
	scheme     *runtime.Scheme
	certClient certv1beta1.CertificateSigningRequestInterface
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
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
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

	for _, resourceFactory := range ownedResources {
		//for _, name := range resourceFactory.getNames() {
		result, err := resourceFactory.ensure(r, request, instance)
		if result != nil {
			return *result, err
		}
		//}
	}

	return reconcile.Result{}, nil
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

// newPodForCR returns a busybox pod with the same name/namespace as the cr
// func newPodForCR(cr *k8sasbackendv1alpha1.K8sAsBackend) *corev1.Pod {
// 	labels := map[string]string{
// 		"app": cr.Name,
// 	}
// 	return &corev1.Pod{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      cr.Name + "-pod",
// 			Namespace: cr.Namespace,
// 			Labels:    labels,
// 		},
// 		Spec: corev1.PodSpec{
// 			Containers: []corev1.Container{
// 				{
// 					Name:    "busybox",
// 					Image:   "busybox",
// 					Command: []string{"sleep", "3600"},
// 				},
// 			},
// 		},
// 	}
// }
