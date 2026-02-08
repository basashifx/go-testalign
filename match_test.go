package testalign

import "testing"

func TestMatchTestFuncs_ExactMatch(t *testing.T) {
	sourceFuncs := []SourceFunc{
		{Name: "Create", ReceiverType: "Service"},
		{Name: "Delete", ReceiverType: "Service"},
	}
	testFuncs := []TestFunc{
		{Name: "TestService_Create"},
		{Name: "TestService_Delete"},
	}

	results := MatchTestFuncs(testFuncs, sourceFuncs)

	if len(results) != 2 {
		t.Fatalf("結果数: got %d, want 2", len(results))
	}

	for i, r := range results {
		if r.SourceFunc == nil {
			t.Errorf("results[%d]: マッチなし、期待はマッチあり", i)
			continue
		}
		if i == 0 && r.SourceFunc.Name != "Create" {
			t.Errorf("results[0].SourceFunc.Name: got %q, want %q", r.SourceFunc.Name, "Create")
		}
		if i == 1 && r.SourceFunc.Name != "Delete" {
			t.Errorf("results[1].SourceFunc.Name: got %q, want %q", r.SourceFunc.Name, "Delete")
		}
	}
}

func TestMatchTestFuncs_SubtestMatch(t *testing.T) {
	sourceFuncs := []SourceFunc{
		{Name: "Create", ReceiverType: "Service"},
		{Name: "Delete", ReceiverType: "Service"},
	}
	testFuncs := []TestFunc{
		{Name: "TestService_Create_Success"},
		{Name: "TestService_Create_Error"},
		{Name: "TestService_Delete_NotFound"},
	}

	results := MatchTestFuncs(testFuncs, sourceFuncs)

	if len(results) != 3 {
		t.Fatalf("結果数: got %d, want 3", len(results))
	}

	for i, r := range results {
		if r.SourceFunc == nil {
			t.Errorf("results[%d]: マッチなし", i)
			continue
		}
	}

	if results[0].SourceFunc.Name != "Create" {
		t.Errorf("results[0]: got %q, want Create", results[0].SourceFunc.Name)
	}
	if results[1].SourceFunc.Name != "Create" {
		t.Errorf("results[1]: got %q, want Create", results[1].SourceFunc.Name)
	}
	if results[2].SourceFunc.Name != "Delete" {
		t.Errorf("results[2]: got %q, want Delete", results[2].SourceFunc.Name)
	}
}

func TestMatchTestFuncs_PackageFunc(t *testing.T) {
	sourceFuncs := []SourceFunc{
		{Name: "MyFunc"},
		{Name: "AnotherFunc"},
	}
	testFuncs := []TestFunc{
		{Name: "TestMyFunc"},
		{Name: "TestAnotherFunc"},
	}

	results := MatchTestFuncs(testFuncs, sourceFuncs)

	if results[0].SourceFunc == nil || results[0].SourceFunc.Name != "MyFunc" {
		t.Errorf("TestMyFunc: マッチ失敗")
	}
	if results[1].SourceFunc == nil || results[1].SourceFunc.Name != "AnotherFunc" {
		t.Errorf("TestAnotherFunc: マッチ失敗")
	}
}

func TestMatchTestFuncs_UnexportedFunc(t *testing.T) {
	sourceFuncs := []SourceFunc{
		{Name: "helper"},
		{Name: "validate"},
	}
	testFuncs := []TestFunc{
		{Name: "Test_helper"},
		{Name: "Test_validate"},
	}

	results := MatchTestFuncs(testFuncs, sourceFuncs)

	if results[0].SourceFunc == nil || results[0].SourceFunc.Name != "helper" {
		t.Errorf("Test_helper: マッチ失敗")
	}
	if results[1].SourceFunc == nil || results[1].SourceFunc.Name != "validate" {
		t.Errorf("Test_validate: マッチ失敗")
	}
}

func TestMatchTestFuncs_NoMatch(t *testing.T) {
	sourceFuncs := []SourceFunc{
		{Name: "Create", ReceiverType: "Service"},
	}
	testFuncs := []TestFunc{
		{Name: "TestUnrelated"},
	}

	results := MatchTestFuncs(testFuncs, sourceFuncs)

	if results[0].SourceFunc != nil {
		t.Errorf("TestUnrelated: マッチすべきでないがマッチした: %q", results[0].SourceFunc.Name)
	}
}

func TestMatchTestFuncs_LongestPrefixWins(t *testing.T) {
	sourceFuncs := []SourceFunc{
		{Name: "Create", ReceiverType: "Service"},
		{Name: "Create_Batch", ReceiverType: "Service"},
	}
	testFuncs := []TestFunc{
		{Name: "TestService_Create_Batch_Success"},
	}

	results := MatchTestFuncs(testFuncs, sourceFuncs)

	if results[0].SourceFunc == nil {
		t.Fatal("マッチなし")
	}
	// "Service_Create_Batch" は "Service_Create" より長いのでこちらにマッチ
	if results[0].SourceFunc.Name != "Create_Batch" {
		t.Errorf("最長プレフィックス一致: got %q, want %q", results[0].SourceFunc.Name, "Create_Batch")
	}
}

func TestMatchTestFuncs_TestOnly(t *testing.T) {
	sourceFuncs := []SourceFunc{
		{Name: "Create", ReceiverType: "Service"},
	}
	testFuncs := []TestFunc{
		{Name: "Test"},
	}

	results := MatchTestFuncs(testFuncs, sourceFuncs)

	if results[0].SourceFunc != nil {
		t.Errorf("Test（ターゲット名が空）: マッチすべきでない")
	}
}
