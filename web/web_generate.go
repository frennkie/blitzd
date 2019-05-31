// +build ignore

package main

import (
	"github.com/shurcooL/vfsgen"
	"log"
	"net/http"
)

var Assets http.FileSystem = http.Dir("assets")

func main() {
	err := vfsgen.Generate(
		Assets,
		vfsgen.Options{
			Filename:     "./web_vfsdata.go",
			PackageName:  "web",
			VariableName: "Assets",
		})
	if err != nil {
		log.Fatalln(err)
	}
}
