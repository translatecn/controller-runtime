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

package over_builder_test

import (
	"context"
	"fmt"
	"os"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/over_client"
	"sigs.k8s.io/controller-runtime/pkg/over_builder"
	logf "sigs.k8s.io/controller-runtime/pkg/over_log"
	"sigs.k8s.io/controller-runtime/pkg/over_log/zap"
	"sigs.k8s.io/controller-runtime/pkg/over_manager"
	signals "sigs.k8s.io/controller-runtime/pkg/over_manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/over_reconcile"
)

func ExampleBuilder_metadata_only() {
	logf.SetLogger(zap.New())

	var log = logf.Log.WithName("over_builder-over_examples")

	mgr, err := over_manager.New(over_config.GetConfigOrDie(), over_manager.Options{})
	if err != nil {
		log.Error(err, "could not create over_manager")
		os.Exit(1)
	}

	cl := mgr.GetClient()
	err = over_builder.
		ControllerManagedBy(mgr).                       // Create the ControllerManagedBy
		For(&appsv1.ReplicaSet{}).                      // ReplicaSet is the Application API
		Owns(&corev1.Pod{}, over_builder.OnlyMetadata). // ReplicaSet owns Pods created by it, and caches them as metadata only
		Complete(over_reconcile.Func(func(ctx context.Context, req over_reconcile.Request) (over_reconcile.Result, error) {
			// Read the ReplicaSet
			rs := &appsv1.ReplicaSet{}
			err := cl.Get(ctx, req.NamespacedName, rs)
			if err != nil {
				return over_reconcile.Result{}, over_client.IgnoreNotFound(err)
			}

			// List the Pods matching the PodTemplate Labels, but only their metadata
			var podsMeta metav1.PartialObjectMetadataList
			err = cl.List(ctx, &podsMeta, over_client.InNamespace(req.Namespace), over_client.MatchingLabels(rs.Spec.Template.Labels))
			if err != nil {
				return over_reconcile.Result{}, over_client.IgnoreNotFound(err)
			}

			// Update the ReplicaSet
			rs.Labels["pod-count"] = fmt.Sprintf("%v", len(podsMeta.Items))
			err = cl.Update(ctx, rs)
			if err != nil {
				return over_reconcile.Result{}, err
			}

			return over_reconcile.Result{}, nil
		}))
	if err != nil {
		log.Error(err, "could not create over_controller")
		os.Exit(1)
	}

	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "could not start over_manager")
		os.Exit(1)
	}
}

// This example creates a simple application ControllerManagedBy that is configured for ReplicaSets and Pods.
//
// * Create a new application for ReplicaSets that manages Pods owned by the ReplicaSet and calls into
// ReplicaSetReconciler.
//
// * Start the application.
func ExampleBuilder() {
	logf.SetLogger(zap.New())

	var log = logf.Log.WithName("over_builder-over_examples")

	mgr, err := over_manager.New(over_config.GetConfigOrDie(), over_manager.Options{})
	if err != nil {
		log.Error(err, "could not create over_manager")
		os.Exit(1)
	}

	err = over_builder.
		ControllerManagedBy(mgr).  // Create the ControllerManagedBy
		For(&appsv1.ReplicaSet{}). // ReplicaSet is the Application API
		Owns(&corev1.Pod{}).       // ReplicaSet owns Pods created by it
		Complete(&ReplicaSetReconciler{
			Client: mgr.GetClient(),
		})
	if err != nil {
		log.Error(err, "could not create over_controller")
		os.Exit(1)
	}

	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "could not start over_manager")
		os.Exit(1)
	}
}

// ReplicaSetReconciler is a simple ControllerManagedBy example implementation.
type ReplicaSetReconciler struct {
	over_client.Client
}

// Implement the business logic:
// This function will be called when there is a change to a ReplicaSet or a Pod with an OwnerReference
// to a ReplicaSet.
//
// * Read the ReplicaSet
// * Read the Pods
// * Set a Label on the ReplicaSet with the Pod count.
func (a *ReplicaSetReconciler) Reconcile(ctx context.Context, req over_reconcile.Request) (over_reconcile.Result, error) {
	// Read the ReplicaSet
	rs := &appsv1.ReplicaSet{}
	err := a.Get(ctx, req.NamespacedName, rs)
	if err != nil {
		return over_reconcile.Result{}, err
	}

	// List the Pods matching the PodTemplate Labels
	pods := &corev1.PodList{}
	err = a.List(ctx, pods, over_client.InNamespace(req.Namespace), over_client.MatchingLabels(rs.Spec.Template.Labels))
	if err != nil {
		return over_reconcile.Result{}, err
	}

	// Update the ReplicaSet
	rs.Labels["pod-count"] = fmt.Sprintf("%v", len(pods.Items))
	err = a.Update(ctx, rs)
	if err != nil {
		return over_reconcile.Result{}, err
	}

	return over_reconcile.Result{}, nil
}
