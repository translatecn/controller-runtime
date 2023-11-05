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

package main

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/utils/pointer"
	"net/http"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	apiutil "sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	signals "sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

var scheme = runtime.NewScheme()

func init() {
	log.SetLogger(zap.New())
	clientgoscheme.AddToScheme(scheme)
}

type reconcileReplicaSet struct {
	client client.Client
}

func (r reconcileReplicaSet) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := log.FromContext(ctx).WithValues("chaospod", request.NamespacedName)
	log.V(1).Info("reconciling chaos pod")
	return reconcile.Result{}, nil
}

func main() {
	entryLog := log.Log.WithName("entrypoint")
	_config := config.GetConfigOrDie()
	_config.WrapTransport = func(rt http.RoundTripper) http.RoundTripper {
		return &LoggingTransport{rt: rt}
	}
	// Setup a Manager
	entryLog.Info("setting up manager")
	entryLog.V(3).Info("setting up manager3")
	u := &unstructured.Unstructured{}
	object, _ := apiutil.GVKForObject(&corev1.Pod{}, scheme)
	u.SetGroupVersionKind(object)
	mgr, err := ctrl.NewManager(_config, ctrl.Options{
		Scheme: scheme,
		Cache: cache.Options{
			//DefaultNamespaces: map[string]cache.Config{
			//	"kube-system": cache.Config{}, // DefaultLabelSelector
			//	"":            cache.Config{},
			//},
			ByObject: map[client.Object]cache.ByObject{
				u: cache.ByObject{
					Namespaces: map[string]cache.Config{
						"kube-system": cache.Config{}, // 依次填充 Label、DefaultLabelSelector
					},
					Label:                 nil, // 填充 DefaultLabelSelector
					Field:                 nil, // 填充 DefaultFieldSelector
					Transform:             nil, // 填充 DefaultTransform
					UnsafeDisableDeepCopy: nil, // 填充 DefaultUnsafeDisableDeepCopy
				},
				&appsv1.ReplicaSet{}: cache.ByObject{},
			},
			DefaultLabelSelector: labels.Everything(),
			DefaultFieldSelector: fields.Everything(), //metadata.namespace 可以在这里指定namespace
			DefaultTransform: func(i interface{}) (interface{}, error) {
				fmt.Println(i) // 从api server 获取到的每一个对象，没有TypeMeta信息
				return i, nil
			},
			DefaultUnsafeDisableDeepCopy: pointer.Bool(true),
		}, // 对ownsInput、forInput 里的结构根据ns 进行watch
		LeaderElection:         false,
		PprofBindAddress:       ":8012",
		HealthProbeBindAddress: ":8013",
		Metrics: metricsserver.Options{
			SecureServing:  false,
			BindAddress:    ":8014",
			ExtraHandlers:  nil,
			FilterProvider: nil,
			CertDir:        "",
			CertName:       "",
			KeyName:        "",
			TLSOpts:        nil,
		},
	}) // 初始化  metricServer、health、readyServer、leaderServer、pprofServer

	// {
	//   gvk:informerCache{Informers:cache.SharedIndexInformer,Namespace},
	//   gvk:multiNamespaceCache{ns1:informerCache{Informers:cache.SharedIndexInformer,Namespace}}
	//  }
	// 每个informerCache可以有多个多个informer ,但在当前场景下,只有一个informer
	if err != nil {
		entryLog.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	// Setup a new controller to reconcile ReplicaSets
	_ = controller.New
	err = ctrl.NewControllerManagedBy(mgr).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: 0,
			CacheSyncTimeout:        0,
			RecoverPanic:            nil,
			NeedLeaderElection:      nil,
			Reconciler:              nil,
			RateLimiter:             nil,
			LogConstructor:          nil,
		}).
		WithEventFilter(predicate.Funcs{
			CreateFunc:  nil,
			DeleteFunc:  nil,
			UpdateFunc:  nil,
			GenericFunc: nil,
		}).
		WithEventFilter(predicate.Funcs{
			CreateFunc:  nil,
			DeleteFunc:  nil,
			UpdateFunc:  nil,
			GenericFunc: nil,
		}).
		For(&appsv1.ReplicaSet{}). // 添加到 forInput 只能有一个
		Watches(&appsv1.ReplicaSet{}, handler.Funcs{
			CreateFunc:  nil,
			UpdateFunc:  nil,
			DeleteFunc:  nil,
			GenericFunc: nil,
		}).                            // 往 watchesInput添加对象
		Owns(&corev1.Pod{}).           // 添加到 ownsInput 可以有很多个；触发 pod owners 中对应的 ReplicaSet ；也会watch  pod
		Complete(&reconcileReplicaSet{ // reconcileReplicaSet 就是 Reconciler
			client: mgr.GetClient(),
		}) // 添加Runnable

	if err != nil {
		entryLog.Error(err, "unable to create controller")
		os.Exit(1)
	}
	mgr.AddHealthzCheck("x", func(req *http.Request) error {
		return nil
	})
	mgr.AddReadyzCheck("f", func(req *http.Request) error {
		return nil
	})
	mgr.Add(RunAble{})
	//hookServer := mgr.GetWebhookServer() // manager.New 初始化了 WebhookServer
	//hookServer.Register("/validate-v1-tokenreview", &authentication.Webhook{Handler: nil})
	//mgr.Add(hookServer)
	entryLog.Info("starting manager")

	go func() {
		mgr.GetCache().WaitForCacheSync(context.Background())
		fmt.Println("over")
		_ = mgr.GetCache().IndexField
		_ = mgr.GetCache().GetInformer
		_ = mgr.GetCache().GetInformerForKind
	}()
	_ = cluster.InternalCluster{}

	if err := mgr.Start(signals.SetupSignalHandler()); err != nil { // 会启动 InternalCluster(newCache.run--> watch)、metricServer、health、readyServer、leaderServer、other、pprof
		entryLog.Error(err, "unable to run manager")
		os.Exit(1)
	}
}

type LoggingTransport struct {
	rt http.RoundTripper
}

func (l *LoggingTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	fmt.Println(request.URL, request.Method)
	return l.rt.RoundTrip(request)
}

type RunAble struct {
}

func (r RunAble) Start(ctx context.Context) error {
	return nil
}
