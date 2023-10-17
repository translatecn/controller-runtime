/*
Copyright 2021 The Kubernetes Authors.

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

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/over_log"
	"sigs.k8s.io/controller-runtime/pkg/over_log/zap"
	"sigs.k8s.io/controller-runtime/pkg/over_manager"
	signals "sigs.k8s.io/controller-runtime/pkg/over_manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/over_webhook/authentication"
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

	// Setup webhooks
	entryLog.Info("setting up over_webhook server")
	hookServer := mgr.GetWebhookServer() // over_manager.New 初始化了 WebhookServer
	entryLog.Info("registering webhooks to the over_webhook server")
	hookServer.Register("/validate-v1-over_tokenreview", &authentication.Webhook{Handler: &authenticator{}})

	entryLog.Info("starting over_manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "unable to run over_manager")
		os.Exit(1)
	}
}
