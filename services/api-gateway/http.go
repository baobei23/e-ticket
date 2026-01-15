package main

import (
	"io"
	"net/http"
	"strconv"

	"github.com/baobei23/e-ticket/services/api-gateway/grpc_clients"
	"github.com/baobei23/e-ticket/shared/proto/auth"
	bookingpb "github.com/baobei23/e-ticket/shared/proto/booking"
	eventpb "github.com/baobei23/e-ticket/shared/proto/event"
	paymentpb "github.com/baobei23/e-ticket/shared/proto/payment"
	"github.com/gin-gonic/gin"
)

type GatewayServer struct {
	eventClient   *grpc_clients.EventServiceClient
	bookingClient *grpc_clients.BookingServiceClient
	paymentClient *grpc_clients.PaymentServiceClient
	authClient    *grpc_clients.AuthServiceClient
}

type createBookingRequest struct {
	EventID  int64 `json:"event_id" binding:"required"`
	Quantity int32 `json:"quantity" binding:"required,min=1"`
}

type authRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type activateRequest struct {
	Token string `json:"token" binding:"required"`
}

func NewGatewayServer(eventClient *grpc_clients.EventServiceClient, bookingClient *grpc_clients.BookingServiceClient, paymentClient *grpc_clients.PaymentServiceClient, authClient *grpc_clients.AuthServiceClient) *GatewayServer {
	return &GatewayServer{
		eventClient:   eventClient,
		bookingClient: bookingClient,
		paymentClient: paymentClient,
		authClient:    authClient,
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

	resp, err := s.eventClient.Client.GetEvents(ctx, &eventpb.GetEventsRequest{
		Pagination: &eventpb.PaginationRequest{
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

	resp, err := s.eventClient.Client.GetEventDetail(ctx, &eventpb.GetEventDetailRequest{
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

	resp, err := s.eventClient.Client.CheckAvailability(ctx, &eventpb.CheckAvailabilityRequest{
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
	var req createBookingRequest // Pastikan struct ini tidak lagi mewajibkan field user_id dari JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// AMBIL UserID DARI CONTEXT (Diset oleh Middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Panggil Booking Service
	resp, err := s.bookingClient.Client.CreateBooking(c.Request.Context(), &bookingpb.CreateBookingRequest{
		UserId:   userID.(int64), // Type assertion
		EventId:  req.EventID,
		Quantity: req.Quantity,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create booking",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, resp)
}
func (s *GatewayServer) GetBookingDetailHandler(c *gin.Context) {
	bookingID := c.Param("id")
	if bookingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking ID is required"})
		return
	}

	// Simulasi ambil UserID dari context/token (karena belum ada auth middleware, kita ambil dari query param dulu)
	userIDStr := c.Query("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing user_id"})
		return
	}

	ctx := c.Request.Context()

	resp, err := s.bookingClient.Client.GetBookingDetail(ctx, &bookingpb.GetBookingDetailRequest{
		BookingId: bookingID,
		UserId:    userID,
	})

	if err != nil {
		// Mapping error gRPC status code bisa ditambahkan di sini
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch booking detail",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (s *GatewayServer) HandleStripeWebhook(c *gin.Context) {
	// 1. Baca Raw Body
	// Penting: Stripe butuh raw bytes untuk verifikasi signature.
	// Gin secara default membaca body stream, jadi kita harus baca manual.
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	sigHeader := c.GetHeader("Stripe-Signature")

	// 2. Kirim ke Payment Service via gRPC
	// Asumsi: client sudah ada method HandleWebhook (hasil generate proto baru)
	_, err = s.paymentClient.Client.HandleWebhook(c.Request.Context(), &paymentpb.HandleWebhookRequest{
		Payload:   payload,
		Signature: sigHeader,
	})

	if err != nil {
		// Log error
		// Return 400/500 agar Stripe tau webhook gagal diproses
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

func (s *GatewayServer) RegisterHandler(c *gin.Context) {
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := s.authClient.Client.Register(c.Request.Context(), &auth.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (s *GatewayServer) ActivateHandler(c *gin.Context) {
	token := c.Param("token")

	_, err := s.authClient.Client.Activate(c.Request.Context(), &auth.ActivateRequest{
		Token: token,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "activated"})
}

func (s *GatewayServer) LoginHandler(c *gin.Context) {
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := s.authClient.Client.Login(c.Request.Context(), &auth.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
