// Copyright (c) 2021, Oracle and/or its affiliates.
//
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl/

package crd

import (
	"fmt"

	ndbv1alpha1 "github.com/mysql/ndb-operator/pkg/apis/ndbcontroller/v1alpha1"
	ndbclientset "github.com/mysql/ndb-operator/pkg/generated/clientset/versioned"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/test/e2e/framework"
)

// NewTestNdbCrd creates a new Ndb object for testing
func NewTestNdbCrd(namespace string, name string, datanodes, replicas, mysqlnodes int32) *ndbv1alpha1.NdbCluster {
	return &ndbv1alpha1.NdbCluster{
		TypeMeta: metav1.TypeMeta{
			APIVersion: ndbv1alpha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: ndbv1alpha1.NdbClusterSpec{
			NodeCount:       datanodes,
			RedundancyLevel: replicas,
			Mysqld: &ndbv1alpha1.NdbMysqldSpec{
				NodeCount: mysqlnodes,
			},
		},
	}
}

func LoadClientset() (*ndbclientset.Clientset, error) {
	config, err := framework.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("error creating ndb client: %v", err.Error())
	}
	ndbc, _ := ndbclientset.NewForConfig(config)
	return ndbc, nil
}
