package build

import (
	"github.com/eryk-vieira/next.go/cli/nextgo/types"
)

type build struct {
	Settings *types.Settings
}

func NewBuilder(settings *types.Settings) *build {
	return &build{
		Settings: settings,
	}
}

func (b *build) Build() []Route {
	builder := routerBuilder{
		Settings: b.Settings,
	}

	routes := builder.Build()

	serverBuilder := serverBuilder{
		Settings: b.Settings,
	}
	serverBuilder.Build(routes)

	return routes
}
