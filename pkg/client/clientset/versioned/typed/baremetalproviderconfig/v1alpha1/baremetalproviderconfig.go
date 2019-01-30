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

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"time"

	v1alpha1 "github.com/redhat-nfvpe/cluster-api-provider-baremetal/pkg/apis/baremetalproviderconfig/v1alpha1"
	scheme "github.com/redhat-nfvpe/cluster-api-provider-baremetal/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// BaremetalProviderConfigsGetter has a method to return a BaremetalProviderConfigInterface.
// A group's client should implement this interface.
type BaremetalProviderConfigsGetter interface {
	BaremetalProviderConfigs(namespace string) BaremetalProviderConfigInterface
}

// BaremetalProviderConfigInterface has methods to work with BaremetalProviderConfig resources.
type BaremetalProviderConfigInterface interface {
	Create(*v1alpha1.BaremetalProviderConfig) (*v1alpha1.BaremetalProviderConfig, error)
	Update(*v1alpha1.BaremetalProviderConfig) (*v1alpha1.BaremetalProviderConfig, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.BaremetalProviderConfig, error)
	List(opts v1.ListOptions) (*v1alpha1.BaremetalProviderConfigList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.BaremetalProviderConfig, err error)
	BaremetalProviderConfigExpansion
}

// baremetalProviderConfigs implements BaremetalProviderConfigInterface
type baremetalProviderConfigs struct {
	client rest.Interface
	ns     string
}

// newBaremetalProviderConfigs returns a BaremetalProviderConfigs
func newBaremetalProviderConfigs(c *BaremetalproviderconfigV1alpha1Client, namespace string) *baremetalProviderConfigs {
	return &baremetalProviderConfigs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the baremetalProviderConfig, and returns the corresponding baremetalProviderConfig object, and an error if there is any.
func (c *baremetalProviderConfigs) Get(name string, options v1.GetOptions) (result *v1alpha1.BaremetalProviderConfig, err error) {
	result = &v1alpha1.BaremetalProviderConfig{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("baremetalproviderconfigs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of BaremetalProviderConfigs that match those selectors.
func (c *baremetalProviderConfigs) List(opts v1.ListOptions) (result *v1alpha1.BaremetalProviderConfigList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.BaremetalProviderConfigList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("baremetalproviderconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested baremetalProviderConfigs.
func (c *baremetalProviderConfigs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("baremetalproviderconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a baremetalProviderConfig and creates it.  Returns the server's representation of the baremetalProviderConfig, and an error, if there is any.
func (c *baremetalProviderConfigs) Create(baremetalProviderConfig *v1alpha1.BaremetalProviderConfig) (result *v1alpha1.BaremetalProviderConfig, err error) {
	result = &v1alpha1.BaremetalProviderConfig{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("baremetalproviderconfigs").
		Body(baremetalProviderConfig).
		Do().
		Into(result)
	return
}

// Update takes the representation of a baremetalProviderConfig and updates it. Returns the server's representation of the baremetalProviderConfig, and an error, if there is any.
func (c *baremetalProviderConfigs) Update(baremetalProviderConfig *v1alpha1.BaremetalProviderConfig) (result *v1alpha1.BaremetalProviderConfig, err error) {
	result = &v1alpha1.BaremetalProviderConfig{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("baremetalproviderconfigs").
		Name(baremetalProviderConfig.Name).
		Body(baremetalProviderConfig).
		Do().
		Into(result)
	return
}

// Delete takes name of the baremetalProviderConfig and deletes it. Returns an error if one occurs.
func (c *baremetalProviderConfigs) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("baremetalproviderconfigs").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *baremetalProviderConfigs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("baremetalproviderconfigs").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched baremetalProviderConfig.
func (c *baremetalProviderConfigs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.BaremetalProviderConfig, err error) {
	result = &v1alpha1.BaremetalProviderConfig{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("baremetalproviderconfigs").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
