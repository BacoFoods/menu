
# Setup environment

1. To run the project locally, create your own `.env` file with your configurations.
2. Start the database: `docker compose up database -d`.
3. Run the project: `bin/loadenv go run cmd/main.go`
4. For hot reload: `bin/loadenv air -c air.toml`


# Update docs:

1. Install `swaggo`: `go install github.com/swaggo/swag/cmd/swag@v1.16.2`
2. Run: `swag init -o pkg/swagger/docs -g cmd/main.go`