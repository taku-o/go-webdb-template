# 英語ドキュメントの用意の設計書

## Overview

### 目的
プロジェクトに英語版のドキュメントを用意し、日本語と英語の両方でドキュメントを提供できるようにする。日英で差分が生じた場合のマスタードキュメントは日本語版とし、内部コメントは従来通り日本語を使用する。

### ユーザー
- **英語話者**: 英語でプロジェクトのドキュメントを読めることを期待する
- **国際的なコントリビューター**: 英語でプロジェクトを理解し、コントリビュートできることを期待する
- **開発者**: 日本語版と英語版の両方を維持できる構造を期待する

### 影響
現在のシステム状態を以下のように変更する：
- `README.md`: 英語版に置き換え（既存は`README.ja.md`にリネーム）
- `docs/*.md`: `docs/ja/*.md`に移動（24ファイル）
- `docs/en/*.md`: 新規作成（24ファイル、英語版）
- `docs/_config.yml`: descriptionを英語に変更
- 既存の`docs/pages/ja/`と`docs/pages/en/`の構造は維持

### Goals
- README.mdの英語版を作成する
- docsディレクトリ内の全ドキュメントの英語版を作成する
- 日本語版と英語版の両方を維持できる構造にする
- 既存のドキュメントリンクが機能し続けるようにする
- GitHub Pagesなどの既存の公開方法が機能し続けるようにする
- 各Markdownファイルに言語切替リンクを追加する（冒頭または末尾）

### Non-Goals
- ドキュメントの内容の大幅な変更・追加
- 既存のドキュメント構造の変更
- コード内のコメントの英語化
- その他の設定ファイルの英語化
- ドキュメントの自動翻訳機能の実装
- 他の言語への対応

## Architecture

### 変更前後のファイル構造比較

#### 変更前のファイル構造
```
/
├── README.md (日本語版)
└── docs/
    ├── _config.yml (description: 日本語)
    ├── Admin.md (日本語版)
    ├── Apache-Superset.md (日本語版)
    ├── API.md (日本語版)
    ├── Architecture.md (日本語版)
    ├── Atlas-Operations.md (日本語版)
    ├── Command-Line-Tool.md (日本語版)
    ├── Database-Viewer.md (日本語版)
    ├── Docker.md (日本語版)
    ├── File-Upload.md (日本語版)
    ├── Generate-Sample-Data.md (日本語版)
    ├── index.md (日本語版)
    ├── Initial-Setup.md (日本語版)
    ├── License-Survey.md (日本語版)
    ├── Logging.md (日本語版)
    ├── Metabase.md (日本語版)
    ├── Partner-Idp-Auth0-Login.md (日本語版)
    ├── Project-Structure.md (日本語版)
    ├── Queue-Job.md (日本語版)
    ├── Rate-Limit.md (日本語版)
    ├── Release-Check.md (日本語版)
    ├── Send-Mail.md (日本語版)
    ├── Sharding.md (日本語版)
    ├── Spec-Driven-Development.md (日本語版)
    ├── Testing.md (日本語版)
    └── pages/
        ├── en/
        │   ├── about.md
        │   ├── index.md
        │   └── setup.md
        └── ja/
            ├── about.md
            ├── index.md
            └── setup.md
```

