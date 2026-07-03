# Chirpy

Chirpy is a minimal, Twitter-like social media API service built in Go and powered by PostgreSQL.

It provides a production-ready template demonstrating modern backend design in Go, featuring secure stateless authentication, database migrations, custom sorting, type-safe SQL query generation, and external webhook integrations.

## Why Chirpy?
* **Dual-Token Authentication Flow**: Implements short-lived JWT Access Tokens combined with long-lived database-backed Refresh Tokens for robust security and session revocation.
* **Go + PostgreSQL Core**: Utilizes standard Go libraries (`net/http`, `database/sql`) for speed, predictability, and minimal dependency overhead.
* **Automated Tooling**: Powered by [Goose](https://github.com/pressly/goose) for database migrations and [sqlc](https://sqlc.dev/) for type-safe SQL queries.
* **Architecture Separation**: Virtually decoupled into:
  * `/app` - A fileserver serving static frontend web assets.
  * `/api` - A REST API managing backend resources.

_This project was completed as a guided curriculum on [Boot.dev](https://boot.dev)._


## Quick Start & Installation

### Prerequisites
Make sure you have the following installed:
* [Go](https://go.dev/) (1.22+)
* [PostgreSQL](https://www.postgresql.org/)
* [Goose](https://github.com/pressly/goose) (for migrations)

### 1. Clone & Setup Configuration
Clone the repository and create your configuration file:
```bash
cp .env.example .env
```
Ensure the variables (`DB_URL`, `JWT_SECRET`, `POLKA_KEY`) are set according to your database configurations.

### 2. Run Database Migrations
Apply the PostgreSQL schemas to your local database:
```bash
cd sql/schema
goose postgres "your-database-connection-string" up
```

### 3. Start the Server
Run the application from the root directory:
```bash
go run .
```
The server will start listening at `http://localhost:8080`.

---


## Authentication Flow
Chirpy implements standard JWT-based authentication using short-lived Access Tokens and long-lived Refresh Tokens.

| Endpoint | Method | Description | Auth Required |
| :--- | :--- | :--- | :--- |
| `POST /api/login` | `POST` | Authenticate credentials and retrieve tokens. | None |
| `POST /api/refresh` | `POST` | Exchange a valid Bearer Refresh Token for a new Access Token. | Bearer Refresh Token |
| `POST /api/revoke` | `POST` | Revoke a Bearer Refresh Token, preventing future refreshes. | Bearer Refresh Token |

### Auth Request Payload Sample
```json
{
  "email": "user@example.com",
  "password": "yourpassword"
}
```

---

## Users
Manage user registrations and profile details.

### Create User
* **Method & Path**: `POST /api/users`
* **Auth**: None
* **Payload**:
```json
{
  "email": "user@example.com",
  "password": "yourpassword"
}
```

### Update User
* **Method & Path**: `PUT /api/users`
* **Auth**: Bearer Access Token
* **Payload**:
```json
{
  "email": "updated@example.com",
  "password": "newpassword"
}
```

---

## Chirps
Submit and retrieve short posts.

### Create Chirp
* **Method & Path**: `POST /api/chirps`
* **Auth**: Bearer Access Token
* **Payload**:
```json
{
  "body": "This is my chirp body."
}
```

### Get All Chirps
* **Method & Path**: `GET /api/chirps`
* **Optional Query Params**: 
  * `author_id` (UUID string) - Filter chirps by a specific author.
  * `sort` (string) - Sort direction, either `asc` (default) or `desc`.

### Get Single Chirp
* **Method & Path**: `GET /api/chirps/{chirpID}`
* **Path Params**: `chirpID` (UUID)

### Delete Chirp
* **Method & Path**: `DELETE /api/chirps/{chirpID}`
* **Auth**: Bearer Access Token (only owner can delete)
* **Path Params**: `chirpID` (UUID)

---

## Webhooks

### Polka Upgrade Webhook
* **Method & Path**: `POST /api/polka/webhooks`
* **Auth**: API Key (sent via `Authorization: ApiKey <key>` header)
* **Payload**:
```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "3311741c-680c-4546-99f3-fc9efac2036c"
  }
}
```

---

## Admin & Miscellaneous

### System Health
* `GET /api/healthz` - Verifies application health check status.

### Reset Database (Dev Mode Only)
* `POST /admin/reset` - Deletes all database rows and resets server metrics.

### Server Metrics
* `GET /admin/metrics` - Renders an HTML page showing total file server page visits.