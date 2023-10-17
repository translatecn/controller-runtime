package main

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	over_apiutil "sigs.k8s.io/controller-runtime/pkg/over_client/apiutil"
)

func main() {
	object, err := over_apiutil.GVKForObject(&v1.Pod{}, clientgoscheme.Scheme)
	fmt.Println(object, err)
}
