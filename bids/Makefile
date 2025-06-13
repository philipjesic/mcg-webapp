# Makefile
# Run commands when running app localy

# Run the whole service locally with all dependencies
run local:
	docker-compose up -d

# Stop the service and its dependencies
stop:
	docker-compose down

# Run the dependencies only
run dependencies:
	docker-compose up -d mongo rabbitmq