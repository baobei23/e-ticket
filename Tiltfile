# Load the restart_process extension

load('ext://restart_process', 'docker_build_with_restart')
### K8s Config ###

# Uncomment to use secrets
k8s_yaml('./infra/development/k8s/secrets.yaml')
k8s_yaml('./infra/development/k8s/app-config.yaml')

### End of K8s Config ###

### RabbitMQ ###
k8s_yaml('./infra/development/k8s/rabbitmq-deployment.yaml')
k8s_resource('rabbitmq', port_forwards=['5672', '15672'], labels='tooling')
### End RabbitMQ ###

### PostgreSQL ###
k8s_yaml('./infra/development/k8s/postgres-deployment.yaml')
k8s_resource('eticket-postgres', port_forwards=['5432'], labels='tooling')
### End PostgreSQL ###

### API Gateway ###

gateway_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/api-gateway ./services/api-gateway'
if os.name == 'nt':
  gateway_compile_cmd = './infra/development/docker/api-gateway-build.bat'

local_resource(
  'api-gateway-compile',
  gateway_compile_cmd,
  deps=['./services/api-gateway', './shared'], labels="compiles")


docker_build_with_restart(
  'e-ticket/api-gateway',
  '.',
  entrypoint=['/app/build/api-gateway'],
  dockerfile='./infra/development/docker/api-gateway.Dockerfile',
  only=[
    './build/api-gateway',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)

k8s_yaml('./infra/development/k8s/api-gateway-deployment.yaml')
k8s_resource('api-gateway', port_forwards=8080, resource_deps=['api-gateway-compile', 'rabbitmq'], labels="services")
### End of API Gateway ###
### Event Service ###

event_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/event-service ./services/event-service/cmd/main.go'
if os.name == 'nt':
 event_compile_cmd = './infra/development/docker/event-build.bat'

local_resource(
  'event-service-compile',
  event_compile_cmd,
  deps=['./services/event-service', './shared'], labels="compiles")

docker_build_with_restart(
  'e-ticket/event-service',
  '.',
  entrypoint=['/app/build/event-service'],
  dockerfile='./infra/development/docker/event-service.Dockerfile',
  only=[
    './build/event-service',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)

k8s_yaml('./infra/development/k8s/event-service-deployment.yaml')
k8s_resource('event-service', resource_deps=['event-service-compile', 'rabbitmq', 'eticket-postgres'], labels="services")

### End of Event Service ###
### Booking Service ###

booking_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/booking-service ./services/booking-service/cmd/main.go'
if os.name == 'nt':
 booking_compile_cmd = './infra/development/docker/booking-build.bat'

local_resource(
  'booking-service-compile',
  booking_compile_cmd,
  deps=['./services/booking-service', './shared'], labels="compiles")

docker_build_with_restart(
  'e-ticket/booking-service',
  '.',
  entrypoint=['/app/build/booking-service'],
  dockerfile='./infra/development/docker/booking-service.Dockerfile',
  only=[
    './build/booking-service',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)

k8s_yaml('./infra/development/k8s/booking-service-deployment.yaml')
k8s_resource('booking-service', resource_deps=['booking-service-compile', 'rabbitmq', 'eticket-postgres'], labels="services")

### End of Booking Service ###
### Payment Service ###

payment_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/payment-service ./services/payment-service/cmd/main.go'
if os.name == 'nt':
 payment_compile_cmd = './infra/development/docker/payment-build.bat'

local_resource(
  'payment-service-compile',
  payment_compile_cmd,
  deps=['./services/payment-service', './shared'], labels="compiles")

docker_build_with_restart(
  'e-ticket/payment-service',
  '.',
  entrypoint=['/app/build/payment-service'],
  dockerfile='./infra/development/docker/payment-service.Dockerfile',
  only=[
    './build/payment-service',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)

k8s_yaml('./infra/development/k8s/payment-service-deployment.yaml')
k8s_resource('payment-service', resource_deps=['payment-service-compile', 'rabbitmq', 'eticket-postgres'], labels="services")

### End of Payment Service ###

### Auth Service ###

auth_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/auth-service ./services/auth-service/cmd/main.go'
if os.name == 'nt':
 auth_compile_cmd = './infra/development/docker/auth-build.bat'

local_resource(
  'auth-service-compile',
  auth_compile_cmd,
  deps=['./services/auth-service', './shared'], labels="compiles")

docker_build_with_restart(
  'e-ticket/auth-service',
  '.',
  entrypoint=['/app/build/auth-service'],
  dockerfile='./infra/development/docker/auth-service.Dockerfile',
  only=[
    './build/auth-service',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)

k8s_yaml('./infra/development/k8s/auth-service-deployment.yaml')
k8s_resource('auth-service', resource_deps=['auth-service-compile', 'eticket-postgres'], labels="services")

### End of Auth Service ###