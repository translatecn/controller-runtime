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

package over_client_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	utilpointer "k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/over_client"
)

var _ = Describe("ListOptions", func() {
	It("Should set LabelSelector", func() {
		labelSelector, err := labels.Parse("a=b")
		Expect(err).NotTo(HaveOccurred())
		o := &over_client.ListOptions{LabelSelector: labelSelector}
		newListOpts := &over_client.ListOptions{}
		o.ApplyToList(newListOpts)
		Expect(newListOpts).To(Equal(o))
	})
	It("Should set FieldSelector", func() {
		o := &over_client.ListOptions{FieldSelector: fields.Nothing()}
		newListOpts := &over_client.ListOptions{}
		o.ApplyToList(newListOpts)
		Expect(newListOpts).To(Equal(o))
	})
	It("Should set Namespace", func() {
		o := &over_client.ListOptions{Namespace: "my-ns"}
		newListOpts := &over_client.ListOptions{}
		o.ApplyToList(newListOpts)
		Expect(newListOpts).To(Equal(o))
	})
	It("Should set Raw", func() {
		o := &over_client.ListOptions{Raw: &metav1.ListOptions{FieldSelector: "Hans"}}
		newListOpts := &over_client.ListOptions{}
		o.ApplyToList(newListOpts)
		Expect(newListOpts).To(Equal(o))
	})
	It("Should set Limit", func() {
		o := &over_client.ListOptions{Limit: int64(1)}
		newListOpts := &over_client.ListOptions{}
		o.ApplyToList(newListOpts)
		Expect(newListOpts).To(Equal(o))
	})
	It("Should set Continue", func() {
		o := &over_client.ListOptions{Continue: "foo"}
		newListOpts := &over_client.ListOptions{}
		o.ApplyToList(newListOpts)
		Expect(newListOpts).To(Equal(o))
	})
	It("Should not set anything", func() {
		o := &over_client.ListOptions{}
		newListOpts := &over_client.ListOptions{}
		o.ApplyToList(newListOpts)
		Expect(newListOpts).To(Equal(o))
	})
})

var _ = Describe("GetOptions", func() {
	It("Should set Raw", func() {
		o := &over_client.GetOptions{Raw: &metav1.GetOptions{ResourceVersion: "RV0"}}
		newGetOpts := &over_client.GetOptions{}
		o.ApplyToGet(newGetOpts)
		Expect(newGetOpts).To(Equal(o))
	})
})

var _ = Describe("CreateOptions", func() {
	It("Should set DryRun", func() {
		o := &over_client.CreateOptions{DryRun: []string{"Hello", "Theodore"}}
		newCreateOpts := &over_client.CreateOptions{}
		o.ApplyToCreate(newCreateOpts)
		Expect(newCreateOpts).To(Equal(o))
	})
	It("Should set FieldManager", func() {
		o := &over_client.CreateOptions{FieldManager: "FieldManager"}
		newCreateOpts := &over_client.CreateOptions{}
		o.ApplyToCreate(newCreateOpts)
		Expect(newCreateOpts).To(Equal(o))
	})
	It("Should set Raw", func() {
		o := &over_client.CreateOptions{Raw: &metav1.CreateOptions{DryRun: []string{"Bye", "Theodore"}}}
		newCreateOpts := &over_client.CreateOptions{}
		o.ApplyToCreate(newCreateOpts)
		Expect(newCreateOpts).To(Equal(o))
	})
	It("Should not set anything", func() {
		o := &over_client.CreateOptions{}
		newCreateOpts := &over_client.CreateOptions{}
		o.ApplyToCreate(newCreateOpts)
		Expect(newCreateOpts).To(Equal(o))
	})
})

var _ = Describe("DeleteOptions", func() {
	It("Should set GracePeriodSeconds", func() {
		o := &over_client.DeleteOptions{GracePeriodSeconds: utilpointer.Int64(42)}
		newDeleteOpts := &over_client.DeleteOptions{}
		o.ApplyToDelete(newDeleteOpts)
		Expect(newDeleteOpts).To(Equal(o))
	})
	It("Should set Preconditions", func() {
		o := &over_client.DeleteOptions{Preconditions: &metav1.Preconditions{}}
		newDeleteOpts := &over_client.DeleteOptions{}
		o.ApplyToDelete(newDeleteOpts)
		Expect(newDeleteOpts).To(Equal(o))
	})
	It("Should set PropagationPolicy", func() {
		policy := metav1.DeletePropagationBackground
		o := &over_client.DeleteOptions{PropagationPolicy: &policy}
		newDeleteOpts := &over_client.DeleteOptions{}
		o.ApplyToDelete(newDeleteOpts)
		Expect(newDeleteOpts).To(Equal(o))
	})
	It("Should set Raw", func() {
		o := &over_client.DeleteOptions{Raw: &metav1.DeleteOptions{}}
		newDeleteOpts := &over_client.DeleteOptions{}
		o.ApplyToDelete(newDeleteOpts)
		Expect(newDeleteOpts).To(Equal(o))
	})
	It("Should set DryRun", func() {
		o := &over_client.DeleteOptions{DryRun: []string{"Hello", "Pippa"}}
		newDeleteOpts := &over_client.DeleteOptions{}
		o.ApplyToDelete(newDeleteOpts)
		Expect(newDeleteOpts).To(Equal(o))
	})
	It("Should not set anything", func() {
		o := &over_client.DeleteOptions{}
		newDeleteOpts := &over_client.DeleteOptions{}
		o.ApplyToDelete(newDeleteOpts)
		Expect(newDeleteOpts).To(Equal(o))
	})
})

