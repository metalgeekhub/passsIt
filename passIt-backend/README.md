# passIt backend

A modern password management and sharing platform built with Go, Keycloak, and PostgreSQL.

## Table of Contents
- [Getting Started](#getting-started)
- [Prerequisites](#prerequisites)
- [Database Migration](#database-migration)
- [Local Development](#local-development)
- [Makefile Commands](#makefile-commands)
- [Project Structure](#project-structure)
- [Environment Variables](#environment-variables)
- [PassIt Authentication](#authentication)
- [API Authentication Guide](API_AUTHENTICATION.md) - Using the API with different clients
- [License](#license)

## Getting Started

Follow these instructions to set up the project for local development and testing.

## Prerequisites
- Go 1.24.0 or higher
- Docker & Docker Compose

## Database Migration

Run database migrations:
```bash
go run ./internal/database/migration/migration.go
```

## Local Development

### Partial development

With docker-compose.yml file a functional application will start
Start Core Services:
- Database
- Keycloak
- Redis
```bash
docker-compose up -d
```
Run the Go application:
```bash
make run
```

### Full developement

Start all Services:
- Database
- Keycloak
- Redis
- Elasticsearch
- Kibana(ElasticSearch UI)
- Filebeat
- PassIt backend(If there is no need for that you can comment it out)
```bash
docker-compose -f docker-compose-full.yml up --build -d
```

## MakeFile

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```
Create DB container
```bash
make docker-run
```

Shutdown DB Container
```bash
make docker-down
```

DB Integrations Test:
```bash
make itest
```

Live reload the application:
```bash
make watch
```

Run the test suite:
```bash
make test
```

Run specific package tests:
```bash
# Test utilities (100% coverage)
go test ./internal/utils -v

# Test all packages
go test ./internal/... -v
```

**Current Test Coverage**:
- ✅ Utils: 100% coverage (DecodeServerInput, InsecureHttpContext)
- ✅ Models: User struct validation, UUID generation
- ✅ Response Codes: Code definitions and uniqueness
- ✅ Configuration: Environment variable handling
- ✅ Server: HTTP helpers and status codes

Clean up binary from the last build:
```bash
make clean
```

## Environment Variables
An .env file in the root directory is needed with the following variables:
```bash
PORT=8080
APP_ENV=local

# Database
DB_HOST=localhost
DB_PORT=5432
DB_DATABASE=passit
DB_USERNAME=melkey
DB_PASSWORD=password1234
DB_SCHEMA=public

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_DATABASE=0
REDIS_USERNAME=default
REDIS_PASSWORD=redis_password

# Keycloak
KEYCLOAK_URL=https://localhost:8443
KEYCLOAK_REALM=passit
KEYCLOAK_CLIENT_ID=passit-backend
KEYCLOAK_CLIENT_SECRET=your_client_secret
KEYCLOAK_ADMIN_USERNAME=passit-admin
KEYCLOAK_ADMIN_PASSWORD=yP7!vR2qLz#eX9sT
REDIRECT_URL=http://localhost:8080/auth/callback
FRONTEND_URL=http://localhost:3000
```

## Architecture

### Service Layer Pattern
The backend follows a clean architecture with separation of concerns:

- **Handlers** (`internal/server/`): HTTP layer - request validation, response formatting
- **Services** (`internal/services/`): Business logic, orchestration, transaction management
- **Database** (`internal/database/`): Data access layer using GORM
- **Auth** (`internal/auth/`): Keycloak integration with interface-based design
- **Middleware** (`internal/middleware/`): Authentication, logging, CORS

### Key Features
- **Service Layer**: Encapsulates business logic with automatic rollback on failures
- **Interface-Based Design**: Loose coupling for easier testing and mocking
- **Dual Authentication**: Supports both session cookies (browser) and Bearer tokens (API clients)
- **PostgreSQL as Source of Truth**: User data stored in PostgreSQL, Keycloak for authentication only

## Authentication

This project uses [Keycloak](https://www.keycloak.org/) as the authentication provider.  
Keycloak handles user login, registration, and token issuance for secure access to protected resources.

- **Login Flow:**  
  The authentication process with Keycloak involves several steps and handlers that work together to secure the application. Users are redirected to Keycloak for authentication. Upon successful login, Keycloak issues tokens that are used to access protected endpoints in the backend.

  ![Login Flow Diagram](../docs/keycloak_auth_flow.png)

  The key points:
  - /auth/login - Initiates the authentication process
  - /auth/callback - Handles the OAuth2 callback from Keycloak

- **Protected Routes:**  
  Most backend routes require a valid access token. The backend verifies tokens using Keycloak’s OIDC endpoints.

- **Session Management:**  
  - **OAuth State**: Stored in secure signed cookies (no Redis needed for OAuth flow)
  - **User Sessions**: Stored in Redis after successful login for scalability and security
  - Session data includes access tokens, user info, and expires after 24 hours

- **Configuration:**  
  Keycloak connection details (URL, realm, client ID, client secret) are set via environment variables in your `.env` file:
  ```env
  KEYCLOAK_URL=https://localhost:8443
  KEYCLOAK_REALM=passit
  KEYCLOAK_CLIENT_ID=passit-backend
  KEYCLOAK_CLIENT_SECRET=your_client_secret
  KEYCLOAK_ADMIN_USERNAME=passit-admin
  KEYCLOAK_ADMIN_PASSWORD=your_admin_password
  ```

- **Login Page:**  
  The login page is served at `/` and provides a "Login with Keycloak" button, which redirects users to the Keycloak login screen.

- **Customizing Authentication:**  
  You can adjust authentication logic in `internal/middleware/auth.go` and route configuration in `internal/server/routes.go`.

For more details on configuring Keycloak, see the [Keycloak documentation](https://www.keycloak.org/documentation).

## API Documentation (Swagger)

The backend provides interactive API documentation using Swagger/OpenAPI.

### Accessing Swagger UI

1. Start the backend server:
   ```bash
   make run
   ```

2. Open your browser and navigate to:
   ```
   http://localhost:8080/swagger/index.html
   ```

### Using Swagger with Authentication

Since the API uses Keycloak for authentication, you need a Bearer token to test protected endpoints in Swagger.

**Steps to get a Bearer token:**

1. **Enable Direct Access Grants in Keycloak:**
   - Go to Keycloak Admin Console: https://localhost:8443
   - Login with admin credentials
   - Select **passit** realm
   - Go to **Clients** → **passit-backend**
   - In **Settings** tab, enable **"Direct access grants"**
   - Disable **"Consent required"** if enabled
   - Click **Save**

2. **Get Access Token using PowerShell:**
   ```powershell
   $body = @{
       client_id = "passit-backend"
       client_secret = "79VK0WOudj"
       grant_type = "password"
       username = "admin@passit.com"
       password = "YourSecurePassword123!"
   }

   $response = Invoke-RestMethod -Uri "https://localhost:8443/realms/passit/protocol/openid-connect/token" -Method Post -Body $body -ContentType "application/x-www-form-urlencoded" -SkipCertificateCheck
   
   # Copy the access token
   $response.access_token
   ```

   Or using curl:
   ```bash
   curl -X POST "https://localhost:8443/realms/passit/protocol/openid-connect/token" \
     -H "Content-Type: application/x-www-form-urlencoded" \
     -d "client_id=passit-backend" \
     -d "client_secret=79VK0WOudj" \
     -d "grant_type=password" \
     -d "username=admin@passit.com" \
     -d "password=YourSecurePassword123!" \
     --insecure
   ```

3. **Use Token in Swagger:**
   - Click the **Authorize** button (green lock icon) in Swagger UI
   - Paste the access token (without "Bearer" prefix)
   - Click **Authorize**
   - Click **Close**
   - Now you can test all protected endpoints

### Regenerating Swagger Documentation

After adding new endpoints or updating annotations, regenerate the docs:

```bash
swag init -g cmd/api/main.go
```

This will update the `docs/` folder with the latest API specifications.

