package rest

import (
	"database/sql"
	"errors"
	"net/http"
	"se-challenge-payroll/internal/api"
	"se-challenge-payroll/internal/service"

	"github.com/gin-gonic/gin"
)

type PayrollHandler struct {
	service *service.Service
}

func NewPayrollHandler(service *service.Service) *PayrollHandler {
	return &PayrollHandler{service: service}
}

// PayrollRoutes define routes on the /v1/payroll
func PayrollRoutes(routes *gin.RouterGroup, s *service.Service) {
	h := NewPayrollHandler(s)

	// swagger:route GET /report report PayrollReportList
	//
	// Payroll report.
	//
	// This will show all available employee payroll report.
	//
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http
	//
	//     Responses:
	//       default: genericError
	//       200: PayrollReportResponse
	routes.GET("/report", h.getPayrollReport)
}

func (h *PayrollHandler) getPayrollReport(ctx *gin.Context) {
	reportRows, err := h.service.PayrollRepository.GetReport(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, api.NewPayrollReportResponse(reportRows))
}
