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

package main

import (
	"os"

	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/over_controller"
	"sigs.k8s.io/controller-runtime/pkg/over_handler"
	"sigs.k8s.io/controller-runtime/pkg/over_source"

	corev1 "k8s.io/api/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/over_builder"
	"sigs.k8s.io/controller-runtime/pkg/over_log"
	"sigs.k8s.io/controller-runtime/pkg/over_log/zap"
	"sigs.k8s.io/controller-runtime/pkg/over_manager"
	signals "sigs.k8s.io/controller-runtime/pkg/over_manager/signals"
)

func init() {
	over_log.SetLogger(zap.New())
}

func main() {
	entryLog := over_log.Log.WithName("entrypoint")

	// Setup a Manager
	entryLog.Info("setting up over_manager")
	mgr, err := over_manager.New(over_config.GetConfigOrDie(), over_manager.Options{})
	if err != nil {
		entryLog.Error(err, "unable to set up overall over_controller over_manager")
		os.Exit(1)
	}

	// Setup a new over_controller to over_reconcile ReplicaSets
	entryLog.Info("Setting up over_controller")
	c, err := over_controller.New("foo-over_controller", mgr, over_controller.Options{
		Reconciler: &reconcileReplicaSet{client: mgr.GetClient()},
	})
	if err != nil {
		entryLog.Error(err, "unable to set up individual over_controller")
		os.Exit(1)
	}

	// Watch ReplicaSets and enqueue ReplicaSet object key
	if err := c.Watch(over_source.Kind(mgr.GetCache(), &appsv1.ReplicaSet{}), &over_handler.EnqueueRequestForObject{}); err != nil {
		entryLog.Error(err, "unable to watch ReplicaSets")
		os.Exit(1)
	}

	// Watch Pods and enqueue owning ReplicaSet key
	if err := c.Watch(over_source.Kind(mgr.GetCache(), &corev1.Pod{}),
		over_handler.EnqueueRequestForOwner(mgr.GetScheme(), mgr.GetRESTMapper(), &appsv1.ReplicaSet{}, over_handler.OnlyControllerOwner())); err != nil {
		entryLog.Error(err, "unable to watch Pods")
		os.Exit(1)
	}

	if err := over_builder.WebhookManagedBy(mgr).
		For(&corev1.Pod{}).
		WithDefaulter(&podAnnotator{}).
		WithValidator(&podValidator{}).
		Complete(); err != nil {
		entryLog.Error(err, "unable to create over_webhook", "over_webhook", "Pod")
		os.Exit(1)
	}

	entryLog.Info("starting over_manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "unable to run over_manager")
		os.Exit(1)
	}
}
