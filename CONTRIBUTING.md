# Contributing to WhereGo

First off, thank you for considering contributing to WhereGo! 🎉

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Commit Messages](#commit-messages)

## Code of Conduct

This project and everyone participating in it is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/WhereGo.git`
3. Add the upstream remote: `git remote add upstream https://github.com/gustavosett/WhereGo.git`
4. Create a new branch: `git checkout -b feature/your-feature-name`

## Development Setup

### Prerequisites

- Go 1.24+
- Docker (optional, for container builds)
- Make (optional, for build automation)

### Building

```bash
# Build the binary
make build

# Run tests
make test

# Run linter
make lint

# Run all checks
make all
```

### Running Locally

```bash
# Run the API server
go run ./cmd/api

# Or use make
make run
```

## How to Contribute

### Reporting Bugs

Before creating bug reports, please check existing issues to avoid duplicates.

When creating a bug report, include:

- **Clear title** describing the issue
- **Steps to reproduce** the behavior
- **Expected behavior** vs actual behavior
- **Environment details** (OS, Go version, etc.)
- **Logs or error messages** if applicable

### Suggesting Features

Feature requests are welcome! Please include:

- **Use case** - Why is this feature needed?
- **Proposed solution** - How do you envision it working?
- **Alternatives considered** - Other approaches you've thought of

### Code Contributions

1. Look for issues labeled `good first issue` or `help wanted`
2. Comment on the issue to let others know you're working on it
3. Submit a PR when ready

## Pull Request Process

1. **Update documentation** if you're changing functionality
2. **Add tests** for new features or bug fixes
3. **Ensure all tests pass**: `make test`
4. **Ensure linter passes**: `make lint`
5. **Update CHANGELOG.md** with your changes
6. **Request review** from maintainers

### PR Checklist

- [ ] My code follows the project's coding standards
- [ ] I have added tests that prove my fix/feature works
- [ ] All new and existing tests pass
- [ ] I have updated the documentation accordingly
- [ ] I have added an entry to CHANGELOG.md

## Coding Standards

### Go Style

- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use `gofmt` for formatting (automatic with most editors)
- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Naming Conventions

- Use camelCase for unexported identifiers
- Use PascalCase for exported identifiers
- Use descriptive names over abbreviations
- Acronyms should be all caps: `HTTPHandler`, `IPAddress`

### Error Handling

```go
// Good
if err != nil {
    return fmt.Errorf("failed to parse IP: %w", err)
}

// Avoid
if err != nil {
    return err // No context
}
```

### Testing

- Write table-driven tests when appropriate
- Use meaningful test names: `TestLookup_InvalidIP_ReturnsError`
- Aim for high test coverage on critical paths

## Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>: <description>

[optional body]

[optional footer]
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `style`: Formatting, missing semicolons, etc.
- `refactor`: Code change that neither fixes a bug nor adds a feature
- `perf`: Performance improvement
- `test`: Adding missing tests
- `chore`: Maintenance tasks

### Examples

```
feat: add rate limiting endpoint

fix: handle IPv6 addresses correctly

docs: update API documentation with new endpoints

perf: optimize memory allocation in hot path
```

## Questions?

Feel free to open an issue with the `question` label or reach out to the maintainers.

Thank you for contributing! 🚀
