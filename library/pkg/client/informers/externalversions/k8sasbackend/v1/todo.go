/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	time "time"

	k8sasbackendv1 "github.com/crixo/k8s-as-backend/library/pkg/apis/k8sasbackend/v1"
	versioned "github.com/crixo/k8s-as-backend/library/pkg/client/clientset/versioned"
	internalinterfaces "github.com/crixo/k8s-as-backend/library/pkg/client/informers/externalversions/internalinterfaces"
	v1 "github.com/crixo/k8s-as-backend/library/pkg/client/listers/k8sasbackend/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// TodoInformer provides access to a shared informer and lister for
// Todos.
type TodoInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.TodoLister
}

type todoInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewTodoInformer constructs a new informer for Todo type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewTodoInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredTodoInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredTodoInformer constructs a new informer for Todo type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredTodoInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.K8sasbackendV1().Todos(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.K8sasbackendV1().Todos(namespace).Watch(options)
			},
		},
		&k8sasbackendv1.Todo{},
		resyncPeriod,
		indexers,
	)
}

func (f *todoInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredTodoInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *todoInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&k8sasbackendv1.Todo{}, f.defaultInformer)
}

func (f *todoInformer) Lister() v1.TodoLister {
	return v1.NewTodoLister(f.Informer().GetIndexer())
}