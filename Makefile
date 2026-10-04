dev:
	@echo "Rump up docker containers (LocalStack)"
	mkdir -p ./volume
	chmod 777 ./volume
	docker compose -f docker-compose.yml -p cloudforge up -d
