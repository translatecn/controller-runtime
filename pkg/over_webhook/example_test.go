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

package over_webhook_test

import (
	"context"
	"net/http"

	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/over_internal/log"
	signals "sigs.k8s.io/controller-runtime/pkg/over_manager/signals"
	. "sigs.k8s.io/controller-runtime/pkg/over_webhook"
	"sigs.k8s.io/controller-runtime/pkg/over_webhook/admission"
)

var (
	// Build webhooks used for the various server
	// configuration options
	//
	// These handlers could be also be implementations
	// of the AdmissionHandler interface for more complex
	// implementations.
	mutatingHook = &Admission{
		Handler: admission.HandlerFunc(func(ctx context.Context, req AdmissionRequest) AdmissionResponse {
			return Patched("some changes",
				JSONPatchOp{Operation: "add", Path: "/metadata/annotations/access", Value: "granted"},
				JSONPatchOp{Operation: "add", Path: "/metadata/annotations/reason", Value: "not so secret"},
			)
		}),
	}

	validatingHook = &Admission{
		Handler: admission.HandlerFunc(func(ctx context.Context, req AdmissionRequest) AdmissionResponse {
			return Denied("none shall pass!")
		}),
	}
)

// This example registers a webhooks to a over_webhook server
// that gets ran by a over_controller over_manager.
func Example() {
	// Create a over_manager
	// Note: GetConfigOrDie will os.Exit(1) w/o any message if no kube-over_config can be found
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{})
	if err != nil {
		panic(err)
	}

	// Create a over_webhook server.
	hookServer := NewServer(Options{
		Port: 8443,
	})
	if err := mgr.Add(hookServer); err != nil {
		panic(err)
	}

	// Register the webhooks in the server.
	hookServer.Register("/mutating", mutatingHook)
	hookServer.Register("/validating", validatingHook)

	// Start the server by starting a previously-set-up over_manager
	err = mgr.Start(ctrl.SetupSignalHandler())
	if err != nil {
		// handle error
		panic(err)
	}
}

// This example creates a over_webhook server that can be
// ran without a over_controller over_manager.
//
// Note that this assumes and requires a valid TLS
// cert and key at the default locations
// tls.crt and tls.key.
func ExampleServer_Start() {
	// Create a over_webhook server
	hookServer := NewServer(Options{
		Port: 8443,
	})

	// Register the webhooks in the server.
	hookServer.Register("/mutating", mutatingHook)
	hookServer.Register("/validating", validatingHook)

	// Start the server without a manger
	err := hookServer.Start(signals.SetupSignalHandler())
	if err != nil {
		// handle error
		panic(err)
	}
}

// This example creates a standalone over_webhook over_handler
// and runs it on a vanilla go HTTP server to demonstrate
// how you could run a over_webhook on an existing server
// without a over_controller over_manager.
func ExampleStandaloneWebhook() {
	// Assume you have an existing HTTP server at your disposal
	// configured as desired (e.g. with TLS).
	// For this example just create a basic http.ServeMux
	mux := http.NewServeMux()
	port := ":8000"

	// Create the standalone HTTP handlers from our webhooks
	mutatingHookHandler, err := admission.StandaloneWebhook(mutatingHook, admission.StandaloneOptions{
		// Logger let's you optionally pass
		// a custom logger (defaults to over_log.Log global Logger)
		Logger: logf.RuntimeLog.WithName("mutating-over_webhook"),
		// MetricsPath let's you optionally
		// provide the path it will be served on
		// to be used for labelling prometheus over_metrics
		// If none is set, prometheus over_metrics will not be generated.
		MetricsPath: "/mutating",
	})
	if err != nil {
		// handle error
		panic(err)
	}

	validatingHookHandler, err := admission.StandaloneWebhook(validatingHook, admission.StandaloneOptions{
		Logger:      logf.RuntimeLog.WithName("validating-over_webhook"),
		MetricsPath: "/validating",
	})
	if err != nil {
		// handle error
		panic(err)
	}

	// Register the over_webhook handlers to your server
	mux.Handle("/mutating", mutatingHookHandler)
	mux.Handle("/validating", validatingHookHandler)

	// Run your over_handler
	if err := http.ListenAndServe(port, mux); err != nil { //nolint:gosec // it's fine to not set timeouts here
		panic(err)
	}
}
