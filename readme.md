# 🔗 URL Shortener API

A simple and secure RESTful URL Shortener API built with **Go (net/http)**, **MySQL** (no ORM), and organized in **3-tier architecture** (handler, service, repository). It allows users to shorten long URLs, retrieve the original URL, update, delete, and view access statistics.

---

## 📦 Features

- Create a new short URL
- Retrieve original URL using short code
- Update long URL for a short code
- Delete a short URL
- View access statistics
- No authentication required
- Clean and testable 3-layer architecture

---

## 🧱 Project Structure
url-shortener/
├── main.go
├── app/database # database connection
├── common/util # handle common custom validations
├── shortenurl/
│ ├── handler/ # HTTP handlers
│ ├── service/ # Business logic
│ ├── repository/ # Database access layer
│ ├── model/ # Request and DB models
├── .env # Environment variables
├── go.mod
└── README.md


## 🚀 API Endpoints

### 1. Create Short URL

**POST /shorten**

#### Request

```json
{
  "url": "https://example.com"
}
Response (201):
{
  "id": 1,
  "url": "https://example.com",
  "shortCode": "abc123",
  "accessCount": 0,
  "createdAt": "2025-06-21T10:00:00Z",
  "updatedAt": "2025-06-21T10:00:00Z"
}

2. Get Original URL
GET /shorten/{code}

Returns the full short URL record and increments accessCount.


3. Update Short URL
PUT /shorten/{code}

{
  "url": "https://example.com/updated"
}
Response
{
  "id": 1,
  "url": "https://example.com/updated",
  "shortCode": "abc123",
  "accessCount": 12,
  "createdAt": "2025-06-21T10:00:00Z",
  "updatedAt": "2025-06-21T12:30:00Z"
}

4. Delete Short URL
DELETE /shorten/{code}

Response: 204 No Content

5. Get Stats
GET /shorten/{code}
Response
{
  "id": 1,
  "url": "https://example.com",
  "shortCode": "abc123",
  "accessCount": 12,
  "createdAt": "2025-06-21T10:00:00Z",
  "updatedAt": "2025-06-21T12:30:00Z"
}

Returns full URL record with current accessCount.

Setup Instructions

1. Clone the repo using git clone { git clone https://github.com/sainayani1/go-url-shortener}
cd url-shortener

Configure Environment
Create a .env file in the root:

DB_USER=root
DB_PASSWORD=yourpassword
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=url_shortener
PORT=8080

3. Start MySQL and Create DB

CREATE DATABASE url_shortener;

CREATE TABLE short_urls (
  id INT AUTO_INCREMENT PRIMARY KEY,
  url TEXT NOT NULL,
  short_code VARCHAR(255) NOT NULL UNIQUE,
  access_count INT DEFAULT 0,
  created_at DATETIME,
  updated_at DATETIME
);

. Run the Server
go mod tidy
go run main.go

Run tests
go test ./shortenurl/handler -v

