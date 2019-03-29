package main

import (
	"flag"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
)

type Config struct {
	Vendor string
	Routes []Route `yaml:"routes"`
}

type Route struct {
	Path     string `yaml:"path"`
	Function string `yaml:"function"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Creates a handler function for the router that will
// invoke the aws lambda
func createHandler(route string, function string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		fmt.Println("ROUTE: ", route)
		fmt.Println("FUNCTION: ", function)
		fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
	}
}

func main() {
	var f = flag.String("config", "./example.yml", "Config file to use")

	flag.Parse()

	dat, err := ioutil.ReadFile(*f)
	check(err)

	config := Config{}

	marshalerr := yaml.Unmarshal([]byte(dat), &config)

	if marshalerr != nil {
		log.Fatalf("cannot unmarshal data: %v", marshalerr)
	}

	fmt.Println(config)

	router := httprouter.New()
	for i := 0; i < len(config.Routes); i++ {
		handler := createHandler(config.Routes[i].Path, config.Routes[i].Function)
		router.GET(config.Routes[i].Path, handler)
	}

	log.Fatal(http.ListenAndServe(":8080", router))
}
