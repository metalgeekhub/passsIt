# passIt-ui

A modern Vue 3 frontend for the PassIt password management and sharing platform.

## Features

- 🔐 **OAuth2 Authentication** via Keycloak
- 📦 **Type-Safe API Client** with TypeScript
- 🎨 **Tailwind CSS** for styling
- ⚡ **Vite** for fast development
- 🧩 **Modular Architecture** with composables and services

## Recommended IDE Setup

[VSCode](https://code.visualstudio.com/) + [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) (and disable Vetur).

## Type Support for `.vue` Imports in TS

TypeScript cannot handle type information for `.vue` imports by default, so we replace the `tsc` CLI with `vue-tsc` for type checking. In editors, we need [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) to make the TypeScript language service aware of `.vue` types.

## Environment Configuration

Create a `.env` file in the root directory (see `.env.example`):

```env
# Backend API URL
VITE_API_BASE_URL=http://localhost:8080

# Frontend Port
VITE_APP_PORT=3000
```

**Important**: All environment variables must be prefixed with `VITE_` to be exposed to the client.

## Project Structure

```
src/
├── api/
│   ├── axios.ts       # Axios instances (public & authenticated)
│   └── services.ts    # API service layer (userApi, authApi, healthApi)
├── components/        # Reusable Vue components
├── composables/       # Vue composables (useAuth, etc.)
├── config/
│   └── env.ts        # Environment configuration
├── models/           # TypeScript interfaces
├── router/           # Vue Router configuration
├── views/            # Page components
└── main.ts           # Application entry point
```

## API Service Usage

The application provides a clean API service layer:

```typescript
import { userApi, authApi, healthApi } from '@/api/services';

// User operations
const users = await userApi.getAll();
const user = await userApi.getById('uuid');
const user = await userApi.getByEmail('user@example.com');
await userApi.create({ email, username, firstName, lastName, password });
await userApi.update('uuid', { firstName: 'John' });

// Authentication
const isAuth = await authApi.checkAuth();
authApi.login();  // Redirect to Keycloak login
authApi.logout(); // Redirect to Keycloak logout

// Health check
const health = await healthApi.check();
```

## Customize configuration

See [Vite Configuration Reference](https://vite.dev/config/).

## Project Setup

```sh
npm install
```

### Compile and Hot-Reload for Development

```sh
npm run dev
```

### Type-Check, Compile and Minify for Production

```sh
npm run build
```

### Run Unit Tests with [Vitest](https://vitest.dev/)

```sh
npm run test:unit
```

## Docker

### Development
```sh
docker build -f Dockerfile_development -t passit-ui:dev .
docker run -p 3000:3000 -v $(pwd):/app passit-ui:dev
```

### Production
```sh
docker build -t passit-ui:latest .
docker run -p 80:80 passit-ui:latest
```

## Backend Integration

This frontend connects to the PassIt backend API. Make sure:

1. Backend is running on `http://localhost:8080` (or update `VITE_API_BASE_URL`)
2. Keycloak is configured and running
3. CORS is enabled for `http://localhost:3000` in backend

See [Backend Documentation](../passIt-backend/README.md) for setup instructions.
