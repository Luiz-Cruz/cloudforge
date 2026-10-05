dev:
	@echo "Rump up docker containers (LocalStack)"
	mkdir -p ./volume
	chmod 777 ./volume
	docker compose -f docker-compose.yml -p cloudforge up -d


.PHONY: build package

build:
	@echo "Building Go binary for AWS Lambda (linux/amd64)..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap main.go

package: build
	@echo "Zipping binary for deployment..."
	zip function.zip bootstrap
	@echo "Package function.zip is ready!"
