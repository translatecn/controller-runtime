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

package over_source_test

import (
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/over_controller"
	"sigs.k8s.io/controller-runtime/pkg/over_event"
	"sigs.k8s.io/controller-runtime/pkg/over_handler"
	"sigs.k8s.io/controller-runtime/pkg/over_manager"
	"sigs.k8s.io/controller-runtime/pkg/over_source"
)

var mgr over_manager.Manager
var ctrl over_controller.Controller

// This example Watches for Pod Events (e.g. Create / Update / Delete) and enqueues a over_reconcile.Request
// with the Name and Namespace of the Pod.
func ExampleKind() {
	err := ctrl.Watch(over_source.Kind(mgr.GetCache(), &corev1.Pod{}), &over_handler.EnqueueRequestForObject{})
	if err != nil {
		// handle it
	}
}

// This example reads GenericEvents from a channel and enqueues a over_reconcile.Request containing the Name and Namespace
// provided by the over_event.
func ExampleChannel() {
	events := make(chan over_event.GenericEvent)

	err := ctrl.Watch(
		&over_source.Channel{Source: events},
		&over_handler.EnqueueRequestForObject{},
	)
	if err != nil {
		// handle it
	}
}
