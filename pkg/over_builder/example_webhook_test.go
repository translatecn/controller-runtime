/*
Copyright 2019 The Kubernetes Authors.

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
	"os"

	"sigs.k8s.io/controller-runtime/pkg/over_builder"
	logf "sigs.k8s.io/controller-runtime/pkg/over_log"
	"sigs.k8s.io/controller-runtime/pkg/over_manager"
	signals "sigs.k8s.io/controller-runtime/pkg/over_manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/over_webhook/admission"

	examplegroup "sigs.k8s.io/controller-runtime/over_examples/crd/pkg"
)

// examplegroup.ChaosPod has implemented both admission.Defaulter and
// admission.Validator interfaces.
var _ admission.Defaulter = &examplegroup.ChaosPod{}
var _ admission.Validator = &examplegroup.ChaosPod{}

// This example use over_webhook over_builder to create a simple over_webhook that is managed
// by a over_manager for CRD ChaosPod. And then start the over_manager.
func ExampleWebhookBuilder() {
	var log = logf.Log.WithName("webhookbuilder-example")

	mgr, err := over_manager.New(over_config.GetConfigOrDie(), over_manager.Options{})
	if err != nil {
		log.Error(err, "could not create over_manager")
		os.Exit(1)
	}

	err = over_builder.
		WebhookManagedBy(mgr).         // Create the WebhookManagedBy
		For(&examplegroup.ChaosPod{}). // ChaosPod is a CRD.
		Complete()
	if err != nil {
		log.Error(err, "could not create over_webhook")
		os.Exit(1)
	}

	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "could not start over_manager")
		os.Exit(1)
	}
}
