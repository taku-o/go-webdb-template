**[日本語](../ja/License-Survey.md) | [English]**

# License Survey Results

This document contains the survey results on licenses, commercial use availability, and pricing for libraries and services used in the go-webdb-template project.

**Survey Date**: January 2025

## Table of Contents

1. [Go Libraries](#go-libraries)
2. [JavaScript/TypeScript Libraries](#javascripttypescript-libraries)
3. [Development Tools](#development-tools)
4. [Docker Container Services](#docker-container-services)
5. [SaaS Services](#saas-services)
6. [Summary](#summary)

---

## Go Libraries

### Web Framework

| Library | Version | License | Commercial Use | Notes |
|---------|---------|---------|----------------|-------|
| `github.com/labstack/echo/v4` | v4.13.3 | MIT | ✅ Yes | No restrictions |
| `github.com/danielgtaylor/huma/v2` | v2.34.1 | MIT | ✅ Yes | No restrictions |
| `github.com/gorilla/mux` | v1.8.1 | BSD 3-Clause | ✅ Yes | No restrictions |

### Database

| Library | Version | License | Commercial Use | Notes |
|---------|---------|---------|----------------|-------|
| `gorm.io/gorm` | v1.25.12 | MIT | ✅ Yes | No restrictions |
| `gorm.io/driver/postgres` | v1.5.9 | MIT | ✅ Yes | No restrictions |
| `gorm.io/driver/mysql` | v1.5.7 | MIT | ✅ Yes | No restrictions |
| `github.com/lib/pq` | v1.10.9 | MIT | ✅ Yes | PostgreSQL driver |

### Authentication & Security

| Library | Version | License | Commercial Use | Notes |
|---------|---------|---------|----------------|-------|
| `github.com/golang-jwt/jwt/v5` | v5.3.0 | MIT | ✅ Yes | No restrictions |
| `github.com/MicahParks/keyfunc/v2` | v2.1.0 | Apache 2.0 | ✅ Yes | No restrictions |

### Cache & Job Queue

| Library | Version | License | Commercial Use | Notes |
|---------|---------|---------|----------------|-------|
| `github.com/redis/go-redis/v9` | v9.17.2 | MIT | ✅ Yes | No restrictions |
| `github.com/hibiken/asynq` | v0.25.1 | MIT | ✅ Yes | No restrictions |

**⚠️ Important**: Redis itself changed its license in March 2024. `github.com/redis/go-redis` is a Redis client library licensed under MIT and is commercially usable. However, when using the Redis server itself, please check the requirements of Redis Ltd.'s new license (RSALv2/SSPLv1).

### Rate Limiting

| Library | Version | License | Commercial Use | Notes |
|---------|---------|---------|----------------|-------|
| `github.com/ulule/limiter/v3` | v3.11.2 | MIT | ✅ Yes | No restrictions |

### File Upload

| Library | Version | License | Commercial Use | Notes |
|---------|---------|---------|----------------|-------|
| `github.com/tus/tusd/v2` | v2.8.0 | MIT | ✅ Yes | No restrictions |

### Admin Panel

| Library | Version | License | Commercial Use | Notes |
|---------|---------|---------|----------------|-------|
| `github.com/GoAdminGroup/go-admin` | v1.2.26 | Apache 2.0 | ✅ Yes | No restrictions |
| `github.com/GoAdminGroup/themes` | v0.0.48 | Apache 2.0 | ✅ Yes | No restrictions |

### AWS SDK

| Library | Version | License | Commercial Use | Notes |
|---------|---------|---------|----------------|-------|
| `github.com/aws/aws-sdk-go-v2` | v1.41.0 | Apache 2.0 | ✅ Yes | No restrictions |
| `github.com/aws/aws-sdk-go-v2/service/s3` | v1.95.0 | Apache 2.0 | ✅ Yes | No restrictions |
| `github.com/aws/aws-sdk-go-v2/service/ses` | v1.34.17 | Apache 2.0 | ✅ Yes | No restrictions |

### Others

| Library | Version | License | Commercial Use | Notes |
|---------|---------|---------|----------------|-------|
| `github.com/google/uuid` | v1.6.0 | BSD-3-Clause | ✅ Yes | No restrictions |
| `github.com/sirupsen/logrus` | v1.9.0 | MIT | ✅ Yes | No restrictions |
| `github.com/spf13/viper` | v1.21.0 | MIT | ✅ Yes | No restrictions |
| `github.com/brianvoe/gofakeit/v6` | v6.28.0 | MIT | ✅ Yes | No restrictions |
| `golang.org/x/crypto` | v0.46.0 | BSD 3-Clause | ✅ Yes | No restrictions |
| `gopkg.in/mail.v2` | v2.3.1 | MIT | ✅ Yes | No restrictions |
| `gopkg.in/natefinch/lumberjack.v2` | v2.2.1 | MIT | ✅ Yes | No restrictions |
| `github.com/avast/retry-go/v4` | v4.7.0 | MIT | ✅ Yes | No restrictions |
| `github.com/stretchr/testify` | v1.11.1 | MIT | ✅ Yes | For testing |

---

## JavaScript/TypeScript Libraries

### Framework

| Library | Version | License | Commercial Use | Notes |
|---------|---------|---------|----------------|-------|
| `next` | ^14.1.0 | MIT | ✅ Yes | No restrictions |
| `react` | ^18.2.0 | MIT | ✅ Yes | No restrictions |
| `react-dom` | ^18.2.0 | MIT | ✅ Yes | No restrictions |

### Authentication

| Library | Version | License | Commercial Use | Notes |
|---------|---------|---------|----------------|-------|
| `@auth0/nextjs-auth0` | ^4.14.0 | MIT | ✅ Yes | SDK is free. Auth0 service fees are separate (see below) |

### File Upload

| Library | Version | License | Commercial Use | Notes |
|---------|---------|---------|----------------|-------|
| `@uppy/core` | ^5.2.0 | MIT | ✅ Yes | No restrictions |
| `@uppy/dashboard` | ^5.1.0 | MIT | ✅ Yes | No restrictions |
| `@uppy/react` | ^5.1.1 | MIT | ✅ Yes | No restrictions |
| `@uppy/tus` | ^5.1.0 | MIT | ✅ Yes | No restrictions |

### Testing

| Library | Version | License | Commercial Use | Notes |
|---------|---------|---------|----------------|-------|
| `@playwright/test` | ^1.41.0 | Apache 2.0 | ✅ Yes | No restrictions |
| `jest` | ^29.7.0 | MIT | ✅ Yes | No restrictions |
| `@storybook/react` | ^7.6.0 | MIT | ✅ Yes | Development tool |
| `@testing-library/react` | ^14.1.2 | MIT | ✅ Yes | For testing |
| `msw` | ^2.0.13 | MIT | ✅ Yes | For testing |

### Others

| Library | Version | License | Commercial Use | Notes |
|---------|---------|---------|----------------|-------|
| `typescript` | ^5.3.3 | Apache 2.0 | ✅ Yes | No restrictions |
| `tailwindcss` | ^3.4.1 | MIT | ✅ Yes | No restrictions |
| `eslint` | ^8.56.0 | MIT | ✅ Yes | Development tool |

---

## Development Tools

### Atlas CLI

| Item | Details |
|------|---------|
| **License** | Apache 2.0 (Open Source) / PRO License (Paid features) |
| **Commercial Use** | ✅ Yes |
| **Pricing** | Open Source: Free / PRO: Paid |
| **Notes** | Database migration management tool (by Ariga).<br><br>**Open Source (Apache 2.0)**:<br>- Basic migration features are free for commercial use<br>- No restrictions under Apache 2.0<br><br>**PRO (Paid License)**:<br>- Advanced features like View (database view testing/validation)<br>- Other PRO features: Migration Linting & Safety Checks, Code Review Guardrails, Atlas Copilot (AI assistance), Schema Policy & Governance, Drift Detection & Monitoring, Audit Trails & Change History<br>- Additional database engine support (SQL Server, ClickHouse, Redshift, Oracle, Spanner, Snowflake, Databricks, etc.)<br><br>**PRO Pricing** (as of January 2025):<br>- Developer seat: $9/month/seat<br>- CI/CD project: $59/month/project (includes 2 databases)<br>- Additional database: $39/month/database<br>- Free trial: 30 days (up to 10 seats)<br><br>*Check [Atlas official site](https://atlasgo.io/pricing/) for latest pricing. |

---

## Docker Container Services

### CloudBeaver

| Item | Details |
|------|---------|
| **License** | Apache 2.0 |
| **Commercial Use** | ✅ Yes (Community Edition) |
| **Pricing** | Free (Community Edition) |
| **Notes** | - Community Edition is commercially usable under Apache 2.0<br>- Enterprise Edition requires paid subscription (with additional features/support) |

### Metabase

| Item | Details |
|------|---------|
| **License** | AGPL v3 (Open Source Edition) |
| **Commercial Use** | ⚠️ Conditional |
| **Pricing** | Free (Open Source Edition) or Paid plans |
| **Notes** | **Important**: AGPL v3 license has the following requirements:<br>- Source code must be disclosed (when providing over network)<br>- Derivative works must also be disclosed under AGPL v3<br>- Commercial License purchase required if you want to avoid source code disclosure<br><br>**Paid Plans (Commercial License)**:<br><br>*Pricing may change. Check [Metabase official site](https://www.metabase.com/pricing/) for latest info.*<br><br>**As of January 2025**:<br>- Pro Plan: $575/month (10 users, +$12/user/month additional)<br>- Enterprise Plan: From $20,000/year (custom pricing available)<br><br>**Past Info (Reference)**:<br>- Starter: $85/month (5 users, +$5/user/month additional)<br>- Pro: $500/month (10 users, +$10/user/month additional)<br>- Enterprise: From $15,000/year |

### Apache Superset

| Item | Details |
|------|---------|
| **License** | Apache 2.0 |
| **Commercial Use** | ✅ Yes |
| **Pricing** | Free |
| **Notes** | - Commercially usable under Apache 2.0 (no restrictions)<br>- No source code disclosure obligation<br>- Modification and distribution freely allowed<br>- Many companies use in production (Airbnb, American Express, Dropbox, Lyft, Netflix, Twitter, Udemy, etc.)<br>- Software is free, but production operation may incur infrastructure/hosting/technical resource costs<br>- Managed services may have separate fees |

### Redis

| Item | Details |
|------|---------|
| **License** | RSALv2 / SSPLv1 (since March 2024) |
| **Commercial Use** | ⚠️ Conditional |
| **Pricing** | Free (self-use) |
| **Notes** | **Important**: License changed in March 2024:<br>- **RSALv2**: Commercial license required when providing Redis as a hosted service<br>- **SSPLv1**: Source code disclosure required when providing as a service<br>- No issues for internal use or embedding in applications<br>- Commercial contract with Redis Ltd. required when providing Redis as a managed service<br><br>**Alternative**: Consider Redict (BSD-licensed fork) |

### Mailpit

| Item | Details |
|------|---------|
| **License** | MIT |
| **Commercial Use** | ✅ Yes |
| **Pricing** | Free |
| **Notes** | Email testing tool for development/test environments |

---

## SaaS Services

### AWS S3 (Simple Storage Service)

| Item | Details |
|------|---------|
| **License** | Proprietary (AWS service) |
| **Commercial Use** | ✅ Yes |
| **Pricing** | Pay-as-you-go |
| **Pricing Details** | **Storage** (varies by region):<br>- S3 Standard: $0.023/GB/month (first 50TB)<br>- S3 Standard-IA: $0.0125/GB/month<br>- S3 One Zone-IA: $0.01/GB/month<br>- S3 Glacier Instant Retrieval: $0.004/GB/month<br>- S3 Glacier Flexible Retrieval: $0.0036/GB/month<br>- S3 Glacier Deep Archive: $0.00099/GB/month<br><br>**Requests**:<br>- PUT, COPY, POST, LIST: $0.005/1,000 requests<br>- GET, SELECT: $0.0004/1,000 requests<br><br>**Data Transfer**:<br>- Internet outbound: First 100GB/month free, then $0.09/GB<br><br>*Pricing varies by region. Check [AWS official site](https://aws.amazon.com/s3/pricing/) for latest. |

### AWS SES (Simple Email Service)

| Item | Details |
|------|---------|
| **License** | Proprietary (AWS service) |
| **Commercial Use** | ✅ Yes |
| **Pricing** | Pay-as-you-go |
| **Pricing Details** | **Email Sending**:<br>- $0.10/1,000 emails ($0.0001/email)<br><br>**Attachments**:<br>- $0.12/GB<br><br>**Dedicated IP**:<br>- $24.95/month/IP<br><br>**Free Tier**:<br>- New customers: Up to 3,000 emails/month free for first 12 months<br><br>*Check [AWS official site](https://aws.amazon.com/ses/pricing/) for latest. |

### Auth0

| Item | Details |
|------|---------|
| **License** | Proprietary (Auth0 service) |
| **Commercial Use** | ✅ Yes (Commercial use allowed even on Free Plan) |
| **Pricing** | Free plan available, Paid plans available |
| **Pricing Details** | **Free Plan** (Updated September 2024):<br>- Monthly Active Users (MAU): Up to 25,000<br>- Unlimited social/Okta connections<br>- Custom domain: 1 (credit card verification required)<br>- Passwordless auth (SMS, Email, Passkey, OTP)<br>- Organizations: Up to 5<br>- SSO feature<br>- Community support<br><br>**Paid Plans**:<br>- Essentials: $35/month (500 MAU, +$0.07/MAU additional)<br>- Professional: $240/month (500 MAU, +$0.07/MAU additional)<br>- Enterprise: Custom pricing<br><br>*Check [Auth0 official site](https://auth0.com/pricing/) for latest. |

---

## Summary

### Commercial Use Summary

#### ✅ Commercially Usable (No Restrictions)

- **Go Libraries**: All MIT, Apache 2.0, BSD-family permissive licenses
- **JavaScript/TypeScript Libraries**: All MIT, Apache 2.0 permissive licenses
- **Development Tools**: Atlas CLI (Apache 2.0, OSS version free, PRO version paid)
- **CloudBeaver**: Apache 2.0 (Community Edition)
- **Apache Superset**: Apache 2.0
- **Mailpit**: MIT
- **AWS S3/SES**: Pay-as-you-go commercial services
- **Auth0**: Commercial use allowed even on Free Plan

#### ⚠️ Requires Attention for Commercial Use

1. **Metabase** (AGPL v3)
   - Has source code disclosure requirements
   - Commercial License purchase required to avoid source code disclosure for commercial use

2. **Redis** (RSALv2/SSPLv1)
   - No issues for internal use or embedding in applications
   - Commercial contract required when providing Redis as a managed service
   - This project uses `github.com/redis/go-redis` (MIT), client library itself is fine

### Pricing Summary

#### Available for Free

- All Go/JavaScript libraries (no license fees)
- Atlas CLI (Open Source, Apache 2.0)
- CloudBeaver Community Edition
- Apache Superset (Apache 2.0)
- Metabase Open Source Edition (when meeting AGPL requirements)
- Mailpit
- Auth0 Free Plan (up to 25,000 MAU)

#### Pay-as-you-go

- **AWS S3**: Based on storage capacity, request count, data transfer
- **AWS SES**: Based on email count (first 12 months: up to 3,000/month free)

#### Paid Plans (Optional)

- **Atlas CLI PRO**: When advanced features like Views are needed (from $9/month/seat)
- **Metabase**: When avoiding source code disclosure, or needing additional features
- **Auth0**: When exceeding Free Plan limits, or needing additional features

### Suggestions

1. **Regarding Metabase Use**
   - Development environment: No issues with AGPL v3
   - Production environment: Can continue with AGPL v3 if source code disclosure is possible
   - Consider Commercial License purchase if source code disclosure is not possible

2. **Regarding Redis Use**
   - This project uses Redis client library (MIT), no issues
   - When operating Redis server internally, check RSALv2/SSPLv1 requirements
   - Commercial contract with Redis Ltd. required when providing as managed service

3. **Regarding AWS Services**
   - Development environment can reduce costs by using local storage and Mailpit
   - Production environment incurs pay-as-you-go charges, conduct cost estimation

4. **Regarding Auth0 Use**
   - Free Plan supports up to 25,000 MAU, sufficient for small commercial applications
   - Credit card verification required for custom domain use

5. **Regarding Apache Superset Use**
   - Commercially usable under Apache 2.0, no source code disclosure obligation
   - Fewer restrictions for commercial use compared to Metabase (AGPL v3)
   - Many companies have production usage track record
   - Software is free, but production operation may incur infrastructure/hosting/technical resource costs

---

## Reference Links

- [AWS S3 Pricing](https://aws.amazon.com/s3/pricing/)
- [AWS SES Pricing](https://aws.amazon.com/ses/pricing/)
- [Auth0 Pricing](https://auth0.com/pricing/)
- [Apache Superset Official Site](https://superset.apache.org/)
- [Apache Superset GitHub](https://github.com/apache/superset)
- [Metabase License](https://www.metabase.com/license/)
- [Metabase Pricing](https://www.metabase.com/pricing/)
- [Redis License Change](https://redis.io/license)
- [Atlas CLI Official Site](https://atlasgo.io/)
- [Atlas CLI Pricing](https://atlasgo.io/pricing/)

---

**Last Updated**: January 2025
