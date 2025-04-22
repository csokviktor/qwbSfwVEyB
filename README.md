# Library Management System - Home Assignment

## Project Setup and Execution

### 1. Prerequisites

#### 1.1 Core Requirements

- **Docker**: Container runtime environment  
  [Installation Guide](https://docs.docker.com/engine/install/ubuntu/)
- **Docker Compose**: Container orchestration  
  [Installation Guide](https://docs.docker.com/compose/install/standalone/)
- **Make**: Build automation tool

#### 1.2 Development Requirements

- **Go**: Backend implementation  
  [Installation Guide](https://go.dev/doc/install)
- **golangci-lint**: Static code analysis  
  [Installation Guide](https://golangci-lithub.io/legacy-v1-doc/)
- **mockgen**: Golang mock generator  
  [Installation Guide](https://github.com/uber-go/mock)
- **Python**: Integration testing framework  
  [Download](https://www.python.org/downloads/)
- **Poetry**: Python dependency management  
  [Installation Guide](https://python-poetry.org/docs/)

### 2. Local Deployment

Execute the complete system stack including services and integration tests:

```bash
cd deployment
make build          # Build all Docker images
make compose-up     # Start services and execute test suite
```

If you do not have `docker-compose` you can run it using only `make`.

#### Run everything with one command

```sh
cd deployment
make build   # Build all Docker images
make run-all # Start services and execute test suite
```

#### Run individual components

```sh
make network    # Create the network
make volume     # Create the volume
make manager    # Start the manager service
make integration # Start the integration service
```

#### Stop and remove everything

```sh
make compose-down # If it was started with compose-up
make clean # If it was started with run-all
```

### 3. Development Workflows

#### 3.1 API Service Development

```sh
make run        # Start development server with hot-reload
make lint       # Execute static code analysis
make mocks      # Generate mocks for tests
make tests      # Run unit tests
make migration name=<migration_name> # Create a new migration
```

#### 3.2 Integration Testing

```sh
make tests  # Run integration tests
```