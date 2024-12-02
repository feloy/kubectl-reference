/*
Copyright 2019 Philippe Martin.

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
package upgrade

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/feloy/kubectl-reference/generators"
	"gopkg.in/yaml.v2"
)

func getTocFile() string {
	return filepath.Join(*generators.GenKubectlDir, *generators.KubernetesVersion, "toc.yaml")
}

func Upgrade() {
	flag.Parse()

	toc := generators.ToC{}
	if len(getTocFile()) < 1 {
		fmt.Printf("Must specify --toc-file.\n")
		os.Exit(2)
	}

	contents, err := ioutil.ReadFile(getTocFile())
	if err != nil {
		fmt.Printf("Failed to read yaml file %s: %v", getTocFile(), err)
	}

	err = yaml.Unmarshal(contents, &toc)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	spec := generators.GetSpec()

	toc.AddMissingCommands(&spec)
	toc.AddMissingOptions(&spec)
	toc.AddMissingUsages(&spec)

	bytes, _ := yaml.Marshal(toc)
	fmt.Print(string(bytes))
}
