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

package over_predicate_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/over_client"
	"sigs.k8s.io/controller-runtime/pkg/over_event"
	"sigs.k8s.io/controller-runtime/pkg/over_predicate"
)

var _ = Describe("Predicate", func() {
	var pod *corev1.Pod
	BeforeEach(func() {
		pod = &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{Namespace: "biz", Name: "baz"},
		}
	})

	Describe("Funcs", func() {
		failingFuncs := over_predicate.Funcs{
			CreateFunc: func(over_event.CreateEvent) bool {
				defer GinkgoRecover()
				Fail("Did not expect CreateFunc to be called.")
				return false
			},
			DeleteFunc: func(over_event.DeleteEvent) bool {
				defer GinkgoRecover()
				Fail("Did not expect DeleteFunc to be called.")
				return false
			},
			UpdateFunc: func(over_event.UpdateEvent) bool {
				defer GinkgoRecover()
				Fail("Did not expect UpdateFunc to be called.")
				return false
			},
			GenericFunc: func(over_event.GenericEvent) bool {
				defer GinkgoRecover()
				Fail("Did not expect GenericFunc to be called.")
				return false
			},
		}

		It("should call Create", func() {
			instance := failingFuncs
			instance.CreateFunc = func(evt over_event.CreateEvent) bool {
				defer GinkgoRecover()
				Expect(evt.Object).To(Equal(pod))
				return false
			}
			evt := over_event.CreateEvent{
				Object: pod,
			}
			Expect(instance.Create(evt)).To(BeFalse())

			instance.CreateFunc = func(evt over_event.CreateEvent) bool {
				defer GinkgoRecover()
				Expect(evt.Object).To(Equal(pod))
				return true
			}
			Expect(instance.Create(evt)).To(BeTrue())

			instance.CreateFunc = nil
			Expect(instance.Create(evt)).To(BeTrue())
		})

		It("should call Update", func() {
			newPod := pod.DeepCopy()
			newPod.Name = "baz2"
			newPod.Namespace = "biz2"

			instance := failingFuncs
			instance.UpdateFunc = func(evt over_event.UpdateEvent) bool {
				defer GinkgoRecover()
				Expect(evt.ObjectOld).To(Equal(pod))
				Expect(evt.ObjectNew).To(Equal(newPod))
				return false
			}
			evt := over_event.UpdateEvent{
				ObjectOld: pod,
				ObjectNew: newPod,
			}
			Expect(instance.Update(evt)).To(BeFalse())

			instance.UpdateFunc = func(evt over_event.UpdateEvent) bool {
				defer GinkgoRecover()
				Expect(evt.ObjectOld).To(Equal(pod))
				Expect(evt.ObjectNew).To(Equal(newPod))
				return true
			}
			Expect(instance.Update(evt)).To(BeTrue())

			instance.UpdateFunc = nil
			Expect(instance.Update(evt)).To(BeTrue())
		})

		It("should call Delete", func() {
			instance := failingFuncs
			instance.DeleteFunc = func(evt over_event.DeleteEvent) bool {
				defer GinkgoRecover()
				Expect(evt.Object).To(Equal(pod))
				return false
			}
			evt := over_event.DeleteEvent{
				Object: pod,
			}
			Expect(instance.Delete(evt)).To(BeFalse())

			instance.DeleteFunc = func(evt over_event.DeleteEvent) bool {
				defer GinkgoRecover()
				Expect(evt.Object).To(Equal(pod))
				return true
			}
			Expect(instance.Delete(evt)).To(BeTrue())

			instance.DeleteFunc = nil
			Expect(instance.Delete(evt)).To(BeTrue())
		})

		It("should call Generic", func() {
			instance := failingFuncs
			instance.GenericFunc = func(evt over_event.GenericEvent) bool {
				defer GinkgoRecover()
				Expect(evt.Object).To(Equal(pod))
				return false
			}
			evt := over_event.GenericEvent{
				Object: pod,
			}
			Expect(instance.Generic(evt)).To(BeFalse())

			instance.GenericFunc = func(evt over_event.GenericEvent) bool {
				defer GinkgoRecover()
				Expect(evt.Object).To(Equal(pod))
				return true
			}
			Expect(instance.Generic(evt)).To(BeTrue())

			instance.GenericFunc = nil
			Expect(instance.Generic(evt)).To(BeTrue())
		})
	})

	Describe("When checking a ResourceVersionChangedPredicate", func() {
		instance := over_predicate.ResourceVersionChangedPredicate{}

		Context("Where the old object doesn't have a ResourceVersion or metadata", func() {
			It("should return false", func() {
				newPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:            "baz",
						Namespace:       "biz",
						ResourceVersion: "1",
					}}

				failEvnt := over_event.UpdateEvent{
					ObjectNew: newPod,
				}
				Expect(instance.Create(over_event.CreateEvent{})).Should(BeTrue())
				Expect(instance.Delete(over_event.DeleteEvent{})).Should(BeTrue())
				Expect(instance.Generic(over_event.GenericEvent{})).Should(BeTrue())
				Expect(instance.Update(failEvnt)).Should(BeFalse())
			})
		})

		Context("Where the new object doesn't have a ResourceVersion or metadata", func() {
			It("should return false", func() {
				oldPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:            "baz",
						Namespace:       "biz",
						ResourceVersion: "1",
					}}

				failEvnt := over_event.UpdateEvent{
					ObjectOld: oldPod,
				}
				Expect(instance.Create(over_event.CreateEvent{})).Should(BeTrue())
				Expect(instance.Delete(over_event.DeleteEvent{})).Should(BeTrue())
				Expect(instance.Generic(over_event.GenericEvent{})).Should(BeTrue())
				Expect(instance.Update(failEvnt)).Should(BeFalse())
				Expect(instance.Update(failEvnt)).Should(BeFalse())
			})
		})

		Context("Where the ResourceVersion hasn't changed", func() {
			It("should return false", func() {
				newPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:            "baz",
						Namespace:       "biz",
						ResourceVersion: "v1",
					}}

				oldPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:            "baz",
						Namespace:       "biz",
						ResourceVersion: "v1",
					}}

				failEvnt := over_event.UpdateEvent{
					ObjectOld: oldPod,
					ObjectNew: newPod,
				}
				Expect(instance.Create(over_event.CreateEvent{})).Should(BeTrue())
				Expect(instance.Delete(over_event.DeleteEvent{})).Should(BeTrue())
				Expect(instance.Generic(over_event.GenericEvent{})).Should(BeTrue())
				Expect(instance.Update(failEvnt)).Should(BeFalse())
				Expect(instance.Update(failEvnt)).Should(BeFalse())
			})
		})

		Context("Where the ResourceVersion has changed", func() {
			It("should return true", func() {
				newPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:            "baz",
						Namespace:       "biz",
						ResourceVersion: "v1",
					}}

				oldPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:            "baz",
						Namespace:       "biz",
						ResourceVersion: "v2",
					}}
				passEvt := over_event.UpdateEvent{
					ObjectOld: oldPod,
					ObjectNew: newPod,
				}
				Expect(instance.Create(over_event.CreateEvent{})).Should(BeTrue())
				Expect(instance.Delete(over_event.DeleteEvent{})).Should(BeTrue())
				Expect(instance.Generic(over_event.GenericEvent{})).Should(BeTrue())
				Expect(instance.Update(passEvt)).Should(BeTrue())
			})
		})

		Context("Where the objects or metadata are missing", func() {

			It("should return false", func() {
				newPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:            "baz",
						Namespace:       "biz",
						ResourceVersion: "v1",
					}}

				oldPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:            "baz",
						Namespace:       "biz",
						ResourceVersion: "v1",
					}}

				failEvt1 := over_event.UpdateEvent{ObjectOld: oldPod}
				failEvt2 := over_event.UpdateEvent{ObjectNew: newPod}
				failEvt3 := over_event.UpdateEvent{ObjectOld: oldPod, ObjectNew: newPod}
				Expect(instance.Create(over_event.CreateEvent{})).Should(BeTrue())
				Expect(instance.Delete(over_event.DeleteEvent{})).Should(BeTrue())
				Expect(instance.Generic(over_event.GenericEvent{})).Should(BeTrue())
				Expect(instance.Update(failEvt1)).Should(BeFalse())
				Expect(instance.Update(failEvt2)).Should(BeFalse())
				Expect(instance.Update(failEvt3)).Should(BeFalse())
			})
		})

	})

	Describe("When checking a GenerationChangedPredicate", func() {
		instance := over_predicate.GenerationChangedPredicate{}
		Context("Where the old object doesn't have a Generation or metadata", func() {
			It("should return false", func() {
				newPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:       "baz",
						Namespace:  "biz",
						Generation: 1,
					}}

				failEvnt := over_event.UpdateEvent{
					ObjectNew: newPod,
				}
				Expect(instance.Create(over_event.CreateEvent{})).To(BeTrue())
				Expect(instance.Delete(over_event.DeleteEvent{})).To(BeTrue())
				Expect(instance.Generic(over_event.GenericEvent{})).To(BeTrue())
				Expect(instance.Update(failEvnt)).To(BeFalse())
			})
		})

		Context("Where the new object doesn't have a Generation or metadata", func() {
			It("should return false", func() {
				oldPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:       "baz",
						Namespace:  "biz",
						Generation: 1,
					}}

				failEvnt := over_event.UpdateEvent{
					ObjectOld: oldPod,
				}
				Expect(instance.Create(over_event.CreateEvent{})).To(BeTrue())
				Expect(instance.Delete(over_event.DeleteEvent{})).To(BeTrue())
				Expect(instance.Generic(over_event.GenericEvent{})).To(BeTrue())
				Expect(instance.Update(failEvnt)).To(BeFalse())
			})
		})

		Context("Where the Generation hasn't changed", func() {
			It("should return false", func() {
				newPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:       "baz",
						Namespace:  "biz",
						Generation: 1,
					}}

				oldPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:       "baz",
						Namespace:  "biz",
						Generation: 1,
					}}

				failEvnt := over_event.UpdateEvent{
					ObjectOld: oldPod,
					ObjectNew: newPod,
				}
				Expect(instance.Create(over_event.CreateEvent{})).To(BeTrue())
				Expect(instance.Delete(over_event.DeleteEvent{})).To(BeTrue())
				Expect(instance.Generic(over_event.GenericEvent{})).To(BeTrue())
				Expect(instance.Update(failEvnt)).To(BeFalse())
			})
		})

		Context("Where the Generation has changed", func() {
			It("should return true", func() {
				newPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:       "baz",
						Namespace:  "biz",
						Generation: 1,
					}}

				oldPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:       "baz",
						Namespace:  "biz",
						Generation: 2,
					}}
				passEvt := over_event.UpdateEvent{
					ObjectOld: oldPod,
					ObjectNew: newPod,
				}
				Expect(instance.Create(over_event.CreateEvent{})).To(BeTrue())
				Expect(instance.Delete(over_event.DeleteEvent{})).To(BeTrue())
				Expect(instance.Generic(over_event.GenericEvent{})).To(BeTrue())
				Expect(instance.Update(passEvt)).To(BeTrue())
			})
		})

		Context("Where the objects or metadata are missing", func() {

			It("should return false", func() {
				newPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:       "baz",
						Namespace:  "biz",
						Generation: 1,
					}}

				oldPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:       "baz",
						Namespace:  "biz",
						Generation: 1,
					}}

				failEvt1 := over_event.UpdateEvent{ObjectOld: oldPod}
				failEvt2 := over_event.UpdateEvent{ObjectNew: newPod}
				failEvt3 := over_event.UpdateEvent{ObjectOld: oldPod, ObjectNew: newPod}
				Expect(instance.Create(over_event.CreateEvent{})).To(BeTrue())
				Expect(instance.Delete(over_event.DeleteEvent{})).To(BeTrue())
				Expect(instance.Generic(over_event.GenericEvent{})).To(BeTrue())
				Expect(instance.Update(failEvt1)).To(BeFalse())
				Expect(instance.Update(failEvt2)).To(BeFalse())
				Expect(instance.Update(failEvt3)).To(BeFalse())
			})
		})

	})

	// AnnotationChangedPredicate has almost identical test cases as LabelChangedPredicates,
	// so the duplication linter should be muted on both two test suites.
	Describe("When checking an AnnotationChangedPredicate", func() {
		instance := over_predicate.AnnotationChangedPredicate{}
		Context("Where the old object is missing", func() {
			It("should return false", func() {
				newPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz",
						Namespace: "biz",
						Annotations: map[string]string{
							"booz": "wooz",
						},
					}}

				failEvnt := over_event.UpdateEvent{
					ObjectNew: newPod,
				}
				Expect(instance.Create(over_event.CreateEvent{})).To(BeTrue())
				Expect(instance.Delete(over_event.DeleteEvent{})).To(BeTrue())
				Expect(instance.Generic(over_event.GenericEvent{})).To(BeTrue())
				Expect(instance.Update(failEvnt)).To(BeFalse())
			})
		})

		Context("Where the new object is missing", func() {
			It("should return false", func() {
				oldPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz",
						Namespace: "biz",
						Annotations: map[string]string{
							"booz": "wooz",
						},
					}}

				failEvnt := over_event.UpdateEvent{
					ObjectOld: oldPod,
				}
				Expect(instance.Create(over_event.CreateEvent{})).To(BeTrue())
				Expect(instance.Delete(over_event.DeleteEvent{})).To(BeTrue())
				Expect(instance.Generic(over_event.GenericEvent{})).To(BeTrue())
				Expect(instance.Update(failEvnt)).To(BeFalse())
			})
		})

		Context("Where the annotations are empty", func() {
			It("should return false", func() {
				newPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz",
						Namespace: "biz",
					}}

				oldPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz",
						Namespace: "biz",
					}}

				failEvnt := over_event.UpdateEvent{
					ObjectOld: oldPod,
					ObjectNew: newPod,
				}
				Expect(instance.Create(over_event.CreateEvent{})).To(BeTrue())
				Expect(instance.Delete(over_event.DeleteEvent{})).To(BeTrue())
				Expect(instance.Generic(over_event.GenericEvent{})).To(BeTrue())
				Expect(instance.Update(failEvnt)).To(BeFalse())
			})
		})

		Context("Where the annotations haven't changed", func() {
			It("should return false", func() {
				newPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz",
						Namespace: "biz",
						Annotations: map[string]string{
							"booz": "wooz",
						},
					}}

				oldPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz",
						Namespace: "biz",
						Annotations: map[string]string{
							"booz": "wooz",
						},
					}}

				failEvnt := over_event.UpdateEvent{
					ObjectOld: oldPod,
					ObjectNew: newPod,
				}
				Expect(instance.Create(over_event.CreateEvent{})).To(BeTrue())
				Expect(instance.Delete(over_event.DeleteEvent{})).To(BeTrue())
				Expect(instance.Generic(over_event.GenericEvent{})).To(BeTrue())
				Expect(instance.Update(failEvnt)).To(BeFalse())
			})
		})

		Context("Where an annotation value has changed", func() {
			It("should return true", func() {
				newPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz",
						Namespace: "biz",
						Annotations: map[string]string{
							"booz": "wooz",
						},
					}}

				oldPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz",
						Namespace: "biz",
						Annotations: map[string]string{
							"booz": "weez",
						},
					}}

				passEvt := over_event.UpdateEvent{
					ObjectOld: oldPod,
					ObjectNew: newPod,
				}
				Expect(instance.Create(over_event.CreateEvent{})).To(BeTrue())
				Expect(instance.Delete(over_event.DeleteEvent{})).To(BeTrue())
				Expect(instance.Generic(over_event.GenericEvent{})).To(BeTrue())
				Expect(instance.Update(passEvt)).To(BeTrue())
			})
		})

		Context("Where an annotation has been added", func() {
			It("should return true", func() {
				newPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz",
						Namespace: "biz",
						Annotations: map[string]string{
							"booz": "wooz",
						},
					}}

				oldPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz",
						Namespace: "biz",
						Annotations: map[string]string{
							"booz": "wooz",
							"zooz": "qooz",
						},
					}}

				passEvt := over_event.UpdateEvent{
					ObjectOld: oldPod,
					ObjectNew: newPod,
				}
				Expect(instance.Create(over_event.CreateEvent{})).To(BeTrue())
				Expect(instance.Delete(over_event.DeleteEvent{})).To(BeTrue())
				Expect(instance.Generic(over_event.GenericEvent{})).To(BeTrue())
				Expect(instance.Update(passEvt)).To(BeTrue())
			})
		})

		Context("Where an annotation has been removed", func() {
			It("should return true", func() {
				newPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz",
						Namespace: "biz",
						Annotations: map[string]string{
							"booz": "wooz",
							"zooz": "qooz",
						},
					}}

				oldPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz",
						Namespace: "biz",
						Annotations: map[string]string{
							"booz": "wooz",
						},
					}}

				passEvt := over_event.UpdateEvent{
					ObjectOld: oldPod,
					ObjectNew: newPod,
				}
				Expect(instance.Create(over_event.CreateEvent{})).To(BeTrue())
				Expect(instance.Delete(over_event.DeleteEvent{})).To(BeTrue())
				Expect(instance.Generic(over_event.GenericEvent{})).To(BeTrue())
				Expect(instance.Update(passEvt)).To(BeTrue())
			})
		})
	})

	// LabelChangedPredicates has almost identical test cases as AnnotationChangedPredicates,
	// so the duplication linter should be muted on both two test suites.
	Describe("When checking a LabelChangedPredicate", func() {
		instance := over_predicate.LabelChangedPredicate{}
		Context("Where the old object is missing", func() {
			It("should return false", func() {
				newPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz",
						Namespace: "biz",
						Labels: map[string]string{
							"foo": "bar",
						},
					}}

				evt := over_event.UpdateEvent{
					ObjectNew: newPod,
				}
				Expect(instance.Create(over_event.CreateEvent{})).To(BeTrue())
				Expect(instance.Delete(over_event.DeleteEvent{})).To(BeTrue())
				Expect(instance.Generic(over_event.GenericEvent{})).To(BeTrue())
				Expect(instance.Update(evt)).To(BeFalse())
			})
		})

		Context("Where the new object is missing", func() {
			It("should return false", func() {
				oldPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz",
						Namespace: "biz",
						Labels: map[string]string{
							"foo": "bar",
						},
					}}

				evt := over_event.UpdateEvent{
					ObjectOld: oldPod,
				}
				Expect(instance.Create(over_event.CreateEvent{})).To(BeTrue())
				Expect(instance.Delete(over_event.DeleteEvent{})).To(BeTrue())
				Expect(instance.Generic(over_event.GenericEvent{})).To(BeTrue())
				Expect(instance.Update(evt)).To(BeFalse())
			})
		})

		Context("Where the labels are empty", func() {
			It("should return false", func() {
				newPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz",
						Namespace: "biz",
					}}

				oldPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz",
						Namespace: "biz",
					}}

				evt := over_event.UpdateEvent{
					ObjectOld: oldPod,
					ObjectNew: newPod,
				}
				Expect(instance.Create(over_event.CreateEvent{})).To(BeTrue())
				Expect(instance.Delete(over_event.DeleteEvent{})).To(BeTrue())
				Expect(instance.Generic(over_event.GenericEvent{})).To(BeTrue())
				Expect(instance.Update(evt)).To(BeFalse())
			})
		})

		Context("Where the labels haven't changed", func() {
			It("should return false", func() {
				newPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz",
						Namespace: "biz",
						Labels: map[string]string{
							"foo": "bar",
						},
					}}

				oldPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz",
						Namespace: "biz",
						Labels: map[string]string{
							"foo": "bar",
						},
					}}

				evt := over_event.UpdateEvent{
					ObjectOld: oldPod,
					ObjectNew: newPod,
				}
				Expect(instance.Create(over_event.CreateEvent{})).To(BeTrue())
				Expect(instance.Delete(over_event.DeleteEvent{})).To(BeTrue())
				Expect(instance.Generic(over_event.GenericEvent{})).To(BeTrue())
				Expect(instance.Update(evt)).To(BeFalse())
			})
		})

		Context("Where a label value has changed", func() {
			It("should return true", func() {
				newPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz",
						Namespace: "biz",
						Labels: map[string]string{
							"foo": "bar",
						},
					}}

				oldPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz",
						Namespace: "biz",
						Labels: map[string]string{
							"foo": "bee",
						},
					}}

				evt := over_event.UpdateEvent{
					ObjectOld: oldPod,
					ObjectNew: newPod,
				}
				Expect(instance.Create(over_event.CreateEvent{})).To(BeTrue())
				Expect(instance.Delete(over_event.DeleteEvent{})).To(BeTrue())
				Expect(instance.Generic(over_event.GenericEvent{})).To(BeTrue())
				Expect(instance.Update(evt)).To(BeTrue())
			})
		})

		Context("Where a label has been added", func() {
			It("should return true", func() {
				newPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz",
						Namespace: "biz",
						Labels: map[string]string{
							"foo": "bar",
						},
					}}

				oldPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz",
						Namespace: "biz",
						Labels: map[string]string{
							"foo": "bar",
							"faa": "bor",
						},
					}}

				evt := over_event.UpdateEvent{
					ObjectOld: oldPod,
					ObjectNew: newPod,
				}
				Expect(instance.Create(over_event.CreateEvent{})).To(BeTrue())
				Expect(instance.Delete(over_event.DeleteEvent{})).To(BeTrue())
				Expect(instance.Generic(over_event.GenericEvent{})).To(BeTrue())
				Expect(instance.Update(evt)).To(BeTrue())
			})
		})

		Context("Where a label has been removed", func() {
			It("should return true", func() {
				newPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz",
						Namespace: "biz",
						Labels: map[string]string{
							"foo": "bar",
							"faa": "bor",
						},
					}}

				oldPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz",
						Namespace: "biz",
						Labels: map[string]string{
							"foo": "bar",
						},
					}}

				evt := over_event.UpdateEvent{
					ObjectOld: oldPod,
					ObjectNew: newPod,
				}
				Expect(instance.Create(over_event.CreateEvent{})).To(BeTrue())
				Expect(instance.Delete(over_event.DeleteEvent{})).To(BeTrue())
				Expect(instance.Generic(over_event.GenericEvent{})).To(BeTrue())
				Expect(instance.Update(evt)).To(BeTrue())
			})
		})
	})

	Context("With a boolean over_predicate", func() {
		funcs := func(pass bool) over_predicate.Funcs {
			return over_predicate.Funcs{
				CreateFunc: func(over_event.CreateEvent) bool {
					return pass
				},
				DeleteFunc: func(over_event.DeleteEvent) bool {
					return pass
				},
				UpdateFunc: func(over_event.UpdateEvent) bool {
					return pass
				},
				GenericFunc: func(over_event.GenericEvent) bool {
					return pass
				},
			}
		}
		passFuncs := funcs(true)
		failFuncs := funcs(false)

		Describe("When checking an And over_predicate", func() {
			It("should return false when one of its predicates returns false", func() {
				a := over_predicate.And(passFuncs, failFuncs)
				Expect(a.Create(over_event.CreateEvent{})).To(BeFalse())
				Expect(a.Update(over_event.UpdateEvent{})).To(BeFalse())
				Expect(a.Delete(over_event.DeleteEvent{})).To(BeFalse())
				Expect(a.Generic(over_event.GenericEvent{})).To(BeFalse())
			})
			It("should return true when all of its predicates return true", func() {
				a := over_predicate.And(passFuncs, passFuncs)
				Expect(a.Create(over_event.CreateEvent{})).To(BeTrue())
				Expect(a.Update(over_event.UpdateEvent{})).To(BeTrue())
				Expect(a.Delete(over_event.DeleteEvent{})).To(BeTrue())
				Expect(a.Generic(over_event.GenericEvent{})).To(BeTrue())
			})
		})
		Describe("When checking an Or over_predicate", func() {
			It("should return true when one of its predicates returns true", func() {
				o := over_predicate.Or(passFuncs, failFuncs)
				Expect(o.Create(over_event.CreateEvent{})).To(BeTrue())
				Expect(o.Update(over_event.UpdateEvent{})).To(BeTrue())
				Expect(o.Delete(over_event.DeleteEvent{})).To(BeTrue())
				Expect(o.Generic(over_event.GenericEvent{})).To(BeTrue())
			})
			It("should return false when all of its predicates return false", func() {
				o := over_predicate.Or(failFuncs, failFuncs)
				Expect(o.Create(over_event.CreateEvent{})).To(BeFalse())
				Expect(o.Update(over_event.UpdateEvent{})).To(BeFalse())
				Expect(o.Delete(over_event.DeleteEvent{})).To(BeFalse())
				Expect(o.Generic(over_event.GenericEvent{})).To(BeFalse())
			})
		})
		Describe("When checking a Not over_predicate", func() {
			It("should return false when its over_predicate returns true", func() {
				n := over_predicate.Not(passFuncs)
				Expect(n.Create(over_event.CreateEvent{})).To(BeFalse())
				Expect(n.Update(over_event.UpdateEvent{})).To(BeFalse())
				Expect(n.Delete(over_event.DeleteEvent{})).To(BeFalse())
				Expect(n.Generic(over_event.GenericEvent{})).To(BeFalse())
			})
			It("should return true when its over_predicate returns false", func() {
				n := over_predicate.Not(failFuncs)
				Expect(n.Create(over_event.CreateEvent{})).To(BeTrue())
				Expect(n.Update(over_event.UpdateEvent{})).To(BeTrue())
				Expect(n.Delete(over_event.DeleteEvent{})).To(BeTrue())
				Expect(n.Generic(over_event.GenericEvent{})).To(BeTrue())
			})
		})
	})

	Describe("NewPredicateFuncs with a namespace filter function", func() {
		byNamespaceFilter := func(namespace string) func(object over_client.Object) bool {
			return func(object over_client.Object) bool {
				return object.GetNamespace() == namespace
			}
		}
		byNamespaceFuncs := over_predicate.NewPredicateFuncs(byNamespaceFilter("biz"))
		Context("Where the namespace is matching", func() {
			It("should return true", func() {
				newPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz",
						Namespace: "biz",
					}}

				oldPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz",
						Namespace: "biz",
					}}
				passEvt1 := over_event.UpdateEvent{ObjectOld: oldPod, ObjectNew: newPod}
				Expect(byNamespaceFuncs.Create(over_event.CreateEvent{Object: newPod})).To(BeTrue())
				Expect(byNamespaceFuncs.Delete(over_event.DeleteEvent{Object: oldPod})).To(BeTrue())
				Expect(byNamespaceFuncs.Generic(over_event.GenericEvent{Object: newPod})).To(BeTrue())
				Expect(byNamespaceFuncs.Update(passEvt1)).To(BeTrue())
			})
		})

		Context("Where the namespace is not matching", func() {
			It("should return false", func() {
				newPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz",
						Namespace: "bizz",
					}}

				oldPod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz",
						Namespace: "biz",
					}}
				failEvt1 := over_event.UpdateEvent{ObjectOld: oldPod, ObjectNew: newPod}
				Expect(byNamespaceFuncs.Create(over_event.CreateEvent{Object: newPod})).To(BeFalse())
				Expect(byNamespaceFuncs.Delete(over_event.DeleteEvent{Object: newPod})).To(BeFalse())
				Expect(byNamespaceFuncs.Generic(over_event.GenericEvent{Object: newPod})).To(BeFalse())
				Expect(byNamespaceFuncs.Update(failEvt1)).To(BeFalse())
			})
		})
	})

	Describe("When checking a LabelSelectorPredicate", func() {
		instance, err := over_predicate.LabelSelectorPredicate(metav1.LabelSelector{MatchLabels: map[string]string{"foo": "bar"}})
		if err != nil {
			Fail("Improper Label Selector passed during over_predicate instantiation.")
		}

		Context("When the Selector does not match the over_event labels", func() {
			It("should return false", func() {
				failMatch := &corev1.Pod{}
				Expect(instance.Create(over_event.CreateEvent{Object: failMatch})).To(BeFalse())
				Expect(instance.Delete(over_event.DeleteEvent{Object: failMatch})).To(BeFalse())
				Expect(instance.Generic(over_event.GenericEvent{Object: failMatch})).To(BeFalse())
				Expect(instance.Update(over_event.UpdateEvent{ObjectNew: failMatch})).To(BeFalse())
			})
		})

		Context("When the Selector matches the over_event labels", func() {
			It("should return true", func() {
				successMatch := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"foo": "bar"},
					},
				}
				Expect(instance.Create(over_event.CreateEvent{Object: successMatch})).To(BeTrue())
				Expect(instance.Delete(over_event.DeleteEvent{Object: successMatch})).To(BeTrue())
				Expect(instance.Generic(over_event.GenericEvent{Object: successMatch})).To(BeTrue())
				Expect(instance.Update(over_event.UpdateEvent{ObjectNew: successMatch})).To(BeTrue())
			})
		})
	})
})
