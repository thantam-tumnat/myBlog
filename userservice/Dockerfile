# Builder stage
FROM golang:1.21.4-alpine3.18 AS builder

WORKDIR /app

# Copy only the necessary files for Go module download and build
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the entire application source code
COPY . .

# Change working directory to the app directory
WORKDIR /app/app

# Build the Go application
RUN go build -o main

# Final stage
FROM scratch

ARG APP_PORT
ENV APP_PORT=${APP_PORT}

EXPOSE ${APP_PORT}

WORKDIR /app/app

# Copy only the necessary artifacts from the builder stage
COPY --from=builder /app/app/main .
COPY  ./config.env /app/config.env

CMD ["/app/app/main"]