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
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jinzhu/copier"
)

func (o *Command) AsDocbook(w io.Writer, config *ToCCommand) {
	refname := o.Name
	if len(o.Path) > 0 {
		refname = strings.Replace(o.Path, "/", " ", 1) + " " + refname
	}
	refpurpose := o.Synopsis
	fmt.Fprintf(w, `    <refentry>
      <refnamediv>
        <refname>%s</refname>

        <refpurpose>%s</refpurpose>
      </refnamediv>

      <refsynopsisdiv><title>Usage</title>

        <cmdsynopsis>
          <command>kubectl %s</command>
`, refname, refpurpose, refname)

	for _, arg := range config.Args {
		if !arg.End {
			arg.AsDocbook(w)
		}
	}
	fmt.Fprint(w, "          <sbr/>\n")

	for _, group := range config.OptionsGroups {
		for _, tocOption := range group.Options {
			option := o.GetOption(tocOption.Name)
			if option == nil {
				option = o.GetInheritedOption(tocOption.Name)
				if option == nil {
					fmt.Printf("option %s of command %s not found\n", tocOption.Name, o.Name)
					os.Exit(1)
				}
			}
			option.AsDocbook(w, &tocOption)
		}
		fmt.Fprint(w, "          <sbr/>\n")
	}

	for _, arg := range config.Args {
		if arg.End {
			arg.AsDocbook(w)
		}
	}

	fmt.Fprint(w, `        </cmdsynopsis>
      </refsynopsisdiv>
`)

	if *ShowUsage {
		// Description
		fmt.Fprintf(w, `      <refsection><title>Original Usage</title>
        <programlisting>%s</programlisting></refsection>
`, escapeXml(o.Usage))

	}

	// Description
	fmt.Fprint(w, `      <refsection>
        <title>Description</title>
`)

	desc := o.Description
	paras := strings.Split(desc, "\n\n")
	for _, para := range paras {
		fmt.Fprintf(w, "          <para>%s</para>\n", escapeXml(para))
	}
	fmt.Fprint(w, `      </refsection>
`)

	// Options
	if len(config.OptionsGroups) > 0 {
		fmt.Fprint(w, `      <refsection>
        <title>Options</title>
`)

		for _, group := range config.OptionsGroups {
			if len(group.Name) > 0 {
				fmt.Fprintf(w, "        <bridgehead renderas=\"sect3\">%s</bridgehead>\n", group.Name)
			}
			fmt.Fprintf(w, "        <variablelist>\n")
			for _, tocOption := range group.Options {
				option := o.GetOption(tocOption.Name)
				if option == nil {
					option = o.GetInheritedOption(tocOption.Name)
					if option == nil {
						fmt.Printf("option %s of command %s not found\n", tocOption.Name, o.Name)
						os.Exit(1)
					}
				}
				option.AsDocbookDetails(w, &tocOption)
			}
			fmt.Fprintf(w, "        </variablelist>\n")
		}

		fmt.Fprint(w, `      </refsection>
`)
	}

	// Examples
	if len(o.Example) > 0 {
		fmt.Fprint(w, `      <refsection>
        <title>Examples</title>
`)

		fmt.Fprint(w, "          <programlisting>\n")
		examples := o.Example
		lines := strings.Split(examples, "\n")
		for _, line := range lines {
			fmt.Fprintf(w, "%s\n", escapeXml(line))
		}

		fmt.Fprint(w, "          </programlisting>\n")
		fmt.Fprint(w, `      </refsection>
`)
	}

	fmt.Fprint(w, `    </refentry>
`)
}

func (o *Arg) AsDocbook(w io.Writer) {
	choice := "plain"
	if o.Choice != nil {
		choice = *o.Choice
	}
	rep := "norepeat"
	if o.Rep != nil {
		rep = *o.Rep
	}
	fmt.Fprintf(w, "          <arg choice=\"%s\" rep=\"%s\"><replaceable>%s</replaceable></arg>\n", choice, rep, string(o.Name))
}

