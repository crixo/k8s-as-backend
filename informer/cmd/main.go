package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	todov1 "github.com/crixo/k8s-as-backend/library/pkg/apis/k8sasbackend/v1"
	clientset "github.com/crixo/k8s-as-backend/library/pkg/client/clientset/versioned"
	todoInformers "github.com/crixo/k8s-as-backend/library/pkg/client/informers/externalversions"

	//"github.com/golang/glog"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/tools/reference"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/kubernetes/pkg/proxy/apis/config/scheme"
)

var (
	logger = zap.NewExample()
)

func main() {
	defer runtime.HandleCrash()
	logger.Info("The todo operator started.")
	stopCh := make(chan struct{})
	defer close(stopCh)
	//queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	//defer queue.ShutDown()
	var kubeconfig *string

	usr, err := user.Current()
	if err != nil {
		panic(err.Error())
	}
	if home := usr.HomeDir; home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	var config *rest.Config
	config, err = rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			logger.Fatal(err.Error())
		}
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.With(zap.Error(err)).Fatal("Error building kubernetes clientset")
	}

	todoClient, err := clientset.NewForConfig(config)
	if err != nil {
		logger.With(zap.Error(err)).Fatal("Error building clientset")
	}

	defaultResync := getEnvAsInt("INFORMER_RESYNC", 0)
	factory := todoInformers.NewSharedInformerFactory(todoClient, (time.Second*time.Duration(defaultResync))) // 0 disable resync // time.Second*30
	todoInformer := factory.K8sasbackend().V1().Todos().Informer()
	todoInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				//queue.Add(key)
				businessLogicAsync(key, "add", todoClient, kubeClient)
			}
		},
		UpdateFunc: func(old interface{}, neww interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(neww)
			if err == nil {
				//queue.Add(key)
				businessLogicAsync(key, "update", todoClient, kubeClient)
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				//queue.Add(key)
				businessLogicAsync(key, "delete", todoClient, kubeClient)
			}
		},
	})
	go todoInformer.Run(stopCh)
	if !cache.WaitForCacheSync(stopCh, todoInformer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync."))
		logger.Fatal("stop execution. Bye.")
	}
	// for processNextItem(queue) {
	// }
	<-stopCh
}

func processNextItem(queue workqueue.RateLimitingInterface) bool {
	// Wait until there is a new item in the working queue
	key, quit := queue.Get()
	if quit {
		return false
	}
	// Tell the queue that we are done with processing this key. This unblocks the key for other workers
	// This allows safe parallel processing because two pods with the same key are never processed in
	// parallel.
	defer queue.Done(key)

	err := businessLogic(key.(string))
	handleErr(err, key, queue)
	return true
}

func businessLogic(key string) error {
	logger.With(zap.String("message_key", key)).Info("new key received")
	return nil
}

func businessLogicAsync(key string, action string, todoClient *clientset.Clientset, kubeClient *kubernetes.Clientset) {
	logger.
		With(
			zap.String("message_key", key),
			zap.String("action", action)).
		Info("new key received")

	result := strings.Split(key, "/")
	ns := result[0]
	name := result[1]

	todo, err := todoClient.K8sasbackendV1().Todos(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		logger.Error("error getting the todo")
	}
	logger.
		With(
			zap.String("crd_name", todo.Name)).
		Info("crd received")

	recorder := eventRecorder(kubeClient)
	todov1.AddToScheme(scheme.Scheme)
	ref, err := reference.GetReference(scheme.Scheme, todo.DeepCopyObject())
	if err != nil {
		//klog.Fatalf("Could not get reference for pod %v: %v\n", todo.Name, err)
		logger.Error(fmt.Sprintf("Could not get reference for pod %v: %v\n", todo.Name, err))
	}
	recorder.Event(ref, v1.EventTypeNormal, "Todo CRD changed", fmt.Sprintf("Todo CRD %s has been %s", todo.Name, action))

}

// handleErr checks if an error happened and makes sure we will retry later.
// from https://github.com/kubernetes/client-go/blob/master/examples/workqueue/main.go#L91
func handleErr(err error, key interface{}, queue workqueue.RateLimitingInterface) {
	if err == nil {
		// Forget about the #AddRateLimited history of the key on every successful synchronization.
		// This ensures that future processing of updates for this key is not delayed because of
		// an outdated error history.
		queue.Forget(key)
		return
	}

	// This controller retries 5 times if something goes wrong. After that, it stops trying.
	if queue.NumRequeues(key) < 5 {
		logger.With(zap.Error(err)).With(zap.String("message_key", key.(string))).Warn("Error syncing todo")

		// Re-enqueue the key rate limited. Based on the rate limiter on the
		// queue and the re-enqueue history, the key will be processed later again.
		queue.AddRateLimited(key)
		return
	}

	queue.Forget(key)
	// Report to an external entity that, even after several retries, we could not successfully process this key
	runtime.HandleError(err)
	logger.With(zap.Error(err)).With(zap.String("message_key", key.(string))).Warn("Dropping message out of the queue")
}

func eventRecorder(kubeClient *kubernetes.Clientset) record.EventRecorder {
	eventBroadcaster := record.NewBroadcaster()
	// regardless glog o klog I get this error when script runs in k8s
	// log: exiting because of error: log: cannot create log: open /tmp/main.todo-crd-65f6cfd66b-gpzhj./root.log.INFO.20200220-115404.1: no such file or directory
	//eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(
		&typedcorev1.EventSinkImpl{
			Interface: kubeClient.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(
		scheme.Scheme,
		v1.EventSource{Component: "todos-informer"})
	return recorder
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
return value
	}

	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}
