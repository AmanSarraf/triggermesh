/*
Copyright 2022 TriggerMesh Inc.

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

package googlecloudstorage

import (
	"context"
	"io/ioutil"
	"net/url"

	. "github.com/onsi/ginkgo/v2" //nolint:stylecheck
	. "github.com/onsi/gomega"    //nolint:stylecheck

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	clientset "k8s.io/client-go/kubernetes"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"

	cloudevents "github.com/cloudevents/sdk-go/v2"

	"github.com/triggermesh/triggermesh/test/e2e/framework"
	"github.com/triggermesh/triggermesh/test/e2e/framework/apps"
	e2ece "github.com/triggermesh/triggermesh/test/e2e/framework/cloudevents"
	"github.com/triggermesh/triggermesh/test/e2e/framework/ducktypes"
	e2egcloud "github.com/triggermesh/triggermesh/test/e2e/framework/gcloud"
	e2estorage "github.com/triggermesh/triggermesh/test/e2e/framework/gcloud/storage"
)

/* This test suite requires:

   - A Google Cloud Service Account key in JSON format, exported in the environment as GCLOUD_SERVICEACCOUNT_KEY
   - The name of the Google Cloud project exported in the environment as GCLOUD_PROJECT
*/

var targetAPIVersion = schema.GroupVersion{
	Group:   "targets.triggermesh.io",
	Version: "v1alpha1",
}

const (
	targetKind     = "GoogleCloudStorageTarget"
	targetResource = "googlecloudstoragetargets"
)

const gcpCredentialsJosnSecretKey = "credentialsJson"

var _ = Describe("Google Cloud Storage target", func() {
	// NOTE: bucket names aren't allowed to contain "goog"
	f := framework.New("gcloudstoragetarget")

	var ns string

	var trgtClient dynamic.ResourceInterface

	var bucketName string
	var gcpSecret *corev1.Secret

	BeforeEach(func() {
		ns = f.UniqueName

		gvr := targetAPIVersion.WithResource(targetResource)
		trgtClient = f.DynamicClient.Resource(gvr).Namespace(ns)
	})

	Context("a target is deployed", func() {
		var trgtURL *url.URL
		var storageClient *storage.Client

		var sentEvent *cloudevents.Event

		BeforeEach(func() {
			var err error
			serviceaccountKey := e2egcloud.ServiceAccountKeyFromEnv()
			gcloudProject := e2egcloud.ProjectNameFromEnv()

			storageClient, err = storage.NewClient(context.Background(), option.WithCredentialsJSON([]byte(serviceaccountKey)))
			Expect(err).ToNot(HaveOccurred())

			gcpSecret = createGCPCredsSecret(f.KubeClient, ns, serviceaccountKey)

			By("creating a Google Cloud Storage Bucket", func() {
				bucketName = e2estorage.CreateBucket(storageClient, gcloudProject, f)

				DeferCleanup(func() {
					By("deleting Google Cloud Storage Bucket "+bucketName, func() {
						e2estorage.DeleteBucket(storageClient, bucketName)
					})
				})
			})
		})

		When("the spec contains default settings", func() {
			BeforeEach(func() {
				By("creating an GoogleCloudStorageTarget object", func() {
					trgt, err := createTarget(trgtClient, ns, "test-",
						withBucketName(bucketName),
						withCredentials(gcpSecret.Name),
					)
					Expect(err).ToNot(HaveOccurred())

					trgt = ducktypes.WaitUntilReady(f.DynamicClient, trgt)

					trgtURL = ducktypes.Address(trgt)
					Expect(trgtURL).ToNot(BeNil())
				})
			})

			When("an event is sent to the target", func() {

				BeforeEach(func() {
					By("sending an event", func() {
						sentEvent = e2ece.NewHelloEvent(f)

						job := e2ece.RunEventSender(f.KubeClient, ns, trgtURL.String(), sentEvent)
						apps.WaitForCompletion(f.KubeClient, job)
					})
				})

				It("creates an object into the bucket", func() {
					var receivedObj []byte

					By("listing the bucket objects", func() {
						var err error

						receivedObjs := e2estorage.GetObjectsReader(storageClient, bucketName)
						Expect(receivedObjs).To(HaveLen(1),
							"Received %d objects instead of 1", len(receivedObjs))

						object := receivedObjs[0]
						receivedObj, err = ioutil.ReadAll(object)
						Expect(err).ToNot(HaveOccurred())
					})

					By("inspecting the object payload", func() {
						Expect(receivedObj).To(Equal(sentEvent.Data()))
					})
				})
			})
		})
	})

	When("a client creates a target object with invalid specs", func() {

		// Those tests do not require a real bucketName or gcpSecret
		BeforeEach(func() {
			bucketName = "test"
			gcpSecret = &corev1.Secret{}
		})

		// Here we use
		//   "Specify: the API server rejects ..., By: setting an invalid ..."
		// instead of
		//   "When: it sets an invalid ..., Specify: the API server rejects ..."
		// to avoid creating a namespace for each spec, due to their simplicity.
		Specify("the API server rejects the creation of that object", func() {

			By("omitting the bucketName", func() {
				_, err := createTarget(trgtClient, ns, "test-no-bucketName",
					withCredentials(gcpSecret.Name),
				)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("spec.bucketName: Required value"))
			})

			By("omitting the credentials", func() {
				_, err := createTarget(trgtClient, ns, "test-nocreds",
					withBucketName(bucketName),
				)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(
					"spec.credentialsJson: Required value"))
			})
		})
	})
})

// createTarget creates an GoogleCloudStorageTarget object initialized with the given options.
func createTarget(trgtClient dynamic.ResourceInterface, namespace, namePrefix string, opts ...targetOption) (*unstructured.Unstructured, error) {

	trgt := &unstructured.Unstructured{}
	trgt.SetAPIVersion(targetAPIVersion.String())
	trgt.SetKind(targetKind)
	trgt.SetNamespace(namespace)
	trgt.SetGenerateName(namePrefix)

	for _, opt := range opts {
		opt(trgt)
	}

	return trgtClient.Create(context.Background(), trgt, metav1.CreateOptions{})
}

type targetOption func(*unstructured.Unstructured)

func withBucketName(bucketName string) targetOption {
	return func(trgt *unstructured.Unstructured) {
		if err := unstructured.SetNestedField(trgt.Object, bucketName, "spec", "bucketName"); err != nil {
			framework.FailfWithOffset(2, "Failed to set spec.bucketName field: %s", err)
		}
	}
}

func withCredentials(secretName string) targetOption {
	credentials := map[string]interface{}{
		"secretKeyRef": map[string]interface{}{
			"name": secretName,
			"key":  gcpCredentialsJosnSecretKey,
		},
	}

	return func(trgt *unstructured.Unstructured) {
		if err := unstructured.SetNestedMap(trgt.Object, credentials, "spec", gcpCredentialsJosnSecretKey); err != nil {
			framework.FailfWithOffset(2, "Failed to set spec.credentialsJson field: %s", err)
		}
	}
}

// createGCPCredsSecret creates a Kubernetes Secret containing GCP credentials.
func createGCPCredsSecret(c clientset.Interface, namespace string, creds string) *corev1.Secret {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    namespace,
			GenerateName: "gcp-creds-",
		},
		StringData: map[string]string{
			gcpCredentialsJosnSecretKey: creds,
		},
	}

	secret, err := c.CoreV1().Secrets(namespace).Create(context.Background(), secret, metav1.CreateOptions{})
	if err != nil {
		framework.FailfWithOffset(2, "Failed to create Secret: %s", err)
	}

	return secret
}