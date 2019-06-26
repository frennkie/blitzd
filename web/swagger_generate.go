package main

import (
	"github.com/shurcooL/vfsgen"
	"log"
	"net/http"
	"path/filepath"
)

var SwaggerFs http.FileSystem = http.Dir(filepath.FromSlash("../third_party/swagger-ui"))

func main() {
	err := vfsgen.Generate(
		SwaggerFs,
		vfsgen.Options{
			Filename:     "./swagger/swagger_vfsdata.go",
			PackageName:  "swagger",
			VariableName: "SwaggerFs",
		})
	if err != nil {
		log.Fatalln(err)
	}
}
