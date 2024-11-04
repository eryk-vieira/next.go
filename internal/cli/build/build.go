package build

import "github.com/eryk-vieira/next.go/internal/cli/types"

type build struct {
	Settings *types.Settings
}

func NewBuilder(settings *types.Settings) *build {
	return &build{
		Settings: settings,
	}
}

func (b *build) Build() ([]Route, []Errors) {
	builder := routerBuilder{
		Settings: b.Settings,
	}

	routes, errorList := builder.Build()

	if len(errorList) > 0 {
		return []Route{}, errorList
	}

	serverBuilder := serverBuilder{
		Settings: b.Settings,
	}

	serverBuilder.Build(routes)

	return routes, errorList
}
