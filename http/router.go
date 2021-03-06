package http

import (
	"one-cd/service"

	"github.com/gin-gonic/gin"
)

var svc *service.Service

// Init init http sever instance.
func Init(s *service.Service) {
	svc = s
}

// Start route
func Start(listen string) {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(cors())
	router.GET("/ping", handler(pingHandler))
	g1 := router.Group("/v1/deployment")
	{
		g1.POST("/deploy", handler(deployHandler))
		g1.POST("/update", handler(updateHandler))
		g1.POST("/undo", handler(undoHandler))
		g1.POST("/scale", handler(scaleHandler))
		g1.POST("/RollBack", handler(rollBackHandler))
		g1.DELETE("/delete", handler(deleteDeploymentHandler))
		g1.GET("/describe", handler(deploymentHandler))
		g1.GET("/rs", handler(replicaSetHandler))
		g1.GET("/scale", handler(getScaleHandler))
	}

	g2 := router.Group("/v1/pod")
	{
		g2.GET("/list", handler(podListHandler))
		g2.GET("/events", handler(podEventsHandler))
		g2.GET("/log", handler(podLogHandler))
	}

	g3 := router.Group("/v1/ingress")
	{
		g3.GET("/describe", handler(ingressHandler))
	}

	router.Run(listen)
}
