# Security Policy

## Supported Versions

We release patches for security vulnerabilities. Currently supported versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take the security of LingProxy seriously. If you have discovered a security vulnerability, we appreciate your help in disclosing it to us in a responsible manner.

### How to Report

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please report them via one of the following methods:

1. **GitHub Security Advisory** (Preferred)
   - Go to the [Security Advisories](https://github.com/lingproxy/lingproxy/security/advisories) page
   - Click "Report a vulnerability"
   - Fill out the form with details about the vulnerability

2. **Email**
   - Send an email to the project maintainers
   - Include "SECURITY" in the subject line
   - Provide detailed information about the vulnerability

### What to Include

Please include the following information in your report:

- Type of vulnerability (e.g., injection, authentication bypass, etc.)
- Full paths of source file(s) related to the vulnerability
- The location of the affected source code (tag/branch/commit or direct URL)
- Any special configuration required to reproduce the issue
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the issue, including how an attacker might exploit it

### Response Timeline

We will respond to security reports according to the following timeline:

| Timeframe | Action |
|-----------|--------|
| 24 hours | Initial response acknowledging receipt |
| 72 hours | Detailed response with next steps |
| 7 days | Target timeframe for fix development |
| 14 days | Maximum timeframe for complex issues |

### Disclosure Policy

- We will acknowledge your email within 24 hours
- We will confirm the vulnerability and determine its severity
- We will work on a fix and release it as soon as possible
- We will notify you when the fix is released
- We will credit you in the security advisory (unless you prefer to remain anonymous)

### Security Best Practices

When deploying LingProxy, we recommend:

1. **Change Default Credentials**
   - Always change the default admin password before deploying to production
   - Use strong, unique passwords

2. **Secure Configuration**
   - Enable authentication in production environments
   - Use HTTPS/TLS for all communications
   - Configure CORS appropriately for your environment

3. **API Key Management**
   - Store API keys securely
   - Rotate API keys periodically
   - Use policies to limit API key permissions

4. **Network Security**
   - Deploy behind a reverse proxy when possible
   - Use firewalls to restrict access
   - Consider using VPN for internal deployments

5. **Regular Updates**
   - Keep LingProxy updated to the latest version
   - Monitor security advisories for updates

### Known Security Considerations

#### API Key Storage

API keys are encrypted at rest using AES-256 encryption. The encryption key is derived from a master key configured in the application. Ensure:

- The master key is kept secure
- The master key is rotated periodically
- Access to configuration files is restricted

#### Password Storage

Admin passwords are hashed using bcrypt before storage. We recommend:

- Using strong passwords (minimum 12 characters)
- Enabling rate limiting to prevent brute-force attacks

#### Request Logging

Request logs may contain sensitive information. We recommend:

- Implementing log rotation
- Securing log file access
- Complying with data retention policies

## Security Updates

Security updates will be released as patch versions (e.g., 1.0.1, 1.0.2). We will:

- Publish security advisories on GitHub
- Include security fixes in release notes
- Credit reporters (with permission)

## Comments on this Policy

If you have suggestions on how this policy could be improved, please submit a pull request or open an issue.

---

Thank you for helping keep LingProxy and its users safe!