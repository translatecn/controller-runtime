package interceptor

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"sigs.k8s.io/controller-runtime/pkg/over_client"
)

var _ = Describe("NewClient", func() {
	wrappedClient := dummyClient{}
	ctx := context.Background()
	It("should call the provided Get function", func() {
		var called bool
		client := NewClient(wrappedClient, Funcs{
			Get: func(ctx context.Context, client over_client.WithWatch, key over_client.ObjectKey, obj over_client.Object, opts ...over_client.GetOption) error {
				called = true
				return nil
			},
		})
		_ = client.Get(ctx, types.NamespacedName{}, nil)
		Expect(called).To(BeTrue())
	})
	It("should call the underlying over_client if the provided Get function is nil", func() {
		var called bool
		client1 := NewClient(wrappedClient, Funcs{
			Get: func(ctx context.Context, client over_client.WithWatch, key over_client.ObjectKey, obj over_client.Object, opts ...over_client.GetOption) error {
				called = true
				return nil
			},
		})
		client2 := NewClient(client1, Funcs{})
		_ = client2.Get(ctx, types.NamespacedName{}, nil)
		Expect(called).To(BeTrue())
	})
	It("should call the provided List function", func() {
		var called bool
		client := NewClient(wrappedClient, Funcs{
			List: func(ctx context.Context, client over_client.WithWatch, list over_client.ObjectList, opts ...over_client.ListOption) error {
				called = true
				return nil
			},
		})
		_ = client.List(ctx, nil)
		Expect(called).To(BeTrue())
	})
	It("should call the underlying over_client if the provided List function is nil", func() {
		var called bool
		client1 := NewClient(wrappedClient, Funcs{
			List: func(ctx context.Context, client over_client.WithWatch, list over_client.ObjectList, opts ...over_client.ListOption) error {
				called = true
				return nil
			},
		})
		client2 := NewClient(client1, Funcs{})
		_ = client2.List(ctx, nil)
		Expect(called).To(BeTrue())
	})
	It("should call the provided Create function", func() {
		var called bool
		client := NewClient(wrappedClient, Funcs{
			Create: func(ctx context.Context, client over_client.WithWatch, obj over_client.Object, opts ...over_client.CreateOption) error {
				called = true
				return nil
			},
		})
		_ = client.Create(ctx, nil)
		Expect(called).To(BeTrue())
	})
	It("should call the underlying over_client if the provided Create function is nil", func() {
		var called bool
		client1 := NewClient(wrappedClient, Funcs{
			Create: func(ctx context.Context, client over_client.WithWatch, obj over_client.Object, opts ...over_client.CreateOption) error {
				called = true
				return nil
			},
		})
		client2 := NewClient(client1, Funcs{})
		_ = client2.Create(ctx, nil)
		Expect(called).To(BeTrue())
	})
	It("should call the provided Delete function", func() {
		var called bool
		client := NewClient(wrappedClient, Funcs{
			Delete: func(ctx context.Context, client over_client.WithWatch, obj over_client.Object, opts ...over_client.DeleteOption) error {
				called = true
				return nil
			},
		})
		_ = client.Delete(ctx, nil)
		Expect(called).To(BeTrue())
	})
	It("should call the underlying over_client if the provided Delete function is nil", func() {
		var called bool
		client1 := NewClient(wrappedClient, Funcs{
			Delete: func(ctx context.Context, client over_client.WithWatch, obj over_client.Object, opts ...over_client.DeleteOption) error {
				called = true
				return nil
			},
		})
		client2 := NewClient(client1, Funcs{})
		_ = client2.Delete(ctx, nil)
		Expect(called).To(BeTrue())
	})
	It("should call the provided DeleteAllOf function", func() {
		var called bool
		client := NewClient(wrappedClient, Funcs{
			DeleteAllOf: func(ctx context.Context, client over_client.WithWatch, obj over_client.Object, opts ...over_client.DeleteAllOfOption) error {
				called = true
				return nil
			},
		})
		_ = client.DeleteAllOf(ctx, nil)
		Expect(called).To(BeTrue())
	})
	It("should call the underlying over_client if the provided DeleteAllOf function is nil", func() {
		var called bool
		client1 := NewClient(wrappedClient, Funcs{
			DeleteAllOf: func(ctx context.Context, client over_client.WithWatch, obj over_client.Object, opts ...over_client.DeleteAllOfOption) error {
				called = true
				return nil
			},
		})
		client2 := NewClient(client1, Funcs{})
		_ = client2.DeleteAllOf(ctx, nil)
		Expect(called).To(BeTrue())
	})
	It("should call the provided Update function", func() {
		var called bool
		client := NewClient(wrappedClient, Funcs{
			Update: func(ctx context.Context, client over_client.WithWatch, obj over_client.Object, opts ...over_client.UpdateOption) error {
				called = true
				return nil
			},
		})
		_ = client.Update(ctx, nil)
		Expect(called).To(BeTrue())
	})
	It("should call the underlying over_client if the provided Update function is nil", func() {
		var called bool
		client1 := NewClient(wrappedClient, Funcs{
			Update: func(ctx context.Context, client over_client.WithWatch, obj over_client.Object, opts ...over_client.UpdateOption) error {
				called = true
				return nil
			},
		})
		client2 := NewClient(client1, Funcs{})
		_ = client2.Update(ctx, nil)
		Expect(called).To(BeTrue())
	})
	It("should call the provided Patch function", func() {
		var called bool
		client := NewClient(wrappedClient, Funcs{
			Patch: func(ctx context.Context, client over_client.WithWatch, obj over_client.Object, patch over_client.Patch, opts ...over_client.PatchOption) error {
				called = true
				return nil
			},
		})
		_ = client.Patch(ctx, nil, nil)
		Expect(called).To(BeTrue())
	})
	It("should call the underlying over_client if the provided Patch function is nil", func() {
		var called bool
		client1 := NewClient(wrappedClient, Funcs{
			Patch: func(ctx context.Context, client over_client.WithWatch, obj over_client.Object, patch over_client.Patch, opts ...over_client.PatchOption) error {
				called = true
				return nil
			},
		})
		client2 := NewClient(client1, Funcs{})
		_ = client2.Patch(ctx, nil, nil)
		Expect(called).To(BeTrue())
	})
	It("should call the provided Watch function", func() {
		var called bool
		client := NewClient(wrappedClient, Funcs{
			Watch: func(ctx context.Context, client over_client.WithWatch, obj over_client.ObjectList, opts ...over_client.ListOption) (watch.Interface, error) {
				called = true
				return nil, nil
			},
		})
		_, _ = client.Watch(ctx, nil)
		Expect(called).To(BeTrue())
	})
	It("should call the underlying over_client if the provided Watch function is nil", func() {
		var called bool
		client1 := NewClient(wrappedClient, Funcs{
			Watch: func(ctx context.Context, client over_client.WithWatch, obj over_client.ObjectList, opts ...over_client.ListOption) (watch.Interface, error) {
				called = true
				return nil, nil
			},
		})
		client2 := NewClient(client1, Funcs{})
		_, _ = client2.Watch(ctx, nil)
		Expect(called).To(BeTrue())
	})
	It("should call the provided SubResource function", func() {
		var called bool
		client := NewClient(wrappedClient, Funcs{
			SubResource: func(client over_client.WithWatch, subResource string) over_client.SubResourceClient {
				called = true
				return nil
			},
		})
		_ = client.SubResource("")
		Expect(called).To(BeTrue())
	})
	It("should call the provided SubResource function with 'status' when calling Status()", func() {
		var called bool
		client := NewClient(wrappedClient, Funcs{
			SubResource: func(client over_client.WithWatch, subResource string) over_client.SubResourceClient {
				if subResource == "status" {
					called = true
				}
				return nil
			},
		})
		_ = client.Status()
		Expect(called).To(BeTrue())
	})
})

