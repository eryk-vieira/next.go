package runner

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"plugin"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/olekukonko/tablewriter"

	"github.com/eryk-vieira/next.go/cli/build"
)

type ServerRunner struct {
	Port   string
	Router string
}

func (runner *ServerRunner) Run() {
	dir, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	jsonRoutes, err := os.ReadFile(filepath.Join(dir, ".dist", "routes.json"))

	if err != nil {
		panic(err)
	}

	var routes []build.Route

	err = json.Unmarshal(jsonRoutes, &routes)

	if err != nil {
		panic(err)
	}

	if runner.Router == "chi" {
		fmt.Println("Running chi server on port: ", runner.Port)

		http.ListenAndServe(":"+runner.Port, runner.runChiServer(routes))
	}

	if runner.Router == "gin" {
		fmt.Println("Running Gin Gonic server on port: ", runner.Port)

		http.ListenAndServe(":"+runner.Port, runner.runGinServer(routes))
	}
}

func (*ServerRunner) runChiServer(routes []build.Route) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	for _, route := range routes {
		if route.RouteType == "handler" {
			plugin, err := plugin.Open(route.FilePath)

			if err != nil {
				panic(err)
			}

			for i := 0; i < len(build.AllowedMethods); i++ {
				symbol, err := plugin.Lookup(build.AllowedMethods[i])

				if err != nil {
					continue
				}

				handler, ok := symbol.(func(w http.ResponseWriter, r *http.Request))

				if !ok {
					log.Fatalf("Method %s at %s is not of type http.HandlerFunc", build.AllowedMethods[i], route.FilePath)
				}

				router.Method(build.AllowedMethods[i], route.Pattern, http.HandlerFunc(handler))
			}
		}

		if route.RouteType == "html" {
			html, err := template.ParseFiles(route.FilePath)

			if err != nil {
				panic(err)
			}

			router.Method("GET", route.Pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Header().Add("Content-Type", "text/html")

				html.Execute(w, nil)
			}))
		}
	}

	printRoutes(router)

	return router
}

func printRoutes(r chi.Routes) {
	data := [][]string{}
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		data = append(data, []string{method, route})
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		fmt.Printf("Logging err: %s\n", err.Error())
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Method", "Route"})

	for _, v := range data {
		table.Append(v)
	}

	go table.Render()
}

func (s *ServerRunner) runGinServer(routes []build.Route) http.Handler {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	for _, route := range routes {
		if route.RouteType == "handler" {
			plugin, err := plugin.Open(route.FilePath)

			if err != nil {
				panic(err)
			}

			for i := 0; i < len(build.AllowedMethods); i++ {
				symbol, err := plugin.Lookup(build.AllowedMethods[i])

				if err != nil {
					continue
				}

				handler, ok := symbol.(func(c *gin.Context))

				if !ok {
					log.Fatalf("Method %s at %s is not of type gin.HandlerFunc", build.AllowedMethods[i], route.FilePath)
				}

				router.Handle(build.AllowedMethods[i], route.Pattern, gin.HandlerFunc(handler))
			}
		}

		if route.RouteType == "html" {
			html, err := template.ParseFiles(route.FilePath)

			if err != nil {
				panic(err)
			}

			router.GET(route.Pattern, func(c *gin.Context) {
				c.Header("Content-Type", "text/html")
				html.Execute(c.Writer, nil)
			})
		}
	}

	printRoutesGin(router)

	return router
}

func printRoutesGin(r *gin.Engine) {
	routes := r.Routes()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Method", "Route"})

	for _, route := range routes {
		table.Append([]string{route.Method, route.Path})
	}

	go table.Render()
}
