package grpc

import (
	"context"

	"github.com/baobei23/e-ticket/services/booking-service/internal/domain"
	bookingpb "github.com/baobei23/e-ticket/shared/proto/booking"
	"google.golang.org/grpc"
)

type BookingHandler struct {
	bookingpb.UnimplementedBookingServiceServer
	service domain.BookingService
}

func NewBookingHandler(server *grpc.Server, service domain.BookingService) *BookingHandler {
	handler := &BookingHandler{service: service}
	bookingpb.RegisterBookingServiceServer(server, handler)
	return handler
}

func (h *BookingHandler) CreateBooking(ctx context.Context, req *bookingpb.CreateBookingRequest) (*bookingpb.CreateBookingResponse, error) {
	booking, paymentURL, err := h.service.CreateBooking(ctx, req.UserId, req.EventId, req.Quantity)
	if err != nil {
		return nil, err
	}

	return &bookingpb.CreateBookingResponse{
		BookingId:      booking.ID,
		PaymentUrl:     paymentURL,
		TimeoutSeconds: 600, // 10 menit
	}, nil
}

func (h *BookingHandler) GetBookingDetail(ctx context.Context, req *bookingpb.GetBookingDetailRequest) (*bookingpb.GetBookingDetailResponse, error) {
	booking, err := h.service.GetBookingDetail(ctx, req.BookingId, req.UserId)
	if err != nil {
		return nil, err
	}

	return &bookingpb.GetBookingDetailResponse{
		Booking: booking.ToProto(),
	}, nil
}
