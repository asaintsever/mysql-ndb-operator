// Copyright (c) 2021, Oracle and/or its affiliates.
//
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl/

package yaml

import (
	"encoding/json"
	"fmt"
	"github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
	"k8s.io/kubernetes/test/e2e/framework"
	"strings"

	"github.com/mysql/ndb-operator/pkg/apis/ndbcontroller/v1alpha1"
)

// K8sObject identifies a K8s object in a yaml file
type K8sObject struct {
	Name, Kind, Version string
}

// ExtractObjectsFromYaml extracts the resource identified by
// the Kind, Version, Name from the given yaml file
func ExtractObjectsFromYaml(path, filename string, k8sObjects []K8sObject, ns string) string {

	// generate a map of k8s object identifiers with
	// key name;kind;version for easier lookup
	k8sObjMap := make(map[string]struct{})
	for _, obj := range k8sObjects {
		k8sObjMap[obj.Name+";"+obj.Kind+";"+obj.Version] = struct{}{}
	}

	// read the given file
	yamlFile := YamlFile(path, filename)

	// split the various docs in the file, read one by one
	// and extract them into a single list
	var objsToBeCreated []string
	yamlDocs := strings.Split(yamlFile, "---")
	for _, yamlDoc := range yamlDocs {
		if len(k8sObjMap) == 0 {
			// Found all objects to be created
			break
		}

		var resource interface{}
		err := yaml.Unmarshal([]byte(yamlDoc), &resource)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		if resource == nil {
			// empty yaml doc
			continue
		}

		// non empty yaml doc - extract the name/kind/version
		resourceMap := resource.(map[interface{}]interface{})
		metadata := resourceMap["metadata"].(map[interface{}]interface{})
		name := metadata["name"].(string)
		kind := resourceMap["kind"].(string)
		apiVersion := resourceMap["apiVersion"].(string)
		// if this object exists in the map, update and append it to output
		objKey := name + ";" + kind + ";" + apiVersion
		if _, exists := k8sObjMap[objKey]; exists {
			// replace namespace if required
			if ns != "" {
				replaceAll(&resource, "namespace", ns)
			}

			// append the resource to the list
			updatedDoc, err := yaml.Marshal(resource)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			objsToBeCreated = append(objsToBeCreated, string(updatedDoc))

			// delete the key from map
			delete(k8sObjMap, objKey)
		}
	}

	if len(k8sObjMap) != 0 {
		// few resources not found
		var errorMsg []string
		for key := range k8sObjMap {
			res := strings.Split(key, ";")
			errorMsg = append(errorMsg, fmt.Sprintf(
				"Requested resource '%s' with kind '%s' and version '%s' not found", res[0], res[1], res[2]))
		}
		framework.Fail(strings.Join(errorMsg, "\n"))
	}

	// Return the docs as a single yaml string
	return strings.Join(objsToBeCreated, "\n---\n")
}

// MarshalNdb marshals the given Ndb object into yaml
func MarshalNdb(ndb *v1alpha1.NdbCluster) []byte {
	// The yaml library will not respect the 'json' flags
	// in the Ndb type. Convert Ndb object to json and
	// back to a map to get the yaml in right format.
	// TODO: Find a better solution or update Ndb flags.
	b, err := json.Marshal(ndb)
	framework.ExpectNoError(err, "MarshalNdb failed")

	var data map[string]interface{}
	err = json.Unmarshal(b, &data)
	framework.ExpectNoError(err, "MarshalNdb failed")

	yamlContent, err := yaml.Marshal(data)
	framework.ExpectNoError(err, "MarshalNdb failed")
	return yamlContent
}