#### 変更後のファイル構造
```
/
├── README.md (英語版、新規作成)
├── README.ja.md (日本語版、既存のREADME.mdからリネーム)
└── docs/
    ├── _config.yml (description: 英語)
    ├── ja/
    │   ├── Admin.md (日本語版、移動)
    │   ├── Apache-Superset.md (日本語版、移動)
    │   ├── API.md (日本語版、移動)
    │   ├── Architecture.md (日本語版、移動)
    │   ├── Atlas-Operations.md (日本語版、移動)
    │   ├── Command-Line-Tool.md (日本語版、移動)
    │   ├── Database-Viewer.md (日本語版、移動)
    │   ├── Docker.md (日本語版、移動)
    │   ├── File-Upload.md (日本語版、移動)
    │   ├── Generate-Sample-Data.md (日本語版、移動)
    │   ├── index.md (日本語版、移動)
    │   ├── Initial-Setup.md (日本語版、移動)
    │   ├── License-Survey.md (日本語版、移動)
    │   ├── Logging.md (日本語版、移動)
    │   ├── Metabase.md (日本語版、移動)
    │   ├── Partner-Idp-Auth0-Login.md (日本語版、移動)
    │   ├── Project-Structure.md (日本語版、移動)
    │   ├── Queue-Job.md (日本語版、移動)
    │   ├── Rate-Limit.md (日本語版、移動)
    │   ├── Release-Check.md (日本語版、移動)
    │   ├── Send-Mail.md (日本語版、移動)
    │   ├── Sharding.md (日本語版、移動)
    │   ├── Spec-Driven-Development.md (日本語版、移動)
    │   └── Testing.md (日本語版、移動)
    ├── en/
    │   ├── Admin.md (英語版、新規作成)
    │   ├── Apache-Superset.md (英語版、新規作成)
    │   ├── API.md (英語版、新規作成)
    │   ├── Architecture.md (英語版、新規作成)
    │   ├── Atlas-Operations.md (英語版、新規作成)
    │   ├── Command-Line-Tool.md (英語版、新規作成)
    │   ├── Database-Viewer.md (英語版、新規作成)
    │   ├── Docker.md (英語版、新規作成)
    │   ├── File-Upload.md (英語版、新規作成)
    │   ├── Generate-Sample-Data.md (英語版、新規作成)
    │   ├── index.md (英語版、新規作成)
    │   ├── Initial-Setup.md (英語版、新規作成)
    │   ├── License-Survey.md (英語版、新規作成)
    │   ├── Logging.md (英語版、新規作成)
    │   ├── Metabase.md (英語版、新規作成)
    │   ├── Partner-Idp-Auth0-Login.md (英語版、新規作成)
    │   ├── Project-Structure.md (英語版、新規作成)
    │   ├── Queue-Job.md (英語版、新規作成)
    │   ├── Rate-Limit.md (英語版、新規作成)
    │   ├── Release-Check.md (英語版、新規作成)
    │   ├── Send-Mail.md (英語版、新規作成)
    │   ├── Sharding.md (英語版、新規作成)
    │   ├── Spec-Driven-Development.md (英語版、新規作成)
    │   └── Testing.md (英語版、新規作成)
    └── pages/
        ├── en/
        │   ├── about.md (既存、維持)
        │   ├── index.md (既存、維持)
        │   └── setup.md (既存、維持)
        └── ja/
            ├── about.md (既存、維持)
            ├── index.md (既存、維持)
            └── setup.md (既存、維持)
```

### ファイル変更の詳細

#### 変更1: README.mdのリネームと英語版作成

**変更前**:
- `README.md`: 日本語版（831行）

**変更後**:
- `README.ja.md`: 日本語版（既存のREADME.mdからリネーム）
- `README.md`: 英語版（新規作成、既存の内容を英語に翻訳）

**作業内容**:
1. 既存の`README.md`を`README.ja.md`にリネーム
2. `README.ja.md`の内容を英語に翻訳して`README.md`を作成
3. リンクや参照先を適切に更新
4. `README.md`と`README.ja.md`の冒頭または末尾に言語切替リンクを追加
   - `README.md`（英語版）: `**[日本語](README.ja.md) | [English]**`
   - `README.ja.md`（日本語版）: `**[日本語]** | [English](README.md)`

#### 変更2: docsディレクトリ内のファイル移動

**対象ファイル一覧（24ファイル）**:
1. Admin.md
2. Apache-Superset.md
3. API.md
4. Architecture.md
5. Atlas-Operations.md
6. Command-Line-Tool.md
7. Database-Viewer.md
8. Docker.md
9. File-Upload.md
10. Generate-Sample-Data.md
11. index.md
12. Initial-Setup.md
13. License-Survey.md
14. Logging.md
15. Metabase.md
16. Partner-Idp-Auth0-Login.md
17. Project-Structure.md
18. Queue-Job.md
19. Rate-Limit.md
20. Release-Check.md
21. Send-Mail.md
22. Sharding.md
23. Spec-Driven-Development.md
24. Testing.md