var _ = Describe("NewSubResourceClient", func() {
	c := dummyClient{}
	ctx := context.Background()
	It("should call the provided Get function", func() {
		var called bool
		c := NewClient(c, Funcs{
			SubResourceGet: func(_ context.Context, client over_client.Client, subResourceName string, obj, subResource over_client.Object, opts ...over_client.SubResourceGetOption) error {
				called = true
				Expect(subResourceName).To(BeEquivalentTo("foo"))
				return nil
			},
		})
		_ = c.SubResource("foo").Get(ctx, nil, nil)
		Expect(called).To(BeTrue())
	})
	It("should call the underlying over_client if the provided Get function is nil", func() {
		var called bool
		client1 := NewClient(c, Funcs{
			SubResourceGet: func(_ context.Context, client over_client.Client, subResourceName string, obj, subResource over_client.Object, opts ...over_client.SubResourceGetOption) error {
				called = true
				Expect(subResourceName).To(BeEquivalentTo("foo"))
				return nil
			},
		})
		client2 := NewClient(client1, Funcs{})
		_ = client2.SubResource("foo").Get(ctx, nil, nil)
		Expect(called).To(BeTrue())
	})
	It("should call the provided Update function", func() {
		var called bool
		client := NewClient(c, Funcs{
			SubResourceUpdate: func(_ context.Context, client over_client.Client, subResourceName string, obj over_client.Object, opts ...over_client.SubResourceUpdateOption) error {
				called = true
				Expect(subResourceName).To(BeEquivalentTo("foo"))
				return nil
			},
		})
		_ = client.SubResource("foo").Update(ctx, nil, nil)
		Expect(called).To(BeTrue())
	})
	It("should call the underlying over_client if the provided Update function is nil", func() {
		var called bool
		client1 := NewClient(c, Funcs{
			SubResourceUpdate: func(_ context.Context, client over_client.Client, subResourceName string, obj over_client.Object, opts ...over_client.SubResourceUpdateOption) error {
				called = true
				Expect(subResourceName).To(BeEquivalentTo("foo"))
				return nil
			},
		})
		client2 := NewClient(client1, Funcs{})
		_ = client2.SubResource("foo").Update(ctx, nil, nil)
		Expect(called).To(BeTrue())
	})
	It("should call the provided Patch function", func() {
		var called bool
		client := NewClient(c, Funcs{
			SubResourcePatch: func(_ context.Context, client over_client.Client, subResourceName string, obj over_client.Object, patch over_client.Patch, opts ...over_client.SubResourcePatchOption) error {
				called = true
				Expect(subResourceName).To(BeEquivalentTo("foo"))
				return nil
			},
		})
		_ = client.SubResource("foo").Patch(ctx, nil, nil)
		Expect(called).To(BeTrue())
	})
	It("should call the underlying over_client if the provided Patch function is nil", func() {
		var called bool
		client1 := NewClient(c, Funcs{
			SubResourcePatch: func(ctx context.Context, client over_client.Client, subResourceName string, obj over_client.Object, patch over_client.Patch, opts ...over_client.SubResourcePatchOption) error {
				called = true
				Expect(subResourceName).To(BeEquivalentTo("foo"))
				return nil
			},
		})
		client2 := NewClient(client1, Funcs{})
		_ = client2.SubResource("foo").Patch(ctx, nil, nil)
		Expect(called).To(BeTrue())
	})
	It("should call the provided Create function", func() {
		var called bool
		client := NewClient(c, Funcs{
			SubResourceCreate: func(_ context.Context, client over_client.Client, subResourceName string, obj, subResource over_client.Object, opts ...over_client.SubResourceCreateOption) error {
				called = true
				Expect(subResourceName).To(BeEquivalentTo("foo"))
				return nil
			},
		})
		_ = client.SubResource("foo").Create(ctx, nil, nil)
		Expect(called).To(BeTrue())
	})
	It("should call the underlying over_client if the provided Create function is nil", func() {
		var called bool
		client1 := NewClient(c, Funcs{
			SubResourceCreate: func(_ context.Context, client over_client.Client, subResourceName string, obj, subResource over_client.Object, opts ...over_client.SubResourceCreateOption) error {
				called = true
				Expect(subResourceName).To(BeEquivalentTo("foo"))
				return nil
			},
		})
		client2 := NewClient(client1, Funcs{})
		_ = client2.SubResource("foo").Create(ctx, nil, nil)
		Expect(called).To(BeTrue())
	})
})

