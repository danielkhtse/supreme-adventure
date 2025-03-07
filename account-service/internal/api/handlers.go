package api

import (
	"net/http"

	"github.com/danielkhtse/supreme-adventure/common/response"
)

// HealthCheckHandler handles the health check endpoint and returns a 200 status
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response.SendSuccess[string](w, response.StatusOK, "Hello, World!", nil)
}
