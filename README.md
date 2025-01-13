# Applications List

A simple web application built with Go (Gin) to manage a list of applications with their logos and links.

## Features

- List applications with their logos and links
- Add new applications with logo upload
- Delete existing applications
- Persistent SQLite database storage

## Prerequisites

- Docker
- Docker Compose

## Running with Docker

1. Clone the repository:
```bash
git clone https://github.com/ed2ti/applicationsList.git
cd applicationsList
```

2. Create required directories:
```bash
mkdir -p images data/db
```

3. Build and run with Docker Compose:
```bash
docker-compose up -d --build
```

4. Access the application at:
```
http://localhost:8080
```

## Project Structure

- `main.go` - Main application code
- `Dockerfile` - Docker container configuration
- `docker-compose.yaml` - Docker Compose configuration
- `templates/` - HTML templates
- `images/` - Uploaded application logos (created at runtime)
- `data/db/` - SQLite database storage (created at runtime)

## License

MIT 