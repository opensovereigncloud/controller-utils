// Copyright 2021 OnMetal authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metautils_test

import (
	"reflect"

	. "github.com/onmetal/controller-utils/metautils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

var _ = Describe("Metautils", func() {
	Describe("ListElementType", func() {
		It("should return the element type of an object list", func() {
			t, err := ListElementType(&appsv1.DeploymentList{})
			Expect(err).NotTo(HaveOccurred())
			Expect(t).To(Equal(reflect.TypeOf(appsv1.Deployment{})))
		})

		It("should error if the list is not a valid object list", func() {
			_, err := ListElementType(&appsv1.Deployment{})
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("GVKForList", func() {
		It("should return the GVK for the list", func() {
			gvk, err := GVKForList(scheme.Scheme, &appsv1.DeploymentList{})
			Expect(err).NotTo(HaveOccurred())
			Expect(gvk).To(Equal(appsv1.SchemeGroupVersion.WithKind("Deployment")))
		})

		It("should error if no gvk could be obtained for the given object", func() {
			_, err := GVKForList(scheme.Scheme, &unstructured.UnstructuredList{})
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("ConvertAndSetList", func() {
		It("should convert the given objects for the list and insert them", func() {
			list := &corev1.ConfigMapList{}
			Expect(ConvertAndSetList(scheme.Scheme, list, []runtime.Object{
				&corev1.ConfigMap{},
				&unstructured.Unstructured{Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "ConfigMap",
				}},
			})).NotTo(HaveOccurred())
			Expect(list).To(Equal(&corev1.ConfigMapList{
				Items: []corev1.ConfigMap{{}, {}},
			}))
		})

		It("should error if the given list is not a list", func() {
			Expect(ConvertAndSetList(scheme.Scheme, &corev1.ConfigMap{}, nil)).To(HaveOccurred())
		})

		It("should error if an element could not be converted", func() {
			Expect(ConvertAndSetList(
				scheme.Scheme,
				&corev1.ConfigMap{},
				[]runtime.Object{&corev1.Secret{}},
			)).To(HaveOccurred())
		})
	})
})