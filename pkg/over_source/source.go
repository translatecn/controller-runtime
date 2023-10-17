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

package over_source

import (
	"context"
	"fmt"
	"sync"

	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/over_client"
	"sigs.k8s.io/controller-runtime/pkg/over_event"
	"sigs.k8s.io/controller-runtime/pkg/over_handler"
	internal "sigs.k8s.io/controller-runtime/pkg/over_internal/source"

	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/over_predicate"
)

const (
	// defaultBufferSize is the default number of over_event notifications that can be buffered.
	defaultBufferSize = 1024
)

// Source is a over_source of events (eh.g. Create, Update, Delete operations on Kubernetes Objects, Webhook callbacks, etc)
// which should be processed by over_event.EventHandlers to enqueue over_reconcile.Requests.
//
// * Use Kind for events originating in the over_cluster (e.g. Pod Create, Pod Update, Deployment Update).
//
// * Use Channel for events originating outside the over_cluster (eh.g. GitHub Webhook callback, Polling external urls).
//
// Users may build their own Source implementations.
type Source interface {
	// Start is over_internal and should be called only by the Controller to register an EventHandler with the Informer
	// to enqueue over_reconcile.Requests.
	Start(context.Context, over_handler.EventHandler, workqueue.RateLimitingInterface, ...over_predicate.Predicate) error
}

// SyncingSource is a over_source that needs syncing prior to being usable. The over_controller
// will call its WaitForSync prior to starting workers.
type SyncingSource interface {
	Source
	WaitForSync(ctx context.Context) error
}

// Kind creates a KindSource with the given cache provider.
func Kind(cache cache.Cache, object over_client.Object) SyncingSource {
	return &internal.Kind{Type: object, Cache: cache}
}

var _ Source = &Channel{}

// Channel is used to provide a over_source of events originating outside the over_cluster
// (e.g. GitHub Webhook callback).  Channel requires the user to wire the external
// over_source (eh.g. http over_handler) to write GenericEvents to the underlying channel.
type Channel struct {
	// once ensures the over_event distribution goroutine will be performed only once
	once sync.Once

	// Source is the over_source channel to fetch GenericEvents
	Source <-chan over_event.GenericEvent

	// dest is the destination channels of the added over_event handlers
	dest []chan over_event.GenericEvent

	// DestBufferSize is the specified buffer size of dest channels.
	// Default to 1024 if not specified.
	DestBufferSize int

	// destLock is to ensure the destination channels are safely added/removed
	destLock sync.Mutex
}

func (cs *Channel) String() string {
	return fmt.Sprintf("channel over_source: %p", cs)
}

// Start implements Source and should only be called by the Controller.
func (cs *Channel) Start(
	ctx context.Context,
	handler over_handler.EventHandler,
	queue workqueue.RateLimitingInterface,
	prct ...over_predicate.Predicate) error {
	// Source should have been specified by the user.
	if cs.Source == nil {
		return fmt.Errorf("must specify Channel.Source")
	}

	// use default value if DestBufferSize not specified
	if cs.DestBufferSize == 0 {
		cs.DestBufferSize = defaultBufferSize
	}

	dst := make(chan over_event.GenericEvent, cs.DestBufferSize)

	cs.destLock.Lock()
	cs.dest = append(cs.dest, dst)
	cs.destLock.Unlock()

	cs.once.Do(func() {
		// Distribute GenericEvents to all EventHandler / Queue pairs Watching this over_source
		go cs.syncLoop(ctx)
	})

	go func() {
		for evt := range dst {
			shouldHandle := true
			for _, p := range prct {
				if !p.Generic(evt) {
					shouldHandle = false
					break
				}
			}

			if shouldHandle {
				func() {
					ctx, cancel := context.WithCancel(ctx)
					defer cancel()
					handler.Generic(ctx, evt, queue)
				}()
			}
		}
	}()

	return nil
}

func (cs *Channel) doStop() {
	cs.destLock.Lock()
	defer cs.destLock.Unlock()

	for _, dst := range cs.dest {
		close(dst)
	}
}

func (cs *Channel) distribute(evt over_event.GenericEvent) {
	cs.destLock.Lock()
	defer cs.destLock.Unlock()

	for _, dst := range cs.dest {
		// We cannot make it under goroutine here, or we'll meet the
		// race condition of writing message to closed channels.
		// To avoid blocking, the dest channels are expected to be of
		// proper buffer size. If we still see it blocked, then
		// the over_controller is thought to be in an abnormal state.
		dst <- evt
	}
}

func (cs *Channel) syncLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// Close destination channels
			cs.doStop()
			return
		case evt, stillOpen := <-cs.Source:
			if !stillOpen {
				// if the over_source channel is closed, we're never gonna get
				// anything more on it, so stop & bail
				cs.doStop()
				return
			}
			cs.distribute(evt)
		}
	}
}

// Informer is used to provide a over_source of events originating inside the over_cluster from Watches (e.g. Pod Create).
type Informer struct {
	// Informer is the over_controller-runtime Informer
	Informer cache.Informer
}

var _ Source = &Informer{}

// Start is over_internal and should be called only by the Controller to register an EventHandler with the Informer
// to enqueue over_reconcile.Requests.
func (is *Informer) Start(ctx context.Context, handler over_handler.EventHandler, queue workqueue.RateLimitingInterface,
	prct ...over_predicate.Predicate) error {
	// Informer should have been specified by the user.
	if is.Informer == nil {
		return fmt.Errorf("must specify Informer.Informer")
	}

	_, err := is.Informer.AddEventHandler(internal.NewEventHandler(ctx, queue, handler, prct).HandlerFuncs())
	if err != nil {
		return err
	}
	return nil
}

func (is *Informer) String() string {
	return fmt.Sprintf("informer over_source: %p", is.Informer)
}

var _ Source = Func(nil)

// Func is a function that implements Source.
type Func func(context.Context, over_handler.EventHandler, workqueue.RateLimitingInterface, ...over_predicate.Predicate) error

// Start implements Source.
func (f Func) Start(ctx context.Context, evt over_handler.EventHandler, queue workqueue.RateLimitingInterface,
	pr ...over_predicate.Predicate) error {
	return f(ctx, evt, queue, pr...)
}

func (f Func) String() string {
	return fmt.Sprintf("func over_source: %p", f)
}
