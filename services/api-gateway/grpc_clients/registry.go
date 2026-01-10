package grpc_clients

type ServiceRegistry struct {
	Event   *EventServiceClient
	Booking *BookingServiceClient
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

	// add other clients here

	return &ServiceRegistry{
		Event:   eventClient,
		Booking: bookingClient,
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
	// close other clients here
}
