package testalign

import (
	"fmt"
	"go/ast"
	"go/token"
	"maps"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Analyzer はテスト関数の順序がソースコードの宣言順序と一致しているかを検証する。
var Analyzer = &analysis.Analyzer{
	Name: "testalign",
	Doc:  "テスト関数の順序がソースコードの宣言順序と一致しているかを検証する",
	Run:  run,
	FactTypes: []analysis.Fact{
		(*SourceOrderFact)(nil),
	},
}

func run(pass *analysis.Pass) (any, error) {
	// テストバイナリのmainパッケージはスキップ
	if pass.Pkg.Name() == "main" {
		return nil, nil
	}

	// ファイルをソースファイルとテストファイルに分類
	sourceFiles := make(map[string]*ast.File)
	testFiles := make(map[string]*ast.File)

	for _, file := range pass.Files {
		pos := pass.Fset.Position(file.Pos())
		fileName := filepath.Base(pos.Filename)
		if IsTestFile(fileName) {
			testFiles[fileName] = file
		} else {
			sourceFiles[fileName] = file
		}
	}

	// ソースファイルから関数を抽出
	allSourceFuncs := make(map[string][]SourceFunc)
	for fileName, file := range sourceFiles {
		funcs := ExtractSourceFuncs(file, pass.Fset)
		if len(funcs) > 0 {
			allSourceFuncs[fileName] = funcs
		}
	}

	// ソース関数情報をFactとしてエクスポート（外部テストパッケージ用）
	if len(allSourceFuncs) > 0 {
		fact := &SourceOrderFact{FileToFuncs: allSourceFuncs}
		pass.ExportPackageFact(fact)
	}

	// 外部テストパッケージの場合、Factからソース関数をインポート
	if len(sourceFiles) == 0 && len(testFiles) > 0 {
		allSourceFuncs = importSourceFuncsFromFact(pass)
	}

	// テストファイルごとに検証
	for testFileName, testFile := range testFiles {
		testFuncs := ExtractTestFuncs(testFile, pass.Fset)
		if len(testFuncs) == 0 {
			continue
		}

		// 対応するソースファイルの関数を収集
		sourceFuncs := collectSourceFuncsForTestFile(testFileName, allSourceFuncs)
		if len(sourceFuncs) == 0 {
			continue
		}

		// マッチング
		matches := MatchTestFuncs(testFuncs, sourceFuncs)

		// 順序検証
		violations := DetectOrderViolations(matches, sourceFuncs)

		// 診断報告
		for _, v := range violations {
			reportViolation(pass, v)
		}
	}

	return nil, nil
}

// collectSourceFuncsForTestFile はテストファイルに対応するソースファイルの関数を収集する。
// 対応ルール: foo_test.go → foo.go
// 対応するソースファイルがない場合は、全ソースファイルの関数を結合して返す。
func collectSourceFuncsForTestFile(testFileName string, allSourceFuncs map[string][]SourceFunc) []SourceFunc {
	sourceFileName := SourceFileForTest(testFileName)

	// まず直接対応を試みる
	if funcs, ok := allSourceFuncs[sourceFileName]; ok {
		return funcs
	}

	// 対応するソースファイルがない場合、全ソースファイルの関数を結合
	var all []SourceFunc
	for _, funcs := range allSourceFuncs {
		all = append(all, funcs...)
	}

	return all
}

// importSourceFuncsFromFact は依存パッケージからFactをインポートし、ソース関数を取得する。
func importSourceFuncsFromFact(pass *analysis.Pass) map[string][]SourceFunc {
	result := make(map[string][]SourceFunc)

	// 外部テストパッケージのパスは "<path>_test" の形式
	// ソースパッケージのパスは "<path>"
	sourcePkgPath := strings.TrimSuffix(pass.Pkg.Path(), "_test")

	for _, imp := range pass.Pkg.Imports() {
		if imp.Path() != sourcePkgPath {
			continue
		}

		var fact SourceOrderFact
		if pass.ImportPackageFact(imp, &fact) {
			maps.Copy(result, fact.FileToFuncs)
		}
	}

	return result
}

// reportViolation は順序違反の診断メッセージを生成・報告する。
func reportViolation(pass *analysis.Pass, v OrderViolation) {
	srcPos := formatSourcePos(pass.Fset, v.SourceFunc)
	precedingPos := ""
	precedingName := ""

	if v.PrecedingTest != nil && v.PrecedingTest.SourceFunc != nil {
		precedingPos = formatSourcePos(pass.Fset, *v.PrecedingTest.SourceFunc)
		precedingName = v.PrecedingTest.TestFunc.Name
	}

	var msg string
	if precedingName != "" {
		msg = fmt.Sprintf(
			"%s corresponds to %s (%s) but appears before %s which corresponds to %s (%s)",
			v.TestFunc.Name,
			formatFuncRef(v.SourceFunc),
			srcPos,
			precedingName,
			formatFuncRef(*v.PrecedingTest.SourceFunc),
			precedingPos,
		)
	} else {
		msg = fmt.Sprintf(
			"%s corresponds to %s (%s) but is out of order",
			v.TestFunc.Name,
			formatFuncRef(v.SourceFunc),
			srcPos,
		)
	}

	pass.Reportf(v.TestFunc.Pos, "%s", msg)
}

// formatFuncRef はソース関数の参照文字列を返す。
// メソッドの場合: "ReceiverType.Name"
// 関数の場合: "Name"
func formatFuncRef(sf SourceFunc) string {
	if sf.ReceiverType != "" {
		return sf.ReceiverType + "." + sf.Name
	}

	return sf.Name
}

// formatSourcePos はソース関数のファイル位置を "filename:line" 形式で返す。
func formatSourcePos(fset *token.FileSet, sf SourceFunc) string {
	if sf.Pos.IsValid() {
		pos := fset.Position(sf.Pos)

		return fmt.Sprintf("%s:%d", filepath.Base(pos.Filename), pos.Line)
	}

	return sf.FileName
}
