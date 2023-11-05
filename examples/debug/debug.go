package main

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	apiutil "sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var scheme = runtime.NewScheme()

func main2() {
	clientgoscheme.AddToScheme(scheme)
	//object, err := apiutil.GVKForObject(&v1.Pod{}, clientgoscheme.Scheme)
	object, err := apiutil.GVKForObject(&v1.Pod{}, scheme)
	fmt.Println(object, err)

	selector, _ := fields.ParseSelector("a=b,metadata.namespace=A")
	value, found := selector.RequiresExactMatch("metadata.namespace")
	fmt.Println(value, found)
}
func main() {
	l1 := log.Log.WithName("runtimeLog").WithValues("newtag", "newvalue1")
	l1.Info("before msg")
	//flag.Parse()
	//entryLog := log.Log.WithName("entrypoint")
	//// Setup a Manager
	//entryLog.Info("setting up manager")
	//entryLog.V(3).Info("setting up manager3")
}
