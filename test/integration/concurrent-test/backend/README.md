# concurrent-test FastAPI Backend

A FastAPI backend service for concurrent-test.

## Features

- FastAPI framework with automatic API documentation
- CORS middleware for cross-origin requests
- Health check endpoint
- Environment-based configuration
- Docker support

## Quick Start

### Prerequisites

- Python 3.11+
- pip

### Local Development

1. Create a virtual environment:
```bash
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate
```

2. Install dependencies:
```bash
pip install -r requirements.txt
```

3. Copy environment variables:
```bash
cp .env.example .env
```

4. Run the development server:
```bash
python main.py
```

The API will be available at http://localhost:8000

### Docker Development

1. Build the Docker image:
```bash
docker build -t concurrent-test-backend .
```

2. Run the container:
```bash
docker run -p 8000:8000 concurrent-test-backend
```

## API Documentation

Once the server is running, you can access:
- Interactive API docs: http://localhost:8000/docs
- ReDoc documentation: http://localhost:8000/redoc

## Available Endpoints

- `GET /` - Welcome message
- `GET /health` - Health check endpoint

## Environment Variables

See `.env.example` for all available environment variables.

## Project Structure

```
.
├── main.py              # Application entry point
├── requirements.txt     # Python dependencies
├── Dockerfile          # Docker configuration
├── .env.example        # Environment variables template
├── .gitignore         # Git ignore rules
└── README.md          # This file
```