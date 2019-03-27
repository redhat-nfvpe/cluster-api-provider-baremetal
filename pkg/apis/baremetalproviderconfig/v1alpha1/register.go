/*
Copyright 2019 The Kubernetes Authors.

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

package v1alpha1

import (
	"bytes"
	"fmt"

	"github.com/ghodss/yaml"
	machinev1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/controller-runtime/pkg/runtime/scheme"

	"github.com/redhat-nfvpe/cluster-api-provider-baremetal/pkg/apis/baremetalproviderconfig"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: baremetalproviderconfig.GroupName, Version: "v1alpha1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}
)

// BaremetalProviderConfigCodec contains encoder/decoder to convert this types from/to serialize data
// +k8s:deepcopy-gen=false
type BaremetalProviderConfigCodec struct {
	encoder runtime.Encoder
	decoder runtime.Decoder
}

// NewScheme creates a new Scheme
func NewScheme() (*runtime.Scheme, error) {
	return SchemeBuilder.Build()
}

// NewCodec returns a encode/decoder for this API
func NewCodec() (*BaremetalProviderConfigCodec, error) {
	scheme, err := NewScheme()
	if err != nil {
		return nil, err
	}
	codecFactory := serializer.NewCodecFactory(scheme)
	encoder, err := newEncoder(&codecFactory)
	if err != nil {
		return nil, err
	}
	codec := BaremetalProviderConfigCodec{
		encoder: encoder,
		decoder: codecFactory.UniversalDecoder(SchemeGroupVersion),
	}
	return &codec, nil
}

// DecodeFromProviderSpec decodes a serialised ProviderConfig into an object
func (codec *BaremetalProviderConfigCodec) DecodeFromProviderSpec(providerConfig machinev1.ProviderSpec, out runtime.Object) error {
	if providerConfig.Value != nil {
		if err := yaml.Unmarshal(providerConfig.Value.Raw, out); err != nil {
			return fmt.Errorf("decoding failure: %v", err)
		}
	}
	return nil
}

// EncodeToProviderSpec encodes an object into a serialised ProviderConfig
func (codec *BaremetalProviderConfigCodec) EncodeToProviderSpec(in runtime.Object) (*machinev1.ProviderSpec, error) {
	var buf bytes.Buffer
	if err := codec.encoder.Encode(in, &buf); err != nil {
		return nil, fmt.Errorf("encoding failed: %v", err)
	}
	return &machinev1.ProviderSpec{
		Value: &runtime.RawExtension{Raw: buf.Bytes()},
	}, nil
}

// EncodeProviderStatus encodes an object into serialised data
func (codec *BaremetalProviderConfigCodec) EncodeProviderStatus(in runtime.Object) (*runtime.RawExtension, error) {
	var buf bytes.Buffer
	if err := codec.encoder.Encode(in, &buf); err != nil {
		return nil, fmt.Errorf("encoding failed: %v", err)
	}

	return &runtime.RawExtension{Raw: buf.Bytes()}, nil
}

// DecodeProviderStatus decodes a serialised providerStatus into an object
func (codec *BaremetalProviderConfigCodec) DecodeProviderStatus(providerStatus *runtime.RawExtension, out runtime.Object) error {
	if providerStatus != nil {
		if err := yaml.Unmarshal(providerStatus.Raw, out); err != nil {
			return fmt.Errorf("decoding failure: %v", err)
		}
	}
	return nil
}

func newEncoder(codecFactory *serializer.CodecFactory) (runtime.Encoder, error) {
	serializerInfos := codecFactory.SupportedMediaTypes()
	if len(serializerInfos) == 0 {
		return nil, fmt.Errorf("unable to find any serlializers")
	}
	encoder := codecFactory.EncoderForVersion(serializerInfos[0].Serializer, SchemeGroupVersion)
	return encoder, nil
}
