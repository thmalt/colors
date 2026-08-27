package codegen

import (
	_ "embed"
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/thmalt/colors/gen/codegen/writer"
)

//go:embed data/css_lv4_named_colors.json
var cssLv4NamedColorData []byte

type namedColor struct {
	Name string  `json:"name"`
	RGBA []uint8 `json:"rgba"`
}

func GenerateNamedPkg(ctx *Context) {
	pkg := ctx.NamedPkg
	if pkg.Name == "" {
		return
	}

	w := newWriter(ctx)

	m := make(map[string]struct{})
	cssLv4Colors := jsonParseNamedColor(cssLv4NamedColorData)
	for _, c := range cssLv4Colors {
		lower := strings.ToLower(c.Name)
		if _, ok := m[lower]; ok {
			log.Fatalln("duplicate named color:", lower)
		}
		m[lower] = struct{}{}
	}

	if _, ok := m["transparent"]; !ok {
		m["transparent"] = struct{}{}
		cssLv4Colors = append(cssLv4Colors, namedColor{Name: "Transparent", RGBA: []uint8{0, 0, 0, 0}})
	}

	emitGoFile(ctx, pkg, w, "css_lv4", func(w *writer.GoWriter) {
		w.Import(ctx.RootPkg.Path)

		genNamedPkgNamedVar(ctx, w, cssLv4Colors)
	})
}

func genNamedPkgNamedVar(ctx *Context, w *writer.GoWriter, colors []namedColor) {
	pkgJoin := func(ident string) string {
		if ctx.NamedPkg.Name == ctx.RootPkg.Name {
			return ident
		}
		return ctx.RootPkg.Join(ident)
	}

	var (
		temp []string
		rgb  string
	)

	w.BeginGroup("var ")
	for i, color := range colors {
		temp = temp[:0]
		for i := range min(4, len(color.RGBA)) {
			temp = append(temp, strconv.Itoa(int(color.RGBA[i])))
		}

		switch len(temp) {
		case 3:
			rgb = "rgb"
		case 4:
			rgb = "rgba"
		default:
			panic("invalid named color channel count")
		}

		if i > 0 {
			w.Separate()
		}

		name := color.Name
		w.Comment(name, " is the CSS named color ", '"', strings.ToLower(name), '"')
		w.Comment('\t', rgb, "(", strings.Join(temp, ", "), ')')
		w.LineWrite(name, " = ")
		fn := "Rgb"
		if len(color.RGBA) > 3 {
			fn += "Alpha"
		}
		w.Write(pkgJoin(fn), '(')
		w.WriteJoin(temp, ", ")
		w.Writeln(')')
	}

	w.End()
}

func jsonParseNamedColor(data []byte) []namedColor {
	var colors []namedColor
	if err := json.Unmarshal(data, &colors); err != nil {
		return nil
	}
	return colors
}
