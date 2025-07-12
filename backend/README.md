<p align="center">
    <h1 align="center">LEADERPRO BACKEND</h1>
</p>
<p align="center">
    <em>
    AI platform backend that amplifies leadership intelligence using Domain-Driven Design (DDD) and Clean Architecture principles.
    </em>
</p>
<p align="center">
   <a href='https://github.com/diegoclair/leaderpro/commits/main'>
	<img src="https://img.shields.io/github/last-commit/diegoclair/leaderpro?style=flat&logo=git&logoColor=white&color=0080ff" alt="last-commit">
   </a>
   <a href="https://github.com/diegoclair/leaderpro/actions">
     <img src="https://github.com/diegoclair/leaderpro/actions/workflows/ci.yaml/badge.svg" alt="build status">
   </a>
  <a href='https://goreportcard.com/badge/github.com/diegoclair/leaderpro'>
     <img src='https://goreportcard.com/badge/github.com/diegoclair/leaderpro' alt='Go Report'/>
    </a>
<p>
<p align="center">
		<em>Developed with the software and tools below.</em>
</p>
<p align="center">
    <img src="https://skillicons.dev/icons?i=githubactions,mysql,redis,go,docker,prometheus,grafana,jaeger" >
</p>

## Table of Contents
- [Description](#description)
- [Project Architecture](#project-architecture)
  - [Directory Structure](#directory-structure)
  - [Dependency Rule](#dependency-rule)
- [Getting Started](#-getting-started)
  - [Prerequisites](#prerequisites-)
  - [Configuration](#configuration)
  - [Launching the Application](#️-launching-the-application)
  - [First Request Example](#first-request-example)
- [Testing](#testing-)
  - [Unit Tests](#unit-tests)
  - [Mocks](#mocks)
  - [Running Tests](#running-tests)
- [API Documentation](#-api-documentation)
  - [Generating Docs](#generating-docs)
- [Contributing](#-contributing)
- [License](#-license)

## Description
LeaderPro Backend is an AI-powered platform that amplifies leadership intelligence using **Domain-Driven Design (DDD)** and **Clean Architecture** principles. The backend provides APIs for team management, 1:1 meetings, and contextual AI coaching features.

The platform maintains perfect memory of every team interaction and suggests contextual actions to help leaders become more effective. Built with Go, it provides a solid foundation for scalable leadership intelligence features.

## Current Architecture Implementation

### Simplified Company-User Relationship
The application uses a simplified ownership model where:
- Each company belongs directly to one user via `user_owner_id`
- Each user can own multiple companies
- Frontend onboarding creates the user's first company automatically

### Company Entity Structure
```sql
CREATE TABLE tab_company (
    company_id INT NOT NULL AUTO_INCREMENT,
    company_uuid CHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    industry VARCHAR(100) NULL,
    size VARCHAR(50) NULL,
    role VARCHAR(200) NULL,
    is_default TINYINT(1) NOT NULL DEFAULT 0,
    user_owner_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    active TINYINT(1) NOT NULL DEFAULT 1,
    PRIMARY KEY (company_id),
    UNIQUE KEY uk_company_uuid (company_uuid),
    KEY fk_company_user (user_owner_id),
    CONSTRAINT fk_company_user FOREIGN KEY (user_owner_id) REFERENCES tab_user (user_id)
)
```

### Onboarding Flow Integration
1. **User Registration**: Creates user account in backend
2. **Auto-login**: Returns authentication tokens immediately
3. **Company Check**: Frontend checks if user has companies via `GET /companies`
4. **Onboarding Wizard**: If no companies, shows 3-step wizard
5. **Company Creation**: `POST /companies` creates first company with `is_default: true`
6. **Dashboard Redirect**: Automatically loads dashboard with new company

## Features

### ✅ Recently Implemented (Backend + Frontend Integration)
- **JWT Authentication**: PASETO-based token system with 15-minute access tokens and 24-hour refresh tokens
- **User Registration**: Streamlined registration with automatic login (single API call)
- **User Profiles**: Complete profile management with update capabilities
- **Session Management**: Secure session tracking with Redis cache
- **Company Management**: Full CRUD operations with simplified user-company ownership model
- **User-Company Association**: Direct ownership via `user_owner_id` field (no junction table)
- **Onboarding Flow**: Frontend wizard integrated with backend company creation
- **Database Integration**: Real company creation and retrieval from MySQL

### 🙧 Planned Core Features  
- **Person Profiles**: Comprehensive team member profile system
- **1:1 Meetings**: Meeting management and note-taking system
- **Feedback Tracking**: Direct and mentioned feedback collection
- **Member Get Member System**: Referral tracking with 50% discounts per valid referral

### 🚧 Planned AI Features
- **Contextual Memory**: Vector database integration for AI-powered insights
- **Smart Suggestions**: AI-driven action recommendations based on team context
- **Meeting Analysis**: Automated insights from 1:1 conversations
- **Performance Tracking**: AI-enhanced performance review capabilities

### 🛠️ Technical Features
- **Clean Architecture**: Separation of concerns with DDD principles
- **Comprehensive Testing**: Unit tests with testcontainers for integration testing
- **API Documentation**: Auto-generated Swagger/OpenAPI documentation
- **Observability**: Prometheus metrics, Grafana dashboards, Jaeger tracing

## Project Architecture
This boilerplate follows the principles of Clean Architecture, organizing code into distinct layers with specific responsibilities.

<div align="center">
    <img src='https://raw.githubusercontent.com/diegoclair/go-boilerplate/main/.github/assets/architecture.png' />
</div>

### Directory Structure
The project aims to follow standard Go project layout conventions where applicable. The main directories are:

*   **`/cmd`**: Contains the application entry point (`cmd/main.go`). Responsible for bootstrapping the application: loading configuration (`config.toml`), setting up dependency injection, and starting the HTTP server.
*   **`/internal`**: Houses the core application logic, not intended for import by external projects. It's organized by Clean Architecture layers:
    *   **`/internal/domain`**: The heart of the application. Contains business logic, entities, value objects, aggregates, and repository *interfaces*. It's independent of external concerns.
    *   **`/internal/application`**: Orchestrates use cases (interactors/services). Depends on `internal/domain` interfaces and defines interfaces for required infrastructure components (like repositories or caches).
    *   **`/internal/transport`**: Handles incoming requests (HTTP via `echo` framework in this case) and outgoing responses. It validates input, calls `internal/application` use cases, and formats the output. Depends on `internal/application`.
*   **`/infra`**: Provides concrete implementations for interfaces defined in `internal/domain` and `internal/application`. This includes database repositories (MySQL), caching (Redis), logging, etc. It depends on interfaces defined in `/internal` and external libraries/drivers.
*   **`/util`**: Contains common utility functions potentially used across different layers (use with caution to avoid creating unwanted coupling).
*   **`/migrator`**: Manages database schema migrations.
*   **`/docs`**: Stores generated API documentation (Swagger/OpenAPI generated by `make docs`).
*   **`/mocks`**: Contains mocks generated using `gomock` (`make mocks`) for testing domain and infrastructure interfaces.
*   **`/goswag`**: A helper tool used internally by the `make docs` command to facilitate Swagger generation.
*   **`.docker`**: Holds Docker volumes data for services like MySQL, Redis and etc.
*   **`.github`**: Contains GitHub Actions workflows for CI/CD and contribution templates (like issue/PR templates).

(Other files like `go.mod`, `go.sum`, `Dockerfile`, `docker-compose.yml`, `Makefile`, `config.toml`, `prometheus.yml`, `.gitignore`, `LICENSE` are standard configuration or project files.)

### Dependency Rule
A fundamental principle of Clean Architecture is the **Dependency Rule**: source code dependencies can only point inwards towards the core business logic.
*   Code within `/internal` follows this rule strictly: `transport` depends on `application`, which depends on `domain`. `domain` has no internal dependencies.
*   The `/infra` layer depends on *interfaces* defined within `internal/domain` and `internal/application`.
*   The `/cmd` layer orchestrates the setup and depends on concrete types from `/internal` and `/infra` during initialization.

This structure ensures that changes in outer layers (like `/infra`, `/transport`, UI frameworks, or databases) have minimal impact on the core business logic.

## 🔐 Authentication System

LeaderPro implements a robust JWT authentication system using PASETO tokens:

### Token Types
- **Access Token**: 15-minute expiration, used for API requests
- **Refresh Token**: 24-hour expiration, used to generate new access tokens

### Authentication Flow
1. **Registration**: `POST /users` - Creates user and automatically logs them in
2. **Login**: `POST /auth/login` - Returns both user data and authentication tokens
3. **Token Refresh**: `POST /auth/refresh` - Generates new access token using refresh token
4. **Protected Routes**: All `/users/profile` endpoints require valid access token

### Session Management
- Sessions are stored in Redis with refresh token mapping
- User-Agent and Client IP tracking for security
- Automatic session cleanup on token expiration

### Security Features
- PASETO symmetric key encryption
- Redis-based token blacklisting capability
- Request context includes user UUID and session UUID
- Middleware-based route protection

## 💻 Getting Started

### Prerequisites ❗
*   Ensure **Docker** is installed on your machine.
*   An installation of **Go 1.22** or later. See [Installing Go](https://go.dev/doc/install).
*   **(Optional)** Recommended VS Code Extensions: Go, Go Test Explorer, EnvFile.

### Configuration
The application loads configuration from `config.toml` by default. You can override settings using environment variables (prefixed with `APP_`). For local development, the `docker-compose.yml` file sets up necessary services (like MySQL, Redis) with default configurations. If you need to customize database credentials or other settings outside of Docker, you can modify `config.toml` or set environment variables.

### ▶️ Launching the Application
To start the application and its dependencies (MySQL, Redis, etc.) using Docker Compose, run the Make command:
```bash
make start
```
Wait for the logs to indicate the services are ready. You should see a message like `your server started on [::]:5000` from the Go application container (`myapp`).

### API Examples

Once the application is running, you can interact with the LeaderPro API:

#### 1. Health Check
```bash
curl -X GET http://localhost:5000/ping
```
```json
{
    "message": "pong"
}
```

#### 2. User Registration (Auto-login)
```bash
curl -X POST http://localhost:5000/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john.doe@example.com",
    "password": "password123",
    "phone": "+1234567890"
  }'
```
```json
{
  "user": {
    "uuid": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "phone": "+1234567890"
  },
  "auth": {
    "access_token": "v2.local.Gdh5kiOTyyaQ3_bNykYDeYHO21Jg2...",
    "access_token_expires_at": "2025-01-07T10:30:00Z",
    "refresh_token": "v2.local.H6Gdh5kiOTyyaQ3_bNykYDeYHO21...",
    "refresh_token_expires_at": "2025-01-08T10:15:00Z"
  }
}
```

#### 3. User Login
```bash
curl -X POST http://localhost:5000/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "password123"
  }'
```

#### 4. Get User Profile (Protected Route)
```bash
curl -X GET http://localhost:5000/users/profile \
  -H "user-token: v2.local.Gdh5kiOTyyaQ3_bNykYDeYHO21Jg2..."
```
```json
{
  "uuid": "550e8400-e29b-41d4-a716-446655440000",
  "name": "John Doe",
  "email": "john.doe@example.com",
  "phone": "+1234567890"
}
```

#### 5. Create Company (Onboarding)
```bash
curl -X POST http://localhost:5000/companies \
  -H "Content-Type: application/json" \
  -H "user-token: v2.local.Gdh5kiOTyyaQ3_bNykYDeYHO21Jg2..." \
  -d '{
    "name": "Tech Startup",
    "industry": "technology",
    "size": "6-15",
    "role": "CTO",
    "is_default": true
  }'
```
```json
{
  "uuid": "company-uuid-here",
  "name": "Tech Startup",
  "industry": "technology",
  "size": "6-15",
  "role": "CTO",
  "is_default": true,
  "created_at": "2025-01-07T10:15:00Z"
}
```

#### 6. List User Companies
```bash
curl -X GET http://localhost:5000/companies \
  -H "user-token: v2.local.Gdh5kiOTyyaQ3_bNykYDeYHO21Jg2..."
```
```json
[
  {
    "uuid": "company-uuid-here",
    "name": "Tech Startup",
    "industry": "technology",
    "size": "6-15",
    "role": "CTO",
    "is_default": true,
    "created_at": "2025-01-07T10:15:00Z"
  }
]
```


## Testing 🧪

### Unit Tests
Unit tests are crucial for ensuring code correctness and maintainability.
*   **Real Dependencies with Testcontainers:** For tests involving databases (MySQL) or caches (Redis), we use [Testcontainers](https://testcontainers.com/). This library spins up real dependencies in Docker containers specifically for the test run, providing a high-fidelity testing environment without mocking the database itself.
*   **Layer Isolation:** Tests are written to respect architectural boundaries. For example, when testing `application` layer use cases, `infra` layer components (like repositories) are typically replaced with mocks.

### Mocks
Mocks are generated using `gomock` to isolate components during testing. They are stored in the `/mocks` directory. To regenerate mocks after modifying interfaces (typically within `/internal/domain/contract` or `/infra/contract` based on the `Makefile`):
```bash
make mocks
```
This command ensures `mockgen` is installed and generates mocks based on the interfaces found in specified contract directories.

### Running Tests
To run all tests (unit and integration tests using testcontainers):
```bash
make tests
```
This command executes `go test` with coverage analysis across all modules.

## 📝 API Documentation
API documentation is generated in Swagger/OpenAPI format and served directly by the application.

*   **Accessing:** Once the application is running (using `make start`), you can access the interactive Swagger UI documentation in your browser at:
    [`http://localhost:5000/swagger/`](http://localhost:5000/swagger/)
*   **Generation:** The documentation is generated automatically from code annotations (specifically in the `/internal/transport` layer handlers) using [goswag](https://github.com/diegoclair/goswag), an open-source tool developed for this purpose.

### Generating Docs
To regenerate the API documentation (Swagger/OpenAPI) using `swag` and the `goswag` helper after making changes to handlers or annotations in the `/internal/transport` layer:
```bash
make docs
```
This updates the files in the `/docs` directory.

## 🤝 Contributing
Contributions are welcome! Improving LeaderPro helps the leadership development community. Here are ways you can contribute:

-   **Submit Pull Requests**: Review open PRs, and submit your own enhancements or bug fixes.
-   **[Join the Discussions](https://github.com/diegoclair/leaderpro/discussions)**: Share insights, provide feedback, or ask questions.
-   **[Report Issues](https://github.com/diegoclair/leaderpro/issues)**: Report bugs or suggest new features.

<details closed>
    <summary>Contributing Guidelines</summary>

1.  **Fork the Repository**: Start by forking the project repository to your GitHub account.
2.  **Clone Locally**: Clone the forked repository to your local machine.
    ```sh
    git clone https://github.com/<your_username>/leaderpro
    ```
3.  **Create a New Branch**: Use a descriptive branch name.
    ```sh
    git checkout -b feature/describe-your-feature
    ```
4.  **Make Changes**: Implement your feature or bug fix following the existing code patterns.
5.  **Run Tests**: Ensure all tests pass before submitting.
    ```sh
    make tests
    ```
6.  **Update Documentation**: Update README or add API documentation if needed.
7.  **Commit Your Changes**: Write clear, descriptive commit messages.
    ```sh
    git commit -m "feat: add new leadership insight feature"
    ```
8.  **Push and Submit PR**: Push your branch and create a pull request.
    ```sh
    git push origin feature/describe-your-feature
    ```

</details>

## 📄 License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.