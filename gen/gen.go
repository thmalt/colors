//go:build ignore

package main

import (
	"os"
	"path/filepath"

	"github.com/thmalt/colors/gen/codegen"
)

func main() {
	var ctx codegen.Context
	ctx.SplitFile = true

	ctx.SetModuleByType(ctx)
	ctx.Directory = findRoot("./")

	ctx.RootPkg.Name = "colors"
	ctx.RootPkg.Path = "/"

	ctx.ConvertPkg.Name = "convert"
	ctx.ConvertPkg.Path = "/convert"

	ctx.SpacePkg.Name = "space"
	ctx.SpacePkg.Path = "/space"

	ctx.AddSpaces(codegen.Spaces[:])
	ctx.AddConvertFunc(codegen.ConvertFuncs[:]...)
	ctx.AddWhitePoint(codegen.WhitePoints[:]...)
	ctx.Build()

	codegen.GenerateConvertPkg(&ctx)
	codegen.GenerateRootPkg(&ctx)
	codegen.GenerateSpacePkg(&ctx)
}

func findRoot(dir string) string {
	if _, err := os.Stat("go.mod"); err == nil {
		return dir
	}

	var (
		prev string
		path string = dir
	)

	if !filepath.IsAbs(path) {
		if p, err := filepath.Abs(path); err == nil {
			path = p
		}
	}

	for {
		path = filepath.Dir(path)
		if prev == path {
			break
		}

		if _, err := os.Stat(filepath.Join(path, "go.mod")); err == nil {
			return path
		}
		prev = path
	}

	return dir
}
