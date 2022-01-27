package proxy

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"html/template"
	"strings"
)

// v1 根据 接口实现 proxy
// TODO 根据 receiver 实现的方法直接实现proxy

var proxyTemp = `
package {{.Package}}

type {{ .StructName }}Proxy struct {
	child *{{ .StructName }}
}

func New{{ .StructName }}Proxy(child *{{ .StructName }}) *{{ .StructName }}Proxy {
	return &{{ .StructName }}Proxy{child: child}
}

{{ range .Methods }}
func (p *{{$.StructName}}Proxy) {{ .Name }} ({{ .Params }}) ({{ .Results }}) {
	// before 这里可能会有一些统计的逻辑
	start := time.Now()

	{{ .ResultNames }} = p.child.{{ .Name }}({{ .ParamNames }})

	// after 这里可能也有一些监控统计的逻辑
	log.Printf("user login cost time: %s", time.Now().Sub(start))

	return {{ .ResultNames }}
}
{{ end }}
`

type GoFile struct {
	Package    string    `json:"package,omitempty"`
	StructName string    `json:"struct_name,omitempty"`
	Methods    []*Method `json:"methods,omitempty"`
}

type Method struct {
	Name        string `json:"name,omitempty"`
	Params      string `json:"params,omitempty"`
	ParamNames  string `json:"param_names,omitempty"`
	Results     string `json:"results,omitempty"`
	ResultNames string `json:"result_names,omitempty"`
}

func generate(file string) []byte {
	// 1.从 go file 中获取参数信息
	data := &GoFile{}

	fileset := token.NewFileSet()
	f, err := parser.ParseFile(fileset, file, nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	// ast.Print(fileset, f)
	data.Package = f.Name.Name

	cmap := ast.NewCommentMap(fileset, f, f.Comments)
	fmt.Println(cmap.String())
	for node, group := range cmap {
		// 从注释 @PROXY 接口名，获取接口名称
		interfaceName := getProxyInterfaceName(group)
		if interfaceName == "" {
			// 只获取这个 @PROXY 注释的对象，其他都跳过
			continue
		}

		// 这是一个 *ast.GenDecl 类型
		structName := node.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Name.Name
		data.StructName = structName

		// 检索接口需要实现的方法
		obj := f.Scope.Lookup(interfaceName)
		interType := obj.Decl.(*ast.TypeSpec).Type.(*ast.InterfaceType)

		for _, method := range interType.Methods.List {
			data.Methods = append(data.Methods, getProxyMethod(method))
		}
	}

	// 2.参数信息注入模版，获取代码
	tpl, err := template.New("").Parse(proxyTemp)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, data); err != nil {
		fmt.Println(err)
		return nil
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return src
}

func getProxyMethod(f *ast.Field) *Method {
	method := &Method{}

	method.Name = f.Names[0].Name
	method.Params, method.ParamNames = getFieldList(f.Type.(*ast.FuncType).Params.List)
	method.Results, method.ResultNames = getFieldList(f.Type.(*ast.FuncType).Results.List)

	return method
}

func getFieldList(fs []*ast.Field) (string, string) {
	var (
		params     []string
		paramNames []string
	)

	for _, p := range fs {
		names := make([]string, 0, len(p.Names))
		for _, name := range p.Names {
			names = append(names, strings.TrimSpace(name.Name))
		}

		paramNames = append(paramNames, names...)
		params = append(params, fmt.Sprintf("%s %s", strings.Join(names, ","), p.Type.(*ast.Ident).Name))
	}

	return strings.Join(params, ","), strings.Join(paramNames, ",")
}

func getProxyInterfaceName(groups []*ast.CommentGroup) string {
	for _, commentGroup := range groups {
		for _, comment := range commentGroup.List {
			if strings.Contains(comment.Text, "@PROXY") {
				interfaceName := strings.TrimLeft(comment.Text, "// @PROXY ")
				return strings.TrimSpace(interfaceName)
			}
		}
	}
	return ""
}
