package main

import (
	"bytes"
	"log"
	"text/template"

	"github.com/golang/protobuf/protoc-gen-go/generator"
	"google.golang.org/protobuf/types/descriptorpb"
)

type netrpcPlugin struct{ *generator.Generator }

func init() {
	generator.RegisterPlugin(new(netrpcPlugin))
}

func (p *netrpcPlugin) Name() string                { return "netrpc" }
func (p *netrpcPlugin) Init(g *generator.Generator) { p.Generator = g }

func (p *netrpcPlugin) GenerateImports(file *generator.FileDescriptor) {
	if len(file.Service) > 0 {
		p.genImportCode(file)
	}
}

func (p *netrpcPlugin) Generate(file *generator.FileDescriptor) {
	for _, svc := range file.Service {
		p.genServiceCode(svc)
	}
}

func (p *netrpcPlugin) genImportCode(file *generator.FileDescriptor) {
	p.P(`import "net/rpc"`)
}

func (p *netrpcPlugin) genServiceCode(svc *descriptorpb.ServiceDescriptorProto) {
	spec := p.buildServiceSpec(svc)

	var buf bytes.Buffer
	t := template.Must(template.New("").Parse(tmplService))
	err := t.Execute(&buf, spec)
	if err != nil {
		log.Fatal(err)
	}

	p.P(buf.String())
}

type ServiceSpec struct {
	ServiceName string
	MethodList  []ServiceMethodSpec
}

type ServiceMethodSpec struct {
	MethodName      string
	InputTypeNames  string
	OutputTypeNames string
}

func (p *netrpcPlugin) buildServiceSpec(
	svc *descriptorpb.ServiceDescriptorProto,
) *ServiceSpec {
	spec := &ServiceSpec{
		ServiceName: generator.CamelCase(svc.GetName()),
	}

	for _, m := range svc.Method {
		spec.MethodList = append(spec.MethodList, ServiceMethodSpec{
			MethodName:      generator.CamelCase(m.GetName()),
			InputTypeNames:  p.TypeName(p.ObjectNamed(m.GetInputType())),
			OutputTypeNames: p.TypeName(p.ObjectNamed(m.GetOutputType())),
		})
	}

	return spec
}
