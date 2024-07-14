package build

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/google/uuid"

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
	Method    string `json:"method"`
	Pattern   string `json:"pattern"`
	RouteType string `json:"route_type"`
	FilePath  string `json:"file_path"`
}

type routerBuilder struct {
	Settings *types.Settings
}

func (builder *routerBuilder) Build() []Route {
	routes, err := builder.scanDirectory(filepath.Join(builder.Settings.RootFolder, "pages"))

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

	osRootDir, _ := os.Getwd()
	err := os.Mkdir(".dist", fs.ModePerm)

	re := regexp.MustCompile(`_(\w+)`)

	if err != nil {
		panic(err)
	}

	var routes []Route = []Route{}

	var wg sync.WaitGroup

	for _, file := range files {
		if file.Name == "index.html" {
			routes = append(routes, Route{
				Method:    http.MethodGet,
				Pattern:   re.ReplaceAllString(file.RelativePath, "{$1}"),
				RouteType: "html",
				FilePath:  filepath.Join(osRootDir, file.Path),
			})
		}

		handlerName := builder.Settings.HTTP.HandlerName + ".go"

		if file.Name == handlerName {
			wg.Add(1)
			pluginId := uuid.NewString()
			pluginDestination := filepath.Join(osRootDir, ".dist", pluginId+".so")

			pluginBuilder := pluginBuilder{
				Dir:       builder.Settings.RootFolder,
				DebugMode: true,
			}

			go func() {
				err := pluginBuilder.Build(filepath.Join("pages", file.RelativePath, file.Name), pluginDestination)

				if err != nil {
					panic(err)
				}

				defer wg.Done()
			}()

			routes = append(routes, Route{
				Method:    http.MethodGet,
				Pattern:   re.ReplaceAllString(file.RelativePath, "{$1}"),
				RouteType: "handler",
				FilePath:  pluginDestination,
			})
		}
	}

	wg.Wait()

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
