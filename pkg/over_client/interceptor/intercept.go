package interceptor

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"sigs.k8s.io/controller-runtime/pkg/over_client"
)

// Funcs contains functions that are called instead of the underlying over_client's methods.
type Funcs struct {
	Get               func(ctx context.Context, client over_client.WithWatch, key over_client.ObjectKey, obj over_client.Object, opts ...over_client.GetOption) error
	List              func(ctx context.Context, client over_client.WithWatch, list over_client.ObjectList, opts ...over_client.ListOption) error
	Create            func(ctx context.Context, client over_client.WithWatch, obj over_client.Object, opts ...over_client.CreateOption) error
	Delete            func(ctx context.Context, client over_client.WithWatch, obj over_client.Object, opts ...over_client.DeleteOption) error
	DeleteAllOf       func(ctx context.Context, client over_client.WithWatch, obj over_client.Object, opts ...over_client.DeleteAllOfOption) error
	Update            func(ctx context.Context, client over_client.WithWatch, obj over_client.Object, opts ...over_client.UpdateOption) error
	Patch             func(ctx context.Context, client over_client.WithWatch, obj over_client.Object, patch over_client.Patch, opts ...over_client.PatchOption) error
	Watch             func(ctx context.Context, client over_client.WithWatch, obj over_client.ObjectList, opts ...over_client.ListOption) (watch.Interface, error)
	SubResource       func(client over_client.WithWatch, subResource string) over_client.SubResourceClient
	SubResourceGet    func(ctx context.Context, client over_client.Client, subResourceName string, obj over_client.Object, subResource over_client.Object, opts ...over_client.SubResourceGetOption) error
	SubResourceCreate func(ctx context.Context, client over_client.Client, subResourceName string, obj over_client.Object, subResource over_client.Object, opts ...over_client.SubResourceCreateOption) error
	SubResourceUpdate func(ctx context.Context, client over_client.Client, subResourceName string, obj over_client.Object, opts ...over_client.SubResourceUpdateOption) error
	SubResourcePatch  func(ctx context.Context, client over_client.Client, subResourceName string, obj over_client.Object, patch over_client.Patch, opts ...over_client.SubResourcePatchOption) error
}

// NewClient returns a new over_interceptor over_client that calls the functions in funcs instead of the underlying over_client's methods, if they are not nil.
func NewClient(interceptedClient over_client.WithWatch, funcs Funcs) over_client.WithWatch {
	return interceptor{
		client: interceptedClient,
		funcs:  funcs,
	}
}

type interceptor struct {
	client over_client.WithWatch
	funcs  Funcs
}

var _ over_client.WithWatch = &interceptor{}

func (c interceptor) GroupVersionKindFor(obj runtime.Object) (schema.GroupVersionKind, error) {
	return c.client.GroupVersionKindFor(obj)
}

func (c interceptor) IsObjectNamespaced(obj runtime.Object) (bool, error) {
	return c.client.IsObjectNamespaced(obj)
}

func (c interceptor) Get(ctx context.Context, key over_client.ObjectKey, obj over_client.Object, opts ...over_client.GetOption) error {
	if c.funcs.Get != nil {
		return c.funcs.Get(ctx, c.client, key, obj, opts...)
	}
	return c.client.Get(ctx, key, obj, opts...)
}

func (c interceptor) List(ctx context.Context, list over_client.ObjectList, opts ...over_client.ListOption) error {
	if c.funcs.List != nil {
		return c.funcs.List(ctx, c.client, list, opts...)
	}
	return c.client.List(ctx, list, opts...)
}

func (c interceptor) Create(ctx context.Context, obj over_client.Object, opts ...over_client.CreateOption) error {
	if c.funcs.Create != nil {
		return c.funcs.Create(ctx, c.client, obj, opts...)
	}
	return c.client.Create(ctx, obj, opts...)
}

