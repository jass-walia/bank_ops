# Build stage.
FROM golang:1.14-alpine AS builder

# Add Maintainer Info.
LABEL maintainer="Jaspreet Singh <jaspreet.surmount@gmail.com>"

# ENV: define runtime variables, these are activated when the container is started.
# Set necessary environmet variables needed for our image.
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

# Move to working directory inside the container.
WORKDIR /app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed.
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container.
COPY . .

# Build the application.
RUN go build -o main .

FROM alpine:latest AS production

# Move to working directory inside the container.
WORKDIR /app

# Copy binary and env file to production build.
COPY --from=builder /app/main ./
COPY --from=builder /app/.env ./

# Command to run
ENTRYPOINT ["sh", "-c", "/app/main -v=2 -logtostderr"]