// Copyright (c) 2021, Oracle and/or its affiliates.
//
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl/

package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mysql/ndb-operator/pkg/apis/ndbcontroller/v1alpha1"
	v1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"testing"
)

func createNewValidatingWebhookConfig(
	t *testing.T, f *fixture,
	validatingWebhookConfigName, serviceName string, numberOfWebhooks uint) {
	t.Helper()

	// create a new vwc with minimal params for testing
	newVwc := &v1.ValidatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: validatingWebhookConfigName,
			Labels: map[string]string{
				"webhook-server": serviceName,
			},
		},
	}

	for i := uint(0); i < numberOfWebhooks; i++ {
		newVwc.Webhooks = append(newVwc.Webhooks, v1.ValidatingWebhook{
			Name: fmt.Sprintf("webhook%d", i+1),
			ClientConfig: v1.WebhookClientConfig{
				Service: &v1.ServiceReference{
					Namespace: "default",
					Name:      serviceName,
				},
			},
		})
	}

	// create it in k8s
	var err error
	vwcInterface := f.kubeclient.AdmissionregistrationV1().ValidatingWebhookConfigurations()
	if newVwc, err = vwcInterface.Create(context.Background(), newVwc, metav1.CreateOptions{}); err != nil {
		t.Fatalf("Error creating validating webhook configs : %s", err.Error())
	}

	// update expected action
	f.expectCreateAction("", "admissionregistration.k8s.io", "v1", "validatingwebhookconfigurations", newVwc)
}

func generateExpectedPatch(
	t *testing.T, serviceName string, numberOfWebhooks uint, cert []byte) []byte {
	t.Helper()

	// Build the webhook config diff
	var webhooks []v1.ValidatingWebhook
	for i := uint(0); i < numberOfWebhooks; i++ {
		webhooks = append(webhooks, v1.ValidatingWebhook{
			Name: fmt.Sprintf("webhook%d", i+1),
			ClientConfig: v1.WebhookClientConfig{
				Service: &v1.ServiceReference{
					Namespace: "default",
					Name:      serviceName,
				},
				CABundle: cert,
			},
		})
	}
	diff := map[string]interface{}{"webhooks": webhooks}

	// Generate the patch for the change
	patch, err := json.Marshal(diff)
	if err != nil {
		t.Fatal("Failed to marshal diff :", err)
	}

	// Re-marshal to sort all the json keys
	var ifce interface{}
	_ = json.Unmarshal(patch, &ifce)
	patch, _ = json.Marshal(ifce)

	return patch
}

func Test_ValidatingWebhook_UpdateWebhookConfigCertificate(t *testing.T) {
	// Create fixture and start informers
	f := newFixture(t, &v1alpha1.NdbCluster{})
	defer f.close()
	f.start()

	// create 2 webhook configs
	serviceName := "test-service"
	createNewValidatingWebhookConfig(
		t, f, "test-validating-webhook-config-1", serviceName, 1)
	createNewValidatingWebhookConfig(
		t, f, "test-validating-webhook-config-2", serviceName, 2)

	// use a simple string as the certificate to be updated in the webhook config
	cert := []byte("CERTIFICATE")

	// Update the webhook configs using the controller
	vwcController := NewValidatingWebhookConfigController(f.kubeclient)
	if !vwcController.UpdateWebhookConfigCertificate(
		context.Background(), "webhook-server="+serviceName, cert) {
		t.Fatal("Failed to update the validating webhook configs")
	}

	// update expected patch action
	f.expectPatchAction("", "validatingwebhookconfigurations",
		"test-validating-webhook-config-1", types.StrategicMergePatchType,
		generateExpectedPatch(t, serviceName, 1, cert))
	f.expectPatchAction("", "validatingwebhookconfigurations",
		"test-validating-webhook-config-2", types.StrategicMergePatchType,
		generateExpectedPatch(t, serviceName, 2, cert))

	// verify everything went as expected
	f.checkActions()
}
