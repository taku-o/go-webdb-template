# ライセンス調査結果

本ドキュメントは、go-webdb-templateプロジェクトで使用しているライブラリ・サービスのライセンスと商用利用可否、利用料金についての調査結果です。

**調査日**: 2025年1月

## 目次

1. [Goライブラリ](#goライブラリ)
2. [JavaScript/TypeScriptライブラリ](#javascripttypescriptライブラリ)
3. [開発ツール](#開発ツール)
4. [Dockerコンテナサービス](#dockerコンテナサービス)
5. [SaaSサービス](#saasサービス)
6. [まとめ](#まとめ)

---

## Goライブラリ

### Webフレームワーク

| ライブラリ | バージョン | ライセンス | 商用利用可否 | 注意事項 |
|-----------|----------|----------|------------|---------|
| `github.com/labstack/echo/v4` | v4.13.3 | MIT | ✅ 可 | 制限なし |
| `github.com/danielgtaylor/huma/v2` | v2.34.1 | MIT | ✅ 可 | 制限なし |
| `github.com/gorilla/mux` | v1.8.1 | BSD 3-Clause | ✅ 可 | 制限なし |

### データベース

| ライブラリ | バージョン | ライセンス | 商用利用可否 | 注意事項 |
|-----------|----------|----------|------------|---------|
| `gorm.io/gorm` | v1.25.12 | MIT | ✅ 可 | 制限なし |
| `gorm.io/driver/postgres` | v1.5.9 | MIT | ✅ 可 | 制限なし |
| `gorm.io/driver/mysql` | v1.5.7 | MIT | ✅ 可 | 制限なし |
| `github.com/lib/pq` | v1.10.9 | MIT | ✅ 可 | PostgreSQLドライバー |

### 認証・セキュリティ

| ライブラリ | バージョン | ライセンス | 商用利用可否 | 注意事項 |
|-----------|----------|----------|------------|---------|
| `github.com/golang-jwt/jwt/v5` | v5.3.0 | MIT | ✅ 可 | 制限なし |
| `github.com/MicahParks/keyfunc/v2` | v2.1.0 | Apache 2.0 | ✅ 可 | 制限なし |

### キャッシュ・ジョブキュー

| ライブラリ | バージョン | ライセンス | 商用利用可否 | 注意事項 |
|-----------|----------|----------|------------|---------|
| `github.com/redis/go-redis/v9` | v9.17.2 | MIT | ✅ 可 | 制限なし |
| `github.com/hibiken/asynq` | v0.25.1 | MIT | ✅ 可 | 制限なし |

**⚠️ 重要**: Redis本体のライセンスは2024年3月に変更されました。`github.com/redis/go-redis`はRedisクライアントライブラリであり、MITライセンスで商用利用可能です。ただし、Redisサーバー自体を使用する場合は、Redis Ltd.の新しいライセンス（RSALv2/SSPLv1）の要件を確認してください。

### レートリミット

| ライブラリ | バージョン | ライセンス | 商用利用可否 | 注意事項 |
|-----------|----------|----------|------------|---------|
| `github.com/ulule/limiter/v3` | v3.11.2 | MIT | ✅ 可 | 制限なし |

### ファイルアップロード

| ライブラリ | バージョン | ライセンス | 商用利用可否 | 注意事項 |
|-----------|----------|----------|------------|---------|
| `github.com/tus/tusd/v2` | v2.8.0 | MIT | ✅ 可 | 制限なし |

### 管理画面

| ライブラリ | バージョン | ライセンス | 商用利用可否 | 注意事項 |
|-----------|----------|----------|------------|---------|
| `github.com/GoAdminGroup/go-admin` | v1.2.26 | Apache 2.0 | ✅ 可 | 制限なし |
| `github.com/GoAdminGroup/themes` | v0.0.48 | Apache 2.0 | ✅ 可 | 制限なし |

### AWS SDK

| ライブラリ | バージョン | ライセンス | 商用利用可否 | 注意事項 |
|-----------|----------|----------|------------|---------|
| `github.com/aws/aws-sdk-go-v2` | v1.41.0 | Apache 2.0 | ✅ 可 | 制限なし |
| `github.com/aws/aws-sdk-go-v2/service/s3` | v1.95.0 | Apache 2.0 | ✅ 可 | 制限なし |
| `github.com/aws/aws-sdk-go-v2/service/ses` | v1.34.17 | Apache 2.0 | ✅ 可 | 制限なし |

### その他

| ライブラリ | バージョン | ライセンス | 商用利用可否 | 注意事項 |
|-----------|----------|----------|------------|---------|
| `github.com/google/uuid` | v1.6.0 | BSD-3-Clause | ✅ 可 | 制限なし |
| `github.com/sirupsen/logrus` | v1.9.0 | MIT | ✅ 可 | 制限なし |
| `github.com/spf13/viper` | v1.21.0 | MIT | ✅ 可 | 制限なし |
| `github.com/brianvoe/gofakeit/v6` | v6.28.0 | MIT | ✅ 可 | 制限なし |
| `golang.org/x/crypto` | v0.46.0 | BSD 3-Clause | ✅ 可 | 制限なし |
| `gopkg.in/mail.v2` | v2.3.1 | MIT | ✅ 可 | 制限なし |
| `gopkg.in/natefinch/lumberjack.v2` | v2.2.1 | MIT | ✅ 可 | 制限なし |
| `github.com/avast/retry-go/v4` | v4.7.0 | MIT | ✅ 可 | 制限なし |
| `github.com/stretchr/testify` | v1.11.1 | MIT | ✅ 可 | テスト用 |

---

## JavaScript/TypeScriptライブラリ

### フレームワーク

| ライブラリ | バージョン | ライセンス | 商用利用可否 | 注意事項 |
|-----------|----------|----------|------------|---------|
| `next` | ^14.1.0 | MIT | ✅ 可 | 制限なし |
| `react` | ^18.2.0 | MIT | ✅ 可 | 制限なし |
| `react-dom` | ^18.2.0 | MIT | ✅ 可 | 制限なし |

### 認証

| ライブラリ | バージョン | ライセンス | 商用利用可否 | 注意事項 |
|-----------|----------|----------|------------|---------|
| `@auth0/nextjs-auth0` | ^4.14.0 | MIT | ✅ 可 | SDK自体は無料。Auth0サービスの利用料金は別途（後述） |

### ファイルアップロード

| ライブラリ | バージョン | ライセンス | 商用利用可否 | 注意事項 |
|-----------|----------|----------|------------|---------|
| `@uppy/core` | ^5.2.0 | MIT | ✅ 可 | 制限なし |
| `@uppy/dashboard` | ^5.1.0 | MIT | ✅ 可 | 制限なし |
| `@uppy/react` | ^5.1.1 | MIT | ✅ 可 | 制限なし |
| `@uppy/tus` | ^5.1.0 | MIT | ✅ 可 | 制限なし |

### テスト

| ライブラリ | バージョン | ライセンス | 商用利用可否 | 注意事項 |
|-----------|----------|----------|------------|---------|
| `@playwright/test` | ^1.41.0 | Apache 2.0 | ✅ 可 | 制限なし |
| `jest` | ^29.7.0 | MIT | ✅ 可 | 制限なし |
| `@storybook/react` | ^7.6.0 | MIT | ✅ 可 | 開発ツール |
| `@testing-library/react` | ^14.1.2 | MIT | ✅ 可 | テスト用 |
| `msw` | ^2.0.13 | MIT | ✅ 可 | テスト用 |

### その他

| ライブラリ | バージョン | ライセンス | 商用利用可否 | 注意事項 |
|-----------|----------|----------|------------|---------|
| `typescript` | ^5.3.3 | Apache 2.0 | ✅ 可 | 制限なし |
| `tailwindcss` | ^3.4.1 | MIT | ✅ 可 | 制限なし |
| `eslint` | ^8.56.0 | MIT | ✅ 可 | 開発ツール |

---

## 開発ツール

### Atlas CLI

| 項目 | 内容 |
|-----|------|
| **ライセンス** | Apache 2.0（オープンソース版）/ PROライセンス（有料機能） |
| **商用利用可否** | ✅ 可 |
| **利用料金** | オープンソース版: 無料 / PRO版: 有料 |
| **注意事項** | データベースマイグレーション管理ツール（Ariga開発）。<br><br>**オープンソース版（Apache 2.0）**:<br>- 基本的なマイグレーション機能は無料で商用利用可能<br>- Apache 2.0ライセンスにより制限なし<br><br>**PRO版（有料ライセンス）**:<br>- View機能（データベースビューのテスト・検証）など高度な機能が利用可能<br>- その他のPRO機能: Migration Linting & Safety Checks、Code Review Guardrails、Atlas Copilot（AI支援）、Schema Policy & Governance、Drift Detection & Monitoring、Audit Trails & Change History<br>- 追加データベースエンジンサポート（SQL Server、ClickHouse、Redshift、Oracle、Spanner、Snowflake、Databricksなど）<br><br>**PRO版の価格**（2025年1月時点）:<br>- 開発者シート: $9/月/シート<br>- CI/CDプロジェクト: $59/月/プロジェクト（データベース2つまで含む）<br>- 追加データベース: $39/月/データベース<br>- 無料トライアル: 30日間（10シートまで）<br><br>※ 最新の価格情報は[Atlas公式サイト](https://atlasgo.io/pricing/)で確認してください。 |

---

## Dockerコンテナサービス

### CloudBeaver

| 項目 | 内容 |
|-----|------|
| **ライセンス** | Apache 2.0 |
| **商用利用可否** | ✅ 可（Community Edition） |
| **利用料金** | 無料（Community Edition） |
| **注意事項** | - Community EditionはApache 2.0で商用利用可能<br>- Enterprise Editionは有料サブスクリプションが必要（追加機能・サポート付き） |

### Metabase

| 項目 | 内容 |
|-----|------|
| **ライセンス** | AGPL v3（Open Source Edition） |
| **商用利用可否** | ⚠️ 条件付き |
| **利用料金** | 無料（Open Source Edition）または有料プラン |
| **注意事項** | **重要**: AGPL v3ライセンスのため、以下の要件があります：<br>- ソースコードを公開する必要がある（ネットワーク経由で提供する場合）<br>- 派生作品もAGPL v3で公開する必要がある<br>- 商用利用でソースコード公開を避けたい場合は、Commercial Licenseの購入が必要<br><br>**有料プラン（Commercial License）**:<br><br>※ 価格情報は変動する可能性があります。最新の情報は[Metabase公式サイト](https://www.metabase.com/pricing/)で確認してください。<br><br>**2025年1月時点の情報**:<br>- Pro Plan: $575/月（10ユーザー、追加$12/ユーザー/月）<br>- Enterprise Plan: $20,000/年から（カスタム価格あり）<br><br>**過去の情報（参考）**:<br>- Starter: $85/月（5ユーザー、追加$5/ユーザー/月）<br>- Pro: $500/月（10ユーザー、追加$10/ユーザー/月）<br>- Enterprise: $15,000/年から |

### Apache Superset

| 項目 | 内容 |
|-----|------|
| **ライセンス** | Apache 2.0 |
| **商用利用可否** | ✅ 可 |
| **利用料金** | 無料 |
| **注意事項** | - Apache 2.0ライセンスにより商用利用可能（制限なし）<br>- ソースコード公開の義務なし<br>- 修正・配布も自由に可能<br>- 多くの企業（Airbnb、American Express、Dropbox、Lyft、Netflix、Twitter、Udemyなど）が商用環境で利用実績あり<br>- ソフトウェア自体は無料だが、本番環境での運用にはインフラ・ホスティング・技術リソースのコストが発生する可能性がある<br>- マネージドサービスを利用する場合は別途料金が発生する可能性がある |

### Redis

| 項目 | 内容 |
|-----|------|
| **ライセンス** | RSALv2 / SSPLv1（2024年3月以降） |
| **商用利用可否** | ⚠️ 条件付き |
| **利用料金** | 無料（自己利用） |
| **注意事項** | **重要**: 2024年3月にライセンスが変更されました：<br>- **RSALv2**: Redisをホストサービスとして提供する場合は商用ライセンスが必要<br>- **SSPLv1**: サービスとして提供する場合はソースコード公開が必要<br>- 自社内での使用やアプリケーションに組み込む場合は問題なし<br>- Redisをマネージドサービスとして提供する場合は、Redis Ltd.との商用契約が必要<br><br>**代替案**: Redict（BSDライセンスのフォーク）を検討可能 |

### Mailpit

| 項目 | 内容 |
|-----|------|
| **ライセンス** | MIT |
| **商用利用可否** | ✅ 可 |
| **利用料金** | 無料 |
| **注意事項** | 開発・テスト環境用のメールテストツール |

---

## SaaSサービス

### AWS S3（Simple Storage Service）

| 項目 | 内容 |
|-----|------|
| **ライセンス** | プロプライエタリ（AWS提供サービス） |
| **商用利用可否** | ✅ 可 |
| **利用料金** | 従量課金制 |
| **料金詳細** | **ストレージ料金**（リージョンにより異なる）:<br>- S3 Standard: $0.023/GB/月（最初の50TB）<br>- S3 Standard-IA: $0.0125/GB/月<br>- S3 One Zone-IA: $0.01/GB/月<br>- S3 Glacier Instant Retrieval: $0.004/GB/月<br>- S3 Glacier Flexible Retrieval: $0.0036/GB/月<br>- S3 Glacier Deep Archive: $0.00099/GB/月<br><br>**リクエスト料金**:<br>- PUT, COPY, POST, LIST: $0.005/1,000リクエスト<br>- GET, SELECT: $0.0004/1,000リクエスト<br><br>**データ転送料金**:<br>- インターネットへの送信: 最初の100GB/月は無料、以降$0.09/GB<br><br>※ 料金はリージョンにより異なります。最新の料金は[AWS公式サイト](https://aws.amazon.com/s3/pricing/)を確認してください。 |

### AWS SES（Simple Email Service）

| 項目 | 内容 |
|-----|------|
| **ライセンス** | プロプライエタリ（AWS提供サービス） |
| **商用利用可否** | ✅ 可 |
| **利用料金** | 従量課金制 |
| **料金詳細** | **メール送信**:<br>- $0.10/1,000通（$0.0001/通）<br><br>**添付ファイル**:<br>- $0.12/GB<br><br>**専用IPアドレス**:<br>- $24.95/月/IP<br><br>**無料枠**:<br>- 新規顧客は最初の12ヶ月間、月3,000通まで無料<br><br>※ 最新の料金は[AWS公式サイト](https://aws.amazon.com/ses/pricing/)を確認してください。 |

### Auth0

| 項目 | 内容 |
|-----|------|
| **ライセンス** | プロプライエタリ（Auth0提供サービス） |
| **商用利用可否** | ✅ 可（Free Planでも商用利用可能） |
| **利用料金** | 無料プランあり、有料プランあり |
| **料金詳細** | **Free Plan**（2024年9月更新）:<br>- 月間アクティブユーザー（MAU）: 25,000人まで<br>- 無制限のソーシャル・Okta接続<br>- カスタムドメイン: 1つ（クレジットカード認証が必要）<br>- パスワードレス認証（SMS、メール、Passkey、OTP）<br>- 組織: 5つまで<br>- SSO機能<br>- コミュニティサポート<br><br>**有料プラン**:<br>- Essentials: $35/月（500 MAU、追加$0.07/MAU）<br>- Professional: $240/月（500 MAU、追加$0.07/MAU）<br>- Enterprise: カスタム価格<br><br>※ 最新の料金は[Auth0公式サイト](https://auth0.com/pricing/)を確認してください。 |

---

## まとめ

### 商用利用可否の総括

#### ✅ 商用利用可能（制限なし）

- **Goライブラリ**: すべてMIT、Apache 2.0、BSD系のパーミッシブライセンス
- **JavaScript/TypeScriptライブラリ**: すべてMIT、Apache 2.0のパーミッシブライセンス
- **開発ツール**: Atlas CLI（Apache 2.0、オープンソース版は無料、PRO版は有料）
- **CloudBeaver**: Apache 2.0（Community Edition）
- **Apache Superset**: Apache 2.0
- **Mailpit**: MIT
- **AWS S3/SES**: 従量課金制の商用サービス
- **Auth0**: Free Planでも商用利用可能

#### ⚠️ 商用利用時に注意が必要

1. **Metabase**（AGPL v3）
   - ソースコード公開要件がある
   - 商用利用でソースコード公開を避けたい場合はCommercial Licenseの購入が必要

2. **Redis**（RSALv2/SSPLv1）
   - 自社内での使用やアプリケーションへの組み込みは問題なし
   - Redisをマネージドサービスとして提供する場合は商用契約が必要
   - 本プロジェクトでは`github.com/redis/go-redis`（MIT）を使用しており、クライアントライブラリ自体は問題なし

### 利用料金の総括

#### 無料で利用可能

- すべてのGo/JavaScriptライブラリ（ライセンス料なし）
- Atlas CLI（オープンソース版、Apache 2.0）
- CloudBeaver Community Edition
- Apache Superset（Apache 2.0）
- Metabase Open Source Edition（AGPL要件を満たす場合）
- Mailpit
- Auth0 Free Plan（25,000 MAUまで）

#### 従量課金制

- **AWS S3**: ストレージ容量・リクエスト数・データ転送量に応じた料金
- **AWS SES**: 送信メール数に応じた料金（最初の12ヶ月は月3,000通まで無料）

#### 有料プラン（オプション）

- **Atlas CLI PRO**: View機能など高度な機能が必要な場合（$9/月/シートから）
- **Metabase**: ソースコード公開を避けたい場合、または追加機能が必要な場合
- **Auth0**: Free Planの制限を超える場合、または追加機能が必要な場合

### 推奨事項

1. **Metabaseの使用について**
   - 開発環境での使用: AGPL v3で問題なし
   - 本番環境での使用: ソースコード公開が可能な場合はAGPL v3で継続利用可能
   - ソースコード公開ができない場合は、Commercial Licenseの購入を検討

2. **Redisの使用について**
   - 本プロジェクトではRedisクライアントライブラリ（MIT）を使用しており、問題なし
   - Redisサーバー自体を自社内で運用する場合は、RSALv2/SSPLv1の要件を確認
   - マネージドサービスとして提供する場合は、Redis Ltd.との商用契約が必要

3. **AWSサービスの使用について**
   - 開発環境ではローカルストレージやMailpitを使用することでコストを抑制可能
   - 本番環境では利用量に応じた従量課金が発生するため、コスト見積もりを実施

4. **Auth0の使用について**
   - Free Planで25,000 MAUまで対応可能なため、小規模な商用アプリケーションでは十分
   - カスタムドメインを使用する場合はクレジットカード認証が必要

5. **Apache Supersetの使用について**
   - Apache 2.0ライセンスにより商用利用可能で、ソースコード公開の義務なし
   - Metabase（AGPL v3）と比較して、商用利用時の制約が少ない
   - 多くの企業が商用環境で利用実績あり
   - ソフトウェア自体は無料だが、本番環境での運用にはインフラ・ホスティング・技術リソースのコストが発生する可能性がある

---

## 参考リンク

- [AWS S3 料金](https://aws.amazon.com/s3/pricing/)
- [AWS SES 料金](https://aws.amazon.com/ses/pricing/)
- [Auth0 料金](https://auth0.com/pricing/)
- [Apache Superset 公式サイト](https://superset.apache.org/)
- [Apache Superset GitHub](https://github.com/apache/superset)
- [Metabase ライセンス](https://www.metabase.com/license/)
- [Metabase 料金](https://www.metabase.com/pricing/)
- [Redis ライセンス変更について](https://redis.io/license)
- [Atlas CLI 公式サイト](https://atlasgo.io/)
- [Atlas CLI 料金](https://atlasgo.io/pricing/)

---

**最終更新**: 2025年1月
