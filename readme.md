
### manager 主要维护了
runnables{
    HTTPServers:    newRunnableGroup(baseContext, errChan), // metricsServer  /metrics
    Webhooks:       newRunnableGroup(baseContext, errChan),
    Caches:         newRunnableGroup(baseContext, errChan),
    Others:         newRunnableGroup(baseContext, errChan),
    LeaderElection: newRunnableGroup(baseContext, errChan), // 异步
}


func (r *runnables) Add(fn Runnable) error {
	switch runnable := fn.(type) {
	case *server:
		//    /readyz <---> AddReadyzCheck       map[string]healthz.Checker
		//    /healthz <---> AddHealthzCheck     map[string]healthz.Checker
		//    /debug/pprof/ <---> pprof.Index
		//    /debug/pprof/cmdline <---> pprof.Cmdline
		//    /debug/pprof/profile <---> pprof.Profile
		//    /debug/pprof/symbol <---> pprof.Symbol
		//    /debug/pprof/trace <---> pprof.Trace
		return r.HTTPServers.Add(fn, nil)
	case hasCache:
		_ = cluster.InternalCluster{}
		return r.Caches.Add(fn, func(ctx context.Context) bool {
			return runnable.GetCache().WaitForCacheSync(ctx)
		})
	case webhook.Server:
		_ = webhook.DefaultServer{}
		return r.Webhooks.Add(fn, nil)
	case LeaderElectionRunnable:
		if !runnable.NeedLeaderElection() {
			return r.Others.Add(fn, nil)
		}
		return r.LeaderElection.Add(fn, nil)
	default:
		return r.LeaderElection.Add(fn, nil)
	}
}




### controller{manager}

- manager.Add(Controller) -----> 会走default 分支
- 每个 Controller 作为一个服务加入到 LeaderElection
- Controller.Start 会将每个watch启动起来,并WaitForSync,都Sync以后运行processNextWorkItem
  - mgr.GetCache().GetInformer().AddEventHandler
  - mgr.GetCache().WaitForCacheSync
  - processNextWorkItem
