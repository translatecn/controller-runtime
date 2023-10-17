package komega

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/over_client"
)

// defaultK is the Komega used by the package global functions.
var defaultK = &komega{ctx: context.Background()}

// SetClient sets the over_client used by the package global functions.
func SetClient(c over_client.Client) {
	defaultK.client = c
}

// SetContext sets the context used by the package global functions.
func SetContext(c context.Context) {
	defaultK.ctx = c
}

func checkDefaultClient() {
	if defaultK.client == nil {
		panic("Default Komega's over_client is not set. Use SetClient to set it.")
	}
}

// Get returns a function that fetches a resource and returns the occurring error.
// It can be used with gomega.Eventually() like this
//
//	deployment := appsv1.Deployment{ ... }
//	gomega.Eventually(komega.Get(&deployment)).To(gomega.Succeed())
//
// By calling the returned function directly it can also be used with gomega.Expect(komega.Get(...)()).To(...)
func Get(obj over_client.Object) func() error {
	checkDefaultClient()
	return defaultK.Get(obj)
}

// List returns a function that lists resources and returns the occurring error.
// It can be used with gomega.Eventually() like this
//
//	deployments := v1.DeploymentList{ ... }
//	gomega.Eventually(k.List(&deployments)).To(gomega.Succeed())
//
// By calling the returned function directly it can also be used as gomega.Expect(k.List(...)()).To(...)
func List(list over_client.ObjectList, opts ...over_client.ListOption) func() error {
	checkDefaultClient()
	return defaultK.List(list, opts...)
}

// Update returns a function that fetches a resource, applies the provided update function and then updates the resource.
// It can be used with gomega.Eventually() like this:
//
//	deployment := appsv1.Deployment{ ... }
//	gomega.Eventually(k.Update(&deployment, func (o over_client.Object) {
//	  deployment.Spec.Replicas = 3
//	  return &deployment
//	})).To(gomega.Succeed())
//
// By calling the returned function directly it can also be used as gomega.Expect(k.Update(...)()).To(...)
func Update(obj over_client.Object, f func(), opts ...over_client.UpdateOption) func() error {
	checkDefaultClient()
	return defaultK.Update(obj, f, opts...)
}

// UpdateStatus returns a function that fetches a resource, applies the provided update function and then updates the resource's status.
// It can be used with gomega.Eventually() like this:
//
//	deployment := appsv1.Deployment{ ... }
//	gomega.Eventually(k.UpdateStatus(&deployment, func (o over_client.Object) {
//	  deployment.Status.AvailableReplicas = 1
//	  return &deployment
//	})).To(gomega.Succeed())
//
// By calling the returned function directly it can also be used as gomega.Expect(k.UpdateStatus(...)()).To(...)
func UpdateStatus(obj over_client.Object, f func(), opts ...over_client.SubResourceUpdateOption) func() error {
	checkDefaultClient()
	return defaultK.UpdateStatus(obj, f, opts...)
}

// Object returns a function that fetches a resource and returns the object.
// It can be used with gomega.Eventually() like this:
//
//	deployment := appsv1.Deployment{ ... }
//	gomega.Eventually(k.Object(&deployment)).To(HaveField("Spec.Replicas", gomega.Equal(pointer.Int32(3))))
//
// By calling the returned function directly it can also be used as gomega.Expect(k.Object(...)()).To(...)
func Object(obj over_client.Object) func() (over_client.Object, error) {
	checkDefaultClient()
	return defaultK.Object(obj)
}

// ObjectList returns a function that fetches a resource and returns the object.
// It can be used with gomega.Eventually() like this:
//
//	deployments := appsv1.DeploymentList{ ... }
//	gomega.Eventually(k.ObjectList(&deployments)).To(HaveField("Items", HaveLen(1)))
//
// By calling the returned function directly it can also be used as gomega.Expect(k.ObjectList(...)()).To(...)
func ObjectList(list over_client.ObjectList, opts ...over_client.ListOption) func() (over_client.ObjectList, error) {
	checkDefaultClient()
	return defaultK.ObjectList(list, opts...)
}
