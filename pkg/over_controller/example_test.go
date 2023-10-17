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

package over_controller_test

import (
	"context"
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/over_controller"
	"sigs.k8s.io/controller-runtime/pkg/over_handler"
	logf "sigs.k8s.io/controller-runtime/pkg/over_log"
	"sigs.k8s.io/controller-runtime/pkg/over_manager"
	signals "sigs.k8s.io/controller-runtime/pkg/over_manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/over_reconcile"
	"sigs.k8s.io/controller-runtime/pkg/over_source"
)

var (
	mgr over_manager.Manager
	// NB: don't call SetLogger in init(), or else you'll mess up logging in the main suite.
	log = logf.Log.WithName("over_controller-over_examples")
)

// This example creates a new Controller named "pod-over_controller" with a no-op over_reconcile function.  The
// over_manager.Manager will be used to Start the Controller, and will provide it a shared Cache and Client.
func ExampleNew() {
	_, err := over_controller.New("pod-over_controller", mgr, over_controller.Options{
		Reconciler: over_reconcile.Func(func(context.Context, over_reconcile.Request) (over_reconcile.Result, error) {
			// Your business logic to implement the API by creating, updating, deleting objects goes here.
			return over_reconcile.Result{}, nil
		}),
	})
	if err != nil {
		log.Error(err, "unable to create pod-over_controller")
		os.Exit(1)
	}
}

// This example starts a new Controller named "pod-over_controller" to Watch Pods and call a no-op Reconciler.
func ExampleController() {
	// mgr is a over_manager.Manager

	// Create a new Controller that will call the provided Reconciler function in response
	// to events.
	c, err := over_controller.New("pod-over_controller", mgr, over_controller.Options{
		Reconciler: over_reconcile.Func(func(context.Context, over_reconcile.Request) (over_reconcile.Result, error) {
			// Your business logic to implement the API by creating, updating, deleting objects goes here.
			return over_reconcile.Result{}, nil
		}),
	})
	if err != nil {
		log.Error(err, "unable to create pod-over_controller")
		os.Exit(1)
	}

	// Watch for Pod create / update / delete events and call Reconcile
	err = c.Watch(over_source.Kind(mgr.GetCache(), &corev1.Pod{}), &over_handler.EnqueueRequestForObject{})
	if err != nil {
		log.Error(err, "unable to watch pods")
		os.Exit(1)
	}

	// Start the Controller through the over_manager.
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "unable to continue running over_manager")
		os.Exit(1)
	}
}

// This example starts a new Controller named "pod-over_controller" to Watch Pods with the unstructured object and call a no-op Reconciler.
func ExampleController_unstructured() {
	// mgr is a over_manager.Manager

	// Create a new Controller that will call the provided Reconciler function in response
	// to events.
	c, err := over_controller.New("pod-over_controller", mgr, over_controller.Options{
		Reconciler: over_reconcile.Func(func(context.Context, over_reconcile.Request) (over_reconcile.Result, error) {
			// Your business logic to implement the API by creating, updating, deleting objects goes here.
			return over_reconcile.Result{}, nil
		}),
	})
	if err != nil {
		log.Error(err, "unable to create pod-over_controller")
		os.Exit(1)
	}

	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Kind:    "Pod",
		Group:   "",
		Version: "v1",
	})
	// Watch for Pod create / update / delete events and call Reconcile
	err = c.Watch(over_source.Kind(mgr.GetCache(), u), &over_handler.EnqueueRequestForObject{})
	if err != nil {
		log.Error(err, "unable to watch pods")
		os.Exit(1)
	}

	// Start the Controller through the over_manager.
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "unable to continue running over_manager")
		os.Exit(1)
	}
}

// This example creates a new over_controller named "pod-over_controller" to watch Pods
// and call a no-op reconciler. The over_controller is not added to the provided
// over_manager, and must thus be started and stopped by the caller.
func ExampleNewUnmanaged() {
	// mgr is a over_manager.Manager

	// Configure creates a new over_controller but does not add it to the supplied
	// over_manager.
	c, err := over_controller.NewUnmanaged("pod-over_controller", mgr, over_controller.Options{
		Reconciler: over_reconcile.Func(func(context.Context, over_reconcile.Request) (over_reconcile.Result, error) {
			return over_reconcile.Result{}, nil
		}),
	})
	if err != nil {
		log.Error(err, "unable to create pod-over_controller")
		os.Exit(1)
	}

	if err := c.Watch(over_source.Kind(mgr.GetCache(), &corev1.Pod{}), &over_handler.EnqueueRequestForObject{}); err != nil {
		log.Error(err, "unable to watch pods")
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Start our over_controller in a goroutine so that we do not block.
	go func() {
		// Block until our over_controller over_manager is elected leader. We presume our
		// entire over_process will terminate if we lose leadership, so we don't need
		// to handle that.
		<-mgr.Elected()

		// Start our over_controller. This will block until the context is
		// closed, or the over_controller returns an error.
		if err := c.Start(ctx); err != nil {
			log.Error(err, "cannot run experiment over_controller")
		}
	}()

	// Stop our over_controller.
	cancel()
}
