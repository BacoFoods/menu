FROM golang AS builder

WORKDIR /src

# Download dependencies
COPY go.mod go.sum /
RUN go mod download

# Add source code
COPY . .
RUN CGO_ENABLED=0 go build -o main ./cmd

# Multi-Stage production build
FROM alpine AS production
RUN apk --no-cache add ca-certificates

WORKDIR /app
# Retrieve the binary from the previous stage
COPY --from=builder /src/main .

COPY entrypoint.sh .

# Set the binary as the entrypoint of the container
EXPOSE 3000
ENTRYPOINT ["/bin/sh"]
CMD ["entrypoint.sh"]
