## Booking Expiry (Auto-Cancel via DLX)

Jika user tidak menyelesaikan pembayaran dalam 30 menit, booking otomatis
dibatalkan dan stok dikembalikan.

### Mekanisme

1. Saat booking dibuat, message `BookingExpiry` dikirim ke queue
   `booking.expiry.pending` dengan per-message TTL 30 menit.
2. Queue ini **tidak memiliki consumer** — message hanya "tidur" sampai TTL
   habis.
3. Setelah TTL habis, RabbitMQ memindahkan message ke `booking_expiry_dlx` (Dead
   Letter Exchange).
4. DLX merouting message ke queue `booking.expiry.process`.
5. `ExpiryConsumer` di booking-service memproses message.

### Flow

```mermaid
sequenceDiagram
  participant MQ_Pending as booking.expiry.pending
  participant DLX as booking_expiry_dlx
  participant MQ_Process as booking.expiry.process
  participant Booking as Booking Service
  participant Event as Event Service
  participant DB as Postgres
  Note over MQ_Pending: Message TTL 30 menit habis
  MQ_Pending -->> DLX: dead letter (expired)
  DLX -->> MQ_Process: route via "BookingExpired"
  MQ_Process -->> Booking: ExpiryConsumer consume
  Booking ->> DB: SELECT booking
  Booking ->> DB: UPDATE status = CANCELLED (WHERE status = PENDING)
  alt Status berhasil diubah
    Booking ->> Event: gRPC ReleaseSeat
    Event ->> DB: available_seats + qty
  else Status bukan PENDING (sudah CONFIRMED/FAILED)
    Note over Booking: Skip — booking sudah diproses
  end
```
