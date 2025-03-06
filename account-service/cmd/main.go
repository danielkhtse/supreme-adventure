package main

import (
	"log"
	"net/http"

	"github.com/danielkhtse/supreme-adventure/account-service/internal/api"
)

func main() {
	r := api.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", r))
}
