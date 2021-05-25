/*
Copyright AppsCode Inc. and Contributors

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
	v1alpha1 "kubedb.dev/apimachinery/apis/autoscaling/v1alpha1"
)

// FakeRedisAutoscalers implements RedisAutoscalerInterface
type FakeRedisAutoscalers struct {
	Fake *FakeAutoscalingV1alpha1
	ns   string
}

var redisautoscalersResource = schema.GroupVersionResource{Group: "autoscaling.kubedb.com", Version: "v1alpha1", Resource: "redisautoscalers"}

var redisautoscalersKind = schema.GroupVersionKind{Group: "autoscaling.kubedb.com", Version: "v1alpha1", Kind: "RedisAutoscaler"}

// Get takes name of the redisAutoscaler, and returns the corresponding redisAutoscaler object, and an error if there is any.
func (c *FakeRedisAutoscalers) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.RedisAutoscaler, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(redisautoscalersResource, c.ns, name), &v1alpha1.RedisAutoscaler{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.RedisAutoscaler), err
}

// List takes label and field selectors, and returns the list of RedisAutoscalers that match those selectors.
func (c *FakeRedisAutoscalers) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.RedisAutoscalerList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(redisautoscalersResource, redisautoscalersKind, c.ns, opts), &v1alpha1.RedisAutoscalerList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.RedisAutoscalerList{ListMeta: obj.(*v1alpha1.RedisAutoscalerList).ListMeta}
	for _, item := range obj.(*v1alpha1.RedisAutoscalerList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested redisAutoscalers.
func (c *FakeRedisAutoscalers) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(redisautoscalersResource, c.ns, opts))

}

// Create takes the representation of a redisAutoscaler and creates it.  Returns the server's representation of the redisAutoscaler, and an error, if there is any.
func (c *FakeRedisAutoscalers) Create(ctx context.Context, redisAutoscaler *v1alpha1.RedisAutoscaler, opts v1.CreateOptions) (result *v1alpha1.RedisAutoscaler, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(redisautoscalersResource, c.ns, redisAutoscaler), &v1alpha1.RedisAutoscaler{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.RedisAutoscaler), err
}

// Update takes the representation of a redisAutoscaler and updates it. Returns the server's representation of the redisAutoscaler, and an error, if there is any.
func (c *FakeRedisAutoscalers) Update(ctx context.Context, redisAutoscaler *v1alpha1.RedisAutoscaler, opts v1.UpdateOptions) (result *v1alpha1.RedisAutoscaler, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(redisautoscalersResource, c.ns, redisAutoscaler), &v1alpha1.RedisAutoscaler{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.RedisAutoscaler), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeRedisAutoscalers) UpdateStatus(ctx context.Context, redisAutoscaler *v1alpha1.RedisAutoscaler, opts v1.UpdateOptions) (*v1alpha1.RedisAutoscaler, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(redisautoscalersResource, "status", c.ns, redisAutoscaler), &v1alpha1.RedisAutoscaler{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.RedisAutoscaler), err
}

// Delete takes name of the redisAutoscaler and deletes it. Returns an error if one occurs.
func (c *FakeRedisAutoscalers) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(redisautoscalersResource, c.ns, name), &v1alpha1.RedisAutoscaler{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeRedisAutoscalers) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(redisautoscalersResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.RedisAutoscalerList{})
	return err
}

// Patch applies the patch and returns the patched redisAutoscaler.
func (c *FakeRedisAutoscalers) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.RedisAutoscaler, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(redisautoscalersResource, c.ns, name, pt, data, subresources...), &v1alpha1.RedisAutoscaler{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.RedisAutoscaler), err
}
