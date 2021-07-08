// Copyright (c) 2020, 2021, Oracle and/or its affiliates.
//
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl/

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/mysql/ndb-operator/pkg/apis/ndbcontroller/v1alpha1"
	"github.com/mysql/ndb-operator/pkg/generated/clientset/versioned/scheme"
	rest "k8s.io/client-go/rest"
)

type MysqlV1alpha1Interface interface {
	RESTClient() rest.Interface
	NdbClustersGetter
}

// MysqlV1alpha1Client is used to interact with features provided by the mysql.oracle.com group.
type MysqlV1alpha1Client struct {
	restClient rest.Interface
}

func (c *MysqlV1alpha1Client) NdbClusters(namespace string) NdbClusterInterface {
	return newNdbClusters(c, namespace)
}

// NewForConfig creates a new MysqlV1alpha1Client for the given config.
func NewForConfig(c *rest.Config) (*MysqlV1alpha1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &MysqlV1alpha1Client{client}, nil
}

// NewForConfigOrDie creates a new MysqlV1alpha1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *MysqlV1alpha1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new MysqlV1alpha1Client for the given RESTClient.
func New(c rest.Interface) *MysqlV1alpha1Client {
	return &MysqlV1alpha1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1alpha1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *MysqlV1alpha1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
