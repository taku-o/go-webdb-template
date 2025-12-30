# cc-sdd

## 概要
cc-sddは仕様駆動開発を行うためのツールです。
Cursor、Claude Code上で動作する。

## コマンド
```
## プロジェクトの解析
## (最初か、プロジェクトの構成が変わった時に実行する。)
/kiro:steering

## プロジェクトの要件・設計・タスクを作成(一括)
/kiro/spec-init "{プロジェクトの要件}"

## プロジェクトの要件を作成
/kiro:spec-requirements "{プロジェクトの要件}"

## 設計を作成
/kiro:spec-design

## タスク一覧を作成
/kiro:spec-tasks

## 指定タスクを実行する
/kiro:spec-impl "{タスク}"
```

## 利用手順

### Cursor上で作業 (ドキュメントの品質が良いため)
```
## プロジェクトを作ったら最初に1回実行
/kiro:steering

## プロジェクトの要件を作成
/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/3 に対応するための要件を作成してください。GitHub CLIは入っています。"
think.

## 要件定義書作成計画が表示された場合は
要件定義書を作成してください。

## 要件定義書の内容を確認したら
要件定義書を承認します。

## 設計を作成
/kiro:spec-design

## 設計書の内容を確認したら
設計書を承認します。

## タスク一覧を作成
/kiro:spec-tasks

## タスク一覧の内容を確認したら
タスクリストを承認します。

この要件の作業用のgitブランチを切ってください。
ここまでの作業をcommitしてください。
```

### Claude Code上で作業 (作業速度が速いため)

```
## 実装開始
/kiro:spec-impl 0003-gorm-introduction
```
