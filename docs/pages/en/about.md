---
layout: default
title: Project Overview
lang: en
---

# Project Overview

Go WebDB Template is a sample project with Go + Next.js + Database Sharding support. It adopts a configuration suitable for large-scale projects.

---

## Purpose

This project was created for the following purposes:

- **Demonstrating Scalability**: Providing an implementation example of horizontal scaling through database sharding
- **Presenting Best Practices**: Implementation examples of layered architecture, testing strategies, and environment-specific configuration management
- **Learning Resource**: Implementation patterns for reference in large-scale application development

---

## Technology Stack

### Server

| Item | Technology |
|------|------------|
| Language | Go 1.21+ |
| Architecture | Layered Architecture |
| Database | PostgreSQL/MySQL |
| ORM | GORM v1.25.12 |
| HTTP Router | Echo v4 |
| API Specification | Huma (OpenAPI auto-generation) |

### Client

| Item | Technology |
|------|------------|
| Framework | Next.js 14 (App Router) |
| Language | TypeScript 5+ |
| UI Components | shadcn/ui |
| Authentication | NextAuth (Auth.js) v5 |
| Styling | Tailwind CSS |

### Testing

| Item | Technology |
|------|------------|
| Server | Go testing, testify |
| Client Unit | Jest, React Testing Library |
| E2E | Playwright |
| API Mocking | MSW |

---

## Key Features

- **Sharding Support**: Table-based sharding (32 partitions) distributing data across multiple DBs
- **GORM Support**: Writer/Reader separation supported
- **GoAdmin Dashboard**: Web-based admin panel for data management
- **Layer Separation**: Clear separation of responsibilities across API, Usecase, Service, Repository, and DB layers
- **Environment-specific Settings**: Configuration switching for develop/staging/production environments
- **Type Safety**: Type definitions with TypeScript
- **Testing**: Unit/Integration/E2E test support
- **Rate Limiting**: API call limiting per IP address
- **Job Queue**: Background job processing using Redis + Asynq
- **Email Sending**: Email sending with stdout, Mailpit, and AWS SES support
- **File Upload**: Large file upload via TUS protocol
- **Logging**: Access logs, email sending logs, and SQL logs
- **Docker Support**: Dockerization of API server, Admin server, and client server

---

## Navigation

- [Home]({{ site.baseurl }}/pages/en/)
- [Setup Guide]({{ site.baseurl }}/pages/en/setup)
- [日本語]({{ site.baseurl }}/pages/ja/about)
