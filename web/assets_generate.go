package main

import (
	"github.com/shurcooL/vfsgen"
	"log"
	"net/http"
	"path/filepath"
)

var Assets http.FileSystem = http.Dir(filepath.FromSlash("assets/src"))

func main() {
	err := vfsgen.Generate(
		Assets,
		vfsgen.Options{
			Filename:     "./assets/assets_vfsdata.go",
			PackageName:  "assets",
			VariableName: "Assets",
		})
	if err != nil {
		log.Fatalln(err)
	}
}
