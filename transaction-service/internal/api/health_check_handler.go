package api

import (
	"net/http"

	"github.com/danielkhtse/supreme-adventure/common/response"
)

// swagger:route GET /health Health healthCheck
// Check API health status
// responses:
//
//	200: description: API is healthy
//	500: description: Internal server error
func (s *Server) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response.SendSuccess[string](w, response.StatusOK, nil)
}
