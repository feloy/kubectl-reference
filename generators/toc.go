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
package generators

import (
	"fmt"
	"os"
)

type ToC struct {
	Categories []*Category `yaml:",omitempty"`
}

type Category struct {
	Name     string        `yaml:",omitempty"`
	Commands []*ToCCommand `yaml:",omitempty"`
	Include  string        `yaml:",omitempty"`
}

type ToCCommand struct {
	Name          string         `yaml:",omitempty"`
	Usage         string         `yaml:",omitempty"`
	Args          []Arg          `yaml:",omitempty"`
	OptionsGroups []OptionsGroup `yaml:"optionsgroups,omitempty"`
}

type Arg struct {
	Name   string  `yaml:",omitempty"`
	End    bool    `yaml:",omitempty"`
	Choice *string `yaml:",omitempty"`
	Rep    *string `yaml:",omitempty"`
}

type OptionsGroup struct {
	Name    string      `yaml:",omitempty"`
	Options []ToCOption `yaml:",omitempty"`
}

type ToCOption struct {
	Name      string  `yaml:",omitempty"`
	Required  bool    `yaml:",omitempty"`
	Type      *string `yaml:",omitempty"`
	Usage     *string `yaml:",omitempty"`
	Shorthand *string `yaml:",omitempty"`
	Default   *string `yaml:",omitempty"`
}

func (o *ToC) GetAllCommandNames() (commands []string) {
	for _, category := range o.Categories {
		for _, command := range category.Commands {
			commands = append(commands, command.Name)
		}
	}
	return
}

func (o *ToCCommand) GetAllOptionNames() (options []string) {
	for _, group := range o.OptionsGroups {
		for _, option := range group.Options {
			options = append(options, option.Name)
		}
	}
	return
}

func (o *ToC) AddMissingCommands(spec *KubectlSpec) {

	commandsInToC := map[string]struct{}{}

	for _, c := range o.GetAllCommandNames() {
		commandsInToC[c] = struct{}{}
	}

	categoryOthers := Category{
		Name: "Other commands",
	}

	for _, c := range spec.GetAllCommandNames() {
		if _, found := commandsInToC[c]; !found {
			fmt.Fprintf(os.Stderr, "command %s not found\n", c)
			categoryOthers.Commands = append(categoryOthers.Commands, &ToCCommand{
				Name: c,
			})
		}
	}

	if len(categoryOthers.Commands) > 0 {
		o.Categories = append(o.Categories, &categoryOthers)
	}

}

func (o *ToC) AddMissingOptions(spec *KubectlSpec) {
	for _, cats := range o.Categories {
		for _, command := range cats.Commands {
			cmd := spec.GetCommand(command.Name)
			if cmd == nil {
				fmt.Fprintln(os.Stderr, command.Name)
				continue
			}
			command.AddMissingOptions(cmd)
		}
	}
}

func (o *ToC) AddMissingUsages(spec *KubectlSpec) {
	for _, cats := range o.Categories {
		for _, command := range cats.Commands {
			cmd := spec.GetCommand(command.Name)
			if cmd == nil {
				fmt.Fprintln(os.Stderr, command.Name)
				continue
			}
			command.AddMissingUsage(cmd)
		}
	}
}

func (o *ToCCommand) AddMissingOptions(spec *Command) {
	optionsInToC := map[string]struct{}{}
	for _, opt := range o.GetAllOptionNames() {
		optionsInToC[opt] = struct{}{}
	}

	newGroup := OptionsGroup{
		Name: "Other options",
	}

	for _, opt := range spec.GetAllOptionNames() {
		if _, found := optionsInToC[opt]; !found {
			fmt.Fprintf(os.Stderr, "option %s not found in %s\n", opt, o.Name)
			newGroup.Options = append(newGroup.Options, ToCOption{
				Name: opt,
			})
		}
	}

	if len(newGroup.Options) > 0 {
		o.OptionsGroups = append(o.OptionsGroups, newGroup)
	}
}

func (o *ToCCommand) AddMissingUsage(spec *Command) {
	if spec == nil {
		return
	}
	o.Usage = spec.Usage
}
