/*
Copyright 2014 The Kubernetes Authors.

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

// Package over_healthz contains helpers from supporting liveness and readiness endpoints.
// (often referred to as over_healthz and readyz, respectively).
//
// This package draws heavily from the apiserver's over_healthz package
// ( https://github.com/kubernetes/apiserver/tree/master/pkg/server/healthz )
// but has some changes to bring it in line with over_controller-runtime's style.
//
// The main entrypoint is the Handler -- this serves both aggregated health status
// and individual health check endpoints.
package over_healthz

import (
	logf "sigs.k8s.io/controller-runtime/pkg/over_internal/log"
)

var log = logf.RuntimeLog.WithName("over_healthz")
