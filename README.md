# Ad Targeting Engine

A service that shows the right ads to the right people.

## What It Does

This service looks at who's asking for ads and figures out which ones to show them. It can target ads based on:

- Which app they're using
- Where they are in the world
- What kind of device they're on (Android, iOS, etc.)

You can set rules to either show ads in certain places or hide them in certain places.

## How to Use It

### Getting Ads

**GET /v1/delivery**

This gives you a list of ads that match what you're looking for.

**What You Need to Tell Us:**

- `app` (required): Which app this is for
- `os` (required): What operating system they're using
- `country` (required): Where they are

**What You'll Get Back:**

- `200 OK`: Here are your matching ads
- `204 No Content`: No ads match what you're looking for
- `400 Bad Request`: You forgot to tell us something
- `500 Internal Server Error`: Something broke on our end

**Example:**

```
GET /v1/delivery?app=com.abc.xyz&country=germany&os=android
```

**What You'll Get:**

```json
[
  {
    "cid": "duolingo",
    "img": "https://somelink2",
    "cta": "Install"
  }
]
```

### Checking If We're Alive

**GET /health**

This tells you if the service is running okay.

**What You'll Get Back:**

- `200 OK`: We're up and running!

## How to Set It Up

You can change how it works using environment variables or a settings file (`configs/config.json`):

### Server Settings
- `PORT`: What port to run on (default: 8080)
- `LOG_LEVEL`: How much to log (default: info)
- `ENABLE_METRICS`: Whether to track metrics (default: false)
- `METRICS_PORT`: What port to run metrics on (default: 9090)
- `ENABLE_HEALTH_CHECK`: Whether to have a health check (default: true)

### Database Settings
- `DB_TYPE`: What kind of database to use, either "memory" or "postgres" (default: memory)
- `POSTGRES_URI`: How to connect to PostgreSQL (default: postgres://localhost:27017)
- `DB_NAME`: What to call the database (default: targeting_engine)

## Where We Store Stuff

### In-Memory Storage
By default, we keep everything in memory. This means if you restart the service, all the data disappears. Good for testing and development.

### PostgreSQL
For real use, you can use PostgreSQL to keep your data safe. We'll create these tables:

- `campaigns`: Where we keep info about ads
- `targeting_rules`: Where we keep the rules for showing ads

To use PostgreSQL, set `DB_TYPE` to "postgres" in your settings.

## How to Run It

### What You Need

- Go 1.16 or newer
- PostgreSQL (if you're using PostgreSQL)

### Building

```bash
go build -o targeting-engine ./cmd/api
```

### Running

```bash
./targeting-engine
```

### Running with Docker

```bash
# Build the Docker image
docker build -t targeting-engine .

# Run with in-memory storage
docker run -p 8080:8080 targeting-engine

# Run with PostgreSQL
docker run -p 8080:8080 -e DB_TYPE=postgres -e POSTGRES_URI=postgres://postgres:postgres@postgres:5432/targeting_engine?sslmode=disable targeting-engine
```

### Running with Docker Compose

We've got a `docker-compose.yml` file to run everything together:

```bash
docker-compose up
```

## Some Examples

### Example 1: Someone in Germany on Android

Request:
```
GET /v1/delivery?app=com.abc.xyz&country=germany&os=android
```

Response:
```json
[
  {
    "cid": "duolingo",
    "img": "https://somelink2",
    "cta": "Install"
  }
]
```

### Example 2: Someone in the US with Ludo King on Android

Request:
```
GET /v1/delivery?app=com.gametion.ludokinggame&country=us&os=android
```

Response:
```json
[
  {
    "cid": "spotify",
    "img": "https://somelink",
    "cta": "Download"
  },
  {
    "cid": "subwaysurfer",
    "img": "https://somelink3",
    "cta": "Play"
  }
]
```

### Example 3: Forgot to Tell Us Something

Request:
```
GET /v1/delivery?country=germany&os=android
```

Response:
```json
{
  "error": "Hey, we need to know which app this is for!"
}
```

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

### Using Docker

1. Build the Docker image:
```bash
docker build -t targeting-engine .
```

2. Run the container:
```bash
docker run -p 8080:8080 targeting-engine
```

The application will be available at `http://localhost:8080` 