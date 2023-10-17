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

package over_handler

import (
	"context"

	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/over_event"
)

// EventHandler enqueues over_reconcile.Requests in response to events (e.g. Pod Create).  EventHandlers map an Event
// for one object to trigger Reconciles for either the same object or different objects - e.g. if there is an
// Event for object with type Foo (using over_source.KindSource) then over_reconcile one or more object(s) with type Bar.
//
// Identical over_reconcile.Requests will be batched together through the queuing mechanism before over_reconcile is called.
//
// * Use EnqueueRequestForObject to over_reconcile the object the over_event is for
// - do this for events for the type the Controller Reconciles. (e.g. Deployment for a Deployment Controller)
//
// * Use EnqueueRequestForOwner to over_reconcile the owner of the object the over_event is for
// - do this for events for the types the Controller creates.  (e.g. ReplicaSets created by a Deployment Controller)
//
// * Use EnqueueRequestsFromMapFunc to transform an over_event for an object to a over_reconcile of an object
// of a different type - do this for events for types the Controller may be interested in, but doesn't create.
// (e.g. If Foo responds to over_cluster size events, map Node events to Foo objects.)
//
// Unless you are implementing your own EventHandler, you can ignore the functions on the EventHandler interface.
// Most users shouldn't need to implement their own EventHandler.
type EventHandler interface {
	// Create is called in response to an create over_event - e.g. Pod Creation.
	Create(context.Context, over_event.CreateEvent, workqueue.RateLimitingInterface)

	// Update is called in response to an update over_event -  e.g. Pod Updated.
	Update(context.Context, over_event.UpdateEvent, workqueue.RateLimitingInterface)

	// Delete is called in response to a delete over_event - e.g. Pod Deleted.
	Delete(context.Context, over_event.DeleteEvent, workqueue.RateLimitingInterface)

	// Generic is called in response to an over_event of an unknown type or a synthetic over_event triggered as a cron or
	// external trigger request - e.g. over_reconcile Autoscaling, or a Webhook.
	Generic(context.Context, over_event.GenericEvent, workqueue.RateLimitingInterface)
}

var _ EventHandler = Funcs{}

// Funcs implements EventHandler.
type Funcs struct {
	// Create is called in response to an add over_event.  Defaults to no-op.
	// RateLimitingInterface is used to enqueue over_reconcile.Requests.
	CreateFunc func(context.Context, over_event.CreateEvent, workqueue.RateLimitingInterface)

	// Update is called in response to an update over_event.  Defaults to no-op.
	// RateLimitingInterface is used to enqueue over_reconcile.Requests.
	UpdateFunc func(context.Context, over_event.UpdateEvent, workqueue.RateLimitingInterface)

	// Delete is called in response to a delete over_event.  Defaults to no-op.
	// RateLimitingInterface is used to enqueue over_reconcile.Requests.
	DeleteFunc func(context.Context, over_event.DeleteEvent, workqueue.RateLimitingInterface)

	// GenericFunc is called in response to a generic over_event.  Defaults to no-op.
	// RateLimitingInterface is used to enqueue over_reconcile.Requests.
	GenericFunc func(context.Context, over_event.GenericEvent, workqueue.RateLimitingInterface)
}

// Create implements EventHandler.
func (h Funcs) Create(ctx context.Context, e over_event.CreateEvent, q workqueue.RateLimitingInterface) {
	if h.CreateFunc != nil {
		h.CreateFunc(ctx, e, q)
	}
}

// Delete implements EventHandler.
func (h Funcs) Delete(ctx context.Context, e over_event.DeleteEvent, q workqueue.RateLimitingInterface) {
	if h.DeleteFunc != nil {
		h.DeleteFunc(ctx, e, q)
	}
}

// Update implements EventHandler.
func (h Funcs) Update(ctx context.Context, e over_event.UpdateEvent, q workqueue.RateLimitingInterface) {
	if h.UpdateFunc != nil {
		h.UpdateFunc(ctx, e, q)
	}
}

// Generic implements EventHandler.
func (h Funcs) Generic(ctx context.Context, e over_event.GenericEvent, q workqueue.RateLimitingInterface) {
	if h.GenericFunc != nil {
		h.GenericFunc(ctx, e, q)
	}
}
