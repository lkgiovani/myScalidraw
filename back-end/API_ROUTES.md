# 🚀 MyScalidraw API Routes

## 📋 Table of Contents

- [Authentication](#authentication)
- [Files](#files)
- [System](#system)
- [Middlewares](#middlewares)
- [Status Codes](#status-codes)

---

## 🔐 Authentication

### Base URL: `/api/auth`

| Method | Route                | Description                    | Middleware           | Body                                              |
| ------ | -------------------- | ------------------------------ | -------------------- | ------------------------------------------------- |
| `POST` | `/create-first-user` | Create first system user       | ❌                   | [CreateFirstUserRequest](#createfirstuserrequest) |
| `GET`  | `/has-users`         | Check if users exist in system | ❌                   | -                                                 |
| `POST` | `/login`             | User login (sets cookie)       | `RequireSystemSetup` | [LoginRequest](#loginrequest)                     |
| `POST` | `/logout`            | User logout (removes cookie)   | `RequireSystemSetup` | -                                                 |
| `GET`  | `/me`                | Get current user data          | `RequireSystemSetup` | -                                                 |
| `GET`  | `/users`             | List all users                 | `RequireSystemSetup` | -                                                 |

---

## 📁 Files

### Base URL: `/api`

| Method   | Route               | Description       | Middleware                           | Body                                    |
| -------- | ------------------- | ----------------- | ------------------------------------ | --------------------------------------- |
| `GET`    | `/ping`             | Health check      | ❌                                   | -                                       |
| `GET`    | `/files`            | List all files    | `RequireSystemSetup` + `RequireAuth` | -                                       |
| `GET`    | `/files/:id`        | Get specific file | `RequireSystemSetup` + `RequireAuth` | -                                       |
| `POST`   | `/files`            | Create new file   | `RequireSystemSetup` + `RequireAuth` | [CreateFileRequest](#createfilerequest) |
| `POST`   | `/files/upload`     | Upload file       | `RequireSystemSetup` + `RequireAuth` | FormData                                |
| `PUT`    | `/files/:id`        | Save file content | `RequireSystemSetup` + `RequireAuth` | JSON Content                            |
| `PUT`    | `/files/:id/rename` | Rename file       | `RequireSystemSetup` + `RequireAuth` | `{"name": "string"}`                    |
| `DELETE` | `/files/:id`        | Delete file       | `RequireSystemSetup` + `RequireAuth` | -                                       |

---

## 📊 System

### Health Check

- **Route:** `GET /api/ping`
- **Description:** Check if API is working
- **Response:** `"pong"`

---

## 🛡️ Middlewares

### 1. RequireSystemSetup

- **Function:** Checks if at least one user exists in the system
- **Failure:** Returns `428 Precondition Required` if no users exist
- **Applied to:** All routes except `/create-first-user` and `/has-users`

### 2. RequireAuth

- **Function:** Verifies authentication via JWT cookie
- **Failure:** Returns `401 Unauthorized` if not authenticated
- **Applied to:** All file routes and some auth routes

### 3. RequireOwnerOrAdmin

- **Function:** Checks if user is owner or admin
- **Failure:** Returns `403 Forbidden` if not authorized

### 4. RequireOwner

- **Function:** Checks if user is owner
- **Failure:** Returns `403 Forbidden` if not owner

---

## 📝 Request Bodies

### CreateFirstUserRequest

```json
{
  "name": "string (min: 2, max: 100)",
  "email": "string (valid email)",
  "password": "string (min: 8 characters)"
}
```

### LoginRequest

```json
{
  "email": "string (valid email)",
  "password": "string"
}
```

### CreateFileRequest

```json
{
  "name": "string",
  "parentId": "string (optional)",
  "parentPath": "string (optional)",
  "isFolder": "boolean",
  "contentType": "string (optional)",
  "content": "string (optional)"
}
```

---

## 📤 Response Examples

### Success - Login

```json
{
  "message": "Login successful",
  "user": {
    "id": "uuid",
    "name": "John Doe",
    "email": "john@email.com",
    "type": "owner",
    "lastActivity": 1640995200000
  }
}
```

### Success - List Files

```json
[
  {
    "id": "uuid",
    "name": "My Drawing.excalidraw",
    "isFolder": false,
    "parentId": "uuid",
    "lastModified": 1640995200000,
    "path": "/My Drawing.excalidraw"
  }
]
```

### Error - System not configured

```json
{
  "error": "System setup required",
  "message": "No users found in system. Please create the first user.",
  "setup_needed": true
}
```

### Error - Not authenticated

```json
{
  "error": "Authentication required"
}
```

### Error - Invalid token

```json
{
  "error": "Invalid or expired token"
}
```

### Error - Access denied

```json
{
  "error": "Owner or admin access required"
}
```

---

## 🔢 Status Codes

| Code  | Description           | When it occurs                |
| ----- | --------------------- | ----------------------------- |
| `200` | OK                    | Successful operation          |
| `201` | Created               | Resource created successfully |
| `400` | Bad Request           | Invalid request data          |
| `401` | Unauthorized          | Missing or invalid token      |
| `403` | Forbidden             | User without permission       |
| `404` | Not Found             | Resource not found            |
| `409` | Conflict              | Conflict (e.g., user exists)  |
| `428` | Precondition Required | System needs configuration    |
| `500` | Internal Server Error | Internal server error         |

---

## 🍪 Cookie Authentication

### Authentication Cookie

- **Name:** `auth_token`
- **Type:** JWT (JSON Web Token)
- **Duration:** 24 hours
- **Settings:**
  - `HttpOnly: true` (not accessible via JavaScript)
  - `Secure: false` (true in production)
  - `SameSite: Lax`

### JWT Payload

```json
{
  "user_id": "uuid",
  "email": "user@email.com",
  "user_type": "owner|admin|guest",
  "exp": 1640995200,
  "iat": 1640908800,
  "nbf": 1640908800
}
```

---

## 🔄 Authentication Flow

1. **First Access:**

   - Check: `GET /api/auth/has-users`
   - If `false`: Create user via `POST /api/auth/create-first-user`

2. **Login:**

   - Login: `POST /api/auth/login`
   - Cookie is set automatically

3. **Route Access:**

   - Cookie is sent automatically
   - Middleware validates the token

4. **Logout:**

   - Logout: `POST /api/auth/logout`
   - Cookie is removed automatically

---

## 🚀 Usage Examples

### 1. Check if system is configured

```bash
curl -X GET http://localhost:8181/api/auth/has-users
```

### 2. Create first user

```bash
curl -X POST http://localhost:8181/api/auth/create-first-user \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Admin",
    "email": "admin@email.com",
    "password": "12345678"
  }'
```

### 3. Login

```bash
curl -X POST http://localhost:8181/api/auth/login \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{
    "email": "admin@email.com",
    "password": "12345678"
  }'
```

### 4. Access files (with cookie)

```bash
curl -X GET http://localhost:8181/api/files \
  -b cookies.txt
```

---

## 📚 User Types

| Type    | Description   | Permissions                         |
| ------- | ------------- | ----------------------------------- |
| `owner` | System owner  | Full access                         |
| `admin` | Administrator | Manage users and files              |
| `guest` | Guest         | Limited access based on permissions |

### Guest Subtypes

- `persistent`: Can access files anytime
- `temporary`: Can only access during active session

---

_Auto-generated documentation - MyScalidraw API v1.0_
