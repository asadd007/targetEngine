
## How to Start the Application

### Using Docker Compose (Recommended)

1. Make sure you have Docker and Docker Compose installed
2. Run the following command:
```bash
docker-compose up
```

This will:
- Start the application on port 8080
- Start PostgreSQL database
- Set up all necessary configurations

## Test With Curl
curl "http://localhost:8080/v1/delivery?app=spotify&os=ios&country=US"
curl "http://localhost:8080/v1/delivery?app=duolingo&os=ios&country=UK"
curl "http://localhost:8080/v1/delivery?app=com.gametion.ludokinggame&os=Android&country=US"