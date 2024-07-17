# syntax=docker/dockerfile:1

FROM golang:1.22 AS build-stage

WORKDIR /app

# Download Go modules
COPY . .
RUN go mod download

# WORKDIR /app/cmd/challenge-02

# # Build
# RUN CGO_ENABLED=0 GOOS=linux go build -o /app/cmd/challenge-02/main

# ## Run the tests in the container
# # FROM build-stage AS run-test-stage
# # RUN go test -v ./...

# # Use a minimal runtime image
# FROM alpine:3.19

# # Copy the executable and the configs directory from the builder stage
# COPY --from=build-stage /app/cmd/challenge-02/main /app/cmd/challenge-02/main
# COPY --from=build-stage /app/internal/postgresql/migrations /app/internal/postgresql/migrations
# COPY --from=build-stage /app/*.env /app

# Set the working directory
# WORKDIR /app
# Run
CMD ["go", "run","./cmd/challenge-02/main.go"]