package main

import (
	"net/http"
	"strconv"

	"github.com/baobei23/e-ticket/services/api-gateway/grpc_clients"
	"github.com/baobei23/e-ticket/shared/proto/event"
	"github.com/gin-gonic/gin"
)

type GatewayServer struct {
	eventClient *grpc_clients.EventServiceClient
}

func NewGatewayServer(eventClient *grpc_clients.EventServiceClient) *GatewayServer {
	return &GatewayServer{
		eventClient: eventClient,
	}
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *GatewayServer) getEventsHandler(c *gin.Context) {

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	ctx := c.Request.Context()

	resp, err := s.eventClient.Client.GetEvents(ctx, &event.GetEventsRequest{
		Pagination: &event.PaginationRequest{
			Page:  int32(page),
			Limit: int32(limit),
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch events",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
