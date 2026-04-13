## End-to-End Booking Flow (Register → Booking → Payment → Confirmed)

```mermaid
sequenceDiagram
  participant User as User
  participant APIGW as API Gateway
  participant Auth as Auth Service
  participant Event as Event Service
  participant Booking as Booking Service
  participant Payment as Payment Service
  participant Notif as Notification Service
  participant Stripe as Stripe API
  participant DB as Postgres
  participant MQ as RabbitMQ
  User ->> APIGW: POST /auth/register
  APIGW ->> Auth: gRPC Register
  Auth ->> DB: INSERT user (is_active=false) + hashed token
  Auth -->> MQ: UserActivationRequested
  Auth -->> APIGW: user_id
  APIGW -->> User: 201 Created
  MQ -->> Notif: consume UserActivationRequested
  Notif ->> User: Activation email (plain token in link)
  User ->> APIGW: GET /auth/activate/:token
  APIGW ->> Auth: gRPC Activate
  Auth ->> DB: hash token, match, set is_active=true
  Auth -->> APIGW: success
  APIGW -->> User: 200 OK
  User ->> APIGW: POST /auth/login
  APIGW ->> Auth: gRPC Login
  Auth ->> DB: SELECT user (is_active=true)
  Auth -->> APIGW: JWT access_token
  APIGW -->> User: 200 OK
  User ->> APIGW: POST /bookings (event_id, qty)
  APIGW ->> Booking: gRPC CreateBooking
  Booking ->> Event: gRPC ReserveSeat (atomic reduce + return price)
  Event ->> DB: UPDATE stock, RETURNING price
  Event -->> Booking: unit_price
  Booking ->> DB: INSERT booking (PENDING, expires_at)
  Booking -->> MQ: BookingExpiry (TTL 30 min, DLX)
  Booking ->> Payment: gRPC CreatePayment
  Payment ->> Stripe: Create checkout session
  Stripe -->> Payment: session + payment_url
  Payment -->> Booking: payment_url
  Booking -->> APIGW: booking + payment_url
  APIGW -->> User: 201 Created
  User ->> Stripe: Complete payment
  Stripe ->> APIGW: Webhook checkout.session.completed
  APIGW ->> Payment: gRPC HandleWebhook
  Payment ->> DB: update payment status (SUCCESS)
  Payment -->> MQ: PaymentSuccess
  MQ -->> Booking: consume PaymentSuccess
  Booking ->> DB: UPDATE booking (CONFIRMED)
```