func (c interceptor) Delete(ctx context.Context, obj over_client.Object, opts ...over_client.DeleteOption) error {
	if c.funcs.Delete != nil {
		return c.funcs.Delete(ctx, c.client, obj, opts...)
	}
	return c.client.Delete(ctx, obj, opts...)
}

func (c interceptor) Update(ctx context.Context, obj over_client.Object, opts ...over_client.UpdateOption) error {
	if c.funcs.Update != nil {
		return c.funcs.Update(ctx, c.client, obj, opts...)
	}
	return c.client.Update(ctx, obj, opts...)
}

func (c interceptor) Patch(ctx context.Context, obj over_client.Object, patch over_client.Patch, opts ...over_client.PatchOption) error {
	if c.funcs.Patch != nil {
		return c.funcs.Patch(ctx, c.client, obj, patch, opts...)
	}
	return c.client.Patch(ctx, obj, patch, opts...)
}

func (c interceptor) DeleteAllOf(ctx context.Context, obj over_client.Object, opts ...over_client.DeleteAllOfOption) error {
	if c.funcs.DeleteAllOf != nil {
		return c.funcs.DeleteAllOf(ctx, c.client, obj, opts...)
	}
	return c.client.DeleteAllOf(ctx, obj, opts...)
}

func (c interceptor) Status() over_client.SubResourceWriter {
	return c.SubResource("status")
}

func (c interceptor) SubResource(subResource string) over_client.SubResourceClient {
	if c.funcs.SubResource != nil {
		return c.funcs.SubResource(c.client, subResource)
	}
	return subResourceInterceptor{
		subResourceName: subResource,
		client:          c.client,
		funcs:           c.funcs,
	}
}

func (c interceptor) Scheme() *runtime.Scheme {
	return c.client.Scheme()
}

func (c interceptor) RESTMapper() meta.RESTMapper {
	return c.client.RESTMapper()
}

func (c interceptor) Watch(ctx context.Context, obj over_client.ObjectList, opts ...over_client.ListOption) (watch.Interface, error) {
	if c.funcs.Watch != nil {
		return c.funcs.Watch(ctx, c.client, obj, opts...)
	}
	return c.client.Watch(ctx, obj, opts...)
}

type subResourceInterceptor struct {
	subResourceName string
	client          over_client.Client
	funcs           Funcs
}

var _ over_client.SubResourceClient = &subResourceInterceptor{}

func (s subResourceInterceptor) Get(ctx context.Context, obj over_client.Object, subResource over_client.Object, opts ...over_client.SubResourceGetOption) error {
	if s.funcs.SubResourceGet != nil {
		return s.funcs.SubResourceGet(ctx, s.client, s.subResourceName, obj, subResource, opts...)
	}
	return s.client.SubResource(s.subResourceName).Get(ctx, obj, subResource, opts...)
}

func (s subResourceInterceptor) Create(ctx context.Context, obj over_client.Object, subResource over_client.Object, opts ...over_client.SubResourceCreateOption) error {
	if s.funcs.SubResourceCreate != nil {
		return s.funcs.SubResourceCreate(ctx, s.client, s.subResourceName, obj, subResource, opts...)
	}
	return s.client.SubResource(s.subResourceName).Create(ctx, obj, subResource, opts...)
}

func (s subResourceInterceptor) Update(ctx context.Context, obj over_client.Object, opts ...over_client.SubResourceUpdateOption) error {
	if s.funcs.SubResourceUpdate != nil {
		return s.funcs.SubResourceUpdate(ctx, s.client, s.subResourceName, obj, opts...)
	}
	return s.client.SubResource(s.subResourceName).Update(ctx, obj, opts...)
}

func (s subResourceInterceptor) Patch(ctx context.Context, obj over_client.Object, patch over_client.Patch, opts ...over_client.SubResourcePatchOption) error {
	if s.funcs.SubResourcePatch != nil {
		return s.funcs.SubResourcePatch(ctx, s.client, s.subResourceName, obj, patch, opts...)
	}
	return s.client.SubResource(s.subResourceName).Patch(ctx, obj, patch, opts...)
}