func (op *Option) AsDocbook(w io.Writer, config *ToCOption) {
	var o Option
	copier.Copy(&o, op)

	if config.Type != nil {
		o.Type = *config.Type
	}
	if config.Usage != nil {
		o.Usage = *config.Usage
	}
	if config.Shorthand != nil {
		o.Shorthand = *config.Shorthand
	}
	if config.Default != nil {
		o.DefaultValue = *config.Default
	}
	choice := "opt"
	if config.Required {
		choice = "plain"
	}

	optionName := "--" + o.Name + "="
	if len(o.Shorthand) > 0 {
		optionName = "-" + o.Shorthand + " "
	}

	switch o.Type {
	case "bool", "tristate":
		var value string
		if len(o.Shorthand) > 0 && o.DefaultValue == "false" {
			value = "-" + o.Shorthand
		} else {
			value = "--" + o.Name
			if o.DefaultValue == "true" {
				value += "=false"
			}
		}
		fmt.Fprintf(w, "          <arg choice=\"%s\">%s</arg>\n", choice, value)

	case "string":
		value := optionName + "<replaceable>value</replaceable>"
		fmt.Fprintf(w, "          <arg choice=\"%s\">%s</arg>\n", choice, value)

	case "int32", "int64", "int", "duration":
		value := optionName + "<replaceable>value</replaceable>"
		fmt.Fprintf(w, "          <arg choice=\"%s\">%s</arg>\n", choice, value)

	case "stringArray":
		if config.Required {
			choice = "req"
		}
		value := optionName + "<replaceable>value</replaceable>"
		fmt.Fprintf(w, "          <arg rep=\"repeat\" choice=\"plain\"><arg choice=\"%s\">%s</arg></arg>\n", choice, value)

	case "stringSlice":
		value := optionName + "<replaceable>value1</replaceable><arg rep=\"repeat\" choice=\"plain\"><arg choice=\"opt\">,<replaceable>valueN</replaceable></arg></arg>"
		fmt.Fprintf(w, "          <arg choice=\"plain\"><arg choice=\"%s\">%s</arg></arg>\n", choice, value)

	case "mapStringString":
		value := optionName + "<replaceable>value</replaceable>"
		fmt.Fprintf(w, "          <arg choice=\"%s\">%s</arg>\n", choice, value)

	default:
		fmt.Fprintf(w, "          <arg>--%s</arg>\n", o.Name)
	}
}

func (op *Option) AsDocbookDetails(w io.Writer, config *ToCOption) {

	var o Option
	copier.Copy(&o, op)

	if config.Type != nil {
		o.Type = *config.Type
	}
	if config.Usage != nil {
		o.Usage = *config.Usage
	}
	if config.Shorthand != nil {
		o.Shorthand = *config.Shorthand
	}
	if config.Default != nil {
		o.DefaultValue = *config.Default
	}

	fmt.Fprintf(w, "          <varlistentry>\n")
	fmt.Fprintf(w, "            <term>")

	value := "--" + o.Name
	if len(o.Shorthand) > 0 {
		value = "-" + o.Shorthand + " | " + value
	}
	fmt.Fprint(w, value)

	var def string
	if len(o.DefaultValue) > 0 && o.DefaultValue != "[]" {
		def = fmt.Sprintf(", defaults to %s", o.DefaultValue)
	}
	fmt.Fprintf(w, " (%s%s)</term>\n", o.Type, def)
	fmt.Fprintf(w, "            <listitem><para>%s</para></listitem>\n", escapeXml(o.Usage))
	fmt.Fprintf(w, "          </varlistentry>\n")
}

func escapeXml(s string) string {
	var b []byte
	buf := bytes.NewBuffer(b)
	xml.EscapeText(buf, []byte(s))
	return buf.String()
}
