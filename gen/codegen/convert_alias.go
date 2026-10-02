package codegen

import "github.com/thmalt/colors/gen/codegen/writer"

func genConvertPkgAliasFunc(ctx *Context, w *writer.GoWriter, alias Pair, target Pair, params, result []string) int {
	if _, ok := ctx.impls[alias]; ok {
		return 0
	}

	ctx.impls[alias] = struct{}{}

	aliasFn := alias.FuncName()
	targetFn := target.FuncName()

	w.Separate()
	w.Comment(aliasFn, " is an alias for ", '[', targetFn, ']')
	w.Func(aliasFn)
	w.FuncParams(joinIdentsWithType(FloatType, params...))
	if ContainsAny(params, result) {
		w.FuncResults(joinRepeatN(FloatType, len(result)))
	} else {
		w.FuncResults(joinIdentsWithType(FloatType, result...))
	}
	w.FuncBody()
	w.ReturnInline()
	w.WriteCallln(targetFn, params...)
	w.End()

	return 1
}
