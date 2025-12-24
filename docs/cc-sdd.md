# cc-sdd

## 概要
使用駆動開発を行うためのツール。
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
```
## プロジェクトを作ったら
/kiro:steering

## プロジェクトの要件を作成
/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/3 に対応するための要件を作成してください。GitHub CLIは入っています。"

## 要件定義書の内容を確認したら
## 設計を作成
/kiro:spec-design

## 設計書の内容を確認したら
## タスク一覧を作成
/kiro:spec-tasks

## タスク一覧の内容を確認したら

この要件の作業用のgitブランチを切ってください。
ここまでの作業をcommitしてください。

## 実装開始
/kiro:spec-impl 0003-gorm-introduction
```
