package testalign

// DetectOrderViolations はマッチ済みテスト関数の順序がソース関数の宣言順序と
// 一致しているかを検証し、違反箇所を返す。
//
// ソース関数のインデックス列が単調非減少であることを検証する。
// 違反箇所: ソースインデックスがそれまでの最大値より小さい位置。
func DetectOrderViolations(matches []MatchResult, sourceFuncs []SourceFunc) []OrderViolation {
	// ソース関数のインデックスマップを構築
	sourceIndex := buildSourceIndex(sourceFuncs)

	var violations []OrderViolation
	maxIdx := -1
	var maxMatch *MatchResult

	for i := range matches {
		m := &matches[i]
		if m.SourceFunc == nil {
			continue
		}

		idx, ok := sourceIndex[m.SourceFunc.Pos]
		if !ok {
			continue
		}

		if idx < maxIdx && maxMatch != nil {
			violations = append(violations, OrderViolation{
				TestFunc:      m.TestFunc,
				SourceFunc:    *m.SourceFunc,
				PrecedingTest: maxMatch,
			})
		}

		if idx > maxIdx {
			maxIdx = idx
			maxMatch = m
		}
	}

	return violations
}

// buildSourceIndex はソース関数のPosからインデックスへのマッピングを構築する。
func buildSourceIndex(sourceFuncs []SourceFunc) map[any]int {
	index := make(map[any]int, len(sourceFuncs))
	for i, sf := range sourceFuncs {
		index[sf.Pos] = i
	}

	return index
}
