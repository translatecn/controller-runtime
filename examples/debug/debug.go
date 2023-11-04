package main

import (
	"fmt"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"

	v1 "k8s.io/api/core/v1"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	apiutil "sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

var scheme = runtime.NewScheme()

func main() {
	clientgoscheme.AddToScheme(scheme)
	//object, err := apiutil.GVKForObject(&v1.Pod{}, clientgoscheme.Scheme)
	object, err := apiutil.GVKForObject(&v1.Pod{}, scheme)
	fmt.Println(object, err)

	selector, _ := fields.ParseSelector("a=b,metadata.namespace=A")
	value, found := selector.RequiresExactMatch("metadata.namespace")
	fmt.Println(value, found)
}
