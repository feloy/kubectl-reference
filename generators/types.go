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

type KubectlSpec struct {
	TopLevelCommandGroups []TopLevelCommands `yaml:",omitempty"`
}

func (o *KubectlSpec) GetCommand(name string) *Command {
	for _, tlCommands := range o.TopLevelCommandGroups {
		for _, command := range tlCommands.Commands {
			if command.MainCommand.Name == name {
				return command.MainCommand
			}
			for _, sub := range command.SubCommands {
				if sub.Path+"/"+sub.Name == name {
					return sub
				}
			}
		}
	}
	return nil
}

func (o *KubectlSpec) GetAllCommandNames() (commands []string) {
	for _, tlCommands := range o.TopLevelCommandGroups {
		for _, command := range tlCommands.Commands {
			commands = append(commands, command.MainCommand.Name)
			for _, sub := range command.SubCommands {
				commands = append(commands, sub.Path+"/"+sub.Name)
			}
		}
	}
	return
}

type TopLevelCommands struct {
	Group    string            `yaml:",omitempty"`
	Commands []TopLevelCommand `yaml:",omitempty"`
}
type TopLevelCommand struct {
	MainCommand *Command `yaml:",omitempty"`
	SubCommands Commands `yaml:",omitempty"`
}

type Options []*Option
type Option struct {
	Name         string `yaml:",omitempty"`
	Shorthand    string `yaml:",omitempty"`
	DefaultValue string `yaml:"default_value,omitempty"`
	Usage        string `yaml:",omitempty"`
	Type         string `yaml:",omitempty"`
}

type Example struct {
	Title   string `yaml:",omitempty"`
	Content string `yaml:",omitempty"`
}

type Commands []*Command
type Command struct {
	Name             string    `yaml:",omitempty"` // done
	Path             string    `yaml:",omitempty"`
	Synopsis         string    `yaml:",omitempty"` // done -> refpurpose
	Description      string    `yaml:",omitempty"` // done -> refsection{Description}
	Options          Options   `yaml:",omitempty"`
	InheritedOptions Options   `yaml:"inherited_options,omitempty"`
	Examples         []Example `yaml:",omitempty"`
	SeeAlso          []string  `yaml:"see_also,omitempty"` // not used
	Usage            string    `yaml:",omitempty"`         // not used
}

type Manifest struct {
	Docs      []Doc  `json:"docs,omitempty"`
	Title     string `json:"title,omitempty"`
	Copyright string `json:"copyright,omitempty"`
}

type Doc struct {
	Filename string `json:"filename,omitempty"`
}

func (o *Command) GetAllOptionNames() (options []string) {
	for _, opt := range o.Options {
		options = append(options, opt.Name)
	}
	return
}

func (o *Command) GetAllInheritedOptionNames() (options []string) {
	for _, opt := range o.InheritedOptions {
		options = append(options, opt.Name)
	}
	return
}

func (o *Command) GetOption(name string) *Option {
	for _, opt := range o.Options {
		if opt.Name == name {
			return opt
		}
	}
	return nil
}

func (o *Command) GetInheritedOption(name string) *Option {
	for _, opt := range o.InheritedOptions {
		if opt.Name == name {
			return opt
		}
	}
	return nil
}
