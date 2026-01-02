# API Documentation

## Base URL

Development: `http://localhost:8080`

## OpenAPI Specification

このAPIはOpenAPI 3.0/3.1仕様に準拠しています。以下のエンドポイントからOpenAPI仕様を取得できます：

- **API Documentation UI**: `http://localhost:8080/docs` (Stoplight Elements)
- **OpenAPI 3.1 (JSON)**: `http://localhost:8080/openapi.json`
- **OpenAPI 3.1 (YAML)**: `http://localhost:8080/openapi.yaml`
- **OpenAPI 3.0.3 (JSON)**: `http://localhost:8080/openapi-3.0.json`

※ OpenAPIドキュメントエンドポイントは認証不要でアクセス可能です。

### フレームワーク

- **Echo v4**: 高性能なHTTPルーター
- **Huma v2**: OpenAPI仕様の自動生成とバリデーション

## Common Headers

```
Content-Type: application/json
Authorization: Bearer <API_TOKEN>
```

## Error Responses

All endpoints may return the following error responses:

### 400 Bad Request
```json
{
  "error": "Invalid request format"
}
```

### 404 Not Found
```json
{
  "error": "Resource not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error"
}
```

---

## User Endpoints

### Create User

**POST** `/api/users`

Creates a new user in the appropriate shard based on the assigned user ID.

**Request Body**:
```json
{
  "name": "John Doe",
  "email": "john@example.com"
}
```

**Response**: `201 Created`
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Sharding Note**: The user will be stored in a shard determined by `hash(id) % shard_count`.

---

### Get All Users

**GET** `/api/users`

Retrieves all users from all shards.

**Response**: `200 OK`
```json
[
  {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  },
  {
    "id": 2,
    "name": "Jane Smith",
    "email": "jane@example.com",
    "created_at": "2024-01-15T11:00:00Z",
    "updated_at": "2024-01-15T11:00:00Z"
  }
]
```

**Sharding Note**: This endpoint queries all shards in parallel and merges the results.

---

### Get User by ID

**GET** `/api/users/{id}`

Retrieves a specific user by ID.

**Path Parameters**:
- `id` (integer): User ID

**Response**: `200 OK`
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Sharding Note**: The request is automatically routed to the correct shard based on `hash(id)`.

---

### Update User

**PUT** `/api/users/{id}`

Updates an existing user.

**Path Parameters**:
- `id` (integer): User ID

**Request Body**:
```json
{
  "name": "John Updated",
  "email": "john.updated@example.com"
}
```

**Response**: `200 OK`
```json
{
  "id": 1,
  "name": "John Updated",
  "email": "john.updated@example.com",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T12:00:00Z"
}
```

**Sharding Note**: The update is performed on the shard where the user exists.

---

### Delete User

**DELETE** `/api/users/{id}`

Deletes a user by ID.

**Path Parameters**:
- `id` (integer): User ID

**Response**: `204 No Content`

**Sharding Note**: The deletion is performed on the shard where the user exists.

---

## Post Endpoints

### Create Post

**POST** `/api/posts`

Creates a new post for a user.

**Request Body**:
```json
{
  "user_id": 1,
  "title": "My First Post",
  "content": "This is the content of my first post."
}
```

**Response**: `201 Created`
```json
{
  "id": 1,
  "user_id": 1,
  "title": "My First Post",
  "content": "This is the content of my first post.",
  "created_at": "2024-01-15T13:00:00Z",
  "updated_at": "2024-01-15T13:00:00Z"
}
```

**Sharding Note**: The post is stored in the same shard as the user (based on `hash(user_id)`).

---

### Get All Posts

**GET** `/api/posts`

Retrieves all posts from all shards.

**Response**: `200 OK`
```json
[
  {
    "id": 1,
    "user_id": 1,
    "title": "My First Post",
    "content": "This is the content of my first post.",
    "created_at": "2024-01-15T13:00:00Z",
    "updated_at": "2024-01-15T13:00:00Z"
  },
  {
    "id": 2,
    "user_id": 2,
    "title": "Another Post",
    "content": "Content from another user.",
    "created_at": "2024-01-15T14:00:00Z",
    "updated_at": "2024-01-15T14:00:00Z"
  }
]
```

**Sharding Note**: This endpoint queries all shards in parallel and merges the results.

---

### Get Post by ID

**GET** `/api/posts/{id}`

Retrieves a specific post by ID.

