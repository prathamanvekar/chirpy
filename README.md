# Chirpy

Chirpy is a minimal, Twitter-like API service built with Go and PostgreSQL.

## Architecture
The application is separated into two components:
* `/app` - A fileserver serving frontend web assets.
* `/api` - A REST API managing backend resources.

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
* **Optional Query Param**: `author_id` (UUID string) to filter by a specific author.

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