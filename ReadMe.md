# PassIt Application

Event ticketing platform with user management, authentication via Keycloak, and role-based access control.

## Quick Start

### Development

**Backend:**
```bash
cd passIt-backend
go run cmd/api/main.go
```

**Frontend:**
```bash
cd passIt-ui
npm run dev
```

### Production Build

**Backend:**
```bash
cd passIt-backend
make build-prod  # Creates passit-server.exe
```

**Frontend:**
```bash
cd passIt-ui
npm run build  # Creates dist/ folder
```

## Documentation

- **Backend Details:** [passIt-backend](./passIt-backend)
- **Frontend Details:** [passIt-ui](./passIt-ui)
- **Deployment Guide:** [DEPLOYMENT.md](./DEPLOYMENT.md)

## Key Features

- ✅ User authentication with Keycloak OAuth2/OIDC
- ✅ Role-based access control (Admin/User)
- ✅ User profile management
- ✅ Soft delete with separate deleted users view
- ✅ Session management with Redis
- ✅ RESTful API with Swagger documentation (dev only)
- ✅ Responsive Vue 3 frontend with TailwindCSS
- ✅ Production-ready builds

## Tech Stack

**Backend:**
- Go 1.23+ with Gin framework
- PostgreSQL database
- Redis for sessions
- Keycloak for authentication
- GORM ORM
- Swagger/OpenAPI

**Frontend:**
- Vue 3 with TypeScript
- Vue Router for SPA routing
- Axios for API calls
- TailwindCSS for styling
- Vite for build tooling

## Environment Configuration

### Development
- Backend: `ENV=development` (enables Swagger)
- Frontend: Development server on `localhost:3000`

### Production
- Backend: `ENV=production` (disables Swagger, release mode)
- Frontend: Static build in `dist/` folder

See [DEPLOYMENT.md](./DEPLOYMENT.md) for complete deployment instructions.