/*
Copyright (c) 2020 - 2024 SUSE LLC

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

// Code generated by main. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/rancher/fleet/pkg/apis/fleet.cattle.io/v1alpha1"
	"github.com/rancher/lasso/pkg/client"
	"github.com/rancher/lasso/pkg/controller"
	"github.com/rancher/wrangler/pkg/generic"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type ContentHandler func(string, *v1alpha1.Content) (*v1alpha1.Content, error)

type ContentController interface {
	generic.ControllerMeta
	ContentClient

	OnChange(ctx context.Context, name string, sync ContentHandler)
	OnRemove(ctx context.Context, name string, sync ContentHandler)
	Enqueue(name string)
	EnqueueAfter(name string, duration time.Duration)

	Cache() ContentCache
}

type ContentClient interface {
	Create(*v1alpha1.Content) (*v1alpha1.Content, error)
	Update(*v1alpha1.Content) (*v1alpha1.Content, error)

	Delete(name string, options *metav1.DeleteOptions) error
	Get(name string, options metav1.GetOptions) (*v1alpha1.Content, error)
	List(opts metav1.ListOptions) (*v1alpha1.ContentList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Content, err error)
}

type ContentCache interface {
	Get(name string) (*v1alpha1.Content, error)
	List(selector labels.Selector) ([]*v1alpha1.Content, error)

	AddIndexer(indexName string, indexer ContentIndexer)
	GetByIndex(indexName, key string) ([]*v1alpha1.Content, error)
}

type ContentIndexer func(obj *v1alpha1.Content) ([]string, error)

type contentController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewContentController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) ContentController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &contentController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromContentHandlerToHandler(sync ContentHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1alpha1.Content
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1alpha1.Content))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *contentController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1alpha1.Content))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateContentDeepCopyOnChange(client ContentClient, obj *v1alpha1.Content, handler func(obj *v1alpha1.Content) (*v1alpha1.Content, error)) (*v1alpha1.Content, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *contentController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *contentController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *contentController) OnChange(ctx context.Context, name string, sync ContentHandler) {
	c.AddGenericHandler(ctx, name, FromContentHandlerToHandler(sync))
}

func (c *contentController) OnRemove(ctx context.Context, name string, sync ContentHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromContentHandlerToHandler(sync)))
}

func (c *contentController) Enqueue(name string) {
	c.controller.Enqueue("", name)
}

func (c *contentController) EnqueueAfter(name string, duration time.Duration) {
	c.controller.EnqueueAfter("", name, duration)
}

func (c *contentController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *contentController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *contentController) Cache() ContentCache {
	return &contentCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *contentController) Create(obj *v1alpha1.Content) (*v1alpha1.Content, error) {
	result := &v1alpha1.Content{}
	return result, c.client.Create(context.TODO(), "", obj, result, metav1.CreateOptions{})
}

func (c *contentController) Update(obj *v1alpha1.Content) (*v1alpha1.Content, error) {
	result := &v1alpha1.Content{}
	return result, c.client.Update(context.TODO(), "", obj, result, metav1.UpdateOptions{})
}

func (c *contentController) Delete(name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), "", name, *options)
}

func (c *contentController) Get(name string, options metav1.GetOptions) (*v1alpha1.Content, error) {
	result := &v1alpha1.Content{}
	return result, c.client.Get(context.TODO(), "", name, result, options)
}

func (c *contentController) List(opts metav1.ListOptions) (*v1alpha1.ContentList, error) {
	result := &v1alpha1.ContentList{}
	return result, c.client.List(context.TODO(), "", result, opts)
}

func (c *contentController) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), "", opts)
}

func (c *contentController) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (*v1alpha1.Content, error) {
	result := &v1alpha1.Content{}
	return result, c.client.Patch(context.TODO(), "", name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type contentCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *contentCache) Get(name string) (*v1alpha1.Content, error) {
	obj, exists, err := c.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v1alpha1.Content), nil
}

func (c *contentCache) List(selector labels.Selector) (ret []*v1alpha1.Content, err error) {

	err = cache.ListAll(c.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Content))
	})

	return ret, err
}

func (c *contentCache) AddIndexer(indexName string, indexer ContentIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1alpha1.Content))
		},
	}))
}

func (c *contentCache) GetByIndex(indexName, key string) (result []*v1alpha1.Content, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1alpha1.Content, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1alpha1.Content))
	}
	return result, nil
}
