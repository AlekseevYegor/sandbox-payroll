package rest

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"se-challenge-payroll/internal/service"
	"strings"
	"time"
)

type (
	API struct {
		*gin.Engine
		service *service.Service
	}
)

// New creates API instance.
func New(s *service.Service) API {
	// create the framework engine
	router := gin.New()
	/* Set Middleware.*/
	router.
		// Include Recovery middleware.
		Use(gin.Recovery())
	// Include CORS.
	router.Use(
		cors.New(cors.Config{
			AllowOrigins:    strings.Split(viper.GetString("CORS_WILDCARD"), ","),
			AllowMethods:    []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodOptions, http.MethodDelete},
			AllowAllOrigins: false,
			AllowHeaders: []string{
				"Authorization",
				"Origin",
				"Content-Length",
				"Content-Type",
				"Accept-Language",
				"Content-Language",
				"http_accept_language",
				"http_content_language",
				"x-trace-id"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			AllowWildcard:    true,
			MaxAge:           12 * time.Hour,
		}))

	api := API{
		router,
		s,
	}

	/*
	 * API version 1.
	 */
	v1 := api.Group("/api/v1")

	TimeTrackingRoutes(v1.Group("/time_tracking"), s)
	PayrollRoutes(v1.Group("/payroll"), s)

	/*
	 *	Docs endpoint.
	 */
	api.StaticFile("/swagger/swagger.json", "static/swagger.json")
	api.Static("/doc/api", "static/swaggerui")

	return api
}
