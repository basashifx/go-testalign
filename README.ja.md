# go-testalign

テスト関数の順序がソースコードの関数宣言順序と一致しているかを検証する Go 静的解析ツールです。

## 動機

テスト関数がソース関数と異なる順序で並んでいると、実装とテストの間を行き来しづらくなります。`go-testalign` はこの不整合を検出・報告し、テストファイルの一貫性と可読性を保ちます。

## インストール

```bash
go install github.com/basashifx/go-testalign/cmd/go-testalign@latest
```

## 使い方

```bash
go-testalign ./...
```

## 使用例

以下のようなソースファイルがあるとします：

```go
// service.go
package service

type Service struct{}

func (s *Service) Create() error { return nil }
func (s *Service) Read() error   { return nil }
func (s *Service) Delete() error { return nil }
```

テストファイルの順序が一致していない場合：

```go
// service_test.go
package service

func TestService_Delete(t *testing.T) {}
func TestService_Create(t *testing.T) {}
func TestService_Read(t *testing.T) {}
```

`go-testalign` は以下のように報告します：

```
service_test.go:5:1: TestService_Create corresponds to Service.Create (service.go:5) but appears before TestService_Delete which corresponds to Service.Delete (service.go:7)
service_test.go:6:1: TestService_Read corresponds to Service.Read (service.go:6) but appears before TestService_Delete which corresponds to Service.Delete (service.go:7)
```

## 仕組み

### テストとソースのマッチング

`go-testalign` はテスト関数名からプレフィックス（`Test`、`Benchmark`、`Fuzz`、`Example`）を除去し、残りの名前をソース関数と照合します：

| テスト関数 | マッチするソース関数 |
|---|---|
| `TestService_Create` | `(s *Service) Create()` |
| `TestService_Create_Success` | `(s *Service) Create()`（サブテスト、最長プレフィックスマッチ） |
| `TestMyFunc` | `func MyFunc()` |
| `Test_helper` | `func helper()` |

マッチしないテスト関数（例：`TestIntegration`、テストヘルパー）は無視されます。

### ファイルの対応付け

テストファイルは命名規則によりソースファイルと対応付けられます：`foo_test.go` は `foo.go` を参照します。対応するソースファイルが存在しない場合は、パッケージ内のすべてのソース関数が対象となります。

### 順序の検証

各テストファイルについて、マッチしたテスト関数にソース宣言のインデックスを割り当て、その並びが単調非減少であることを検証します。テスト関数のソースが、先行するテストのソースよりも前に宣言されている場合に違反として報告されます。

### 外部テストパッケージ

外部テストパッケージ（`package foo_test`）の場合、ソース関数の順序は `analysis.Fact` を介して伝達されるため、テストパッケージがソースパッケージと分離していても正しく動作します。

## 対応パターン

- ポインタレシーバおよび値レシーバのメソッド
- ジェネリクス型（`Container[T]`、`Pair[K, V]`）
- サブテストの命名（`TestFoo_Success`、`TestFoo_Error`）
- 非公開関数（`Test_helper` → `helper()`）
- パッケージ内の複数ソースファイル
- 単一ファイル内の複数レシーバ型
- `Test`、`Benchmark`、`Fuzz`、`Example` プレフィックス
- 外部テストパッケージ（`package foo_test`）

## 要件

- Go 1.22 以降

## ライセンス

MIT
