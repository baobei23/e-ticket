package grpc_clients

type ServiceRegistry struct {
	Event   *EventServiceClient
	Booking *BookingServiceClient
	Payment *PaymentServiceClient
	Auth    *AuthServiceClient
	// add other clients here
}

func NewServiceRegistry() (*ServiceRegistry, error) {
	eventClient, err := NewEventServiceClient()
	if err != nil {
		return nil, err
	}

	bookingClient, err := NewBookingServiceClient()
	if err != nil {
		return nil, err
	}

	paymentClient, err := NewPaymentServiceClient()
	if err != nil {
		return nil, err
	}

	authClient, err := NewAuthServiceClient()
	if err != nil {
		return nil, err
	}

	// add other clients here

	return &ServiceRegistry{
		Event:   eventClient,
		Booking: bookingClient,
		Payment: paymentClient,
		Auth:    authClient,
		// add other clients here
	}, nil
}

// Helper to close all connections
func (r *ServiceRegistry) Close() {
	if r.Event != nil {
		r.Event.Close()
	}
	if r.Booking != nil {
		r.Booking.Close()
	}
	if r.Payment != nil {
		r.Payment.Close()
	}
	if r.Auth != nil {
		r.Auth.Close()
	}
	// close other clients here
}
