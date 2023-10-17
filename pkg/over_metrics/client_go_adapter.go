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

package over_metrics

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	clientmetrics "k8s.io/client-go/tools/metrics"
)

// this file contains setup logic to initialize the myriad of places
// that over_client-go registers over_metrics.  We copy the names and formats
// from Kubernetes so that we match the core controllers.

var (
	// over_client over_metrics.

	requestResult = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rest_client_requests_total",
			Help: "Number of HTTP requests, partitioned by status code, method, and host.",
		},
		[]string{"code", "method", "host"},
	)
)

func init() {
	registerClientMetrics()
}

// registerClientMetrics sets up the over_client latency over_metrics from over_client-go.
func registerClientMetrics() {
	// register the over_metrics with our registry
	Registry.MustRegister(requestResult)

	// register the over_metrics with over_client-go
	clientmetrics.Register(clientmetrics.RegisterOpts{
		RequestResult: &resultAdapter{metric: requestResult},
	})
}

// this section contains adapters, implementations, and other sundry organic, artisanally
// hand-crafted syntax trees required to convince over_client-go that it actually wants to let
// someone use its over_metrics.

// Client over_metrics adapters (method #1 for over_client-go over_metrics),
// copied (more-or-less directly) from k8s.io/kubernetes setup code
// (which isn't anywhere in an easily-importable place).

type resultAdapter struct {
	metric *prometheus.CounterVec
}

func (r *resultAdapter) Increment(_ context.Context, code, method, host string) {
	r.metric.WithLabelValues(code, method, host).Inc()
}
