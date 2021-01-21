/*
 * ACME Reviews - PSIDI
 *
 * Swagger proposed server for the Review Product infrastructure at ACME .Inc
 *
 * API version: 0.1
 * Contact: 1171071@isep.ipp.pt
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nvellon/hal"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var tpl []string
var met []string

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		tpl = append(tpl, path)
		mets, _ := route.GetMethods()
		met = append(met, mets[0])
		return nil
	})
	return router
}

type Root struct {
}

func (r Root) GetMap() hal.Entry {
	return hal.Entry{}
}

func Index(w http.ResponseWriter, r *http.Request) {

	root := hal.NewResource(Root{}, "https://"+r.Host+tpl[0])

	// Customers
	customerURL := "https://psidi-customers.herokuapp.com"
	resp, err := http.Get(customerURL + "/v1/routes")
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	type Links struct {
		S string `json:"href"`
	}

	type Hyper struct {
		Links map[string]Links `json:"_links"`
	}

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}

		hyper := &Hyper{}
		json.Unmarshal(bodyBytes, &hyper)

		customers := hal.NewResource(Root{}, customerURL)

		for key, value := range hyper.Links {
			if key != "self" {
				customers.AddNewLink(hal.Relation(key), customerURL+value.S)
			}
		}

		root.Embed("customers", customers)

	}

	hypermedia, _ := json.MarshalIndent(root, "", "  ")
	w.Write([]byte(hypermedia))
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/v1/api",
		Index,
	},

	Route{
		"GetCustomer",
		strings.ToUpper("Get"),
		"/v1/customers/{customerId}",
		Index,
	},

	Route{
		"RegisterCustomer",
		strings.ToUpper("Post"),
		"/v1/customers",
		Index,
	},
	/*
		Route{
			"UpdateCustomer",
			strings.ToUpper("Patch"),
			"/v1/customers/{customerId}",
			UpdateCustomer,
		},
	*/
}
