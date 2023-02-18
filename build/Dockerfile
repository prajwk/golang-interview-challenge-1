############################
# STEP 1 build executable binary
############################
# Download the base image for building Go apps
FROM golang:1.19-alpine as builder

RUN apk add --no-cache --virtual build-deps \
    make \
    gcc \
    build-base

# Set the working directory
WORKDIR /app
# Copy the code into the container
COPY . .
# Build the Go binary
ARG GO_BINARY
ENV GO_BINARY ${GO_BINARY:-./cmd/go-interview-challenge}
RUN CGO_ENABLED=1 go build -mod vendor -o ./bin/app $GO_BINARY

############################
# STEP 2 build a small image
############################
FROM alpine:3
COPY --from=builder /app/bin/app /app
ENTRYPOINT ["/app"]
