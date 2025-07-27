package main

import (
	"flag"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	generateClient()
}

func generateClient() {
	protogen.Options{ParamFunc: flag.Set}.Run(func(plugin *protogen.Plugin) error {
		for _, file := range plugin.Files {
			if !file.Generate {
				continue
			}

			for _, service := range file.Services {
				err := generateClientCode(plugin, service, file)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func generateClientCode(plugin *protogen.Plugin, service *protogen.Service, file *protogen.File) error {
	clientFile := plugin.NewGeneratedFile(file.GeneratedFilenamePrefix+"_provider.go", file.GoImportPath)

	clientFile.P("package ", file.GoPackageName)
	clientFile.P()
	clientFile.P("import (")
	clientFile.P("\t\"github.com/hgyowan/go-pkg-library/envs\"")
	clientFile.P("\tpkgLibrary \"github.com/hgyowan/go-pkg-library/grpc-library/grpc\"")
	clientFile.P(")")

	clientFile.P("func ", service.GoName, "ClientProvider() ", service.GoName, "Client {")
	clientFile.P("\tconn := pkgLibrary.MustNewGRPCClient(envs.CFMAPIHost)")
	clientFile.P("\treturn New", service.GoName, "Client(conn)")
	clientFile.P("}")

	return nil
}
