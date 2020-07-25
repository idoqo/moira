package routesgen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
)

type apiRoute struct {
	Path    string   `json:"path"`
	Methods []string `json:"methods"`
}

type apiRouteList struct {
	List []apiRoute `json:"list"`
}

func GenerateRoutes(router chi.Router) {
	f, err := os.Create("apiroutes.json")
	if err != nil {
		fmt.Printf("error while creating file: %v", err)
	}
	defer f.Close()

	routes := make(map[string][]string)
	routeList := apiRouteList{}

	walker := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)

		routes[route] = append(routes[route], method)

		return nil
	}

	if err := chi.Walk(router, walker); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	for k, v := range routes {
		r := apiRoute{
			Path:    k,
			Methods: v,
		}
		routeList.List = append(routeList.List, r)
		fmt.Println(r.Path)
	}
	b, err := json.MarshalIndent(routeList, "", "\t")
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
	ioutil.WriteFile("apiroutes.json", b, 0644)

}
