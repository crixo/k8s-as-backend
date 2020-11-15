package k8sasbackend

import (
	"context"
	"fmt"
	"os"

	k8sasbackendv1alpha1 "github.com/crixo/k8s-as-backend/operator/pkg/apis/k8sasbackend/v1alpha1"
	authz "github.com/crixo/k8s-as-backend/operator/pkg/controller/k8sasbackend/authz"
	clusterdep "github.com/crixo/k8s-as-backend/operator/pkg/controller/k8sasbackend/clusterdep"
	common "github.com/crixo/k8s-as-backend/operator/pkg/controller/k8sasbackend/common"
	todoapp "github.com/crixo/k8s-as-backend/operator/pkg/controller/k8sasbackend/todoapp"
	webhookserver "github.com/crixo/k8s-as-backend/operator/pkg/controller/k8sasbackend/webhookserver"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	certv1beta1 "k8s.io/client-go/kubernetes/typed/certificates/v1beta1"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var (
	clusterDependencies *clusterdep.ClusterDependencies
	authorization       *authz.Authz
	webhookServer       *webhookserver.WebhookServer
	todoApp             *todoapp.TodoApp
	log                 logr.Logger = common.Log
	pemFolder                       = common.GetEnv("PEM_FOLDER", os.TempDir()) //pflag.String("pem-folder", "/tmp", "Folder where pem files will be stored during the container lifetime")
	//ingressHost                     = common.GetEnv("INGRESS_HOST", "")
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

	inClusterConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Info("ClientConfig", "InClusterConfig", err.Error())
		common.AppState.ClientConfig = mgr.GetConfig()
	} else {
		common.AppState.ClientConfig = inClusterConfig
	}

	if common.AppState.ClientConfig == nil {
		panic("Unable to load ClientConfig")
	}

	log.Info("ClientConfig", "ClientConfig.CAData", len(common.AppState.ClientConfig.CAData))
	log.Info("ClientConfig", "TLSClientConfig.CAFile", common.AppState.ClientConfig.TLSClientConfig.CAFile)

	reconciler := &ReconcileK8sAsBackend{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		//certClient: certClient.CertificateSigningRequests(),
	}

	clusterDependencies = clusterdep.NewClusterDependencies(reconciler.Client, reconciler.Scheme)

	authorization = authz.NewAuthz(mgr)

	log.Info("Reconciler configuration", "pemFolder", pemFolder)

	certCl, _ := certv1beta1.NewForConfig(mgr.GetConfig())
	webhookServer = webhookserver.NewWebhookServer(reconciler.Client,
		reconciler.Scheme,
		// path.Join(pemFolder, "server-cert.pem"),
		// path.Join(pemFolder, "server-key.pem"),
		certCl.CertificateSigningRequests(),
	)

	todoApp = todoapp.NewTodoApp(mgr)

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

	// TODO:
	// for cluster-wide resource w/ dependecies on primary respurce eg. ValidatingWebhookConfiguration
	// use EnqueueRequestsFromMapFunc but you have to track all main resources created by the oprator to set as "ToRequests"
	// err = c.Watch(&source.Kind{Type: &arv1beta1.ValidatingWebhookConfiguration{}}, &handler.EnqueueRequestsFromMapFunc{
	// 	ToRequests: handler.ToRequestsFunc(func(a handler.MapObject) []reconcile.Request {
	// 		return []reconcile.Request{
	// 			{
	// 				NamespacedName: types.NamespacedName{Namespace: "", Name: webhookserver.ValidationWebhookName},
	// 			},
	// 			// {
	// 			// 	NamespacedName: types.NamespacedName{Namespace: "biz", Name: "baz"},
	// 			// },
	// 		}
	// 	}),
	// })
	// if err != nil {
	// 	return err
	// }

	watchedObjects := []runtime.Object{}
	watchedObjects = append(watchedObjects, authorization.GetWatchedResources()...)
	watchedObjects = append(watchedObjects, webhookServer.GetWatchedResources()...)
	watchedObjects = append(watchedObjects, todoApp.GetWatchedResources()...)
	log.Info("Watching", "watchedObjects", len(watchedObjects))
	watchedObjects = removeDuplicates(watchedObjects)
	log.Info("Watching after remove duplicates", "watchedObjects", len(watchedObjects))
	for _, obj := range watchedObjects {
		log.Info("Watching secondary resources", "resource", fmt.Sprintf("%T", obj))
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
	reqLogger := log.WithValues(
		"Request.NamespacedName", request.NamespacedName,
		"Request.Namespace", request.Namespace,
		"Request.Name", request.Name)
	reqLogger.Info("Reconciling K8sAsBackend")

	// TODO: cluster-wide resources limitation
	// cluster-wide resources are watched
	// cluster-wide resources are garbage-collected due to resource factory hack forcing primary resource ns
	// cluster-wide resources cannot be reconcile because the request created by its change
	//   does not contain the ns needed to load the primary resource
	// if you create a cluster-wide resources w/o the primary resource that resource cannot be owned by the primary resource
	// request example for cluster-wide resource owned by primary resource thruogh res factory hack
	// "msg":"Reconciling K8sAsBackend","Request.NamespacedName":"/example-k8sasbackend","Request.Namespace":"","Request.Name":"example-k8sasbackend"}
	// HACK -> force Request.NamespacedName used to load primary resource getting the ns value from a global settings
	//   or from a previous reconcile request strored within the state of *ReconcileK8sAsBackend
	result, err := clusterDependencies.Reconcile()
	if result != nil {
		return *result, err
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Fetch the K8sAsBackend instance
	instance := &k8sasbackendv1alpha1.K8sAsBackend{}
	err = r.Client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("main crd not found")
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	result, err = authorization.Reconcile(instance)
	if result != nil {
		return *result, err
	} else if err != nil {
		return reconcile.Result{}, err
	}

	result, err = webhookServer.Reconcile(instance)
	if result != nil {
		return *result, err
	} else if err != nil {
		return reconcile.Result{}, err
	}

	result, err = todoApp.Reconcile(instance)
	if result != nil {
		return *result, err
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// == Finish ==========
	// Everything went fine, don't requeue
	log.Info("Reconcile says Everything went fine, don't requeue")
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
