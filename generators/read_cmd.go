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
	// "io/ioutil"
	// "os"

	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"k8s.io/kubernetes/pkg/kubectl/cmd"
	// cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

func GetSpec() KubectlSpec {
	// Initialize a kubectl command that we can use to get the help documentation
	kubectl := cmd.NewDefaultKubectlCommand()

	// Create the structural representation
	return NewKubectlSpec(kubectl)
}

func NewKubectlSpec(c *cobra.Command) KubectlSpec {
	return KubectlSpec{
		TopLevelCommandGroups: []TopLevelCommands{NewTopLevelCommands(c.Commands())},
	}
}

func NewTopLevelCommands(cs []*cobra.Command) TopLevelCommands {
	tlc := TopLevelCommands{}
	for _, c := range cs {
		tlc.Commands = append(tlc.Commands, NewTopLevelCommand(c))
	}
	sort.Sort(tlc)
	return tlc
}

func NewTopLevelCommand(c *cobra.Command) TopLevelCommand {
	result := TopLevelCommand{
		MainCommand: NewCommand(c, ""),
	}
	for _, sub := range c.Commands() {
		result.SubCommands = append(result.SubCommands, NewSubCommands(sub, c.Name())...)
	}
	sort.Sort(result.SubCommands)
	return result
}

// Parse the Options
func NewOptions(flags *pflag.FlagSet) Options {
	result := Options{}
	flags.VisitAll(func(flag *pflag.Flag) {
		opt := &Option{
			Name:         flag.Name,
			Shorthand:    flag.Shorthand,
			DefaultValue: flag.DefValue,
			Usage:        flag.Usage,
			Type:         flag.Value.Type(),
		}
		result = append(result, opt)
	})
	return result
}

// Parse the Commands
func NewSubCommands(c *cobra.Command, path string) Commands {
	subCommands := Commands{NewCommand(c, path)}
	for _, subCommand := range c.Commands() {
		subCommands = append(subCommands, NewSubCommands(subCommand, path+"/"+c.Name())...)
	}
	return subCommands
}

func NewCommand(c *cobra.Command, path string) *Command {
	return &Command{
		Name:             c.Name(),
		Path:             path,
		Description:      c.Long,
		Synopsis:         c.Short,
		Examples:         SplitExamples(c.Example),
		Options:          NewOptions(c.NonInheritedFlags()),
		InheritedOptions: NewOptions(c.InheritedFlags()),
		Usage:            c.Use,
	}
}

func (a Options) Len() int      { return len(a) }
func (a Options) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a Options) Less(i, j int) bool {
	return a[i].Name < a[j].Name
}

func (a TopLevelCommands) Len() int      { return len(a.Commands) }
func (a TopLevelCommands) Swap(i, j int) { a.Commands[i], a.Commands[j] = a.Commands[j], a.Commands[i] }
func (a TopLevelCommands) Less(i, j int) bool {
	return a.Commands[i].MainCommand.Path < a.Commands[j].MainCommand.Path
}

func (a Commands) Len() int      { return len(a) }
func (a Commands) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a Commands) Less(i, j int) bool {
	return a[i].Path < a[j].Path
}

const (
	start = iota
	title
	content
)

func SplitExamples(examples string) (result []Example) {
	lines := strings.Split(examples, "\n")
	pos := start
	currentExample := Example{}
	for _, line := range lines {
		line = escapeXml(strings.Trim(line, " "))
		if len(line) == 0 {
			continue
		}

		if line[0] == '#' {
			if pos == title {
				line = strings.TrimLeft(line, "# ")
				currentExample.Title += "\n " + line
			} else {
				if len(currentExample.Title) > 0 || len(currentExample.Content) > 0 {
					result = append(result, currentExample)
				}
				currentExample = Example{}
				line = strings.TrimLeft(line, "# ")
				currentExample.Title = line
				pos = title
			}
		} else {
			if pos == content {
				currentExample.Content += "\n" + line
			} else {
				currentExample.Content += line
				pos = content
			}
		}
	}
	if len(currentExample.Title) > 0 || len(currentExample.Content) > 0 {
		result = append(result, currentExample)
	}
	return
}
