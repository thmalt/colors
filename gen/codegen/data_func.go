package codegen

import (
	"github.com/thmalt/colors/gen/codegen/model"
)

func channelName(name, ident, symbol, displayName string) model.Channel {
	return model.Channel{
		Name:        name,
		Ident:       ident,
		Symbol:      symbol,
		DisplayName: displayName,
	}
}

func extendChannel(ch model.Channel, funcs ...func(ch *model.Channel)) model.Channel {
	for _, fn := range funcs {
		fn(&ch)
	}
	return ch
}

func numberChannel(min, max float64, prec int) func(ch *model.Channel) {
	return func(ch *model.Channel) {
		ch.Min = min
		ch.Max = max
		ch.Precision = prec
		ch.Unit = model.UnitNumber
	}
}

func percentChannel(min, max float64, prec int) func(ch *model.Channel) {
	return func(ch *model.Channel) {
		ch.Min = min
		ch.Max = max
		ch.Precision = prec
		ch.Unit = model.UnitPercent
	}
}

func unrestrictedChannel(ch *model.Channel) {
	ch.Unrestricted = true
}

func circularChannel(ch *model.Channel) {
	ch.Circular = true
}

func degreeChannel(ch *model.Channel) {
	ch.Min = 0
	ch.Max = 360
	ch.Precision = 4
	ch.Circular = true
	ch.Unit = model.UnitDegree
}

func precisionChannel(prec int) func(ch *model.Channel) {
	return func(ch *model.Channel) {
		ch.Precision = prec
	}
}
