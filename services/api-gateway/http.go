package main

import (
	"net/http"
	"strconv"

	"github.com/baobei23/e-ticket/services/api-gateway/grpc_clients"
	"github.com/baobei23/e-ticket/shared/proto/booking"
	"github.com/baobei23/e-ticket/shared/proto/event"
	"github.com/gin-gonic/gin"
)

type GatewayServer struct {
	eventClient   *grpc_clients.EventServiceClient
	bookingClient *grpc_clients.BookingServiceClient
}

type createBookingRequest struct {
	EventID  int64 `json:"event_id" binding:"required"`
	Quantity int32 `json:"quantity" binding:"required,min=1"`
	UserID   int64 `json:"user_id" binding:"required"` // Nanti diambil dari Token/Context Auth
}

func NewGatewayServer(eventClient *grpc_clients.EventServiceClient, bookingClient *grpc_clients.BookingServiceClient) *GatewayServer {
	return &GatewayServer{
		eventClient:   eventClient,
		bookingClient: bookingClient,
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

func (s *GatewayServer) getEventDetailHandler(c *gin.Context) {
	eventIdStr := c.Param("id")
	eventId, _ := strconv.ParseInt(eventIdStr, 10, 64)

	ctx := c.Request.Context()

	resp, err := s.eventClient.Client.GetEventDetail(ctx, &event.GetEventDetailRequest{
		EventId: eventId,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch event detail",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (s *GatewayServer) checkAvailabilityHandler(c *gin.Context) {

	eventIdStr := c.Param("id")
	quantityStr := c.Query("quantity")

	eventId, err := strconv.ParseInt(eventIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing event_id"})
		return
	}

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil || quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing quantity"})
		return
	}

	ctx := c.Request.Context()

	resp, err := s.eventClient.Client.CheckAvailability(ctx, &event.CheckAvailabilityRequest{
		EventId:  eventId,
		Quantity: int32(quantity),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to check availability",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"is_available": resp.IsAvailable,
		"unit_price":   resp.UnitPrice,
	})
}

func (s *GatewayServer) CreateBookingHandler(c *gin.Context) {
	var req createBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	// Panggil Booking Service via gRPC
	resp, err := s.bookingClient.Client.CreateBooking(ctx, &booking.CreateBookingRequest{
		UserId:   req.UserID,
		EventId:  req.EventID,
		Quantity: req.Quantity,
	})

	if err != nil {
		// Bisa tambah mapping error code gRPC ke HTTP di sini (misal Unavailable -> 400/409)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create booking",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, resp)
}
