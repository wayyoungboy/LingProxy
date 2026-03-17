# Contributing to LingProxy

First off, thank you for considering contributing to LingProxy! It's people like you that make LingProxy such a great tool.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Coding Standards](#coding-standards)
- [Commit Guidelines](#commit-guidelines)
- [Pull Request Process](#pull-request-process)
- [Reporting Issues](#reporting-issues)

## Code of Conduct

This project and everyone participating in it is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## Getting Started

### Prerequisites

- **Go**: Version 1.21 or higher
- **Node.js**: Version 18 or higher
- **npm** or **yarn**: For frontend dependencies
- **Git**: For version control
- **Make**: For using the Makefile commands

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/lingproxy.git
   cd lingproxy
   ```
3. Add the upstream repository:
   ```bash
   git remote add upstream https://github.com/lingproxy/lingproxy.git
   ```

## Development Setup

### Backend Setup

```bash
# Navigate to backend directory
cd backend

# Install Go dependencies
go mod download
go mod tidy

# Run the backend server
go run cmd/main.go
```

The backend will start at `http://localhost:8080`.

### Frontend Setup

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies
npm install

# Run development server
npm run dev
```

The frontend will be available at `http://localhost:3000`.

### Configuration

1. Copy the example configuration:
   ```bash
   cp backend/configs/config.yaml.example backend/configs/config.yaml
   ```
2. Edit `config.yaml` with your settings
3. **Important**: Change the default admin password before running in production!

## How to Contribute

### Types of Contributions

We welcome many types of contributions:

- **Bug fixes**: Fix issues in the codebase
- **New features**: Add new functionality
- **Documentation**: Improve or translate documentation
- **Tests**: Add or improve test coverage
- **Code refactoring**: Improve code quality without changing behavior
- **Performance improvements**: Make the code faster or more efficient

### Before You Start

1. Check if there's an existing issue or discussion related to your contribution
2. For significant changes, open an issue first to discuss the approach
3. Make sure your changes align with the project's goals

## Coding Standards

### Go Code

- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines
- Run `gofmt` before committing: `gofmt -s -w .`
- Run `go vet` to catch common errors: `go vet ./...`
- Use meaningful variable and function names
- Add comments for exported functions and types
- Write unit tests for new functionality

### Frontend Code (Vue.js)

- Follow the [Vue.js Style Guide](https://vuejs.org/style-guide/)
- Use Composition API with `<script setup>` syntax
- Keep components small and focused
- Use TypeScript-style JSDoc comments for complex logic
- Follow the existing project structure

### General Guidelines

- Keep lines under 100 characters when possible
- Use meaningful commit messages
- Update documentation when changing behavior
- Add tests for new features

## Commit Guidelines

We follow [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

### Types

- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `perf`: A code change that improves performance
- `test`: Adding missing tests or correcting existing tests
- `chore`: Changes to the build process or auxiliary tools

### Examples

```
feat(api): add support for custom retry policies
fix(proxy): resolve memory leak in connection pool
docs(readme): update installation instructions
test(handler): add unit tests for API key validation
```

## Pull Request Process

1. **Create a branch** from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following the coding standards above

3. **Write/update tests** for your changes

4. **Update documentation** if needed

5. **Run tests and linters**:
   ```bash
   # Backend tests
   cd backend && go test -v -race ./...

   # Backend linting
   golangci-lint run

   # Frontend linting
   cd frontend && npm run lint
   ```

6. **Commit your changes** following the commit guidelines

7. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

8. **Open a Pull Request** on GitHub

### PR Checklist

- [ ] Code compiles correctly
- [ ] All tests pass
- [ ] New code has test coverage
- [ ] Documentation is updated
- [ ] Commit messages follow the guidelines
- [ ] Branch is up to date with `main`

### Review Process

1. At least one maintainer must approve the PR
2. All CI checks must pass
3. No merge conflicts
4. PR must be open for at least 24 hours for non-trivial changes (exceptions for minor fixes)

## Reporting Issues

### Bug Reports

When reporting a bug, please include:

1. **Description**: Clear description of the bug
2. **Steps to reproduce**: Detailed steps to reproduce the issue
3. **Expected behavior**: What you expected to happen
4. **Actual behavior**: What actually happened
5. **Environment**:
   - OS and version
   - Go version
   - Node.js version
   - LingProxy version
6. **Logs**: Relevant log output
7. **Screenshots**: If applicable

### Feature Requests

For feature requests, please include:

1. **Problem**: What problem does this feature solve?
2. **Solution**: Your proposed solution
3. **Alternatives**: Other solutions you've considered
4. **Impact**: How would this benefit other users?

## Getting Help

- Open a [GitHub Discussion](https://github.com/lingproxy/lingproxy/discussions) for questions
- Join our community (links to be added)
- Check existing documentation in the `docs/` directory

## Recognition

Contributors are recognized in our README and release notes. We appreciate every contribution, no matter how small!

---

Thank you for contributing to LingProxy! 🎉