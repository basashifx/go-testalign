package testalign

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func parseSource(t *testing.T, src string) (*ast.File, *token.FileSet) {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "example.go", src, 0)
	if err != nil {
		t.Fatalf("パース失敗: %v", err)
	}

	return file, fset
}

func parseTestFile(t *testing.T, src string) (*ast.File, *token.FileSet) {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "example_test.go", src, 0)
	if err != nil {
		t.Fatalf("パース失敗: %v", err)
	}

	return file, fset
}

func TestExtractSourceFuncs(t *testing.T) {
	src := `package example

func init() {}

func PublicFunc() {}

func helper() {}

type Service struct{}

func (s *Service) Create() {}

func (s Service) Delete() {}
`
	file, fset := parseSource(t, src)
	funcs := ExtractSourceFuncs(file, fset)

	// initは除外されるので4つ
	if got := len(funcs); got != 4 {
		t.Fatalf("関数数: got %d, want 4", got)
	}

	expected := []struct {
		name         string
		receiverType string
	}{
		{"PublicFunc", ""},
		{"helper", ""},
		{"Create", "Service"},
		{"Delete", "Service"},
	}

	for i, e := range expected {
		if funcs[i].Name != e.name {
			t.Errorf("funcs[%d].Name: got %q, want %q", i, funcs[i].Name, e.name)
		}
		if funcs[i].ReceiverType != e.receiverType {
			t.Errorf("funcs[%d].ReceiverType: got %q, want %q", i, funcs[i].ReceiverType, e.receiverType)
		}
	}
}

func TestExtractSourceFuncs_Generics(t *testing.T) {
	src := `package example

type Container[T any] struct{}

func (c *Container[T]) Get() T { var zero T; return zero }

type Pair[K comparable, V any] struct{}

func (p Pair[K, V]) Keys() []K { return nil }
`
	file, fset := parseSource(t, src)
	funcs := ExtractSourceFuncs(file, fset)

	if got := len(funcs); got != 2 {
		t.Fatalf("関数数: got %d, want 2", got)
	}

	if funcs[0].ReceiverType != "Container" {
		t.Errorf("funcs[0].ReceiverType: got %q, want %q", funcs[0].ReceiverType, "Container")
	}
	if funcs[1].ReceiverType != "Pair" {
		t.Errorf("funcs[1].ReceiverType: got %q, want %q", funcs[1].ReceiverType, "Pair")
	}
}

func TestExtractTestFuncs(t *testing.T) {
	src := `package example

import "testing"

func TestCreate(t *testing.T) {}

func TestDelete(t *testing.T) {}

func BenchmarkCreate(b *testing.B) {}

func FuzzCreate(f *testing.F) {}

func ExampleCreate() {}

func Test_helper(t *testing.T) {}

func helperFunc() {}

func TestMain(m *testing.M) {}
`
	file, fset := parseTestFile(t, src)
	funcs := ExtractTestFuncs(file, fset)

	expected := []string{
		"TestCreate",
		"TestDelete",
		"BenchmarkCreate",
		"FuzzCreate",
		"ExampleCreate",
		"Test_helper",
		"TestMain",
	}

	if got := len(funcs); got != len(expected) {
		t.Fatalf("関数数: got %d, want %d", got, len(expected))
	}

	for i, name := range expected {
		if funcs[i].Name != name {
			t.Errorf("funcs[%d].Name: got %q, want %q", i, funcs[i].Name, name)
		}
	}
}

func TestExtractTestFuncs_ExcludesNonTest(t *testing.T) {
	src := `package example

import "testing"

func setup() {}

func teardown() {}

func Testable() {}

func Testing(t *testing.T) {}

func TestValid(t *testing.T) {}
`
	file, fset := parseTestFile(t, src)
	funcs := ExtractTestFuncs(file, fset)

	if got := len(funcs); got != 1 {
		t.Fatalf("関数数: got %d, want 1 (TestValid only)", got)
	}

	if funcs[0].Name != "TestValid" {
		t.Errorf("Name: got %q, want %q", funcs[0].Name, "TestValid")
	}
}

func TestIsTestFile(t *testing.T) {
	tests := []struct {
		filename string
		want     bool
	}{
		{"foo_test.go", true},
		{"foo.go", false},
		{"test.go", false},
		{"foo_test_test.go", true},
	}

	for _, tt := range tests {
		if got := IsTestFile(tt.filename); got != tt.want {
			t.Errorf("IsTestFile(%q): got %v, want %v", tt.filename, got, tt.want)
		}
	}
}

func TestSourceFileForTest(t *testing.T) {
	tests := []struct {
		testFile string
		want     string
	}{
		{"foo_test.go", "foo.go"},
		{"service_test.go", "service.go"},
	}

	for _, tt := range tests {
		if got := SourceFileForTest(tt.testFile); got != tt.want {
			t.Errorf("SourceFileForTest(%q): got %q, want %q", tt.testFile, got, tt.want)
		}
	}
}
