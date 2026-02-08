package testalign

import (
	"go/token"
	"testing"
)

func TestDetectOrderViolations_NoViolation(t *testing.T) {
	sourceFuncs := []SourceFunc{
		{Name: "Create", ReceiverType: "Service", Pos: token.Pos(10)},
		{Name: "Delete", ReceiverType: "Service", Pos: token.Pos(20)},
	}
	matches := []MatchResult{
		{TestFunc: TestFunc{Name: "TestService_Create"}, SourceFunc: &sourceFuncs[0]},
		{TestFunc: TestFunc{Name: "TestService_Delete"}, SourceFunc: &sourceFuncs[1]},
	}

	violations := DetectOrderViolations(matches, sourceFuncs)

	if len(violations) != 0 {
		t.Errorf("違反数: got %d, want 0", len(violations))
	}
}

func TestDetectOrderViolations_SingleViolation(t *testing.T) {
	sourceFuncs := []SourceFunc{
		{Name: "Create", ReceiverType: "Service", Pos: token.Pos(10)},
		{Name: "Delete", ReceiverType: "Service", Pos: token.Pos(20)},
	}
	// テスト順序が逆
	matches := []MatchResult{
		{TestFunc: TestFunc{Name: "TestService_Delete"}, SourceFunc: &sourceFuncs[1]},
		{TestFunc: TestFunc{Name: "TestService_Create"}, SourceFunc: &sourceFuncs[0]},
	}

	violations := DetectOrderViolations(matches, sourceFuncs)

	if len(violations) != 1 {
		t.Fatalf("違反数: got %d, want 1", len(violations))
	}

	v := violations[0]
	if v.TestFunc.Name != "TestService_Create" {
		t.Errorf("TestFunc.Name: got %q, want %q", v.TestFunc.Name, "TestService_Create")
	}
	if v.SourceFunc.Name != "Create" {
		t.Errorf("SourceFunc.Name: got %q, want %q", v.SourceFunc.Name, "Create")
	}
	if v.PrecedingTest == nil || v.PrecedingTest.TestFunc.Name != "TestService_Delete" {
		t.Errorf("PrecedingTest: 不正な値")
	}
}

func TestDetectOrderViolations_SkipsUnmatched(t *testing.T) {
	sourceFuncs := []SourceFunc{
		{Name: "Create", ReceiverType: "Service", Pos: token.Pos(10)},
		{Name: "Delete", ReceiverType: "Service", Pos: token.Pos(20)},
	}
	matches := []MatchResult{
		{TestFunc: TestFunc{Name: "TestService_Delete"}, SourceFunc: &sourceFuncs[1]},
		{TestFunc: TestFunc{Name: "TestUnrelated"}, SourceFunc: nil}, // マッチなし
		{TestFunc: TestFunc{Name: "TestService_Create"}, SourceFunc: &sourceFuncs[0]},
	}

	violations := DetectOrderViolations(matches, sourceFuncs)

	if len(violations) != 1 {
		t.Fatalf("違反数: got %d, want 1", len(violations))
	}
}

func TestDetectOrderViolations_SameSourceFunc(t *testing.T) {
	sourceFuncs := []SourceFunc{
		{Name: "Create", ReceiverType: "Service", Pos: token.Pos(10)},
		{Name: "Delete", ReceiverType: "Service", Pos: token.Pos(20)},
	}
	// 同じソース関数へのサブテスト（単調非減少なのでOK）
	matches := []MatchResult{
		{TestFunc: TestFunc{Name: "TestService_Create_Success"}, SourceFunc: &sourceFuncs[0]},
		{TestFunc: TestFunc{Name: "TestService_Create_Error"}, SourceFunc: &sourceFuncs[0]},
		{TestFunc: TestFunc{Name: "TestService_Delete"}, SourceFunc: &sourceFuncs[1]},
	}

	violations := DetectOrderViolations(matches, sourceFuncs)

	if len(violations) != 0 {
		t.Errorf("違反数: got %d, want 0", len(violations))
	}
}

func TestDetectOrderViolations_MultipleViolations(t *testing.T) {
	sourceFuncs := []SourceFunc{
		{Name: "A", Pos: token.Pos(10)},
		{Name: "B", Pos: token.Pos(20)},
		{Name: "C", Pos: token.Pos(30)},
	}
	// C, A, B の順（CがmaxになるのでAもBも違反）
	matches := []MatchResult{
		{TestFunc: TestFunc{Name: "TestC"}, SourceFunc: &sourceFuncs[2]},
		{TestFunc: TestFunc{Name: "TestA"}, SourceFunc: &sourceFuncs[0]},
		{TestFunc: TestFunc{Name: "TestB"}, SourceFunc: &sourceFuncs[1]},
	}

	violations := DetectOrderViolations(matches, sourceFuncs)

	if len(violations) != 2 {
		t.Fatalf("違反数: got %d, want 2", len(violations))
	}
}
