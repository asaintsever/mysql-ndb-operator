// Copyright (c) 2021, Oracle and/or its affiliates.
//
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl/

package secret

import (
	"context"
	"fmt"

	"github.com/mysql/ndb-operator/pkg/apis/ndbcontroller/v1alpha1"
	"github.com/mysql/ndb-operator/pkg/resources"
	"github.com/onsi/gomega"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/test/e2e/framework"
)

const (
	testRootPassword = "ndbpass"
)

func CreateSecretForMySQLRootAccount(clientset kubernetes.Interface, secretName, namespace string) {
	// build Secret, create it in K8s and return
	rootPassSecret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		Data: map[string][]byte{v1.BasicAuthPasswordKey: []byte(testRootPassword)},
		Type: v1.SecretTypeBasicAuth,
	}

	_, err := clientset.CoreV1().Secrets(namespace).Create(context.TODO(), rootPassSecret, metav1.CreateOptions{})
	framework.ExpectNoError(err, "failed to create the custom secret")
}

func DeleteSecret(clientset kubernetes.Interface, secretName, namespace string) {
	err := clientset.CoreV1().Secrets(namespace).Delete(context.TODO(), secretName, metav1.DeleteOptions{})
	framework.ExpectNoError(err, "failed to delete the custom secret")
}

// GetMySQLRootPassword returns the root password for the MySQL Servers maintained by the given NdbCluster.
func GetMySQLRootPassword(ctx context.Context, clientset kubernetes.Interface, nc *v1alpha1.NdbCluster) string {
	gomega.Expect(nc.GetMySQLServerNodeCount()).NotTo(
		gomega.BeZero(), fmt.Sprintf("No MySQL Servers configured for NdbCluster %q", nc.Name))
	// Retrieve the Secret
	secretName, _ := resources.GetMySQLRootPasswordSecretName(nc)
	secret, err := clientset.CoreV1().Secrets(nc.Namespace).Get(ctx, secretName, metav1.GetOptions{})
	framework.ExpectNoError(err, "failed to retrieve the MySQL root password secret")

	// Extract the password
	password := secret.Data[v1.BasicAuthPasswordKey]
	gomega.Expect(password).NotTo(
		gomega.BeEmpty(), "MySQL root password was not found in secret")
	return string(password)
}
