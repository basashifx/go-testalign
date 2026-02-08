package testalign

import (
	"go/token"
	"strings"
)

// SourceFunc はソースファイル内の関数またはメソッド宣言を表す。
type SourceFunc struct {
	Name         string    // 関数名
	ReceiverType string    // レシーバー型名（関数の場合は空）
	Pos          token.Pos // 宣言位置
	FileName     string    // ファイル名
}

// QualifiedName はレシーバー型を含む修飾名を返す。
// メソッドの場合は "ReceiverType_Name"、関数の場合は "Name" を返す。
// 非公開関数の場合は "_name" の形式を返す。
func (sf SourceFunc) QualifiedName() string {
	if sf.ReceiverType != "" {
		return sf.ReceiverType + "_" + sf.Name
	}

	// 非公開関数の場合、先頭が小文字なので "_name" 形式にする
	if len(sf.Name) > 0 && sf.Name[0] >= 'a' && sf.Name[0] <= 'z' {
		return "_" + sf.Name
	}

	return sf.Name
}

// TestFunc はテストファイル内のテスト関数宣言を表す。
type TestFunc struct {
	Name     string    // テスト関数名（例: "TestService_Create"）
	Pos      token.Pos // 宣言位置
	FileName string    // ファイル名
}

// TargetName はテストプレフィックス（Test/Benchmark/Fuzz/Example）を除去した名前を返す。
func (tf TestFunc) TargetName() string {
	for _, prefix := range []string{"Test", "Benchmark", "Fuzz", "Example"} {
		if strings.HasPrefix(tf.Name, prefix) {
			rest := tf.Name[len(prefix):]
			// "Test" のみの場合は空文字を返す
			if rest == "" {
				return ""
			}
			// "Test_helper" のようなケースはそのまま返す
			if rest[0] == '_' {
				return rest
			}

			return rest
		}
	}

	return tf.Name
}

// MatchResult はテスト関数とソース関数の対応を表す。
type MatchResult struct {
	TestFunc   TestFunc
	SourceFunc *SourceFunc // nilの場合は対応するソース関数なし
}

// OrderViolation は順序違反の情報を表す。
type OrderViolation struct {
	TestFunc      TestFunc     // 順序違反のテスト関数
	SourceFunc    SourceFunc   // 対応するソース関数
	PrecedingTest *MatchResult // ソース順序的に後にあるべきテスト関数
}

// SourceOrderFact は外部テストパッケージ用のFactとしてエクスポートされる。
// ソースファイルごとの関数一覧を保持する。
type SourceOrderFact struct {
	FileToFuncs map[string][]SourceFunc
}

func (*SourceOrderFact) AFact() {}

func (*SourceOrderFact) String() string { return "testalign source order" }
