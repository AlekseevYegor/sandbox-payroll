package rest

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	timetrackingServ "se-challenge-payroll/internal/domain/timetracking"
	"se-challenge-payroll/internal/entity/timetracking"
	"se-challenge-payroll/internal/repository"
	"se-challenge-payroll/internal/service"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRouter_postTimeTrackerUploadCSV(t *testing.T) {
	// init mock service
	ttRepo := repository.NewTimeTrackerRepoMock()
	payrollRepo := repository.NewPayrollRepoMock(ttRepo)
	svc := service.NewService(timetrackingServ.New(ttRepo), payrollRepo)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	TimeTrackingRoutes(router.Group("/time_tracking"), &svc)

	type args struct {
		csvFilePath string
	}
	tests := []struct {
		name         string
		args         args
		waitHttpCode int
		prepareFunc  func(s service.Service)
	}{
		{
			name:         "Successful",
			args:         args{csvFilePath: "testdata/time-report-40.csv"},
			waitHttpCode: http.StatusCreated,
		},
		{
			name:         "Failed file name validation",
			args:         args{csvFilePath: "testdata/times-reports-41.csv"},
			waitHttpCode: http.StatusBadRequest,
		},
		{
			name:         "Failed file name validation - report id",
			args:         args{csvFilePath: "testdata/times-reports-100.csv"},
			waitHttpCode: http.StatusBadRequest,
		},
		{
			name: "Failed - duplicated report id",
			args: args{csvFilePath: "testdata/time-report-42.csv"},
			prepareFunc: func(s service.Service) { //Prepare data in db for this case
				ctx := gin.Context{}
				_ = s.TimeTrackingService.BulkCreate(&ctx, []timetracking.TimeTracking{{
					TimeReportID: 42,
					Date:         time.Now(),
					WorkedHours:  1,
					JobGroupID:   "A",
					EmployeeID:   "123444",
				}})
			},
			waitHttpCode: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepareFunc != nil {
				tt.prepareFunc(svc)
			}

			file, err := os.Open(tt.args.csvFilePath)
			if err != nil {
				t.Fatalf("Error opening file: %s", err.Error())
			}
			defer file.Close()
			// Create a buffer to store the request body
			body := &bytes.Buffer{}

			// Create a multipart writer to compose the request body
			writer := multipart.NewWriter(body)

			// Create a form file field
			part, err := writer.CreateFormFile("file", file.Name())
			if err != nil {
				t.Fatalf("Error creating form file: %s", err.Error())
			}

			// Copy file content into the form file field
			_, err = io.Copy(part, file)
			if err != nil {
				t.Fatalf("Error copying file content: %s", err.Error())
			}

			// Close the multipart writer to finalize the body
			err = writer.Close()
			if err != nil {
				t.Fatalf("Error closing multipart writer: %s", err.Error())
			}

			rr := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/time_tracking/upload", body)
			if err != nil {
				t.Fatalf("failed to create request: %s", err.Error())
			}

			req.Header.Set("Accept", "text/csv")
			req.Header.Set("Content-Type", writer.FormDataContentType())

			router.ServeHTTP(rr, req)

			// Check the response status code
			if status := rr.Code; status != tt.waitHttpCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.waitHttpCode)
			}
		})
	}
}

func Test_extractReportID(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Success",
			args: args{name: "time-report-40.csv"},
			want: 40,
		},
		{
			name: "Failed wrong file name",
			args: args{name: "times-reports-42.csv"},
			want: 0,
		},
		{
			name: "Failed wrong job reports id",
			args: args{name: "time-report-100.csv"},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractReportID(tt.args.name); got != tt.want {
				t.Errorf("extractReportID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateFileName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Success",
			args: args{name: "time-report-42.csv"},
			want: true,
		},
		{
			name: "Failed wrong file name",
			args: args{name: "times-reports-42.csv"},
			want: false,
		},
		{
			name: "Failed wrong job reports id",
			args: args{name: "time-report-100.csv"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateFileName(tt.args.name); got != tt.want {
				t.Errorf("validateFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}
