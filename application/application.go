package main

import (
	"flag"
	fjv "github.com/morras/firebaseJwtValidator"
	api "github.com/morras/neutrinoapi"
	"log"
	"net/http"
)

// Command line arguments are purely to make it possible to run the server locally.
// For instance the web app assumes that the API endpoint is on the same hostname using
// the same port, so to be able to both run a web server and an api server, we allow the
// api server to serve the static website as well using a dev override.
func main() {
	// TODO add a flag for the stage, test / prod
	frontendLocation := flag.String("dev", "", "Relative location for front end if using a dev server.")
	portArg := flag.String("port", api.DEFAULT_PORT, "Port for the server to listen on.")
	flag.Parse()

	if *frontendLocation != "" {
		fileHandler := http.FileServer(http.Dir(*frontendLocation))
		http.Handle("/", fileHandler)
		log.Printf("Running in development mode using frontend location %s\n", *frontendLocation)
	}

	port := *portArg
	if port == "" {
		port = api.DEFAULT_PORT
	}

	requestParser := api.NewRequestParser(fjv.NewDefaultTokenValidator(api.FIREBASE_PROJECT_ID))
	newGameEndpoint := api.NewNewGameEndpoint(requestParser, nil) //TODO fix up a real data store

	http.Handle("/newGame", newGameEndpoint)

	log.Printf("Listening on port %s\n\n", port)
	http.ListenAndServe(":"+port, nil)
}
