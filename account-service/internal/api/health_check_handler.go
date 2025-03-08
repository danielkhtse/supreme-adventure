package api

import (
	"net/http"

	"github.com/danielkhtse/supreme-adventure/common/response"
)

// @Summary Check API health status
// @Description Check if the API is healthy
// @Tags Health
// @Accept json
// @Produce json
// @Success 200
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /health-check [get]
func (s *Server) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response.SendSuccess[string](w, response.StatusOK, nil)
}
