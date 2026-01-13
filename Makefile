include .env
PROTO_DIR := proto
PROTO_SRC := $(wildcard $(PROTO_DIR)/*.proto)
GO_OUT := .

# Helper URL builder
DB_DSN = postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)

.PHONY: generate-proto
generate-proto:
	protoc \
		--proto_path=. \
		--go_out=$(GO_OUT) \
		--go-grpc_out=$(GO_OUT) \
		$(PROTO_SRC)

# --- MIGRATIONS ---

# Usage: make migrate-create service=event name=create_events_table
.PHONY: migrate-create
migrate-create:
	@if [ -z "$(service)" ] || [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create service=<service_name> name=<migration_name>"; \
		exit 1; \
	fi
	@mkdir -p services/$(service)-service/cmd/migrations
	@migrate create -seq -ext sql -dir services/$(service)-service/cmd/migrations $(name)

# Usage: make migrate-up service=event
.PHONY: migrate-up
migrate-up:
	@if [ -z "$(service)" ]; then \
		echo "Migrating ALL services..."; \
		$(MAKE) migrate-up-event; \
		$(MAKE) migrate-up-booking; \
		$(MAKE) migrate-up-payment; \
		$(MAKE) migrate-up-auth; \
	else \
		echo "Migrating $(service)-service..."; \
		migrate -path services/$(service)-service/cmd/migrations -database "$(DB_DSN)/$(service)_service?sslmode=$(DB_SSL)" up; \
	fi

# Usage: make migrate-down service=event
.PHONY: migrate-down
migrate-down:
	@if [ -z "$(service)" ]; then \
		echo "Rolling back ALL services..."; \
		$(MAKE) migrate-down-event; \
		$(MAKE) migrate-down-booking; \
		$(MAKE) migrate-down-payment; \
		$(MAKE) migrate-down-auth; \
	else \
		echo "Rolling back $(service)-service..."; \
		migrate -path services/$(service)-service/cmd/migrations -database "$(DB_DSN)/$(service)_service?sslmode=$(DB_SSL)" down; \
	fi
# --- Individual Service Shortcuts (Optional) ---
.PHONY: migrate-up-event migrate-up-booking migrate-up-payment migrate-up-auth
.PHONY: migrate-down-event migrate-down-booking migrate-down-payment migrate-down-auth

# --- Migrate Up Shortcuts ---
migrate-up-event:
	@migrate -path services/event-service/cmd/migrations -database "$(DB_DSN)/event_service?sslmode=$(DB_SSL)" up

migrate-up-booking:
	@migrate -path services/booking-service/cmd/migrations -database "$(DB_DSN)/booking_service?sslmode=$(DB_SSL)" up

migrate-up-payment:
	@migrate -path services/payment-service/cmd/migrations -database "$(DB_DSN)/payment_service?sslmode=$(DB_SSL)" up

migrate-up-auth:
	@migrate -path services/auth-service/cmd/migrations -database "$(DB_DSN)/auth_service?sslmode=$(DB_SSL)" up

# --- Migrate Down Shortcuts ---
migrate-down-event:
	@migrate -path services/event-service/cmd/migrations -database "$(DB_DSN)/event_service?sslmode=$(DB_SSL)" down

migrate-down-booking:
	@migrate -path services/booking-service/cmd/migrations -database "$(DB_DSN)/booking_service?sslmode=$(DB_SSL)" down

migrate-down-payment:
	@migrate -path services/payment-service/cmd/migrations -database "$(DB_DSN)/payment_service?sslmode=$(DB_SSL)" down

migrate-down-auth:
	@migrate -path services/auth-service/cmd/migrations -database "$(DB_DSN)/auth_service?sslmode=$(DB_SSL)" down