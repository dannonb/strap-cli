# concurrent-test

This microservice project was generated using the Microservice Bootstrapper.

## Architecture Overview

This project follows a microservices architecture with the following components:


### Backend Service
- **Technology**: fastapi
- **Port**: 8000
- **Location**: `./backend/`






## Getting Started

### Prerequisites

Make sure you have the following installed:
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Quick Start

1. **Clone and navigate to the project**:
   ```bash
   cd concurrent-test
   ```

2. **Set up environment variables**:
   ```bash
   cp .env.example .env
   ```
   Edit `.env` file with your specific configuration values.

3. **Start all services**:
   ```bash
   docker-compose up -d
   ```

4. **Verify services are running**:
   ```bash
   docker-compose ps
   ```

### Service URLs

Once all services are running, you can access them at:


- **Backend API**: http://localhost:8000
  - API Documentation: http://localhost:8000/docs








## Development

### Starting Services

To start all services in development mode:
```bash
docker-compose up
```

To start services in the background:
```bash
docker-compose up -d
```

### Stopping Services

To stop all services:
```bash
docker-compose down
```

To stop and remove volumes (⚠️ **This will delete database data**):
```bash
docker-compose down -v
```

### Viewing Logs

To view logs from all services:
```bash
docker-compose logs -f
```

To view logs from a specific service:
```bash
docker-compose logs -f backend


```

### Rebuilding Services

If you make changes to Dockerfiles or dependencies:
```bash
docker-compose build
docker-compose up -d
```

## Project Structure

```
concurrent-test/
├── docker-compose.yml          # Docker orchestration configuration
├── .env.example               # Environment variables template
├── .env                       # Your local environment variables (git-ignored)
├── .gitignore                # Git ignore rules
├── README.md                 # This file
├── backend/                   # Backend service
│   ├── Dockerfile            # Backend container configuration
│   ├── requirements.txt      # Python dependencies
│   ├── main.py              # FastAPI application entry point
│   └── app/                 # Application source code


└── docs/                      # Project documentation
    └── setup.md              # Detailed setup instructions
```

## Environment Configuration

The project uses environment variables for configuration. Key variables include:




### Backend Configuration
- `BACKEND_PORT`: 8000




See `.env.example` for a complete list of configurable variables.

## API Documentation



### FastAPI Backend
- **Interactive API Docs**: http://localhost:8000/docs
- **OpenAPI Schema**: http://localhost:8000/openapi.json







## Troubleshooting

### Common Issues

1. **Port conflicts**: If you get port binding errors, check if the ports are already in use:
   ```bash
   lsof -i :8000
   
   
   ```

2. **Docker issues**: Make sure Docker is running and you have sufficient resources allocated.

3. **Permission issues**: On Linux/macOS, you might need to run Docker commands with `sudo` or add your user to the docker group.

### Getting Help

- Check service logs: `docker-compose logs [service-name]`
- Restart services: `docker-compose restart [service-name]`
- Rebuild containers: `docker-compose build [service-name]`

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes
4. Test your changes: `docker-compose up --build`
5. Commit your changes: `git commit -am 'Add feature'`
6. Push to the branch: `git push origin feature-name`
7. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.