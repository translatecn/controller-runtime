/*
Copyright 2018 The Kubernetes Authors.

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

package over_handler_test

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/over_client"
	"sigs.k8s.io/controller-runtime/pkg/over_controller"
	"sigs.k8s.io/controller-runtime/pkg/over_event"
	"sigs.k8s.io/controller-runtime/pkg/over_handler"
	"sigs.k8s.io/controller-runtime/pkg/over_manager"
	"sigs.k8s.io/controller-runtime/pkg/over_reconcile"
	"sigs.k8s.io/controller-runtime/pkg/over_source"
)

var mgr over_manager.Manager
var c over_controller.Controller

// This example watches Pods and enqueues Requests with the Name and Namespace of the Pod from
// the Event (i.e. change caused by a Create, Update, Delete).
func ExampleEnqueueRequestForObject() {
	// over_controller is a over_controller.over_controller
	err := c.Watch(
		over_source.Kind(mgr.GetCache(), &corev1.Pod{}),
		&over_handler.EnqueueRequestForObject{},
	)
	if err != nil {
		// handle it
	}
}

// This example watches ReplicaSets and enqueues a Request containing the Name and Namespace of the
// owning (direct) Deployment responsible for the creation of the ReplicaSet.
func ExampleEnqueueRequestForOwner() {
	// over_controller is a over_controller.over_controller
	err := c.Watch(
		over_source.Kind(mgr.GetCache(), &appsv1.ReplicaSet{}),
		over_handler.EnqueueRequestForOwner(mgr.GetScheme(), mgr.GetRESTMapper(), &appsv1.Deployment{}, over_handler.OnlyControllerOwner()),
	)
	if err != nil {
		// handle it
	}
}

// This example watches Deployments and enqueues a Request contain the Name and Namespace of different
// objects (of Type: MyKind) using a mapping function defined by the user.
func ExampleEnqueueRequestsFromMapFunc() {
	// over_controller is a over_controller.over_controller
	err := c.Watch(
		over_source.Kind(mgr.GetCache(), &appsv1.Deployment{}),
		over_handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, a over_client.Object) []over_reconcile.Request {
			return []over_reconcile.Request{
				{NamespacedName: types.NamespacedName{
					Name:      a.GetName() + "-1",
					Namespace: a.GetNamespace(),
				}},
				{NamespacedName: types.NamespacedName{
					Name:      a.GetName() + "-2",
					Namespace: a.GetNamespace(),
				}},
			}
		}),
	)
	if err != nil {
		// handle it
	}
}

// This example implements over_handler.EnqueueRequestForObject.
func ExampleFuncs() {
	// over_controller is a over_controller.over_controller
	err := c.Watch(
		over_source.Kind(mgr.GetCache(), &corev1.Pod{}),
		over_handler.Funcs{
			CreateFunc: func(ctx context.Context, e over_event.CreateEvent, q workqueue.RateLimitingInterface) {
				q.Add(over_reconcile.Request{NamespacedName: types.NamespacedName{
					Name:      e.Object.GetName(),
					Namespace: e.Object.GetNamespace(),
				}})
			},
			UpdateFunc: func(ctx context.Context, e over_event.UpdateEvent, q workqueue.RateLimitingInterface) {
				q.Add(over_reconcile.Request{NamespacedName: types.NamespacedName{
					Name:      e.ObjectNew.GetName(),
					Namespace: e.ObjectNew.GetNamespace(),
				}})
			},
			DeleteFunc: func(ctx context.Context, e over_event.DeleteEvent, q workqueue.RateLimitingInterface) {
				q.Add(over_reconcile.Request{NamespacedName: types.NamespacedName{
					Name:      e.Object.GetName(),
					Namespace: e.Object.GetNamespace(),
				}})
			},
			GenericFunc: func(ctx context.Context, e over_event.GenericEvent, q workqueue.RateLimitingInterface) {
				q.Add(over_reconcile.Request{NamespacedName: types.NamespacedName{
					Name:      e.Object.GetName(),
					Namespace: e.Object.GetNamespace(),
				}})
			},
		},
	)
	if err != nil {
		// handle it
	}
}
