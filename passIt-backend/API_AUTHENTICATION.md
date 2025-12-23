# API Authentication Guide

> 📚 **Main Documentation**: See [README.md](README.md) for general setup and configuration.

Your PassIt API now supports **two authentication methods**:

## 1. Browser-Based Authentication (OAuth Flow)

**For**: Vue.js frontend, web applications

**Flow**:
1. User visits `/auth/login` → Redirects to Keycloak
2. User logs in → Keycloak redirects to `/auth/callback`
3. Backend sets `session_id` cookie
4. Frontend makes requests with cookie automatically included

**Usage**:
```javascript
// Frontend automatically sends cookie
fetch('http://localhost:8080/api/users', {
  credentials: 'include'  // Include cookies
})
```

## 2. API Token Authentication (Bearer Token)

**For**: Mobile apps, CLI tools, backend services, third-party integrations

**Flow**:
1. Get access token from Keycloak (direct authentication)
2. Send token in `Authorization` header with each request

**Usage**:
```bash
# Get token from Keycloak directly
curl -X POST "https://localhost:8443/realms/passit/protocol/openid-connect/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=password" \
  -d "client_id=passit-backend" \
  -d "client_secret=YOUR_CLIENT_SECRET" \
  -d "username=john.doe" \
  -d "password=SecurePassword123!"

# Response includes access_token
# {
#   "access_token": "eyJhbGc...",
#   "expires_in": 300,
#   "refresh_token": "...",
#   ...
# }

# Use the token in API requests
curl -X GET "http://localhost:8080/api/users" \
  -H "Authorization: Bearer eyJhbGc..."
```

**Example in Python**:
```python
import requests

# Get token
token_response = requests.post(
    'https://localhost:8443/realms/passit/protocol/openid-connect/token',
    data={
        'grant_type': 'password',
        'client_id': 'passit-backend',
        'client_secret': 'YOUR_CLIENT_SECRET',
        'username': 'john.doe',
        'password': 'SecurePassword123!'
    },
    verify=False  # Only for self-signed certs
)
access_token = token_response.json()['access_token']

# Make API request
response = requests.get(
    'http://localhost:8080/api/users',
    headers={'Authorization': f'Bearer {access_token}'}
)
```

## API Endpoints

## API Routes

### Public Endpoints (No auth required)
- `GET /` - Login page
- `GET /health` - Health check

### OAuth Flow (Browser-based)
- `GET /auth/login` - Start OAuth login (redirects to Keycloak)
- `GET /auth/callback` - OAuth callback (sets session cookie)
- `GET /auth/logout` - Logout (clears session)

### Protected API Endpoints
All endpoints under `/api/*` require authentication (either session cookie or Bearer token):

**User Management**:
- `POST /api/users` - Create user (requires valid user data)
- `GET /api/users` - Get all users
- `GET /api/users/find?id=<uuid>` - Find user by ID
- `GET /api/users/by-email?email=<email>` - Find user by email
- `PUT /api/users/:id` - Update user by ID

## Keycloak Configuration

For **Password Grant** (direct token authentication), you need to enable it in Keycloak:

1. Go to Keycloak Admin Console
2. Select your realm (`passit`)
3. Go to Clients → `passit-backend`
4. Enable "Direct Access Grants" (Resource Owner Password Credentials)
5. Save

## Security Notes

- Browser clients use **secure, httpOnly cookies** (protected from XSS)
- API clients use **Bearer tokens** in headers
- Both methods validate tokens using Keycloak's OIDC provider
- Tokens are signed and verified using JWT
- **OAuth State**: Stored in secure signed cookies (5 min TTL, SameSite=Lax)
- **User Sessions**: Stored in Redis with 24-hour TTL for scalability
- Cookie security: `httpOnly=true`, `secure=false` (localhost), `SameSite=Lax`

## Why This Architecture?

**OAuth State in Cookies**:
- Simpler than Redis for temporary state (5 minutes)
- Signed cookies prevent tampering
- No cross-server sharing needed
- Automatic cleanup via cookie expiration

**Sessions in Redis**:
- Sensitive data (access tokens) not in cookies
- Scalable across multiple backend instances
- Centralized session invalidation on logout
- Automatic cleanup with TTL