**Path Parameters**:
- `id` (integer): Post ID

**Query Parameters**:
- `user_id` (integer, required): User ID who owns the post

**Response**: `200 OK`
```json
{
  "id": 1,
  "user_id": 1,
  "title": "My First Post",
  "content": "This is the content of my first post.",
  "created_at": "2024-01-15T13:00:00Z",
  "updated_at": "2024-01-15T13:00:00Z"
}
```

**Sharding Note**: Requires `user_id` to route to the correct shard.

---

### Get User Posts (JOIN)

**GET** `/api/user-posts`

Retrieves all posts with their associated user information (cross-shard JOIN).

**Response**: `200 OK`
```json
[
  {
    "user_id": 1,
    "user_name": "John Doe",
    "user_email": "john@example.com",
    "post_id": 1,
    "post_title": "My First Post",
    "post_content": "This is the content of my first post.",
    "created_at": "2024-01-15T13:00:00Z"
  },
  {
    "user_id": 2,
    "user_name": "Jane Smith",
    "user_email": "jane@example.com",
    "post_id": 2,
    "post_title": "Another Post",
    "post_content": "Content from another user.",
    "created_at": "2024-01-15T14:00:00Z"
  }
]
```

**Sharding Note**: This endpoint demonstrates cross-shard JOIN operations:
1. Fetches users from all shards
2. Fetches posts from all shards
3. Joins the data in application memory
4. Returns the merged result

---

### Update Post

**PUT** `/api/posts/{id}`

Updates an existing post.

**Path Parameters**:
- `id` (integer): Post ID

**Request Body**:
```json
{
  "user_id": 1,
  "title": "Updated Title",
  "content": "Updated content."
}
```

**Response**: `200 OK`
```json
{
  "id": 1,
  "user_id": 1,
  "title": "Updated Title",
  "content": "Updated content.",
  "created_at": "2024-01-15T13:00:00Z",
  "updated_at": "2024-01-15T15:00:00Z"
}
```

**Sharding Note**: Requires `user_id` to route to the correct shard.

---

### Delete Post

**DELETE** `/api/posts/{id}`

Deletes a post by ID.

**Path Parameters**:
- `id` (integer): Post ID

**Query Parameters**:
- `user_id` (integer, required): User ID who owns the post

**Response**: `204 No Content`

**Sharding Note**: Requires `user_id` to route to the correct shard for deletion.

---

## CORS Configuration

The API allows cross-origin requests from the following origins:
- `http://localhost:3000` (development client)

Allowed methods:
- GET, POST, PUT, DELETE, OPTIONS

Allowed headers:
- Content-Type

---

## Rate Limiting

APIレートリミット機能は実装済みです。IPアドレス単位でリクエスト数を制限します。

詳細は [Rate-Limit.md](Rate-Limit.md) を参照してください。

---

## Pagination

Currently not implemented. All list endpoints return all results.

**Future Enhancement**: Add pagination support with query parameters:
```
GET /api/users?page=1&limit=10
```

---

## Filtering and Sorting

Currently not implemented.

**Future Enhancement**: Add filtering and sorting:
```
GET /api/posts?user_id=1&sort=created_at&order=desc
```

---

## API Versioning

Currently not versioned. All endpoints are at the root `/api/` path.

**Future Enhancement**: Add versioning:
```
/api/v1/users
/api/v2/users
```

---

## Example Usage

### Creating a User and Post

```bash
# Create a user
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice", "email": "alice@example.com"}'

# Response: {"id": 3, "name": "Alice", ...}

# Create a post for that user
curl -X POST http://localhost:8080/api/posts \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 3,
    "title": "Hello World",
    "content": "My first post!"
  }'

# Get user posts with JOIN
curl http://localhost:8080/api/user-posts
```

### Error Handling Example

```bash
# Try to get a non-existent user
curl http://localhost:8080/api/users/999

# Response: 404 Not Found
# {"error": "User not found"}
```

---

## Client Integration

The Next.js client provides a typed API client in `client/src/lib/api.ts`:

```typescript
import { apiClient } from '@/lib/api'

// Create user
const user = await apiClient.createUser({
  name: 'Alice',
  email: 'alice@example.com'
})

// Get all users
const users = await apiClient.getUsers()

// Create post
const post = await apiClient.createPost({
  user_id: user.id,
  title: 'Hello',
  content: 'World'
})
```

See [Testing.md](./Testing.md) for API testing examples.
