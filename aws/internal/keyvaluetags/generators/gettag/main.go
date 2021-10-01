//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"go/format"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"

	tftags "github.com/hashicorp/terraform-provider-aws/aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

const filename = `get_tag_gen.go`

var serviceNames = []string{
	"autoscaling",
	"batch",
	"dynamodb",
	"ec2",
	"ecs",
	"route53resolver",
}

type TemplateData struct {
	ServiceNames []string
}

func main() {
	// Always sort to reduce any potential generation churn
	sort.Strings(serviceNames)

	templateData := TemplateData{
		ServiceNames: serviceNames,
	}
	templateFuncMap := template.FuncMap{
		"ClientType":                        tftags.ServiceClientType,
		"ListTagsFunction":                  tftags.ServiceListTagsFunction,
		"ListTagsInputFilterIdentifierName": tftags.ServiceListTagsInputFilterIdentifierName,
		"ListTagsOutputTagsField":           tftags.ServiceListTagsOutputTagsField,
		"TagPackage":                        tftags.ServiceTagPackage,
		"TagResourceTypeField":              tftags.ServiceTagResourceTypeField,
		"TagTypeAdditionalBoolFields":       tftags.ServiceTagTypeAdditionalBoolFields,
		"TagTypeIdentifierField":            tftags.ServiceTagTypeIdentifierField,
		"Title":                             strings.Title,
	}

	tmpl, err := template.New("gettag").Funcs(templateFuncMap).Parse(templateBody)

	if err != nil {
		log.Fatalf("error parsing template: %s", err)
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, templateData)

	if err != nil {
		log.Fatalf("error executing template: %s", err)
	}

	generatedFileContents, err := format.Source(buffer.Bytes())

	if err != nil {
		log.Fatalf("error formatting generated file: %s", err)
	}

	f, err := os.Create(filename)

	if err != nil {
		log.Fatalf("error creating file (%s): %s", filename, err)
	}

	defer f.Close()

	_, err = f.Write(generatedFileContents)

	if err != nil {
		log.Fatalf("error writing to file (%s): %s", filename, err)
	}
}

var templateBody = `
// Code generated by generators/gettag/main.go; DO NOT EDIT.

package keyvaluetags

import (
	"github.com/aws/aws-sdk-go/aws"
{{- range .ServiceNames }}
	"github.com/aws/aws-sdk-go/service/{{ . }}"
{{- end }}
    "github.com/hashicorp/terraform-provider-aws/aws/internal/tfresource"
)

{{- range .ServiceNames }}

// {{ . | Title }}GetTag fetches an individual {{ . }} service tag for a resource.
// Returns whether the key value and any errors. A NotFoundError is used to signal that no value was found.
// This function will optimise the handling over {{ . | Title }}ListTags, if possible.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
{{- if or ( . | TagTypeIdentifierField ) ( . | TagTypeAdditionalBoolFields) }}
func {{ . | Title }}GetTag(conn {{ . | ClientType }}, identifier string{{ if . | TagResourceTypeField }}, resourceType string{{ end }}, key string) (*TagData, error) {
{{- else }}
func {{ . | Title }}GetTag(conn {{ . | ClientType }}, identifier string{{ if . | TagResourceTypeField }}, resourceType string{{ end }}, key string) (*string, error) {
{{- end }}
	{{- if . | ListTagsInputFilterIdentifierName }}
	input := &{{ . | TagPackage  }}.{{ . | ListTagsFunction }}Input{
		Filters: []*{{ . | TagPackage  }}.Filter{
			{
				Name:   aws.String("{{ . | ListTagsInputFilterIdentifierName }}"),
				Values: []*string{aws.String(identifier)},
			},
			{
				Name:   aws.String("key"),
				Values: []*string{aws.String(key)},
			},
		},
	}

	output, err := conn.{{ . | ListTagsFunction }}(input)

	if err != nil {
		return nil, err
	}

	listTags := {{ . | Title }}KeyValueTags(output.{{ . | ListTagsOutputTagsField }}{{ if . | TagTypeIdentifierField }}, identifier{{ if . | TagResourceTypeField }}, resourceType{{ end }}{{ end }})
	{{- else }}
	listTags, err := {{ . | Title }}ListTags(conn, identifier{{ if . | TagResourceTypeField }}, resourceType{{ end }})

	if err != nil {
		return nil, err
	}
	{{- end }}

	if !listTags.KeyExists(key) {
		return nil, tfresource.NewEmptyResultError(nil)
	}

	{{ if or ( . | TagTypeIdentifierField ) ( . | TagTypeAdditionalBoolFields) }}
	return listTags.KeyTagData(key), nil
	{{- else }}
	return listTags.KeyValue(key), nil
	{{- end }}
}
{{- end }}
`
