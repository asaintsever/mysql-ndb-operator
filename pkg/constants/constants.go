// Copyright (c) 2020, Oracle and/or its affiliates.
//
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl/

package constants

import "github.com/mysql/ndb-operator/pkg/apis/ndbcontroller"

// ClusterLabel is applied to all components of a Ndb cluster
const ClusterLabel = ndbcontroller.GroupName + "/v1alpha1"
const ClusterNodeTypeLabel = ndbcontroller.GroupName + "/nodetype"

const DataDir = "/var/lib/ndb"