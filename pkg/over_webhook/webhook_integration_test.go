/*
Copyright 2021 The Kubernetes Authors.

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

package over_webhook_test

import (
	"context"
	"crypto/tls"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/over_metrics/server"

	"sigs.k8s.io/controller-runtime/pkg/over_client"
	"sigs.k8s.io/controller-runtime/pkg/over_manager"
	"sigs.k8s.io/controller-runtime/pkg/over_webhook"
	"sigs.k8s.io/controller-runtime/pkg/over_webhook/admission"
)

var _ = Describe("Webhook", func() {
	var c over_client.Client
	var obj *appsv1.Deployment
	BeforeEach(func() {
		Expect(cfg).NotTo(BeNil())
		var err error
		c, err = over_client.New(cfg, over_client.Options{})
		Expect(err).NotTo(HaveOccurred())

		obj = &appsv1.Deployment{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-deployment",
				Namespace: "default",
			},
			Spec: appsv1.DeploymentSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"foo": "bar"},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"foo": "bar"}},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "nginx",
								Image: "nginx",
							},
						},
					},
				},
			},
		}
	})
	Context("when running a over_webhook server with a over_manager", func() {
		It("should reject create request for over_webhook that rejects all requests", func() {
			m, err := over_manager.New(cfg, over_manager.Options{
				WebhookServer: over_webhook.NewServer(over_webhook.Options{
					Port:    testenv.WebhookInstallOptions.LocalServingPort,
					Host:    testenv.WebhookInstallOptions.LocalServingHost,
					CertDir: testenv.WebhookInstallOptions.LocalServingCertDir,
					TLSOpts: []func(*tls.Config){func(config *tls.Config) {}},
				}),
			}) // we need over_manager here just to leverage over_manager.SetFields
			Expect(err).NotTo(HaveOccurred())
			server := m.GetWebhookServer()
			server.Register("/failing", &over_webhook.Admission{Handler: &rejectingValidator{d: admission.NewDecoder(testenv.Scheme)}})

			ctx, cancel := context.WithCancel(context.Background())
			go func() {
				err := server.Start(ctx)
				Expect(err).NotTo(HaveOccurred())
			}()

			Eventually(func() bool {
				err := c.Create(context.TODO(), obj)
				return err != nil && strings.HasSuffix(err.Error(), "Always denied") && apierrors.ReasonForError(err) == metav1.StatusReasonForbidden
			}, 1*time.Second).Should(BeTrue())

			cancel()
		})
		It("should reject create request for multi-over_webhook that rejects all requests", func() {
			m, err := over_manager.New(cfg, over_manager.Options{
				Metrics: metricsserver.Options{BindAddress: "0"},
				WebhookServer: over_webhook.NewServer(over_webhook.Options{
					Port:    testenv.WebhookInstallOptions.LocalServingPort,
					Host:    testenv.WebhookInstallOptions.LocalServingHost,
					CertDir: testenv.WebhookInstallOptions.LocalServingCertDir,
					TLSOpts: []func(*tls.Config){func(config *tls.Config) {}},
				}),
			}) // we need over_manager here just to leverage over_manager.SetFields
			Expect(err).NotTo(HaveOccurred())
			server := m.GetWebhookServer()
			server.Register("/failing", &over_webhook.Admission{Handler: admission.MultiValidatingHandler(&rejectingValidator{d: admission.NewDecoder(testenv.Scheme)})})

			ctx, cancel := context.WithCancel(context.Background())
			go func() {
				err = server.Start(ctx)
				Expect(err).NotTo(HaveOccurred())
			}()

			Eventually(func() bool {
				err = c.Create(context.TODO(), obj)
				return err != nil && strings.HasSuffix(err.Error(), "Always denied") && apierrors.ReasonForError(err) == metav1.StatusReasonForbidden
			}, 1*time.Second).Should(BeTrue())

			cancel()
		})
	})
	Context("when running a over_webhook server without a over_manager", func() {
		It("should reject create request for over_webhook that rejects all requests", func() {
			server := over_webhook.NewServer(over_webhook.Options{
				Port:    testenv.WebhookInstallOptions.LocalServingPort,
				Host:    testenv.WebhookInstallOptions.LocalServingHost,
				CertDir: testenv.WebhookInstallOptions.LocalServingCertDir,
			})
			server.Register("/failing", &over_webhook.Admission{Handler: &rejectingValidator{d: admission.NewDecoder(testenv.Scheme)}})

			ctx, cancel := context.WithCancel(context.Background())
			go func() {
				err := server.Start(ctx)
				Expect(err).NotTo(HaveOccurred())
			}()

			Eventually(func() bool {
				err := c.Create(context.TODO(), obj)
				return err != nil && strings.HasSuffix(err.Error(), "Always denied") && apierrors.ReasonForError(err) == metav1.StatusReasonForbidden
			}, 1*time.Second).Should(BeTrue())

			cancel()
		})
	})
})

type rejectingValidator struct {
	d *admission.Decoder
}

func (v *rejectingValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	var obj appsv1.Deployment
	if err := v.d.Decode(req, &obj); err != nil {
		return admission.Denied(err.Error())
	}
	return admission.Denied("Always denied")
}
