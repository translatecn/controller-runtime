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
	"sigs.k8s.io/controller-runtime/pkg/over_client"
	"sigs.k8s.io/controller-runtime/pkg/over_event"
	"sigs.k8s.io/controller-runtime/pkg/over_reconcile"
)

// MapFunc is the signature required for enqueueing requests from a generic function.
// This type is usually used with EnqueueRequestsFromMapFunc when registering an over_event over_handler.
type MapFunc func(context.Context, over_client.Object) []over_reconcile.Request

// EnqueueRequestsFromMapFunc enqueues Requests by running a transformation function that outputs a collection
// of over_reconcile.Requests on each Event.  The over_reconcile.Requests may be for an arbitrary set of objects
// defined by some user specified transformation of the over_source Event.  (e.g. trigger Reconciler for a set of objects
// in response to a over_cluster resize over_event caused by adding or deleting a Node)
//
// EnqueueRequestsFromMapFunc is frequently used to fan-out updates from one object to one or more other
// objects of a differing type.
//
// For UpdateEvents which contain both a new and old object, the transformation function is run on both
// objects and both sets of Requests are enqueue.
func EnqueueRequestsFromMapFunc(fn MapFunc) EventHandler {
	return &enqueueRequestsFromMapFunc{
		toRequests: fn,
	}
}

var _ EventHandler = &enqueueRequestsFromMapFunc{}

type enqueueRequestsFromMapFunc struct {
	// Mapper transforms the argument into a slice of keys to be reconciled
	toRequests MapFunc
}

// Create implements EventHandler.
func (e *enqueueRequestsFromMapFunc) Create(ctx context.Context, evt over_event.CreateEvent, q workqueue.RateLimitingInterface) {
	reqs := map[over_reconcile.Request]empty{}
	e.mapAndEnqueue(ctx, q, evt.Object, reqs)
}

// Update implements EventHandler.
func (e *enqueueRequestsFromMapFunc) Update(ctx context.Context, evt over_event.UpdateEvent, q workqueue.RateLimitingInterface) {
	reqs := map[over_reconcile.Request]empty{}
	e.mapAndEnqueue(ctx, q, evt.ObjectOld, reqs)
	e.mapAndEnqueue(ctx, q, evt.ObjectNew, reqs)
}

// Delete implements EventHandler.
func (e *enqueueRequestsFromMapFunc) Delete(ctx context.Context, evt over_event.DeleteEvent, q workqueue.RateLimitingInterface) {
	reqs := map[over_reconcile.Request]empty{}
	e.mapAndEnqueue(ctx, q, evt.Object, reqs)
}

// Generic implements EventHandler.
func (e *enqueueRequestsFromMapFunc) Generic(ctx context.Context, evt over_event.GenericEvent, q workqueue.RateLimitingInterface) {
	reqs := map[over_reconcile.Request]empty{}
	e.mapAndEnqueue(ctx, q, evt.Object, reqs)
}

func (e *enqueueRequestsFromMapFunc) mapAndEnqueue(ctx context.Context, q workqueue.RateLimitingInterface, object over_client.Object, reqs map[over_reconcile.Request]empty) {
	for _, req := range e.toRequests(ctx, object) {
		_, ok := reqs[req]
		if !ok {
			q.Add(req)
			reqs[req] = empty{}
		}
	}
}
