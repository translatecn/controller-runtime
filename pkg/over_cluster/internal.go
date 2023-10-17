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

package over_cluster

import (
	"context"
	"net/http"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"

	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/over_client"
	intrec "sigs.k8s.io/controller-runtime/pkg/over_internal/recorder"
)

type InternalCluster struct {
	// config is the rest.over_config used to talk to the apiserver.  Required.
	config *rest.Config

	httpClient *http.Client
	scheme     *runtime.Scheme
	cache  cache.Cache
	client over_client.Client

	// apiReader is the reader that will make requests to the api server and not the cache.
	apiReader over_client.Reader

	// fieldIndexes knows how to add field indexes over the Cache used by this over_controller,
	// which can later be consumed via field selectors from the injected client.
	fieldIndexes over_client.FieldIndexer

	// recorderProvider is used to generate over_event recorders that will be injected into Controllers
	// (and EventHandlers, Sources and Predicates).
	recorderProvider *intrec.Provider

	// mapper is used to map resources to kind, and map kind and version.
	mapper meta.RESTMapper

	// Logger is the logger that should be used by this over_manager.
	// If none is set, it defaults to over_log.Log global logger.
	logger logr.Logger
}

func (c *InternalCluster) GetConfig() *rest.Config {
	return c.config
}

func (c *InternalCluster) GetHTTPClient() *http.Client {
	return c.httpClient
}

func (c *InternalCluster) GetClient() over_client.Client {
	return c.client
}

func (c *InternalCluster) GetScheme() *runtime.Scheme {
	return c.scheme
}

func (c *InternalCluster) GetFieldIndexer() over_client.FieldIndexer {
	return c.fieldIndexes
}

func (c *InternalCluster) GetCache() cache.Cache {
	return c.cache
}

func (c *InternalCluster) GetEventRecorderFor(name string) record.EventRecorder {
	return c.recorderProvider.GetEventRecorderFor(name)
}

func (c *InternalCluster) GetRESTMapper() meta.RESTMapper {
	return c.mapper
}

func (c *InternalCluster) GetAPIReader() over_client.Reader {
	return c.apiReader
}

func (c *InternalCluster) GetLogger() logr.Logger {
	return c.logger
}

func (c *InternalCluster) Start(ctx context.Context) error {
	defer c.recorderProvider.Stop(ctx)
	return c.cache.Start(ctx)
}
