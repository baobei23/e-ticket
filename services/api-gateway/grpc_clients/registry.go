package grpc_clients

type ServiceRegistry struct {
	Event *EventServiceClient
	// add other clients here
}

func NewServiceRegistry() (*ServiceRegistry, error) {
	eventClient, err := NewEventServiceClient()
	if err != nil {
		return nil, err
	}

	// add other clients here

	return &ServiceRegistry{
		Event: eventClient,
		// add other clients here
	}, nil
}

// Helper to close all connections
func (r *ServiceRegistry) Close() {
	if r.Event != nil {
		r.Event.Close()
	}
	// close other clients here
}
