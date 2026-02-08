package testalign

import "strings"

// MatchTestFuncs はテスト関数をソース関数にマッチングする。
// マッチングルール（優先度順）:
// 1. 完全一致: TargetName == QualifiedName
// 2. サブテストマッチ: TargetNameがQualifiedName + "_"で始まる最長一致
// 3. マッチなしの場合はSourceFuncがnilのMatchResultを返す
func MatchTestFuncs(testFuncs []TestFunc, sourceFuncs []SourceFunc) []MatchResult {
	results := make([]MatchResult, 0, len(testFuncs))

	for _, tf := range testFuncs {
		matched := matchTestToSource(tf.TargetName(), sourceFuncs)
		results = append(results, MatchResult{
			TestFunc:   tf,
			SourceFunc: matched,
		})
	}

	return results
}

// matchTestToSource はテスト関数のターゲット名に対応するソース関数を探す。
func matchTestToSource(targetName string, sourceFuncs []SourceFunc) *SourceFunc {
	if targetName == "" {
		return nil
	}

	// 1. 完全一致を試行
	for i := range sourceFuncs {
		if sourceFuncs[i].QualifiedName() == targetName {
			return &sourceFuncs[i]
		}
	}

	// 2. サブテストマッチ: targetNameがQualifiedName + "_"で始まる最長一致
	var bestMatch *SourceFunc
	bestLen := 0

	for i := range sourceFuncs {
		qname := sourceFuncs[i].QualifiedName()
		prefix := qname + "_"
		if strings.HasPrefix(targetName, prefix) && len(qname) > bestLen {
			bestMatch = &sourceFuncs[i]
			bestLen = len(qname)
		}
	}

	return bestMatch
}
