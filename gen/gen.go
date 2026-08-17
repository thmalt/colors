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
	var enableDebug = os.Getenv("DEBUG") != ""

	ctx := codegen.NewContext(codegen.Options{
		EmbedMatrix: true,

		// SeparateAfterComment: true,
	})

	if !enableDebug {
		ctx.Opts.FormatSource = true
	}

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

	ctx.MixerPkg.Name = "mixer"
	ctx.MixerPkg.Path = "/mixer"

	ctx.AddSpace(codegen.Spaces[:]...)
	ctx.AddConvertFunc(codegen.ConvertFuncs[:]...)
	ctx.AddWhitePoint(codegen.WhitePoints[:]...)

	beg := time.Now()

	fmt.Println("Building...")
	ctx.Build()

	fmt.Println()

	if enableDebug {
		fmt.Println("Spaces order:")
		for i, space := range ctx.BuiltSpaces {
			fmt.Printf("%d\t%s\n", i, space.Name)
		}
		logGraphPaths(ctx)
	}

	fmt.Println("Generating convert...")
	codegen.GenerateConvertPkg(ctx)

	fmt.Println("Generating space...")
	codegen.GenerateSpacePkg(ctx)

	fmt.Println("Generating colors...")
	codegen.GenerateRootPkg(ctx)

	fmt.Println("Generating mixer...")
	codegen.GenerateMixerPkg(ctx)

	end := time.Now()

	if !ctx.Opts.FormatSource {
		fmt.Println()
		fmt.Println("NOTICE: Source formatting is disabled")
	}

	fmt.Println()
	fmt.Printf("INFO: Spaces: %d, Conversion: %d\n", len(ctx.BuiltSpaces), ctx.TotalConversionGenerated)
	fmt.Printf("Completed in %v.\n", end.Sub(beg))
}

func logGraphPaths(ctx *codegen.Context) {
	fmt.Println()
	fmt.Println("Conversion path counts:")
	for i, s := range ctx.BuiltSpaces {
		for j := i + 1; j < len(ctx.BuiltSpaces); j++ {
			to := ctx.BuiltSpaces[j]

			allPath := ctx.Graph.FindAllPath(s, to)
			fmt.Println(s.Name, "->", to.Name, len(allPath))

			allPath = ctx.Graph.FindAllPath(to, s)
			fmt.Println(to.Name, "->", s.Name, len(allPath))

			fmt.Println()
		}
	}
	fmt.Println()
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
