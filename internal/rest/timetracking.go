package rest

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"se-challenge-payroll/internal/entity/timetracking"
	"se-challenge-payroll/internal/service"
	"se-challenge-payroll/pkg/util"
	"strconv"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

// TimeTrackingRoutes define routes on the /v1/time_tracking
func TimeTrackingRoutes(routes *gin.RouterGroup, s *service.Service) {
	h := NewHandler(s)

	routes.POST("/upload", h.postTimeTrackerUploadCSV)
}

// upload csv file with time tracking by employee
func (h *Handler) postTimeTrackerUploadCSV(ctx *gin.Context) {
	var (
		accept = ctx.GetHeader("Accept")
	)

	if len(accept) == 0 || accept == "*/*" {
		accept = "text/csv"
	}

	if !util.ContainsString(accept, "text/csv") {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf("accept type '%v' is not supported. Content-type must be text/csv", accept))
		return
	}

	csvHeader, openErr := ctx.FormFile("file")
	if openErr != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, openErr.Error())
		return
	}

	// File name Validate
	if !validateFileName(csvHeader.Filename) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf("file name has unexpected name format: %s", csvHeader.Filename))
		return
	}

	reportID := extractReportID(csvHeader.Filename)
	if reportID == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf("file name has unexpected report ID: %s", csvHeader.Filename))
		return
	}

	//Check if the report number was uploaded before
	tr, err := h.service.TimeTrackingService.GetByFilter(ctx, timetracking.TimeTrackingFilter{TimeReportID: &reportID})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	// If found rows with the same report ID return conflict
	if len(tr) > 0 {
		ctx.AbortWithStatusJSON(http.StatusConflict,
			fmt.Sprintf("Report ID has already been uploaded. Re-uploading a file with the same report ID is not allowed: %s",
				csvHeader.Filename))
		return
	}

	file, err := csvHeader.Open()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprintf("failed to open file: %s", openErr.Error()))
		return
	}

	defer func() {
		_ = file.Close()
	}()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprintf("failed to read file: %s", openErr.Error()))
		return
	}

	timeTrackingList, errs := timetracking.CreateTimeTrackingList(records, reportID)
	if len(errs) > 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errs)
		return
	}

	err = h.service.TimeTrackingService.BulkCreate(ctx, timeTrackingList)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Header("total_row_processed", strconv.Itoa(len(timeTrackingList)))
	ctx.Header("success_row_processed", strconv.Itoa(len(timeTrackingList)))

	ctx.Status(http.StatusCreated)
}

func validateFileName(name string) bool {
	re := regexp.MustCompile(`^time-report-([1-9]|[1-9][0-9])\.csv$`)
	return re.MatchString(name)
}

func extractReportID(name string) int {
	re := regexp.MustCompile(`time-report-([1-9]|[1-9][0-9])\.csv`)
	match := re.FindStringSubmatch(name)

	if len(match) > 1 {
		res, err := strconv.Atoi(match[1])
		if err != nil {
			return 0
		}

		return res
	}

	return 0
}
