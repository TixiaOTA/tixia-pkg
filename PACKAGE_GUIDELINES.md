# Package Development Guidelines

## Overview

This document outlines the structure, conventions, and best practices for creating new packages in the `tixia-pkg` repository. These guidelines apply to both provider packages (API clients) and utility packages (shared functionality).

## Package Structure

### Provider Packages (e.g., `providers/mg`, `providers/vervotech`)

**Required Files:**

- `client.go` - Client struct, Config struct, and `New()` constructor function
- `models.go` - Request and response type definitions
- `errors.go` - Custom error types and error code constants
- `endpoint.go` - API endpoint methods (methods on Client)
- `helpers.go` - Internal utility functions
- `README.md` - Package documentation with usage examples
- `*_test.go` - Test files (e.g., `endpoint_test.go`, `client_test.go`)

**Optional Files:**

- `usecase.go` - High-level business logic combining multiple endpoints
- `usecase_test.go` - Tests for usecase methods

### Utility Packages (e.g., `logger`, `redis`, `channel`)

**Required Files:**

- Main package file (e.g., `logger.go`, `manager.go`) - Core functionality
- `config.go` - Configuration structs (if needed)
- `README.md` - Package documentation with usage examples
- `*_test.go` - Test files

**Optional Files:**

- Subdirectories for organized functionality (e.g., `redis/cache/`, `redis/streams/`)
- `types.go` - Type definitions (if models.go pattern not used)

## Naming Conventions

- Package names should be lowercase, short, and descriptive
- Client structs should be named `Client`
- Config structs should be named `Config`
- Constructor functions should be named `New()`
- Request types should end with `Request` (e.g., `SearchHotelRequest`)
- Response types should end with `Response` (e.g., `SearchHotelResponse`)
- Error types should end with `Error` (e.g., `MGError`, `APIError`)
- Error code constants should be descriptive (e.g., `InternalServerError`, `JRVXML161`)

## Common Patterns

### Client Initialization

- Use a `Config` struct for configuration
- Provide sensible defaults for optional configuration
- Validate required fields in `New()` constructor
- Return errors for invalid configuration

### HTTP Client Usage

- Use `github.com/TixiaOTA/tixia-pkg/http/client` for HTTP operations
- Configure retry behavior via `RetryConfig` in Config
- Set timeouts per endpoint when needed using `httpClient.SetTimeout()`
- Use `PostJSON()` for POST requests with JSON bodies
- Use `Get()` for GET requests with query parameters

### Logging

- Accept `*logger.Logger` as a dependency in Client struct (for provider packages)
- Use `logger.Sanitize()` for sensitive data before logging
- Include trace context using `logger.WithTrace()` with span ID
- Log at appropriate levels (Info for operations, Error for failures)
- Log request/response data when appropriate (after sanitization)

### OpenTelemetry Tracing

- Use `github.com/TixiaOTA/tixia-pkg/otel` for distributed tracing
- Start spans with `otelpkg.StartSpanWithID()` at the beginning of endpoint methods
- Use span names following pattern: `"provider", "<package-name>", "<method-name>"`
- Always defer `span.End()`
- Capture span ID for logging correlation

### Error Handling

