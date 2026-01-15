**[日本語](../ja/Partner-Idp-Auth0-Login.md) | [English]**

# Auth0 External ID Federation - Setup & Development Guide

## 1. Overview and Architecture

This document describes the configuration steps for implementing login (single sign-on) to **MINE (our service)** using accounts from **PARTNER (external service)**.

### Architecture

While PARTNER's environment is incomplete, we set up a **"Mock PARTNER (IdP)"** within Auth0 and connect Auth0 instances via OIDC to enable early development.

* **MINE**: Our service (the one requesting authentication / Service Provider)
* **PARTNER**: Partner's service (the one providing authentication / Identity Provider)
* **Auth0**: Hub mediating both parties

---

## 2. Development: Creating Mock PARTNER (IdP Role)

Set up a dummy configuration to provide a "login screen" instead of PARTNER.

1. **Create Application**
* `Applications` > `Applications` > **[Create Application]**
* **Name**: `Mock-Partner-System`
* **Type**: `Regular Web Applications`


2. **Get Authentication Credentials (Note these down)**
* Note the **Domain**, **Client ID**, **Client Secret** from the Settings tab.


3. **Create Test User**
* `User Management` > `Users` > **[Create User]**
* Create with an email address that might exist in PARTNER (e.g., `test@partner.com`).



---

## 3. MINE-side Connection Settings (Enterprise Connection)

Settings for redirecting from MINE to PARTNER.

1. **Create OIDC Connection**
* `Authentication` > `Enterprise` > `OpenID Connect` > **[Create New]**


2. **Enter Basic Information**
* **Connection Name**: `partner-oidc` (internal identifier)
* **Display Name**: `Login with PARTNER account` (button display text)


3. **Endpoint Settings** (use values noted in step 2)
* **Discovery URL**: `https://{Domain from step 2}/.well-known/openid-configuration`
* **Client ID / Secret**: Enter values noted in step 2
* **Token Endpoint Auth Method**: `Post`


4. **Specify Scopes**
* `openid profile email` (required for user attribute retrieval)


5. **Copy Information (Important)**
* After saving, copy the **`Callback URL`** shown at the bottom of the details screen (e.g., `https://.../login/callback`).



---

## 4. Mutual Authentication Permission Settings

Allow communication between MINE and Mock PARTNER.

1. Go to `Applications` > `Applications` > **Mock-Partner-System**.
2. Paste the **URL copied in step 3-5** into the **Allowed Callback URLs** field.
3. Click **[Save Changes]** at the bottom.

---

## 5. Linking to Application (MINE itself)

1. Go to `Applications` > `Applications` > **(your MINE app)**.
2. In the `Connections` tab > `Enterprise` section, turn **`partner-oidc`** **ON**.

---

## 5.1. Callback URL Settings for Development Environment

Set up callback URLs for the Next.js application (PARTNER) to integrate with Auth0.

1. Go to `Applications` > `Applications` > **(PARTNER app)**.
2. Open the `Settings` tab.
3. Add to **Allowed Callback URLs**:
   ```
   http://localhost:3000/api/auth/callback/auth0
   ```
4. Add to **Allowed Logout URLs**:
   ```
   http://localhost:3000
   ```
5. Click **[Save Changes]** at the bottom.

### URL List by Environment

| Environment | Callback URL | Logout URL |
| --- | --- | --- |
| Development | `http://localhost:3000/api/auth/callback/auth0` | `http://localhost:3000` |
| Staging | `https://staging.example.com/api/auth/callback/auth0` | `https://staging.example.com` |
| Production | `https://example.com/api/auth/callback/auth0` | `https://example.com` |

**Note**: When setting multiple URLs, add them comma-separated.

---

## 5.2. API Settings (For JWT Format Access Token)

Create an API to obtain JWT format access tokens from Auth0. Without this setting, Auth0 returns opaque tokens (random strings), making server-side JWT verification impossible.

### API Creation Steps in Auth0 Dashboard

1. Log in to Auth0 Dashboard
2. Go to `Applications` > `APIs`
3. Click **[+ Create API]**
4. Enter the following:
   - **Name**: `go-webdb-template API`
   - **Identifier**: `https://go-webdb-template/api` (any identifier, URL format recommended though not required)
   - **Signing Algorithm**: `RS256`
5. Click **[Create]**
6. Open `Settings` tab
7. Turn **ON** `Allow Offline Access`
8. Click `Save Changes` at the bottom

### Environment Variable Settings

Set the created API's Identifier as an environment variable.

Add to `client/.env.local` (when using NextAuth (Auth.js) v5):
```
AUTH0_AUDIENCE=https://go-webdb-template/api
AUTH0_SCOPE='openid profile email offline_access'
```

### Configuration Examples by Environment

| Environment | AUTH0_AUDIENCE | AUTH0_SCOPE |
| --- | --- | --- |
| Development | `https://go-webdb-template/api` | `openid profile email offline_access` |
| Staging | `https://go-webdb-template/api` | `openid profile email offline_access` |
| Production | `https://go-webdb-template/api` | `openid profile email offline_access` |

**Notes**:
- The Audience value must exactly match the API's Identifier.
- Including `offline_access` scope enables refresh token retrieval.

---

## 6. Testing (Try Connection)

1. Go to `Authentication` > `Enterprise` > `OpenID Connect`.
2. Click **[Try] (eye icon)** on the right side of `partner-oidc`.
3. When login screen appears in another window, log in with the dummy user from step 2-3.
4. If **"It works!"** is displayed and user info JSON is returned, it's successful.

---

## 7. Adding Development Members & Team Management

Invite members involved in the project and set up a collaborative development environment.

### Member Invitation Steps

1. Click **[Settings]** (gear icon) at bottom left of dashboard > **[Tenant Members]**.
2. Click **[+ Add Member]**.
3. **Email**: Enter the member's email address to invite.
4. **Role (Permissions)**:
* `Admin`: Full permissions including settings changes and member management.
* `Editor`: Can change app settings but cannot add administrators.


5. Complete when they approve via the link in their email.

### About Auth0 Teams

When operating multiple environments (development, staging, production), using the "Teams" feature to group and manage tenants centrally is advised.

---

## 8. Switching Items for Production (When PARTNER is Ready)

When official information arrives from PARTNER, update the following values in `Enterprise Connection`.

| Item | Change Content |
| --- | --- |
| **Discovery URL** | Change to official OIDC document URL provided by PARTNER |
| **Client ID / Secret** | Change to official credentials issued by PARTNER |
| **Allowed Callback URL** | Request PARTNER to register MINE's Callback URL |
