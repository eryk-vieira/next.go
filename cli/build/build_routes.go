package build

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/eryk-vieira/next.go/cli/parsers"
	"github.com/eryk-vieira/next.go/cli/types"
)

var AllowedMethods = [...]string{
	http.MethodGet,
	http.MethodPut,
	http.MethodDelete,
	http.MethodPost,
	http.MethodPatch,
	http.MethodPut,
}

type file struct {
	Name         string
	Path         string
	RelativePath string
}

type Route struct {
	Method      string `json:"method"`
	Pattern     string `json:"pattern"`
	RouteType   string `json:"route_type"`
	FilePath    string `json:"file_path"`
	Handler     string `json:"handler"`
	PackageName string
}

type routerBuilder struct {
	Settings *types.Settings
}

func (builder *routerBuilder) Build() []Route {
	path := filepath.Join(builder.Settings.RootFolder, "src", "pages")

	routes, err := builder.scanDirectory(path)

	if err != nil {
		panic(err)
	}

	return builder.registerRoutes(routes)
}

func (*routerBuilder) scanDirectory(base_path string) ([]file, error) {
	var files []file

	err := filepath.Walk(base_path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			files = append(files, file{
				Name:         info.Name(),
				Path:         path,
				RelativePath: strings.TrimSuffix(strings.TrimPrefix(path, base_path), info.Name()),
			})
		}

		return nil
	})

	if err != nil {
		return []file{}, err
	}

	return files, nil
}

func (builder *routerBuilder) registerRoutes(files []file) []Route {
	os.RemoveAll(".dist")

	err := os.Mkdir(".dist", fs.ModePerm)

	re := regexp.MustCompile(`_(\w+)`)

	if err != nil {
		panic(err)
	}

	var routes []Route = []Route{}

	for _, file := range files {
		handlerName := builder.Settings.HTTP.HandlerName + ".go"

		if file.Name == handlerName {
			parser := parsers.HandlerParser{}

			signature, err := parser.Parse(file.Path)

			if err != nil {
				panic(err)
			}

			for _, f := range signature.Functions {
				var knownMethod = false

				for i := 0; i < len(AllowedMethods); i++ {
					if AllowedMethods[i] == f.FunctionName {
						knownMethod = true
						break
					}
				}

				if knownMethod {
					routes = append(routes, Route{
						Method:      f.FunctionName,
						Pattern:     re.ReplaceAllString(file.RelativePath, "{$1}"),
						RouteType:   "handler",
						FilePath:    filepath.Join(builder.Settings.Package, file.Path),
						PackageName: signature.PackageName,
					})
				}
			}
		}
	}

	file, err := os.Create(filepath.Join(".dist", "routes.json"))

	if err != nil {
		panic(err)
	}

	defer file.Close()

	jsonData, err := json.Marshal(&routes)

	if err != nil {
		panic(err)
	}

	file.Write(jsonData)
	file.Close()

	return routes
}
