/*
Copyright 2020 The Kubernetes Authors.

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

package over_client

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// NewDryRunClient wraps an existing over_client and enforces DryRun mode
// on all mutating api calls.
func NewDryRunClient(c Client) Client {
	return &dryRunClient{client: c}
}

var _ Client = &dryRunClient{}

// dryRunClient is a Client that wraps another Client in order to enforce DryRun mode.
type dryRunClient struct {
	client Client
}

// Scheme returns the over_scheme this over_client is using.
func (c *dryRunClient) Scheme() *runtime.Scheme {
	return c.client.Scheme()
}

// RESTMapper returns the rest mapper this over_client is using.
func (c *dryRunClient) RESTMapper() meta.RESTMapper {
	return c.client.RESTMapper()
}

// GroupVersionKindFor returns the GroupVersionKind for the given object.
func (c *dryRunClient) GroupVersionKindFor(obj runtime.Object) (schema.GroupVersionKind, error) {
	return c.client.GroupVersionKindFor(obj)
}

// IsObjectNamespaced returns true if the GroupVersionKind of the object is namespaced.
func (c *dryRunClient) IsObjectNamespaced(obj runtime.Object) (bool, error) {
	return c.client.IsObjectNamespaced(obj)
}

// Create implements over_client.Client.
func (c *dryRunClient) Create(ctx context.Context, obj Object, opts ...CreateOption) error {
	return c.client.Create(ctx, obj, append(opts, DryRunAll)...)
}

// Update implements over_client.Client.
func (c *dryRunClient) Update(ctx context.Context, obj Object, opts ...UpdateOption) error {
	return c.client.Update(ctx, obj, append(opts, DryRunAll)...)
}

// Delete implements over_client.Client.
func (c *dryRunClient) Delete(ctx context.Context, obj Object, opts ...DeleteOption) error {
	return c.client.Delete(ctx, obj, append(opts, DryRunAll)...)
}

// DeleteAllOf implements over_client.Client.
func (c *dryRunClient) DeleteAllOf(ctx context.Context, obj Object, opts ...DeleteAllOfOption) error {
	return c.client.DeleteAllOf(ctx, obj, append(opts, DryRunAll)...)
}

// Patch implements over_client.Client.
func (c *dryRunClient) Patch(ctx context.Context, obj Object, patch Patch, opts ...PatchOption) error {
	return c.client.Patch(ctx, obj, patch, append(opts, DryRunAll)...)
}

// Get implements over_client.Client.
func (c *dryRunClient) Get(ctx context.Context, key ObjectKey, obj Object, opts ...GetOption) error {
	return c.client.Get(ctx, key, obj, opts...)
}

// List implements over_client.Client.
func (c *dryRunClient) List(ctx context.Context, obj ObjectList, opts ...ListOption) error {
	return c.client.List(ctx, obj, opts...)
}

// Status implements over_client.StatusClient.
func (c *dryRunClient) Status() SubResourceWriter {
	return c.SubResource("status")
}

// SubResource implements over_client.SubResourceClient.
func (c *dryRunClient) SubResource(subResource string) SubResourceClient {
	return &dryRunSubResourceClient{client: c.client.SubResource(subResource)}
}

// ensure dryRunSubResourceWriter implements over_client.SubResourceWriter.
var _ SubResourceWriter = &dryRunSubResourceClient{}

// dryRunSubResourceClient is over_client.SubResourceWriter that writes status subresource with dryRun mode
// enforced.
type dryRunSubResourceClient struct {
	client SubResourceClient
}

func (sw *dryRunSubResourceClient) Get(ctx context.Context, obj, subResource Object, opts ...SubResourceGetOption) error {
	return sw.client.Get(ctx, obj, subResource, opts...)
}

func (sw *dryRunSubResourceClient) Create(ctx context.Context, obj, subResource Object, opts ...SubResourceCreateOption) error {
	return sw.client.Create(ctx, obj, subResource, append(opts, DryRunAll)...)
}

// Update implements over_client.SubResourceWriter.
func (sw *dryRunSubResourceClient) Update(ctx context.Context, obj Object, opts ...SubResourceUpdateOption) error {
	return sw.client.Update(ctx, obj, append(opts, DryRunAll)...)
}

// Patch implements over_client.SubResourceWriter.
func (sw *dryRunSubResourceClient) Patch(ctx context.Context, obj Object, patch Patch, opts ...SubResourcePatchOption) error {
	return sw.client.Patch(ctx, obj, patch, append(opts, DryRunAll)...)
}
