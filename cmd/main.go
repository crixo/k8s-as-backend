package main

import (
	"flag"
	"fmt"
	"os/user"
	"path/filepath"
	"time"

	clientset "github.com/crixo/k8s-as-backend/pkg/client/clientset/versioned"
	todoInformers "github.com/crixo/k8s-as-backend/pkg/client/informers/externalversions"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
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

	todoClient, err := clientset.NewForConfig(config)
	if err != nil {
		logger.With(zap.Error(err)).Fatal("Error building clientset")
	}

	factory := todoInformers.NewSharedInformerFactory(todoClient, time.Second*30)
	todoInformer := factory.K8sasbackend().V1().Todos().Informer()
	todoInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				//queue.Add(key)
				businessLogicAsync(key, "add")
			}
		},
		UpdateFunc: func(old interface{}, neww interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(neww)
			if err == nil {
				//queue.Add(key)
				businessLogicAsync(key, "update")
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				//queue.Add(key)
				businessLogicAsync(key, "delete")
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
	<- stopCh
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

func businessLogicAsync(key string, action string) {
	logger.
	With(
		zap.String("message_key", key), 
		zap.String("action", action)).
	Info("new key received")
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