- Define custom error types in `errors.go` (or `models.go` if they're part of API responses)
- Error types should be structs that implement the `error` interface with an `Error()` method
- Use error code constants for API errors (e.g., `MGErrorCode` type)
- Return errors that implement the `error` interface
- Include both error code and error message in error responses
- Handle HTTP error status codes appropriately
- Try to unmarshal error responses from APIs when possible

### Request/Response Models

- Define request types with JSON tags matching API specification
- Define response types with JSON tags matching API specification
- Use pointers for optional fields
- Group related types together (requests, responses, shared types)
- Add comments for request/response structs describing their purpose

## Testing Requirements

### Test File Organization

- Create test files alongside source files (e.g., `endpoint_test.go` for `endpoint.go`)
- Use `*_test.go` naming convention
- Place tests in the same package (package-level tests)

### Test Coverage

- Test all public methods and functions
- Test error cases and edge cases
- Test configuration validation
- Test HTTP client interactions (use mocks or test servers)
- Test timeout and retry behavior when applicable

### Test Structure

- Use table-driven tests for multiple test cases
- Use descriptive test function names: `TestClient_MethodName_Scenario`
- Use `github.com/stretchr/testify` for assertions when appropriate
- Test both success and failure paths

## Dependencies

### Standard Dependencies

- `github.com/TixiaOTA/tixia-pkg/logger` - For logging (provider packages)
- `github.com/TixiaOTA/tixia-pkg/otel` - For distributed tracing (provider packages)
- `github.com/TixiaOTA/tixia-pkg/http/client` - For HTTP operations (provider packages)

### External Dependencies

- Prefer using existing tixia-pkg utilities over external libraries
- When external libraries are needed, document the reason
- Keep dependencies minimal and well-justified

## Code Style

- Follow Go standard formatting (`gofmt`) (https://go.dev/doc/effective_go)
- Avoid one-liner code for readability
- Write comments for complex functions and code blocks

## Examples

### Provider Package Example Structure

```

providers/example/
├── client.go # Client, Config, New()
├── models.go # Request/Response types
├── errors.go # Error types and codes
├── endpoint.go # API methods
├── helpers.go # Internal utilities
├── endpoint_test.go # Tests
├── client_test.go # Tests
└── README.md # Documentation

```

### Utility Package Example Structure

```

example/
├── example.go # Core functionality
├── config.go # Configuration
├── example_test.go # Tests
└── README.md # Documentation

```

## Documentation Requirements

### README.md Structure

Every package must include a comprehensive README.md file that follows a consistent structure. The README serves as the primary documentation for users of the package.

#### Required Sections

1. **Title and Brief Description**

   - Package name as H1 heading
   - One-line description of what the package does
   - Brief overview paragraph

2. **Overview**

   - Bulleted list of main capabilities/features
   - What problems the package solves
   - Key use cases

3. **Installation**

   - `go get` command with full import path
   - Example: `go get github.com/TixiaOTA/tixia-pkg/providers/mg` or `go get github.com/TixiaOTA/tixia-pkg/logger`

4. **Configuration**

   - **Client Configuration** (for provider packages) or **Configuration** (for utility packages): Complete code example showing how to create a client instance or configure/initialize the package
   - **Configuration Fields**: Detailed table or list describing each configuration field:
     - Field name and type
     - Required vs optional
     - Default values (if applicable)
     - Brief description of what it does

5. **Usage**

   - Separate subsection for each main operation/method (provider packages) or function (utility packages)
   - Each subsection should include:
     - Brief description of what the operation does
     - Complete, runnable code example
     - Comments explaining important parts
     - Expected output or result

6. **API Reference**

   - **Client Methods** (for provider packages) or **Public Functions/Types** (for utility packages): Documentation for all public methods or functions
   - For each method/function, include:
     - Function signature
     - **Request Parameters** (for provider methods) or **Parameters** (for utility functions): List of all parameters with types and descriptions
     - **Response** (for provider methods) or **Return Values** (for utility functions): Structure of the response/return values with field descriptions
     - **Timeout**: If applicable (typically for provider packages), mention timeout duration
     - **Retries**: If applicable (typically for provider packages), mention retry behavior
     - **Special Notes**: Any important behavior or side effects (e.g., "No retries for BookHotel to prevent duplicate bookings")

7. **Error Handling** (if applicable)

   - Document all error types
   - List error codes with descriptions (if applicable)
   - Provide examples of error handling
   - Explain when each error type occurs

8. **Data Models** (if applicable)

   - Document key data structures
   - Show struct definitions with field descriptions
   - Reference source files (e.g., `models.go`) for complete definitions

9. **Examples**

   - Complete, end-to-end example showing a common workflow
   - Should be copy-paste runnable (with placeholder values)
   - Demonstrate best practices

10. **License**
    - Standard license notice: "This package is part of the TixiaOTA project."

#### Optional Sections

- **Pagination** - If the package supports pagination
- **Advanced Usage** - For complex use cases
- **Performance Considerations** - If relevant
- **Troubleshooting** - Common issues and solutions

### README Content Guidelines

#### Code Examples

- **Always use complete, runnable examples** - Include necessary imports, error handling, and context
- **Use realistic placeholder values** - Use descriptive placeholders like `"your-api-key"` instead of `"xxx"`
- **Add inline comments** - Explain important parts of the code
- **Show error handling** - Demonstrate proper error checking
- **Format dates/values clearly** - Include format specifications (e.g., `"2025-12-01" // Format: YYYY-MM-DD`)

#### Writing Style

- **Be clear and concise** - Use simple language, avoid jargon
- **Be consistent** - Use the same terminology throughout
- **Use active voice** - "The package enables you to..." not "Users are enabled to..."
- **Use bullet points** - For lists of features, parameters, etc.
- **Use code formatting** - Use backticks for code elements, file names, and package names

#### Configuration Documentation

- **Group related fields** - Use nested structures in examples
- **Mark required vs optional** - Clearly indicate which fields are required
- **Show defaults** - Document default values for optional fields
- **Explain complex fields** - Provide additional context for fields that need explanation

#### API Reference Format

Use a consistent format for documenting methods and functions:

**For Methods (Provider Packages):**

````markdown
#### MethodName

Brief description of what the method does.

```go
func (c *Client) MethodName(ctx context.Context, req RequestType) (*ResponseType, error)
```
````

**Request Parameters:**

- `FieldName` (type, required/optional): Description

**Response:**

- `FieldName` (type): Description

**Timeout:** X seconds (if applicable)

**Retries:** X retries (if applicable)

**For Functions (Utility Packages):**

````markdown
#### FunctionName

Brief description of what the function does.

```go
func FunctionName(param1 Type1, param2 Type2) (ReturnType, error)
```

**Parameters:**

- `param1` (Type1, required/optional): Description
- `param2` (Type2, required/optional): Description

**Return Values:**

- `ReturnType`: Description
- `error`: Description of error conditions

**Special Notes:** Any important behavior or side effects
````

#### Error Documentation

- **List all error codes** - Include constants from `errors.go`
- **Explain when errors occur** - Describe the conditions that trigger each error
- **Provide handling examples** - Show how to check and handle specific errors
- **Include error structure** - Show the error type definition

### README Best Practices

1. **Start with the basics** - Assume readers are new to the package
2. **Progressive disclosure** - Start simple, then show advanced usage
3. **Link to code** - Reference source files when appropriate (e.g., "See `models.go` for complete model definitions")
4. **Keep examples updated** - Ensure examples match current API
5. **Test your examples** - Verify that code examples compile and work
6. **Use consistent formatting** - Follow the same structure as other package READMEs
7. **Include real-world scenarios** - Show how the package is used in practice

### README Template

````markdown
# PackageName Package

A Go [client/utility] library for [brief description]. This package provides [main features].

## Overview

The PackageName package enables you to:

- Feature 1
- Feature 2
- Feature 3

## Installation

```bash
go get github.com/TixiaOTA/tixia-pkg/path/to/package
```

## Configuration

### Client Configuration

<!-- For provider packages, use "Client Configuration" -->
<!-- For utility packages, use "Configuration" -->

[Complete code example]

### Configuration Fields

- **FieldName** (type, required/optional): Description

## Usage

### Operation Name

[Description and code example]

## API Reference

### Client Methods

<!-- For provider packages, use "Client Methods" -->
<!-- For utility packages, use "Public Functions" or "Public Types" -->

#### MethodName

<!-- For provider packages: methods on Client -->
<!-- For utility packages: standalone functions -->

[Method/Function documentation]

## Error Handling

[Error types and handling examples]

## Data Models

[Key data structures]

## Examples

### Complete Example

[End-to-end example]

## License

This package is part of the TixiaOTA project.
````

### Code Comments

- Add comments for complex functions only (per user rules)
- Document exported types, functions, and methods
- Include parameter and return value descriptions for public APIs
- Add section comments to organize code (e.g., `// --------------------------------- REQUESTS ---------------------------------`)
