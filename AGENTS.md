# AGENTS.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

Goのテスト関数の順序がソースコードの宣言順序と一致しているかを検証する静的解析ツール（linter）。`golang.org/x/tools/go/analysis`フレームワークに基づく。

## コマンド

[Task](https://taskfile.dev/)を使用。

- `task build` - ビルド
- `task test` - 全テスト実行
- `task lint` - 全静的解析（fmt:check, vet, staticcheck）
- `task check` - lint + test 一括実行
- `task fmt` - goimportsでコード整形

単一テスト実行: `go test -run TestName ./...`

## アーキテクチャ

### パイプライン構成

Analyzerの`run`関数（`analyzer.go`）が以下の4段階パイプラインを実行する:

1. **Extract**（`extract.go`）- ASTからソース関数（`SourceFunc`）とテスト関数（`TestFunc`）を抽出
2. **Match**（`match.go`）- テスト関数名からソース関数への対応付け（完全一致 → サブテスト最長プレフィックス一致）
3. **Order**（`order.go`）- マッチ結果のソースインデックス列が単調非減少かを検証
4. **Report**（`analyzer.go`内）- 順序違反の診断メッセージを生成

型定義は`types.go`に集約。CLIエントリポイントは`cmd/go-testalign/main.go`（`singlechecker.Main`呼び出しのみ）。

### 外部テストパッケージ対応

`package foo_test`のような外部テストパッケージでは、ソースパッケージの関数情報に直接アクセスできない。`analysis.Fact`（`SourceOrderFact`）を使い、ソースパッケージのpassで関数情報をエクスポートし、外部テストパッケージのpassでインポートする。

### テスト構成

- **ユニットテスト**（`extract_test.go`, `match_test.go`, `order_test.go`）: 各パイプライン段階を個別にテスト。内部パッケージ（`package testalign`）。
- **統合テスト**（`analyzer_test.go`）: `analysistest.Run`で`testdata/src/`配下のテストケースを実行。外部テストパッケージ（`package testalign_test`）。
- **テストデータ**（`testdata/src/`）: `analysistest`用のフィクスチャ。`// want "..."`コメントで期待される診断を記述。

### analysistest の注意点

- `// want "regex"` でその行の診断を期待、`// want package:"regex"` でパッケージFactを期待
- テストバイナリのmainパッケージでもrunが呼ばれるため、`pass.Pkg.Name() == "main"`でスキップが必要
- 外部テストパッケージのパスは `<path>_test` の形式。`strings.TrimSuffix`で`_test`を除去してソースパッケージを特定する
