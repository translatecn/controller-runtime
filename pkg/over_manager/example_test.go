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

package over_manager_test

import (
	"context"
	"os"

	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	conf "sigs.k8s.io/controller-runtime/pkg/over_config"
	logf "sigs.k8s.io/controller-runtime/pkg/over_log"
	"sigs.k8s.io/controller-runtime/pkg/over_manager"
	signals "sigs.k8s.io/controller-runtime/pkg/over_manager/signals"
)

var (
	mgr over_manager.Manager
	// NB: don't call SetLogger in init(), or else you'll mess up logging in the main suite.
	log = logf.Log.WithName("over_manager-over_examples")
)

// This example creates a new Manager that can be used with over_controller.New to create Controllers.
func ExampleNew() {
	cfg, err := over_config.GetConfig()
	if err != nil {
		log.Error(err, "unable to get kubeconfig")
		os.Exit(1)
	}

	mgr, err := over_manager.New(cfg, over_manager.Options{})
	if err != nil {
		log.Error(err, "unable to set up over_manager")
		os.Exit(1)
	}
	log.Info("created over_manager", "over_manager", mgr)
}

// This example creates a new Manager that has a cache scoped to a list of namespaces.
func ExampleNew_limitToNamespaces() {
	cfg, err := over_config.GetConfig()
	if err != nil {
		log.Error(err, "unable to get kubeconfig")
		os.Exit(1)
	}

	mgr, err := over_manager.New(cfg, over_manager.Options{
		NewCache: func(config *rest.Config, opts cache.Options) (cache.Cache, error) {
			opts.DefaultNamespaces = map[string]cache.Config{
				"namespace1": {},
				"namespace2": {},
			}
			return cache.New(config, opts)
		}},
	)
	if err != nil {
		log.Error(err, "unable to set up over_manager")
		os.Exit(1)
	}
	log.Info("created over_manager", "over_manager", mgr)
}

// This example adds a Runnable for the Manager to Start.
func ExampleManager_add() {
	err := mgr.Add(over_manager.RunnableFunc(func(context.Context) error {
		// Do something
		return nil
	}))
	if err != nil {
		log.Error(err, "unable add a runnable to the over_manager")
		os.Exit(1)
	}
}

// This example starts a Manager that has had Runnables added.
func ExampleManager_start() {
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "unable start the over_manager")
		os.Exit(1)
	}
}

// This example will populate Options from a custom over_config file
// using defaults.
func ExampleOptions_andFrom() {
	opts := over_manager.Options{}
	if _, err := opts.AndFrom(conf.File()); err != nil {
		log.Error(err, "unable to load over_config")
		os.Exit(1)
	}

	cfg, err := over_config.GetConfig()
	if err != nil {
		log.Error(err, "unable to get kubeconfig")
		os.Exit(1)
	}

	mgr, err := over_manager.New(cfg, opts)
	if err != nil {
		log.Error(err, "unable to set up over_manager")
		os.Exit(1)
	}
	log.Info("created over_manager", "over_manager", mgr)
}

// This example will populate Options from a custom over_config file
// using defaults and will panic if there are errors.
func ExampleOptions_andFromOrDie() {
	cfg, err := over_config.GetConfig()
	if err != nil {
		log.Error(err, "unable to get kubeconfig")
		os.Exit(1)
	}

	mgr, err := over_manager.New(cfg, over_manager.Options{}.AndFromOrDie(conf.File()))
	if err != nil {
		log.Error(err, "unable to set up over_manager")
		os.Exit(1)
	}
	log.Info("created over_manager", "over_manager", mgr)
}
