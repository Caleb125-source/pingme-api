# Contributing to PingMe API

Thank you for considering contributing to PingMe API! This document provides guidelines for contributing to this project.

## 🎯 How Can I Contribute?

### Reporting Bugs

If you find a bug, please create an issue with:
- A clear, descriptive title
- Steps to reproduce the problem
- Expected behavior
- Actual behavior
- Your environment (OS, Go version, etc.)

### Suggesting Enhancements

Enhancement suggestions are welcome! Please create an issue with:
- A clear description of the enhancement
- Why this would be useful
- Example usage if applicable

### Pull Requests

1. **Fork the repository**
2. **Create a feature branch** (`git checkout -b feature/amazing-feature`)
3. **Make your changes**
4. **Test thoroughly**
5. **Commit your changes** (`git commit -m 'Add amazing feature'`)
6. **Push to the branch** (`git push origin feature/amazing-feature`)
7. **Open a Pull Request**

## 📝 Code Guidelines

### Go Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` to format your code
- Write clear, descriptive variable and function names
- Add comments for exported functions and complex logic

### Code Structure

- Keep handlers focused and single-purpose
- Use the existing `Response` structure for consistency
- Validate all inputs
- Return appropriate HTTP status codes
- Include error messages in responses

### Example Handler

```go
func myHandler(w http.ResponseWriter, r *http.Request) {
    // 1. Validate HTTP method
    if r.Method != http.MethodPost {
        respondJSON(w, http.StatusMethodNotAllowed, Response{
            Success: false,
            Error:   "Method not allowed. Use POST.",
        })
        return
    }

    // 2. Validate headers
    // 3. Parse and validate input
    // 4. Process request
    // 5. Return response
    
    respondJSON(w, http.StatusOK, Response{
        Success: true,
        Message: "Operation successful",
        Data:    result,
    })
}
```

## 🧪 Testing

- Add test cases for new functionality
- Update `tests/api-tests.sh` with new test scenarios
- Ensure all existing tests pass
- Test error cases, not just happy paths

## 📚 Documentation

When adding features, please update:
- `README.md` - Overview and feature list
- `API_DOCUMENTATION.md` - Endpoint details
- `tests/example-requests.md` - Example requests
- Code comments - Explain complex logic

## 🔀 Branch Naming

- `feature/` - New features
- `bugfix/` - Bug fixes
- `docs/` - Documentation changes
- `refactor/` - Code refactoring

Examples:
- `feature/add-authentication`
- `bugfix/fix-json-parsing`
- `docs/update-readme`

## 💬 Commit Messages

Write clear commit messages:
- Use present tense ("Add feature" not "Added feature")
- Keep first line under 50 characters
- Add detailed description if needed

Good examples:
```
Add rate limiting middleware
Fix JSON parsing for empty messages
Update API documentation for echo endpoint
```

## 🎨 What We're Looking For

Great contributions might include:
- Additional endpoints with good documentation
- Middleware (logging, rate limiting, CORS, auth)
- Database integration examples
- Unit tests
- Performance improvements
- Documentation improvements
- Bug fixes

## ❌ What to Avoid

- Breaking changes without discussion
- Removing existing functionality
- Adding dependencies without good reason
- Ignoring existing code style
- Insufficient testing

## 🤔 Questions?

If you're unsure about anything:
1. Check existing issues and PRs
2. Create an issue to discuss your idea
3. Ask in the PR comments

## 📜 License

By contributing, you agree that your contributions will be licensed under the MIT License.

## 🙏 Thank You!

Every contribution, no matter how small, helps make PingMe API better for everyone!