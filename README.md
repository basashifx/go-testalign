# go-testalign

A Go static analysis tool that verifies test function order matches the declaration order of the corresponding source functions.

## Motivation

When test functions are ordered differently from their corresponding source functions, navigating between implementation and tests becomes harder. `go-testalign` detects this misalignment and reports it, helping keep your test files consistent and easy to follow.

## Installation

```bash
go install github.com/basashifx/go-testalign/cmd/go-testalign@latest
```

## Usage

```bash
go-testalign ./...
```

## Example

Given a source file:

```go
// service.go
package service

type Service struct{}

func (s *Service) Create() error { return nil }
func (s *Service) Read() error   { return nil }
func (s *Service) Delete() error { return nil }
```

And a test file with mismatched order:

```go
// service_test.go
package service

func TestService_Delete(t *testing.T) {}
func TestService_Create(t *testing.T) {}
func TestService_Read(t *testing.T) {}
```

`go-testalign` reports:

```
service_test.go:5:1: TestService_Create corresponds to Service.Create (service.go:5) but appears before TestService_Delete which corresponds to Service.Delete (service.go:7)
service_test.go:6:1: TestService_Read corresponds to Service.Read (service.go:6) but appears before TestService_Delete which corresponds to Service.Delete (service.go:7)
```

## How it works

### Test-to-source matching

`go-testalign` strips the test prefix (`Test`, `Benchmark`, `Fuzz`, `Example`) and matches the remaining name against source functions:

| Test function | Matched source |
|---|---|
| `TestService_Create` | `(s *Service) Create()` |
| `TestService_Create_Success` | `(s *Service) Create()` (subtest, longest prefix match) |
| `TestMyFunc` | `func MyFunc()` |
| `Test_helper` | `func helper()` |

Unmatched test functions (e.g. `TestIntegration`, test helpers) are silently skipped.

### File pairing

Test files are paired with source files by naming convention: `foo_test.go` checks against `foo.go`. If no matching source file exists, all source functions in the package are used.

### Order verification

For each test file, the tool assigns source declaration indices to matched test functions and verifies the sequence is monotonically non-decreasing. A violation is reported when a test function's source appears earlier than the source of a preceding test.

### External test packages

For external test packages (`package foo_test`), source function order is communicated via `analysis.Fact`, so the tool works correctly even when the test package is separate from the source package.

## Supported patterns

- Methods with pointer and value receivers
- Generic types (`Container[T]`, `Pair[K, V]`)
- Subtest naming (`TestFoo_Success`, `TestFoo_Error`)
- Unexported functions (`Test_helper` -> `helper()`)
- Multiple source files per package
- Multiple receiver types in a single file
- `Test`, `Benchmark`, `Fuzz`, and `Example` prefixes
- External test packages (`package foo_test`)

## Requirements

- Go 1.22 or later

## License

MIT
