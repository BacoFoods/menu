
# Setup environment

1. To run the project locally, create your own `.env` file with your configurations.
2. Start the database: `docker compose up database`.
3. Run the project: `bin/loadenv go run cmd/main.go`
4. For hot reload: `bin/loadenv air -c air.toml`