var _ = Describe("UpdateOptions", func() {
	It("Should set DryRun", func() {
		o := &over_client.UpdateOptions{DryRun: []string{"Bye", "Pippa"}}
		newUpdateOpts := &over_client.UpdateOptions{}
		o.ApplyToUpdate(newUpdateOpts)
		Expect(newUpdateOpts).To(Equal(o))
	})
	It("Should set FieldManager", func() {
		o := &over_client.UpdateOptions{FieldManager: "Hello Boris"}
		newUpdateOpts := &over_client.UpdateOptions{}
		o.ApplyToUpdate(newUpdateOpts)
		Expect(newUpdateOpts).To(Equal(o))
	})
	It("Should set Raw", func() {
		o := &over_client.UpdateOptions{Raw: &metav1.UpdateOptions{}}
		newUpdateOpts := &over_client.UpdateOptions{}
		o.ApplyToUpdate(newUpdateOpts)
		Expect(newUpdateOpts).To(Equal(o))
	})
	It("Should not set anything", func() {
		o := &over_client.UpdateOptions{}
		newUpdateOpts := &over_client.UpdateOptions{}
		o.ApplyToUpdate(newUpdateOpts)
		Expect(newUpdateOpts).To(Equal(o))
	})
})

var _ = Describe("PatchOptions", func() {
	It("Should set DryRun", func() {
		o := &over_client.PatchOptions{DryRun: []string{"Bye", "Boris"}}
		newPatchOpts := &over_client.PatchOptions{}
		o.ApplyToPatch(newPatchOpts)
		Expect(newPatchOpts).To(Equal(o))
	})
	It("Should set Force", func() {
		o := &over_client.PatchOptions{Force: utilpointer.Bool(true)}
		newPatchOpts := &over_client.PatchOptions{}
		o.ApplyToPatch(newPatchOpts)
		Expect(newPatchOpts).To(Equal(o))
	})
	It("Should set FieldManager", func() {
		o := &over_client.PatchOptions{FieldManager: "Hello Julian"}
		newPatchOpts := &over_client.PatchOptions{}
		o.ApplyToPatch(newPatchOpts)
		Expect(newPatchOpts).To(Equal(o))
	})
	It("Should set Raw", func() {
		o := &over_client.PatchOptions{Raw: &metav1.PatchOptions{}}
		newPatchOpts := &over_client.PatchOptions{}
		o.ApplyToPatch(newPatchOpts)
		Expect(newPatchOpts).To(Equal(o))
	})
	It("Should not set anything", func() {
		o := &over_client.PatchOptions{}
		newPatchOpts := &over_client.PatchOptions{}
		o.ApplyToPatch(newPatchOpts)
		Expect(newPatchOpts).To(Equal(o))
	})
})

var _ = Describe("DeleteAllOfOptions", func() {
	It("Should set ListOptions", func() {
		o := &over_client.DeleteAllOfOptions{ListOptions: over_client.ListOptions{Raw: &metav1.ListOptions{}}}
		newDeleteAllOfOpts := &over_client.DeleteAllOfOptions{}
		o.ApplyToDeleteAllOf(newDeleteAllOfOpts)
		Expect(newDeleteAllOfOpts).To(Equal(o))
	})
	It("Should set DeleleteOptions", func() {
		o := &over_client.DeleteAllOfOptions{DeleteOptions: over_client.DeleteOptions{GracePeriodSeconds: utilpointer.Int64(44)}}
		newDeleteAllOfOpts := &over_client.DeleteAllOfOptions{}
		o.ApplyToDeleteAllOf(newDeleteAllOfOpts)
		Expect(newDeleteAllOfOpts).To(Equal(o))
	})
})

