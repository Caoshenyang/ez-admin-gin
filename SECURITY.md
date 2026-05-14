# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.1.x   | :white_check_mark: |
| 1.0.x   | :x:                |
| < 1.0   | :x:                |

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

## Security Best Practices

When deploying EZ Admin Gin to production, make sure to follow these security requirements:

- **JWT Secret**: Replace the default secret with a random string (at least 32 characters). Use `openssl rand -hex 32`.
- **Database Password**: Use a strong, unique password for PostgreSQL/MySQL.
- **Redis Password**: Set a Redis password in production.
- **CORS Origins**: Set `EZ_CORS_ALLOWED_ORIGINS` to your actual domain. Never use `*` in production.
- **HTTPS**: Always enable HTTPS in production.
- **Swagger**: Disable Swagger UI in production (`EZ_SWAGGER_ENABLED=false`).
- **Default Password**: Change the default admin password immediately after first login.

See the [Production Checklist](docs/deployment/production-checklist.md) for a complete pre-launch security checklist.

## Security Features

EZ Admin Gin includes the following built-in security measures:

- **Dual Token Authentication**: Access token (stateless) + Refresh token (HttpOnly Secure Cookie with rotation)
- **Token Blacklisting**: Server-side session revocation via Redis
- **RBAC + Casbin**: Role-based access control with fine-grained API permission enforcement
- **Data Scope**: Row-level data permissions (all / dept / dept_and_children / custom_dept / self)
- **Security Headers**: X-Content-Type-Options, X-Frame-Options, Referrer-Policy, HSTS
- **Upload Security**: File extension whitelist, double-extension rejection, path traversal prevention
- **Login Rate Limiting**: Account-level lockout after failed attempts
- **Production Config Validation**: Refuses to start if insecure defaults are detected in production
