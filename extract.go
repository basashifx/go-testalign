package testalign

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"
)

// テスト関数のプレフィックス一覧
var testPrefixes = []string{"Test", "Benchmark", "Fuzz", "Example"}

// ExtractSourceFuncs はASTファイルからソース関数/メソッドを抽出する。
// init関数は除外される。
func ExtractSourceFuncs(file *ast.File, fset *token.FileSet) []SourceFunc {
	fileName := filepath.Base(fset.Position(file.Pos()).Filename)
	var funcs []SourceFunc

	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// init関数を除外
		if funcDecl.Name.Name == "init" {
			continue
		}

		sf := SourceFunc{
			Name:     funcDecl.Name.Name,
			Pos:      funcDecl.Pos(),
			FileName: fileName,
		}

		// レシーバー型を取得
		if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
			sf.ReceiverType = extractReceiverType(funcDecl.Recv.List[0].Type)
		}

		funcs = append(funcs, sf)
	}

	return funcs
}

// ExtractTestFuncs はASTファイルからテスト関数を抽出する。
// Test*/Benchmark*/Fuzz*/Example* プレフィックスの関数のみ抽出する。
func ExtractTestFuncs(file *ast.File, fset *token.FileSet) []TestFunc {
	fileName := filepath.Base(fset.Position(file.Pos()).Filename)
	var funcs []TestFunc

	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// メソッドは対象外
		if funcDecl.Recv != nil {
			continue
		}

		if !isTestFunc(funcDecl.Name.Name) {
			continue
		}

		funcs = append(funcs, TestFunc{
			Name:     funcDecl.Name.Name,
			Pos:      funcDecl.Pos(),
			FileName: fileName,
		})
	}

	return funcs
}

// IsTestFile はファイル名がテストファイルかどうかを判定する。
func IsTestFile(filename string) bool {
	return strings.HasSuffix(filename, "_test.go")
}

// SourceFileForTest はテストファイルに対応するソースファイル名を返す。
// 例: "foo_test.go" → "foo.go"
func SourceFileForTest(testFile string) string {
	return strings.TrimSuffix(testFile, "_test.go") + ".go"
}

// isTestFunc はテスト関数プレフィックスを持つか判定する。
func isTestFunc(name string) bool {
	for _, prefix := range testPrefixes {
		if strings.HasPrefix(name, prefix) {
			rest := name[len(prefix):]
			// プレフィックスだけの場合（例: "Test"）も有効
			if rest == "" {
				return true
			}
			// プレフィックスの後はアンダースコアまたは大文字で始まる必要がある
			if rest[0] == '_' || (rest[0] >= 'A' && rest[0] <= 'Z') {
				return true
			}
		}
	}

	return false
}

// extractReceiverType はレシーバーの型表現から型名を取得する。
// ポインタレシーバー（*Type）とジェネリクス（Type[T]）に対応する。
func extractReceiverType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.StarExpr:
		return extractReceiverType(t.X)
	case *ast.Ident:
		return t.Name
	case *ast.IndexExpr:
		// ジェネリクス: Type[T]
		return extractReceiverType(t.X)
	case *ast.IndexListExpr:
		// ジェネリクス: Type[T1, T2]
		return extractReceiverType(t.X)
	}

	return ""
}