type dummyClient struct{}

var _ over_client.WithWatch = &dummyClient{}

func (d dummyClient) Get(ctx context.Context, key over_client.ObjectKey, obj over_client.Object, opts ...over_client.GetOption) error {
	return nil
}

func (d dummyClient) List(ctx context.Context, list over_client.ObjectList, opts ...over_client.ListOption) error {
	return nil
}

func (d dummyClient) Create(ctx context.Context, obj over_client.Object, opts ...over_client.CreateOption) error {
	return nil
}

func (d dummyClient) Delete(ctx context.Context, obj over_client.Object, opts ...over_client.DeleteOption) error {
	return nil
}

func (d dummyClient) Update(ctx context.Context, obj over_client.Object, opts ...over_client.UpdateOption) error {
	return nil
}

func (d dummyClient) Patch(ctx context.Context, obj over_client.Object, patch over_client.Patch, opts ...over_client.PatchOption) error {
	return nil
}

func (d dummyClient) DeleteAllOf(ctx context.Context, obj over_client.Object, opts ...over_client.DeleteAllOfOption) error {
	return nil
}

func (d dummyClient) Status() over_client.SubResourceWriter {
	return d.SubResource("status")
}

func (d dummyClient) SubResource(subResource string) over_client.SubResourceClient {
	return nil
}

func (d dummyClient) Scheme() *runtime.Scheme {
	return nil
}

func (d dummyClient) RESTMapper() meta.RESTMapper {
	return nil
}

func (d dummyClient) GroupVersionKindFor(obj runtime.Object) (schema.GroupVersionKind, error) {
	return schema.GroupVersionKind{}, nil
}

func (d dummyClient) IsObjectNamespaced(obj runtime.Object) (bool, error) {
	return false, nil
}

func (d dummyClient) Watch(ctx context.Context, obj over_client.ObjectList, opts ...over_client.ListOption) (watch.Interface, error) {
	return nil, nil
}
