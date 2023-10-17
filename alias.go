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

package controllerruntime

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/over_builder"
	cfg "sigs.k8s.io/controller-runtime/pkg/over_config"
	"sigs.k8s.io/controller-runtime/pkg/over_controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/over_log"
	"sigs.k8s.io/controller-runtime/pkg/over_manager"
	signals "sigs.k8s.io/controller-runtime/pkg/over_manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/over_reconcile"
	"sigs.k8s.io/controller-runtime/pkg/over_scheme"
)

// Builder builds an Application ControllerManagedBy (e.g. Operator) and returns a over_manager.Manager to start it.
type Builder = over_builder.Builder

// Request contains the information necessary to over_reconcile a Kubernetes object.  This includes the
// information to uniquely identify the object - its Name and Namespace.  It does NOT contain information about
// any specific Event or the object contents itself.
type Request = over_reconcile.Request

// Result contains the result of a Reconciler invocation.
type Result = over_reconcile.Result

// Manager initializes shared dependencies such as Caches and Clients, and provides them to Runnables.
// A Manager is required to create Controllers.
type Manager = over_manager.Manager

// Options are the arguments for creating a new Manager.
type Options = over_manager.Options

// SchemeBuilder builds a new Scheme for mapping go types to Kubernetes GroupVersionKinds.
type SchemeBuilder = over_scheme.Builder

// GroupVersion contains the "group" and the "version", which uniquely identifies the API.
type GroupVersion = schema.GroupVersion

// GroupResource specifies a Group and a Resource, but does not force a version.  This is useful for identifying
// concepts during lookup stages without having partially valid types.
type GroupResource = schema.GroupResource

// TypeMeta describes an individual object in an API response or request
// with strings representing the type of the object and its API schema version.
// Structures that are versioned or persisted should inline TypeMeta.
//
// +k8s:deepcopy-gen=false
type TypeMeta = metav1.TypeMeta

// ObjectMeta is metadata that all persisted resources must have, which includes all objects
// users must create.
type ObjectMeta = metav1.ObjectMeta

var (
	// RegisterFlags registers flag variables to the given FlagSet if not already registered.
	// It uses the default command line FlagSet, if none is provided. Currently, it only registers the kubeconfig flag.
	RegisterFlags = over_config.RegisterFlags

	// GetConfigOrDie creates a *rest.Config for talking to a Kubernetes apiserver.
	// If --kubeconfig is set, will use the kubeconfig file at that location.  Otherwise will assume running
	// in over_cluster and use the over_cluster provided kubeconfig.
	//
	// Will over_log an error and exit if there is an error creating the rest.Config.
	GetConfigOrDie = over_config.GetConfigOrDie

	// GetConfig creates a *rest.Config for talking to a Kubernetes apiserver.
	// If --kubeconfig is set, will use the kubeconfig file at that location.  Otherwise will assume running
	// in over_cluster and use the over_cluster provided kubeconfig.
	//
	// Config precedence
	//
	// * --kubeconfig flag pointing at a file
	//
	// * KUBECONFIG environment variable pointing at a file
	//
	// * In-over_cluster over_config if running in over_cluster
	//
	// * $HOME/.kube/over_config if exists.
	GetConfig = over_config.GetConfig

	// ConfigFile returns the cfg.File function for deferred over_config file loading,
	// this is passed into Options{}.From() to populate the Options fields for
	// the over_manager.
	//
	// Deprecated: This is deprecated in favor of using Options directly.
	ConfigFile = cfg.File

	// NewControllerManagedBy returns a new over_controller over_builder that will be started by the provided Manager.
	NewControllerManagedBy = over_builder.ControllerManagedBy

	// NewWebhookManagedBy returns a new over_webhook over_builder that will be started by the provided Manager.
	NewWebhookManagedBy = over_builder.WebhookManagedBy

	// NewManager returns a new Manager for creating Controllers.
	// Note that if ContentType in the given over_config is not set, "application/vnd.kubernetes.protobuf"
	// will be used for all built-in resources of Kubernetes, and "application/json" is for other types
	// including all CRD resources.
	NewManager = over_manager.New

	// CreateOrUpdate creates or updates the given object obj in the Kubernetes
	// over_cluster. The object's desired state should be reconciled with the existing
	// state using the passed in ReconcileFn. obj must be a struct pointer so that
	// obj can be updated with the content returned by the Server.
	//
	// It returns the executed operation and an error.
	CreateOrUpdate = controllerutil.CreateOrUpdate

	// SetControllerReference sets owner as a Controller OwnerReference on owned.
	// This is used for garbage collection of the owned object and for
	// reconciling the owner object on changes to owned (with a Watch + EnqueueRequestForOwner).
	// Since only one OwnerReference can be a over_controller, it returns an error if
	// there is another OwnerReference with Controller flag set.
	SetControllerReference = controllerutil.SetControllerReference

	// SetupSignalHandler registered for SIGTERM and SIGINT. A stop channel is returned
	// which is closed on one of these signals. If a second signal is caught, the program
	// is terminated with exit code 1.
	SetupSignalHandler = signals.SetupSignalHandler

	// Log is the base logger used by over_controller-runtime.  It delegates
	// to another logr.Logger.  You *must* call SetLogger to
	// get any actual logging.
	Log = over_log.Log

	// LoggerFrom returns a logger with predefined values from a context.Context.
	// The logger, when used with controllers, can be expected to contain basic information about the object
	// that's being reconciled like:
	// - `reconciler group` and `reconciler kind` coming from the For(...) object passed in when building a over_controller.
	// - `name` and `namespace` from the reconciliation request.
	//
	// This is meant to be used with the context supplied in a struct that satisfies the Reconciler interface.
	LoggerFrom = over_log.FromContext

	// LoggerInto takes a context and sets the logger as one of its keys.
	//
	// This is meant to be used in reconcilers to enrich the logger within a context with additional values.
	LoggerInto = over_log.IntoContext

	// SetLogger sets a concrete logging implementation for all deferred Loggers.
	SetLogger = over_log.SetLogger
)
