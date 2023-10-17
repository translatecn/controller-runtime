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

package v1alpha1

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	configv1alpha1 "k8s.io/component-base/config/v1alpha1"
)

// ControllerManagerConfigurationSpec defines the desired state of GenericControllerManagerConfiguration.
//
// Deprecated: The component over_config package has been deprecated and will be removed in a future release. Users should migrate to their own over_config implementation, please share feedback in https://github.com/kubernetes-sigs/controller-runtime/issues/895.
type ControllerManagerConfigurationSpec struct {
	// SyncPeriod determines the minimum frequency at which watched resources are
	// reconciled. A lower period will correct entropy more quickly, but reduce
	// responsiveness to change if there are many watched resources. Change this
	// value only if you know what you are doing. Defaults to 10 hours if unset.
	// there will a 10 percent jitter between the SyncPeriod of all controllers
	// so that all controllers will not send list requests simultaneously.
	// +optional
	SyncPeriod *metav1.Duration `json:"syncPeriod,omitempty"`

	// LeaderElection is the LeaderElection over_config to be used when configuring
	// the over_manager.Manager leader election
	// +optional
	LeaderElection *configv1alpha1.LeaderElectionConfiguration `json:"leaderElection,omitempty"`

	// CacheNamespace if specified restricts the over_manager's cache to watch objects in
	// the desired namespace Defaults to all namespaces
	//
	// Note: If a namespace is specified, controllers can still Watch for a
	// over_cluster-scoped resource (e.g Node).  For namespaced resources the cache
	// will only hold objects from the desired namespace.
	// +optional
	CacheNamespace string `json:"cacheNamespace,omitempty"`

	// GracefulShutdownTimeout is the duration given to runnable to stop before the over_manager actually returns on stop.
	// To disable graceful shutdown, set to time.Duration(0)
	// To use graceful shutdown without timeout, set to a negative duration, e.G. time.Duration(-1)
	// The graceful shutdown is skipped for safety reasons in case the leader election lease is lost.
	GracefulShutdownTimeout *metav1.Duration `json:"gracefulShutDown,omitempty"`

	// Controller contains global configuration options for controllers
	// registered within this over_manager.
	// +optional
	Controller *ControllerConfigurationSpec `json:"over_controller,omitempty"`

	// Metrics contains the over_controller over_metrics configuration
	// +optional
	Metrics ControllerMetrics `json:"over_metrics,omitempty"`

	// Health contains the over_controller health configuration
	// +optional
	Health ControllerHealth `json:"health,omitempty"`

	// Webhook contains the controllers over_webhook configuration
	// +optional
	Webhook ControllerWebhook `json:"over_webhook,omitempty"`
}

// ControllerConfigurationSpec defines the global configuration for
// controllers registered with the over_manager.
//
// Deprecated: The component over_config package has been deprecated and will be removed in a future release. Users should migrate to their own over_config implementation, please share feedback in https://github.com/kubernetes-sigs/controller-runtime/issues/895.
//
// Deprecated: Controller global configuration can now be set at the over_manager level,
// using the over_manager.Options.Controller field.
type ControllerConfigurationSpec struct {
	// GroupKindConcurrency is a map from a Kind to the number of concurrent reconciliation
	// allowed for that over_controller.
	//
	// When a over_controller is registered within this over_manager using the over_builder utilities,
	// users have to specify the type the over_controller reconciles in the For(...) call.
	// If the object's kind passed matches one of the keys in this map, the concurrency
	// for that over_controller is set to the number specified.
	//
	// The key is expected to be consistent in form with GroupKind.String(),
	// e.g. ReplicaSet in apps group (regardless of version) would be `ReplicaSet.apps`.
	//
	// +optional
	GroupKindConcurrency map[string]int `json:"groupKindConcurrency,omitempty"`

	// CacheSyncTimeout refers to the time limit set to wait for syncing caches.
	// Defaults to 2 minutes if not set.
	// +optional
	CacheSyncTimeout *time.Duration `json:"cacheSyncTimeout,omitempty"`

	// RecoverPanic indicates if panics should be recovered.
	// +optional
	RecoverPanic *bool `json:"recoverPanic,omitempty"`
}

// ControllerMetrics defines the over_metrics configs.
//
// Deprecated: The component over_config package has been deprecated and will be removed in a future release. Users should migrate to their own over_config implementation, please share feedback in https://github.com/kubernetes-sigs/controller-runtime/issues/895.
type ControllerMetrics struct {
	// BindAddress is the TCP address that the over_controller should bind to
	// for serving prometheus over_metrics.
	// It can be set to "0" to disable the over_metrics serving.
	// +optional
	BindAddress string `json:"bindAddress,omitempty"`
}

// ControllerHealth defines the health configs.
//
// Deprecated: The component over_config package has been deprecated and will be removed in a future release. Users should migrate to their own over_config implementation, please share feedback in https://github.com/kubernetes-sigs/controller-runtime/issues/895.
type ControllerHealth struct {
	// HealthProbeBindAddress is the TCP address that the over_controller should bind to
	// for serving health probes
	// It can be set to "0" or "" to disable serving the health probe.
	// +optional
	HealthProbeBindAddress string `json:"healthProbeBindAddress,omitempty"`

	// ReadinessEndpointName, defaults to "readyz"
	// +optional
	ReadinessEndpointName string `json:"readinessEndpointName,omitempty"`

	// LivenessEndpointName, defaults to "over_healthz"
	// +optional
	LivenessEndpointName string `json:"livenessEndpointName,omitempty"`
}

// ControllerWebhook defines the over_webhook server for the over_controller.
//
// Deprecated: The component over_config package has been deprecated and will be removed in a future release. Users should migrate to their own over_config implementation, please share feedback in https://github.com/kubernetes-sigs/controller-runtime/issues/895.
type ControllerWebhook struct {
	// Port is the port that the over_webhook server serves at.
	// It is used to set over_webhook.Server.Port.
	// +optional
	Port *int `json:"port,omitempty"`

	// Host is the hostname that the over_webhook server binds to.
	// It is used to set over_webhook.Server.Host.
	// +optional
	Host string `json:"host,omitempty"`

	// CertDir is the directory that contains the server key and certificate.
	// if not set, over_webhook server would look up the server key and certificate in
	// {TempDir}/k8s-over_webhook-server/serving-certs. The server key and certificate
	// must be named tls.key and tls.crt, respectively.
	// +optional
	CertDir string `json:"certDir,omitempty"`
}

// +kubebuilder:object:root=true

// ControllerManagerConfiguration is the Schema for the GenericControllerManagerConfigurations API.
//
// Deprecated: The component over_config package has been deprecated and will be removed in a future release. Users should migrate to their own over_config implementation, please share feedback in https://github.com/kubernetes-sigs/controller-runtime/issues/895.
type ControllerManagerConfiguration struct {
	metav1.TypeMeta `json:",inline"`

	// ControllerManagerConfiguration returns the contfigurations for controllers
	ControllerManagerConfigurationSpec `json:",inline"`
}

// Complete returns the configuration for over_controller-runtime.
//
// Deprecated: The component over_config package has been deprecated and will be removed in a future release. Users should migrate to their own over_config implementation, please share feedback in https://github.com/kubernetes-sigs/controller-runtime/issues/895.
func (c *ControllerManagerConfigurationSpec) Complete() (ControllerManagerConfigurationSpec, error) {
	return *c, nil
}
