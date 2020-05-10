/*
Copyright 2016 The Kubernetes Authors.
Copyright 2019 Philippe Martin

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

package generators

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var KubernetesVersion = flag.String("kubernetes-version", "", "Version of Kubernetes to generate docs for.")

var GenKubectlDir = flag.String("gen-kubectl-dir", "generators", "Directory containing kubectl files")

var ShowUsage = flag.Bool("show-usage", false, "Show original usage (for debugging)")

func getTocFile() string {
	return filepath.Join(*GenKubectlDir, *KubernetesVersion, "toc.yaml")
}

func getStaticIncludesDir() string {
	return filepath.Join(*GenKubectlDir, *KubernetesVersion, "static_includes")
}

func AsDocbook() {

	f, err := os.Create("build/index.xml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Fprintf(f, `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE book PUBLIC "-//OASIS//DTD DocBook XML V4.5//EN"
"http://www.oasis-open.org/docbook/xml/4.5/docbookx.dtd">
<book>
  <bookinfo>
    <title>Kubectl Reference</title>

    <subtitle>v1.18.0</subtitle>

    <releaseinfo>By the Kubernetes Authors</releaseinfo>

    <releaseinfo>Edited and published by Philippe Martin</releaseinfo>

    <copyright>
      <year>2020</year>

      <holder>The Kubernetes Authors</holder>
    </copyright>

    <legalnotice>
      <para>Permission is granted to copy, distribute and/or modify this
      document under the terms of the Apache License version 2. A copy of the
      license is included in <xref linkend="license"/>.</para>
    </legalnotice>

    <legalnotice>
      <para>The tool used to generate this document is available at
      https://github.com/feloy/kubectl-reference</para>
    </legalnotice>
  </bookinfo>
`)

	spec := GetSpec()

	toc := ToC{}
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

	for _, category := range toc.Categories {
		fmt.Fprintf(f, "  <reference><title>%s</title>\n", category.Name)

		for _, tocCommand := range category.Commands {
			command := spec.GetCommand(tocCommand.Name)
			if command == nil {
				fmt.Printf("command %s not found", tocCommand.Name)
				os.Exit(1)
			}
			command.AsDocbook(f, tocCommand)
		}
		fmt.Fprintf(f, `</reference>`)
	}

	addLicense(f)

	fmt.Fprintf(f, `</book>`)
}

func addLicense(w io.Writer) {
	f, err := os.Open("./static/license.xml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		w.Write(scanner.Bytes())
		w.Write([]byte("\n"))
	}
}
