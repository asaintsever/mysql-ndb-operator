package yaml

import (
	"os"
	"path/filepath"

	"k8s.io/klog"

	"k8s.io/kubernetes/test/e2e/framework"
	"k8s.io/kubernetes/test/e2e/framework/testfiles"
)

// YamlFile reads the content of a single yamle file as string
// Read happens from a test files path which needs to first
// be registered with testfiles.AddFileSource()
func YamlFile(test, file string) string {
	from := filepath.Join(test, file+".yaml")
	data, err := testfiles.Read(from)
	if err != nil {
		dir, _ := os.Getwd()
		klog.Infof("Maybe in wrong directory %s", dir)
		framework.Fail(err.Error())
	}
	return string(data)
}
