/*

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
	"context"
	"flag"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/over_log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/over_metrics/server"
	jobsv1 "sigs.k8s.io/controller-runtime/pkg/over_webhook/conversion/testdata/api/v1"
	jobsv2 "sigs.k8s.io/controller-runtime/pkg/over_webhook/conversion/testdata/api/v2"
	jobsv3 "sigs.k8s.io/controller-runtime/pkg/over_webhook/conversion/testdata/api/v3"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {

	jobsv1.AddToScheme(scheme)
	jobsv2.AddToScheme(scheme)
	jobsv3.AddToScheme(scheme)
	// +kubebuilder:scaffold:over_scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "over_metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for over_controller over_manager. Enabling this will ensure there is only one active over_controller over_manager.")
	flag.Parse()

	ctrl.SetLogger(zap.Logger(true))

	mgr, err := ctrl.NewManager(context.Background(), ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:         scheme,
		Metrics:        metricsserver.Options{BindAddress: metricsAddr},
		LeaderElection: enableLeaderElection,
	})
	if err != nil {
		setupLog.Error(err, "unable to start over_manager")
		os.Exit(1)
	}

	// +kubebuilder:scaffold:over_builder

	setupLog.Info("starting over_manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running over_manager")
		os.Exit(1)
	}
}