**作業内容**:
1. `docs/ja/`ディレクトリを作成（存在しない場合）
2. 上記24ファイルを`docs/`から`docs/ja/`に移動
3. 各ファイルの内容を英語に翻訳して`docs/en/`に作成
4. 各ファイル（日本語版と英語版）の冒頭または末尾に言語切替リンクを追加
   - 日本語版（`docs/ja/*.md`）: `**[日本語]** | [English](../en/{filename})`
   - 英語版（`docs/en/*.md`）: `**[日本語](../ja/{filename}) | [English]**`

#### 変更3: docs/_config.ymlの更新

**変更前**:
```yaml
description: "大量ユーザー、大量アクセスの運用に耐えうるGo APIサーバー、データベース構成のテンプレートプロジェクト"
```

**変更後**:
```yaml
description: "Template project for Go API server and database configuration that can withstand large-scale users and high traffic operations"
```

**作業内容**:
1. `docs/_config.yml`を開く
2. `description`フィールドを英語に変更
3. その他の設定は変更しない

### 技術的詳細

#### 翻訳方針
- **マスタードキュメント**: 日本語版をマスターとする
- **翻訳品質**: 技術用語の一貫性を保つ
- **コード例**: コード例やコマンド例はそのまま維持
- **構造**: ドキュメントの構造、セクション、見出しレベルを維持
- **リンク**: 相対パスや絶対パスを適切に更新

#### リンクの更新
- ドキュメント内の相互参照リンクを適切に更新
- `docs/ja/`や`docs/en/`へのパスを考慮
- GitHub PagesのURL構造を考慮
- README.md内のリンクを更新（`docs/ja/`や`docs/en/`への参照）

