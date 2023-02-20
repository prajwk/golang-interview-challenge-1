.PHONY: test
ROOT_DIR = $(shell pwd)

#################################################################################
# RUN COMMANDS
#################################################################################
run:
	go mod vendor
	docker-compose --file ./build/docker-compose.yml --file ./build/docker-compose.dev.yml --project-directory . up --build; \
	docker-compose --file ./build/docker-compose.yml --file ./build/docker-compose.dev.yml --project-directory . down --volumes; \
	rm -rf vendor

#################################################################################
# LINT COMMANDS
#################################################################################
tidy:
	goimports -w .
	gofmt -s -w .
	go mod tidy

#################################################################################
# BUILD COMMANDS
#################################################################################
protoc:
	rm -rf pkg/*
	protoc -I=proto --go_out=. proto/*.proto

#################################################################################
# TEST COMMANDS
#################################################################################
generate-messages:
	go run test/generate_messages.go
