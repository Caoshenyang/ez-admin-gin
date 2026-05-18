# Security Policy

## Reporting a Vulnerability

If you discover a security vulnerability, please report it responsibly:

1. **Do not** open a public GitHub Issue for security vulnerabilities.
2. Send an email to the maintainer or use [GitHub Security Advisories](https://github.com/Caoshenyang/ez-admin-gin/security/advisories/new).
3. Include the following information:
   - Vulnerability type (e.g., XSS, SQL injection, authentication bypass)
   - Steps to reproduce
   - Affected version
   - Potential impact
   - Any suggested fix

We aim to respond within **48 hours** and provide a fix or mitigation within **7 days** for confirmed vulnerabilities.

## Production Security Checklist

When deploying EZ Admin Gin to production, you **must** complete these steps:

- **JWT Secret**: Replace the default secret with a random string (at least 32 characters). Use `openssl rand -hex 32`. Set via `EZ_AUTH_JWT_SECRET`.
- **Database Password**: Use a strong, unique password for PostgreSQL/MySQL.
- **Redis Password**: Set a Redis password in production.
- **CORS Origins**: Set `EZ_CORS_ALLOWED_ORIGINS` to your actual domain. The server will refuse to start in production if this is empty.
- **HTTPS**: Always enable HTTPS in production.
- **Swagger**: Disable Swagger UI in production (`EZ_SWAGGER_ENABLED=false`). The server warns if Swagger is enabled in production mode.
- **Default Password**: The setup script creates an admin account with `EzAdmin@123456`. Change it immediately after first login.

The server includes **production config validation** that refuses to start if insecure defaults are detected (JWT secret not changed, CORS origins empty, etc.).

See the [Production Checklist](https://caoshenyang.github.io/ez-admin-gin/deployment/production-checklist) for a complete pre-launch checklist.

## Security Features

### Implemented

| Feature | Mechanism | Notes |
|---------|-----------|-------|
| JWT Authentication | HS256 signed tokens | Access token sent via `Authorization: Bearer` header |
| Refresh Token | HttpOnly cookie + Redis storage with rotation | Cookie: `ez_admin_refresh_token`, path `/api/v1/auth`, HttpOnly=true, Secure=prod |
| Token Blacklisting | Redis-based session revocation | On logout, refresh token is deleted from Redis |
| RBAC | Casbin policy engine | Role × URL pattern × HTTP method matching |
| Data Scope | Row-level data permissions in Repository layer | 5 levels: all / dept / dept_and_children / self / custom_dept |
| Security Headers | `X-Content-Type-Options`, `X-Frame-Options`, `Referrer-Policy`, `X-XSS-Protection`, HSTS (prod) | Applied via middleware |
| Upload Security | Extension whitelist, double-extension rejection, size limit, path traversal prevention | Configurable via `upload.allowed_exts` |
| Login Rate Limiting | Per-IP + per-account lockout | Redis sliding window; configurable threshold and duration |
| Production Validation | Refuses to start with insecure defaults | Checks JWT secret, CORS origins, Swagger toggle |

### Implementation Details

**Token storage and transmission:**

- Access token: Stored in browser `localStorage` / `sessionStorage` (user chooses "remember me"), sent as `Authorization: Bearer <token>` header.
- Refresh token: Stored as HttpOnly cookie by the server. The frontend never reads this cookie — it's sent automatically via `withCredentials: true` during token refresh requests.
- Token refresh uses rotation: each refresh issues a new refresh token and invalidates the old one.

**Rate limiting configuration:**

| Setting | Default | Environment Variable |
|---------|---------|---------------------|
| Per-IP limit | 10 req/60s | `EZ_RATE_LIMIT_LOGIN_MAX_REQUESTS` |
| Account lockout threshold | 5 failures | `EZ_RATE_LIMIT_LOGIN_LOCKOUT_THRESHOLD` |
| Account lockout duration | 300s | `EZ_RATE_LIMIT_LOGIN_LOCKOUT_SEC` |

**CORS behavior:**

- Development mode: `localhost` origins are automatically allowed.
- Production mode: `CORS_ALLOWED_ORIGINS` must be explicitly set; wildcard `*` is rejected.

### Planned / Not Yet Implemented

| Feature | Status |
|---------|--------|
| Content Security Policy (CSP) headers | Not implemented |
| Account lockout notification | Not implemented |
| Password complexity enforcement | Not enforced server-side (relies on frontend validation) |
| Two-factor authentication (2FA) | Not implemented |
| Audit log tamper protection | Not implemented (operation logs stored in database) |
