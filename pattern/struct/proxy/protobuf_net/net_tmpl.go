package main

// need args
// ServiceName
//   MethodList
//     Method
//       MethodName
//       InputTypeNames
//       OutputTypeNames
const tmplService = `
{{$root := .}}

type {{.ServiceName}}Interface interface {
    {{- range $_, $m := .MethodList}}
    {{$m.MethodName}}(*{{$m.InputTypeNames}}) (*{{$m.OutputTypeNames}}, error)
    {{- end}}
}

func Register{{.ServiceName}}(
    srv *rpc.Server, x {{.ServiceName}}Interface,
) error {
    if err := srv.RegisterName("{{.ServiceName}}", x); err != nil {
        return err
    }
    return nil
}

type {{.ServiceName}}Client struct {
    *rpc.Client
}

var _ {{.ServiceName}}Interface = (*{{.ServiceName}}Client)(nil)

func Dial{{.ServiceName}}(network, address string) (
    *{{.ServiceName}}Client, error,
) {
    c, err := rpc.Dial(network, address)
    if err != nil {
        return nil, err
    }
    return &{{.ServiceName}}Client{Client: c}, nil
}

{{range .MethodList}}
func (p *{{$root.ServiceName}}Client) {{.MethodName}}(in *{{.InputTypeNames}}) (out *{{$m.OutputTypeNames}},err error) {
    err = p.Client.Call("{{$root.ServiceName}}.{{.MethodName}}", in, out)
    return
}
{{end}}
`
