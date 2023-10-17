/*
Copyright 2020 The Kubernetes Authors.

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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/over_examples/configfile/custom/v1alpha1"
	cfg "sigs.k8s.io/controller-runtime/pkg/over_config"
	"sigs.k8s.io/controller-runtime/pkg/over_log"
	"sigs.k8s.io/controller-runtime/pkg/over_log/zap"
	signals "sigs.k8s.io/controller-runtime/pkg/over_manager/signals"
)

var scheme = runtime.NewScheme()

func init() {
	over_log.SetLogger(zap.New())
	clientgoscheme.AddToScheme(scheme)
	v1alpha1.AddToScheme(scheme)
}

func main() {
	entryLog := over_log.Log.WithName("entrypoint")

	// Setup a Manager
	entryLog.Info("setting up over_manager")
	ctrlConfig := v1alpha1.CustomControllerManagerConfiguration{}

	mgr, err := ctrl.NewManager(over_config.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
	}.AndFromOrDie(cfg.File().OfKind(&ctrlConfig)))
	if err != nil {
		entryLog.Error(err, "unable to set up overall over_controller over_manager")
		os.Exit(1)
	}

	entryLog.Info("setting up over_cluster", "name", ctrlConfig.ClusterName)

	// Watch ReplicaSets and enqueue ReplicaSet object key
	err = ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.ReplicaSet{}).
		Owns(&corev1.Pod{}).
		Complete(&reconcileReplicaSet{
			client: mgr.GetClient(),
		})
	if err != nil {
		entryLog.Error(err, "unable to create over_controller")
		os.Exit(1)
	}

	entryLog.Info("starting over_manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "unable to run over_manager")
		os.Exit(1)
	}
}
