build:
	./build.sh

compose-up:
	docker-compose -f docker-compose.yaml -p lib_manager up -d --remove-orphans

compose-down:
	docker-compose -f docker-compose.yaml -p lib_manager down

network:
	docker network create applications

volume:
	docker volume create database

manager:
	docker run -d \
		--name manager \
		--restart always \
		--network applications \
		-p 8080:8080 \
		-v database:/app \
		lib_manager/manager:0.0.1

integration:
	docker run -d \
		--name integration \
		--network applications \
		lib_manager/integration:0.0.1

run-all: network volume manager integration
	@echo "All services are up and running!"

clean:
	docker stop manager integration || true
	docker rm manager integration || true
	docker network rm applications || true
	docker volume rm database || true
	@echo "Cleaned up all containers, network, and volume."