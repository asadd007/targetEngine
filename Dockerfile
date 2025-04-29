# First stage: build the app
FROM golang:1.22-alpine AS builder

# Set up our workspace
WORKDIR /app

# Copy dependency files
COPY go.mod ./

# Get all the dependencies
RUN go mod download

# Copy our code
COPY . .

# Build the app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o targeting-engine ./cmd/api

# Second stage: run the app
FROM alpine:3.14

# Set up our workspace
WORKDIR /app

# Add HTTPS support
RUN apk --no-cache add ca-certificates

# Copy the app from the builder
COPY --from=builder /app/targeting-engine /app/
COPY --from=builder /app/configs/config.json /app/configs/

# Tell Docker what port we use
EXPOSE 8080

# Run the app
CMD ["/app/targeting-engine"] 