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

package internal_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/over_event"
	"sigs.k8s.io/controller-runtime/pkg/over_handler"
	internal "sigs.k8s.io/controller-runtime/pkg/over_internal/source"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/over_controller/controllertest"
	"sigs.k8s.io/controller-runtime/pkg/over_predicate"
)

var _ = Describe("Internal", func() {
	var ctx = context.Background()
	var instance *internal.EventHandler
	var funcs, setfuncs *over_handler.Funcs
	var set bool
	BeforeEach(func() {
		funcs = &over_handler.Funcs{
			CreateFunc: func(context.Context, over_event.CreateEvent, workqueue.RateLimitingInterface) {
				defer GinkgoRecover()
				Fail("Did not expect CreateEvent to be called.")
			},
			DeleteFunc: func(context.Context, over_event.DeleteEvent, workqueue.RateLimitingInterface) {
				defer GinkgoRecover()
				Fail("Did not expect DeleteEvent to be called.")
			},
			UpdateFunc: func(context.Context, over_event.UpdateEvent, workqueue.RateLimitingInterface) {
				defer GinkgoRecover()
				Fail("Did not expect UpdateEvent to be called.")
			},
			GenericFunc: func(context.Context, over_event.GenericEvent, workqueue.RateLimitingInterface) {
				defer GinkgoRecover()
				Fail("Did not expect GenericEvent to be called.")
			},
		}

		setfuncs = &over_handler.Funcs{
			CreateFunc: func(context.Context, over_event.CreateEvent, workqueue.RateLimitingInterface) {
				set = true
			},
			DeleteFunc: func(context.Context, over_event.DeleteEvent, workqueue.RateLimitingInterface) {
				set = true
			},
			UpdateFunc: func(context.Context, over_event.UpdateEvent, workqueue.RateLimitingInterface) {
				set = true
			},
			GenericFunc: func(context.Context, over_event.GenericEvent, workqueue.RateLimitingInterface) {
				set = true
			},
		}
		instance = internal.NewEventHandler(ctx, &controllertest.Queue{}, funcs, nil)
	})

	Describe("EventHandler", func() {
		var pod, newPod *corev1.Pod

		BeforeEach(func() {
			pod = &corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{Name: "test", Image: "test"}},
				},
			}
			newPod = pod.DeepCopy()
			newPod.Labels = map[string]string{"foo": "bar"}
		})

		It("should create a CreateEvent", func() {
			funcs.CreateFunc = func(ctx context.Context, evt over_event.CreateEvent, q workqueue.RateLimitingInterface) {
				defer GinkgoRecover()
				Expect(evt.Object).To(Equal(pod))
			}
			instance.OnAdd(pod)
		})

		It("should used Predicates to filter CreateEvents", func() {
			instance = internal.NewEventHandler(ctx, &controllertest.Queue{}, setfuncs, []over_predicate.Predicate{
				over_predicate.Funcs{CreateFunc: func(over_event.CreateEvent) bool { return false }},
			})
			set = false
			instance.OnAdd(pod)
			Expect(set).To(BeFalse())

			set = false
			instance = internal.NewEventHandler(ctx, &controllertest.Queue{}, setfuncs, []over_predicate.Predicate{
				over_predicate.Funcs{CreateFunc: func(over_event.CreateEvent) bool { return true }},
			})
			instance.OnAdd(pod)
			Expect(set).To(BeTrue())

			set = false
			instance = internal.NewEventHandler(ctx, &controllertest.Queue{}, setfuncs, []over_predicate.Predicate{
				over_predicate.Funcs{CreateFunc: func(over_event.CreateEvent) bool { return true }},
				over_predicate.Funcs{CreateFunc: func(over_event.CreateEvent) bool { return false }},
			})
			instance.OnAdd(pod)
			Expect(set).To(BeFalse())

			set = false
			instance = internal.NewEventHandler(ctx, &controllertest.Queue{}, setfuncs, []over_predicate.Predicate{
				over_predicate.Funcs{CreateFunc: func(over_event.CreateEvent) bool { return false }},
				over_predicate.Funcs{CreateFunc: func(over_event.CreateEvent) bool { return true }},
			})
			instance.OnAdd(pod)
			Expect(set).To(BeFalse())

			set = false
			instance = internal.NewEventHandler(ctx, &controllertest.Queue{}, setfuncs, []over_predicate.Predicate{
				over_predicate.Funcs{CreateFunc: func(over_event.CreateEvent) bool { return true }},
				over_predicate.Funcs{CreateFunc: func(over_event.CreateEvent) bool { return true }},
			})
			instance.OnAdd(pod)
			Expect(set).To(BeTrue())
		})

		It("should not call Create EventHandler if the object is not a runtime.Object", func() {
			instance.OnAdd(&metav1.ObjectMeta{})
		})

		It("should not call Create EventHandler if the object does not have metadata", func() {
			instance.OnAdd(FooRuntimeObject{})
		})

		It("should create an UpdateEvent", func() {
			funcs.UpdateFunc = func(ctx context.Context, evt over_event.UpdateEvent, q workqueue.RateLimitingInterface) {
				defer GinkgoRecover()
				Expect(evt.ObjectOld).To(Equal(pod))
				Expect(evt.ObjectNew).To(Equal(newPod))
			}
			instance.OnUpdate(pod, newPod)
		})

		It("should used Predicates to filter UpdateEvents", func() {
			set = false
			instance = internal.NewEventHandler(ctx, &controllertest.Queue{}, setfuncs, []over_predicate.Predicate{
				over_predicate.Funcs{UpdateFunc: func(updateEvent over_event.UpdateEvent) bool { return false }},
			})
			instance.OnUpdate(pod, newPod)
			Expect(set).To(BeFalse())

			set = false
			instance = internal.NewEventHandler(ctx, &controllertest.Queue{}, setfuncs, []over_predicate.Predicate{
				over_predicate.Funcs{UpdateFunc: func(over_event.UpdateEvent) bool { return true }},
			})
			instance.OnUpdate(pod, newPod)
			Expect(set).To(BeTrue())

			set = false
			instance = internal.NewEventHandler(ctx, &controllertest.Queue{}, setfuncs, []over_predicate.Predicate{
				over_predicate.Funcs{UpdateFunc: func(over_event.UpdateEvent) bool { return true }},
				over_predicate.Funcs{UpdateFunc: func(over_event.UpdateEvent) bool { return false }},
			})
			instance.OnUpdate(pod, newPod)
			Expect(set).To(BeFalse())

			set = false
			instance = internal.NewEventHandler(ctx, &controllertest.Queue{}, setfuncs, []over_predicate.Predicate{
				over_predicate.Funcs{UpdateFunc: func(over_event.UpdateEvent) bool { return false }},
				over_predicate.Funcs{UpdateFunc: func(over_event.UpdateEvent) bool { return true }},
			})
			instance.OnUpdate(pod, newPod)
			Expect(set).To(BeFalse())

			set = false
			instance = internal.NewEventHandler(ctx, &controllertest.Queue{}, setfuncs, []over_predicate.Predicate{
				over_predicate.Funcs{CreateFunc: func(over_event.CreateEvent) bool { return true }},
				over_predicate.Funcs{CreateFunc: func(over_event.CreateEvent) bool { return true }},
			})
			instance.OnUpdate(pod, newPod)
			Expect(set).To(BeTrue())
		})

		It("should not call Update EventHandler if the object is not a runtime.Object", func() {
			instance.OnUpdate(&metav1.ObjectMeta{}, &corev1.Pod{})
			instance.OnUpdate(&corev1.Pod{}, &metav1.ObjectMeta{})
		})

		It("should not call Update EventHandler if the object does not have metadata", func() {
			instance.OnUpdate(FooRuntimeObject{}, &corev1.Pod{})
			instance.OnUpdate(&corev1.Pod{}, FooRuntimeObject{})
		})

		It("should create a DeleteEvent", func() {
			funcs.DeleteFunc = func(ctx context.Context, evt over_event.DeleteEvent, q workqueue.RateLimitingInterface) {
				defer GinkgoRecover()
				Expect(evt.Object).To(Equal(pod))
			}
			instance.OnDelete(pod)
		})

		It("should used Predicates to filter DeleteEvents", func() {
			set = false
			instance = internal.NewEventHandler(ctx, &controllertest.Queue{}, setfuncs, []over_predicate.Predicate{
				over_predicate.Funcs{DeleteFunc: func(over_event.DeleteEvent) bool { return false }},
			})
			instance.OnDelete(pod)
			Expect(set).To(BeFalse())

			set = false
			instance = internal.NewEventHandler(ctx, &controllertest.Queue{}, setfuncs, []over_predicate.Predicate{
				over_predicate.Funcs{DeleteFunc: func(over_event.DeleteEvent) bool { return true }},
			})
			instance.OnDelete(pod)
			Expect(set).To(BeTrue())

			set = false
			instance = internal.NewEventHandler(ctx, &controllertest.Queue{}, setfuncs, []over_predicate.Predicate{
				over_predicate.Funcs{DeleteFunc: func(over_event.DeleteEvent) bool { return true }},
				over_predicate.Funcs{DeleteFunc: func(over_event.DeleteEvent) bool { return false }},
			})
			instance.OnDelete(pod)
			Expect(set).To(BeFalse())

			set = false
			instance = internal.NewEventHandler(ctx, &controllertest.Queue{}, setfuncs, []over_predicate.Predicate{
				over_predicate.Funcs{DeleteFunc: func(over_event.DeleteEvent) bool { return false }},
				over_predicate.Funcs{DeleteFunc: func(over_event.DeleteEvent) bool { return true }},
			})
			instance.OnDelete(pod)
			Expect(set).To(BeFalse())

			set = false
			instance = internal.NewEventHandler(ctx, &controllertest.Queue{}, setfuncs, []over_predicate.Predicate{
				over_predicate.Funcs{DeleteFunc: func(over_event.DeleteEvent) bool { return true }},
				over_predicate.Funcs{DeleteFunc: func(over_event.DeleteEvent) bool { return true }},
			})
			instance.OnDelete(pod)
			Expect(set).To(BeTrue())
		})

		It("should not call Delete EventHandler if the object is not a runtime.Object", func() {
			instance.OnDelete(&metav1.ObjectMeta{})
		})

		It("should not call Delete EventHandler if the object does not have metadata", func() {
			instance.OnDelete(FooRuntimeObject{})
		})

		It("should create a DeleteEvent from a tombstone", func() {

			tombstone := cache.DeletedFinalStateUnknown{
				Obj: pod,
			}
			funcs.DeleteFunc = func(ctx context.Context, evt over_event.DeleteEvent, q workqueue.RateLimitingInterface) {
				defer GinkgoRecover()
				Expect(evt.Object).To(Equal(pod))
				Expect(evt.DeleteStateUnknown).Should(BeTrue())
			}

			instance.OnDelete(tombstone)
		})

		It("should ignore tombstone objects without meta", func() {
			tombstone := cache.DeletedFinalStateUnknown{Obj: Foo{}}
			instance.OnDelete(tombstone)
		})
		It("should ignore objects without meta", func() {
			instance.OnAdd(Foo{})
			instance.OnUpdate(Foo{}, Foo{})
			instance.OnDelete(Foo{})
		})
	})
})

type Foo struct{}

var _ runtime.Object = FooRuntimeObject{}

type FooRuntimeObject struct{}

func (FooRuntimeObject) GetObjectKind() schema.ObjectKind { return nil }
func (FooRuntimeObject) DeepCopyObject() runtime.Object   { return nil }
