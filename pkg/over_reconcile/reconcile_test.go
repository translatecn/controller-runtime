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

package over_reconcile_test

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/over_reconcile"
)

var _ = Describe("over_reconcile", func() {
	Describe("Result", func() {
		It("IsZero should return true if empty", func() {
			var res *over_reconcile.Result
			Expect(res.IsZero()).To(BeTrue())
			res2 := &over_reconcile.Result{}
			Expect(res2.IsZero()).To(BeTrue())
			res3 := over_reconcile.Result{}
			Expect(res3.IsZero()).To(BeTrue())
		})

		It("IsZero should return false if Requeue is set to true", func() {
			res := over_reconcile.Result{Requeue: true}
			Expect(res.IsZero()).To(BeFalse())
		})

		It("IsZero should return false if RequeueAfter is set to true", func() {
			res := over_reconcile.Result{RequeueAfter: 1 * time.Second}
			Expect(res.IsZero()).To(BeFalse())
		})
	})

	Describe("Func", func() {
		It("should call the function with the request and return a nil error.", func() {
			request := over_reconcile.Request{
				NamespacedName: types.NamespacedName{Name: "foo", Namespace: "bar"},
			}
			result := over_reconcile.Result{
				Requeue: true,
			}

			instance := over_reconcile.Func(func(_ context.Context, r over_reconcile.Request) (over_reconcile.Result, error) {
				defer GinkgoRecover()
				Expect(r).To(Equal(request))

				return result, nil
			})
			actualResult, actualErr := instance.Reconcile(context.Background(), request)
			Expect(actualResult).To(Equal(result))
			Expect(actualErr).NotTo(HaveOccurred())
		})

		It("should call the function with the request and return an error.", func() {
			request := over_reconcile.Request{
				NamespacedName: types.NamespacedName{Name: "foo", Namespace: "bar"},
			}
			result := over_reconcile.Result{
				Requeue: false,
			}
			err := fmt.Errorf("hello world")

			instance := over_reconcile.Func(func(_ context.Context, r over_reconcile.Request) (over_reconcile.Result, error) {
				defer GinkgoRecover()
				Expect(r).To(Equal(request))

				return result, err
			})
			actualResult, actualErr := instance.Reconcile(context.Background(), request)
			Expect(actualResult).To(Equal(result))
			Expect(actualErr).To(Equal(err))
		})

		It("should allow unwrapping inner error from terminal error", func() {
			inner := apierrors.NewGone("")
			terminalError := over_reconcile.TerminalError(inner)

			Expect(apierrors.IsGone(terminalError)).To(BeTrue())
		})

		It("should handle nil terminal errors properly", func() {
			err := over_reconcile.TerminalError(nil)
			Expect(err.Error()).To(Equal("nil terminal error"))
		})
	})
})
