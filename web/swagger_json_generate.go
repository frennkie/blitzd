package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Reads all .json files in the specified srcPath folder
// and encodes them as strings literals in dstPath
func main() {
	srcPath := filepath.FromSlash("../api/swagger/v1/blitzd.swagger.json")
	dstPath := filepath.FromSlash("../pkg/api/v1/swagger.pb.go")

	out, _ := os.Create(dstPath)

	fmt.Println("writing", dstPath)

	_, _ = out.Write([]byte("package v1\n\nconst (\n"))

	_, _ = out.Write([]byte("	Swagger = `"))
	f, _ := os.Open(srcPath)
	_, _ = io.Copy(out, f)
	_, _ = out.Write([]byte("`\n"))

	_, _ = out.Write([]byte(")\n"))
}
