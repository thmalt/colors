//go:build ignore

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/thmalt/colors/gen/codegen"
)

func main() {
	var ctx codegen.Context
	ctx.Optimization = codegen.OptimizeSpeed
	ctx.FormatSource = true

	ctx.SetModuleByType(ctx)
	ctx.Directory = findRoot("./")

	ctx.RootPkg.Name = "colors"
	ctx.RootPkg.Path = "/"

	ctx.ConvertPkg.Name = "convert"
	ctx.ConvertPkg.Path = "/convert"

	ctx.SpacePkg.Name = "space"
	ctx.SpacePkg.Path = "/space"

	ctx.InterpPkg.Name = "interp"
	ctx.InterpPkg.Path = "/interp"

	ctx.AddSpaces(codegen.Spaces[:])
	ctx.AddConvertFunc(codegen.ConvertFuncs[:]...)
	ctx.AddWhitePoint(codegen.WhitePoints[:]...)

	beg := time.Now()

	fmt.Println("Building...")
	ctx.Build()

	fmt.Println()
	fmt.Println("Spaces order:")
	for i, space := range ctx.BuildSpaces {
		fmt.Printf("%d\t%s\n", i, space.Name)
	}

	fmt.Println()

	fmt.Println("Generating convert...")
	codegen.GenerateConvertPkg(&ctx)

	fmt.Println("Generating space...")
	codegen.GenerateSpacePkg(&ctx)

	fmt.Println("Generating colors...")
	codegen.GenerateRootPkg(&ctx)

	fmt.Println("Generating interp...")
	codegen.GenerateInterpPkg(&ctx)

	end := time.Now()

	fmt.Println()
	fmt.Printf("Completed in %v.\n", end.Sub(beg))
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
