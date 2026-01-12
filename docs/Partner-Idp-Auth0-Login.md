# Auth0 外部ID連携 導入・開発ガイド

## 1. 概要と構成

本ドキュメントは、 **MINE（自社サービス）** において、 **PARTNER（外部サービス）** のアカウントを用いたログイン（シングルサインオン）を実現するための設定手順です。

### アーキテクチャ構成

PARTNERの環境が未完成の間、Auth0内に **「疑似PARTNER（IdP）」** を用意し、Auth0同士をOIDCで接続することで先行開発を可能にします。

* **MINE**: 我々のサービス（認証を依頼する側 / Service Provider）
* **PARTNER**: パートナーのサービス（認証を提供する側 / Identity Provider）
* **Auth0**: 両者を仲介するハブ

---

## 2. 開発用：疑似PARTNER（IdP役）の作成

PARTNERの代わりに「ログイン画面」を提供するためのダミー設定を行います。

1. **Application作成**
* `Applications` > `Applications` > **[Create Application]**
* **Name**: `Mock-Partner-System`
* **Type**: `Regular Web Applications`


2. **認証情報の取得（メモしておく）**
* Settingsタブにある **Domain**, **Client ID**, **Client Secret** を手元に控えます。


3. **テストユーザーの作成**
* `User Management` > `Users` > **[Create User]**
* PARTNERに存在しそうなメールアドレス（例：`test@partner.com`）で作成します。



---

## 3. MINE側の接続設定 (Enterprise Connection)

MINEからPARTNERへリダイレクトするための設定です。

1. **OIDC接続の作成**
* `Authentication` > `Enterprise` > `OpenID Connect` > **[Create New]**


2. **基本情報の入力**
* **Connection Name**: `partner-oidc`（内部識別子）
* **Display Name**: `PARTNERアカウントでログイン`（ボタンの表示名）


3. **エンドポイント等の設定**（手順2で控えた内容を使用）
* **Discovery URL**: `https://{手順2のDomain}/.well-known/openid-configuration`
* **Client ID / Secret**: 手順2で控えたものを入力
* **Token Endpoint Auth Method**: `Post`


4. **Scopesの指定**
* `openid profile email` （ユーザー属性取得に必須）


5. **情報のコピー（重要）**
* 保存後、詳細画面の下部に表示される **`Callback URL`**（例: `https://.../login/callback`）をコピーします。



---

## 4. 相互認証の許可設定

MINEと疑似PARTNER間の通信を許可します。

1. `Applications` > `Applications` > **Mock-Partner-System** を開く。
2. **Allowed Callback URLs** の欄に、**手順3-5でコピーしたURL** を貼り付けます。
3. 画面下部の **[Save Changes]** をクリック。

---

## 5. アプリケーション（MINE本体）への紐付け

1. `Applications` > `Applications` > **(自分のMINEアプリ)** を選択。
2. `Connections` タブ > `Enterprise` セクションの **`partner-oidc`** を **ON** にします。

---

## 5.1. 開発環境用のコールバックURL設定

Next.jsアプリケーション（PARTNER）がAuth0と連携するために、コールバックURLを設定します。

1. `Applications` > `Applications` > **(PARTNERアプリ)** を選択。
2. `Settings` タブを開く。
3. **Allowed Callback URLs** に以下を追加:
   ```
   http://localhost:3000/api/auth/callback/auth0
   ```
4. **Allowed Logout URLs** に以下を追加:
   ```
   http://localhost:3000
   ```
5. 画面下部の **[Save Changes]** をクリック。

### 各環境のURL一覧

| 環境 | Callback URL | Logout URL |
| --- | --- | --- |
| 開発環境 | `http://localhost:3000/api/auth/callback/auth0` | `http://localhost:3000` |
| ステージング | `https://staging.example.com/api/auth/callback/auth0` | `https://staging.example.com` |
| 本番 | `https://example.com/api/auth/callback/auth0` | `https://example.com` |

**注意**: 複数のURLを設定する場合は、カンマ区切りで追加できます。

---

## 5.2. API設定（JWT形式のアクセストークン取得用）

Auth0からJWT形式のアクセストークンを取得するために、APIを作成します。この設定がないと、Auth0はopaque token（ランダム文字列）を返すため、サーバー側でのJWT検証ができません。

### Auth0ダッシュボードでのAPI作成手順

1. Auth0ダッシュボードにログイン
2. `Applications` > `APIs` に移動
3. **[+ Create API]** をクリック
4. 以下の情報を入力:
   - **Name**: `go-webdb-template API`
   - **Identifier**: `https://go-webdb-template/api` （任意の識別子、URLである必要はないがURL形式が推奨）
   - **Signing Algorithm**: `RS256`
5. **[Create]** をクリック
6. `Settings` タブを開く
7. `Allow Offline Access` を **ON** にする
8. 画面下部の `Save Changes` をクリック

### 環境変数の設定

作成したAPIのIdentifierを環境変数に設定します。

`client/.env.local` に以下を追加（NextAuth (Auth.js) v5を使用する場合）:
```
AUTH0_AUDIENCE=https://go-webdb-template/api
AUTH0_SCOPE='openid profile email offline_access'
```

### 各環境の設定例

| 環境 | AUTH0_AUDIENCE | AUTH0_SCOPE |
| --- | --- | --- |
| 開発環境 | `https://go-webdb-template/api` | `openid profile email offline_access` |
| ステージング | `https://go-webdb-template/api` | `openid profile email offline_access` |
| 本番 | `https://go-webdb-template/api` | `openid profile email offline_access` |

**注意**:
- Audienceの値はAPIのIdentifierと完全に一致する必要があります。
- `offline_access`スコープを含めることで、リフレッシュトークンが取得できます。

---

## 6. 動作確認 (Try Connection)

1. `Authentication` > `Enterprise` > `OpenID Connect` を開く。
2. `partner-oidc` の右側にある **[Try]（目のアイコン）** をクリック。
3. 別ウィンドウでログイン画面が出たら、手順2-3のダミーユーザーでログイン。
4. **"It works!"** と表示され、ユーザー情報のJSONが返ってくれば成功です。

---

## 7. 開発メンバーの追加・チーム管理

プロジェクトに関わるメンバーを招待し、共同開発体制を整えます。

### メンバーの招待手順

1. ダッシュボード左下 **[Settings]** (歯車アイコン) > **[Tenant Members]** を選択。
2. **[+ Add Member]** をクリック。
3. **Email**: 招待するメンバーのメールアドレスを入力。
4. **Role (権限)**:
* `Admin`: 設定変更やメンバー管理が可能なフル権限。
* `Editor`: アプリ設定の変更は可能だが、管理者の追加は不可。


5. 相手に届いたメールのリンクを承認すれば完了です。

### Auth0 Teams について

複数の環境（開発、ステージング、本番）を運用する際は、「Teams」機能を使用して各テナントをグループ化し、一括管理することを推奨します。

---

## 8. 本番（PARTNER完成時）の切り替え項目

PARTNERから正式な情報が届いたら、`Enterprise Connection` の以下の値を更新します。

| 項目 | 変更内容 |
| --- | --- |
| **Discovery URL** | PARTNERが提供する正式なOIDCドキュメントのURLへ変更 |
| **Client ID / Secret** | PARTNERから発行された正式な認証情報へ変更 |
| **Allowed Callback URL** | MINEのCallback URLをPARTNER側へ登録依頼する |