var _ = Describe("MatchingLabels", func() {
	It("Should produce an invalid selector when given invalid input", func() {
		matchingLabels := over_client.MatchingLabels(map[string]string{"k": "axahm2EJ8Phiephe2eixohbee9eGeiyees1thuozi1xoh0GiuH3diewi8iem7Nui"})
		listOpts := &over_client.ListOptions{}
		matchingLabels.ApplyToList(listOpts)

		r, _ := listOpts.LabelSelector.Requirements()
		_, err := labels.NewRequirement(r[0].Key(), r[0].Operator(), r[0].Values().List())
		Expect(err).To(HaveOccurred())
		expectedErrMsg := `values[0][k]: Invalid value: "axahm2EJ8Phiephe2eixohbee9eGeiyees1thuozi1xoh0GiuH3diewi8iem7Nui": must be no more than 63 characters`
		Expect(err.Error()).To(Equal(expectedErrMsg))
	})

	It("Should add matchingLabels to existing selector", func() {
		listOpts := &over_client.ListOptions{}

		matchingLabels := over_client.MatchingLabels(map[string]string{"k": "v"})
		matchingLabels2 := over_client.MatchingLabels(map[string]string{"k2": "v2"})

		matchingLabels.ApplyToList(listOpts)
		Expect(listOpts.LabelSelector.String()).To(Equal("k=v"))

		matchingLabels2.ApplyToList(listOpts)
		Expect(listOpts.LabelSelector.String()).To(Equal("k=v,k2=v2"))
	})
})

var _ = Describe("FieldOwner", func() {
	It("Should apply to PatchOptions", func() {
		o := &over_client.PatchOptions{FieldManager: "bar"}
		t := over_client.FieldOwner("foo")
		t.ApplyToPatch(o)
		Expect(o.FieldManager).To(Equal("foo"))
	})
	It("Should apply to CreateOptions", func() {
		o := &over_client.CreateOptions{FieldManager: "bar"}
		t := over_client.FieldOwner("foo")
		t.ApplyToCreate(o)
		Expect(o.FieldManager).To(Equal("foo"))
	})
	It("Should apply to UpdateOptions", func() {
		o := &over_client.UpdateOptions{FieldManager: "bar"}
		t := over_client.FieldOwner("foo")
		t.ApplyToUpdate(o)
		Expect(o.FieldManager).To(Equal("foo"))
	})
	It("Should apply to SubResourcePatchOptions", func() {
		o := &over_client.SubResourcePatchOptions{PatchOptions: over_client.PatchOptions{FieldManager: "bar"}}
		t := over_client.FieldOwner("foo")
		t.ApplyToSubResourcePatch(o)
		Expect(o.FieldManager).To(Equal("foo"))
	})
	It("Should apply to SubResourceCreateOptions", func() {
		o := &over_client.SubResourceCreateOptions{CreateOptions: over_client.CreateOptions{FieldManager: "bar"}}
		t := over_client.FieldOwner("foo")
		t.ApplyToSubResourceCreate(o)
		Expect(o.FieldManager).To(Equal("foo"))
	})
	It("Should apply to SubResourceUpdateOptions", func() {
		o := &over_client.SubResourceUpdateOptions{UpdateOptions: over_client.UpdateOptions{FieldManager: "bar"}}
		t := over_client.FieldOwner("foo")
		t.ApplyToSubResourceUpdate(o)
		Expect(o.FieldManager).To(Equal("foo"))
	})
})

var _ = Describe("ForceOwnership", func() {
	It("Should apply to PatchOptions", func() {
		o := &over_client.PatchOptions{}
		t := over_client.ForceOwnership
		t.ApplyToPatch(o)
		Expect(*o.Force).To(BeTrue())
	})
	It("Should apply to SubResourcePatchOptions", func() {
		o := &over_client.SubResourcePatchOptions{PatchOptions: over_client.PatchOptions{}}
		t := over_client.ForceOwnership
		t.ApplyToSubResourcePatch(o)
		Expect(*o.Force).To(BeTrue())
	})
})

var _ = Describe("HasLabels", func() {
	It("Should produce hasLabels in given order", func() {
		listOpts := &over_client.ListOptions{}

		hasLabels := over_client.HasLabels([]string{"labelApe", "labelFox"})
		hasLabels.ApplyToList(listOpts)
		Expect(listOpts.LabelSelector.String()).To(Equal("labelApe,labelFox"))
	})

	It("Should add hasLabels to existing hasLabels selector", func() {
		listOpts := &over_client.ListOptions{}

		hasLabel := over_client.HasLabels([]string{"labelApe"})
		hasLabel.ApplyToList(listOpts)

		hasOtherLabel := over_client.HasLabels([]string{"labelFox"})
		hasOtherLabel.ApplyToList(listOpts)
		Expect(listOpts.LabelSelector.String()).To(Equal("labelApe,labelFox"))
	})

	It("Should add hasLabels to existing matchingLabels", func() {
		listOpts := &over_client.ListOptions{}

		matchingLabels := over_client.MatchingLabels(map[string]string{"k": "v"})
		matchingLabels.ApplyToList(listOpts)

		hasLabel := over_client.HasLabels([]string{"labelApe"})
		hasLabel.ApplyToList(listOpts)
		Expect(listOpts.LabelSelector.String()).To(Equal("k=v,labelApe"))
	})
})