#### 言語切替リンクの実装
- **配置**: 各Markdownファイルの冒頭または末尾に追加
- **形式**: `**[日本語] | [English]**`
- **実装方法**:
  - README.md（英語版）: `**[日本語](README.ja.md) | [English]**`
  - README.ja.md（日本語版）: `**[日本語]** | [English](README.md)`
  - docs/ja/*.md（日本語版）: `**[日本語]** | [English](../en/{filename})`
  - docs/en/*.md（英語版）: `**[日本語](../ja/{filename}) | [English]**`
- **表示**: 現在の言語は太字（`**`）で表示し、リンクは通常のテキストで表示

#### 既存構造の維持
- `docs/pages/ja/`と`docs/pages/en/`は既に存在し、維持する
- Jekyllの設定（`_config.yml`）は既存の構造を維持
- GitHub Pagesの設定は既存の構造を維持

### アーキテクチャの整合性
- **既存パターンの維持**: 既存のドキュメント構造とパターンを維持
- **多言語対応の構造**: 将来的に他の言語を追加しやすい構造
- **保守性**: 日本語版と英語版の両方を維持できる構造

## Implementation Details

### 実装手順

#### Step 1: ディレクトリ構造の準備
1. `docs/ja/`ディレクトリを作成（存在しない場合）
2. `docs/en/`ディレクトリを作成（存在しない場合）
3. ディレクトリが正しく作成されたことを確認

#### Step 2: README.mdのリネーム
1. 既存の`README.md`を`README.ja.md`にリネーム
2. リネームが成功したことを確認
3. `README.ja.md`の内容が正しいことを確認

#### Step 3: docsディレクトリ内のファイル移動
1. 24個のMarkdownファイルを`docs/`から`docs/ja/`に移動
2. 各ファイルが正しく移動されたことを確認
3. ファイル数が24個であることを確認

#### Step 4: README.mdの英語版作成
1. `README.ja.md`の内容を英語に翻訳
2. 構造、セクション、コード例を維持
3. リンクや参照先を適切に更新
4. 冒頭または末尾に言語切替リンクを追加: `**[日本語](README.ja.md) | [English]**`
5. `README.ja.md`の冒頭または末尾にも言語切替リンクを追加: `**[日本語]** | [English](README.md)`
6. `README.md`として保存

#### Step 5: docs/en/内の英語版ファイル作成
1. `docs/ja/`内の各ファイルの内容を英語に翻訳
2. 24個のファイルすべてを翻訳
3. 各ファイルの構造、セクション、コード例を維持
4. リンクや参照先を適切に更新
5. 各ファイル（英語版）の冒頭または末尾に言語切替リンクを追加: `**[日本語](../ja/{filename}) | [English]**`
6. `docs/ja/`内の各ファイル（日本語版）の冒頭または末尾にも言語切替リンクを追加: `**[日本語]** | [English](../en/{filename})`
7. `docs/en/`に保存

#### Step 6: docs/_config.ymlの更新
1. `docs/_config.yml`を開く
2. `description`フィールドを英語に変更
3. その他の設定は変更しない
4. YAMLの構文が正しいことを確認

#### Step 7: リンクの更新と動作確認
1. README.md内のリンクを更新（`docs/ja/`や`docs/en/`への参照）
2. docs/en/内のファイル間の相互参照リンクを更新
3. docs/ja/内のファイル間の相互参照リンクを更新（必要に応じて）
4. 各ドキュメントファイルが正常に表示されることを確認
5. リンクが正しく機能することを確認

#### Step 8: ビルドと表示の確認
1. GitHub Pagesのビルドが正常に動作することを確認（該当する場合）
2. Jekyllのビルドが正常に動作することを確認（該当する場合）
3. ビルドエラーが発生しないことを確認
4. 英語版と日本語版の両方が正常に表示されることを確認

### エラーハンドリング
- ファイル移動時のエラー: ファイルが存在しない、権限エラーなどの確認
- 翻訳時のエラー: 翻訳品質の確認、技術用語の一貫性の確認
- リンク更新時のエラー: リンクが正しく機能することを確認
- ビルド時のエラー: YAML構文エラー、Markdown構文エラーの確認

### パフォーマンスへの影響
- ファイル数の増加: 24ファイルが48ファイルになる（日本語版24 + 英語版24）
- ビルド時間: ファイル数が増えるため、ビルド時間が若干増加する可能性がある
- ストレージ: ファイルサイズが約2倍になる
- **影響**: 軽微（ドキュメントファイルは比較的小さい）

## Testing Strategy

### ファイル構造の確認
- [ ] `README.md`が英語版になっている
- [ ] `README.ja.md`が日本語版として存在する
- [ ] `docs/ja/`ディレクトリに24個の日本語版Markdownファイルが存在する
- [ ] `docs/en/`ディレクトリに24個の英語版Markdownファイルが存在する
- [ ] 既存の`docs/pages/ja/`と`docs/pages/en/`の構造は維持されている

### 内容の確認
- [ ] README.mdの英語版が適切に翻訳されている
- [ ] docs/en/内の全ファイルが適切に翻訳されている
- [ ] 技術用語の一貫性が保たれている
- [ ] コード例やコマンド例が正しく動作する
- [ ] ドキュメントの構造が維持されている

### リンクの確認
- [ ] README.md内のリンクが正しく機能する
- [ ] docs/en/内のファイル間の相互参照リンクが正しく機能する
- [ ] docs/ja/内のファイル間の相互参照リンクが正しく機能する
- [ ] 外部リンクが正しく機能する
- [ ] README.mdとREADME.ja.mdの言語切替リンクが正しく機能する
- [ ] docs/ja/内の全ファイルの言語切替リンクが正しく機能する
- [ ] docs/en/内の全ファイルの言語切替リンクが正しく機能する

### 設定ファイルの確認
- [ ] `docs/_config.yml`のdescriptionが英語になっている
- [ ] 既存の設定が壊れていない
- [ ] YAMLの構文が正しい

### ビルドと表示の確認
- [ ] GitHub Pagesが正常に動作する（該当する場合）
- [ ] Jekyllのビルドが正常に動作する（該当する場合）
- [ ] ビルドエラーが発生しない
- [ ] 英語版と日本語版の両方が正常に表示される

### 手動テスト
1. **ファイル構造の確認**
   - すべてのファイルが正しい場所に存在することを確認
   - ファイル数が正しいことを確認

2. **内容の確認**
   - 英語版の翻訳品質を確認
   - 技術用語の一貫性を確認
   - コード例やコマンド例が正しいことを確認

3. **リンクの確認**
   - 各ドキュメント内のリンクが正しく機能することを確認
   - 相互参照リンクが正しく機能することを確認
   - 言語切替リンクが正しく機能することを確認（README.md、README.ja.md、docs/ja/*.md、docs/en/*.md）

4. **ビルドの確認**
   - GitHub PagesやJekyllのビルドが正常に動作することを確認
   - エラーが発生しないことを確認

## Migration/Rollback Strategy

### 移行戦略
この実装はファイルの移動、リネーム、新規作成が主な作業のため、以下の手順で実装する：

1. **バックアップの作成**（推奨）
   - 既存のREADME.mdとdocsディレクトリをバックアップ
   - Gitでコミット前の状態を確認

2. **段階的な実装**
   - Step 1: ディレクトリ構造の準備
   - Step 2: README.mdのリネーム
   - Step 3: docsディレクトリ内のファイル移動
   - Step 4: README.mdの英語版作成
   - Step 5: docs/en/内の英語版ファイル作成
   - Step 6: docs/_config.ymlの更新
   - Step 7: リンクの更新と動作確認
   - Step 8: ビルドと表示の確認

3. **各ステップでの確認**
   - 各ステップの完了後に動作確認
   - 問題が発生した場合は、そのステップで停止

### ロールバック戦略
変更がファイルの移動、リネーム、新規作成のため、ロールバックも比較的簡単：

1. **Gitを使用したロールバック**（推奨）
   - `git reset --hard HEAD`で変更を元に戻す
   - または、変更前のコミットに戻る

2. **手動でのロールバック**
   - `README.ja.md`を`README.md`にリネーム
   - `docs/ja/`内のファイルを`docs/`に移動
   - `docs/en/`ディレクトリを削除
   - `docs/_config.yml`のdescriptionを日本語に戻す

### リスク評価
- **低リスク**: ファイルの移動、リネーム、新規作成のみ
- **影響範囲**: ドキュメントファイルのみ（コードファイルへの影響なし）
- **既存機能への影響**: なし（ドキュメントの表示のみ）
- **データ損失のリスク**: 低（Gitで管理されているため）

## Security Considerations

### セキュリティへの影響
- 既存のセキュリティ機能に影響を与えない
- ドキュメントファイルのみの変更のため、セキュリティリスクは低い
- コードファイルへの影響はない

### セキュリティの確認事項
- ドキュメント内に機密情報が含まれていないことを確認
- リンクが正しく機能し、不正なリダイレクトが発生しないことを確認

## Performance Considerations

### パフォーマンスへの影響
- **ファイル数の増加**: 24ファイルが48ファイルになる
- **ビルド時間**: ファイル数が増えるため、ビルド時間が若干増加する可能性がある
- **ストレージ**: ファイルサイズが約2倍になる
- **影響**: 軽微（ドキュメントファイルは比較的小さい）

### パフォーマンスの確認事項
- ビルド時間が許容範囲内であることを確認
- ファイルサイズの増加が許容範囲内であることを確認

## Documentation

### ドキュメント更新
- 既存のドキュメント構造を維持
- 新しいドキュメント構造（`docs/ja/`と`docs/en/`）を追加
- 既存の`docs/pages/ja/`と`docs/pages/en/`は維持

### コードコメント
- コード内のコメントは従来通り日本語を使用
- ドキュメントの英語化はコード内コメントには影響しない

## Dependencies

### 依存関係
- **Markdown**: ドキュメント形式（既存のまま）
- **Jekyll**: GitHub Pagesのビルドツール（該当する場合、既存のまま）
- **GitHub Pages**: ドキュメント公開プラットフォーム（該当する場合、既存のまま）

### 依存関係の変更
- 依存関係の変更は不要
- 新しい依存関係の追加は不要

## Future Considerations

### 将来の拡張性
- 将来的に他の言語を追加しやすい構造になっている
- `docs/ja/`と`docs/en/`の構造を他の言語にも適用可能
- 例: `docs/fr/`（フランス語）、`docs/de/`（ドイツ語）など

### 保守性
- 日本語版と英語版の両方を維持できる構造
- マスタードキュメント（日本語版）を優先的に更新し、英語版はそれに追従
- 翻訳の品質を維持するためのプロセスを確立

### 自動化の可能性
- 将来的に、翻訳の自動化ツールを導入する可能性がある
- ただし、現時点では手動翻訳を想定
