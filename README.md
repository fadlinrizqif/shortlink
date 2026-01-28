# Shortlink Service

## Overview
Shortlink service adalah projek yang saya buat menggunakan Bahasa pemrograman Go. Aplikasi ini memungkinkan user untuk membuat short URL dan tetap akan mengarahkan ke URL yang asli. Aplikasi ini cocok untuk user yang perlu mengirimkan URL namun khawatir URL tersebut terlalu panjang maka aplikasi ini cocok untuk masalah tersebut.

## Features
- User Authentication (Login/Logout)
- JWT-based authentication with refresh token
- Generate short URL
- Redirect the URL to original URL
- Dashboard (server-side rendered)
- Middleware-base authentication

## Tech Stack
- Go (net/http)
- PostgreSQL
- Goose
- SQLC
- JWT
- HTML template Engine
- Cookies base authentication

## Project Stucture
```text
├── app                 # HTML Template folder
│   ├── dashboard.html  # Dashboard page for shorten the URL
│   ├── sigup.html      # Register page for new user
│   └── login.html      # Login page 
├── internal
│   ├── auth          # JWT & refresh token logic
│   ├── database      # query database (SQLC)
│   ├── handler       # HTTP handlers
│   ├── middleware    # authentication middleware
│   └── app           # shared config
├── sql
    ├── queries       # query to database 
│   └── schema        # migration for database using Goose
└── static
    └── css           # styling
```

## How to run 
```bash
git clone https://github.com/fadlinrizqif/shortlink.git
cd shortlink
go mod tidy
go run cmd/server/main.go
```

## Environment Variables
| Variable | Description |
|--------|------------|
| DB_URL | PostgreSQL connection string |
| SERVER_SECRET | JWT signing secret |

## Routes
| Method | Path | Description |
|------|------|------------|
| GET | /signup | Register page |
| POST | /signup | Register New user |
| GET | /login | Login page |
| POST | /login | Authenticate user |
| GET | /dashboard | User dashboard |
| POST | /dashboard | Create shortlink |
| GET | /{code} | Redirect to original URL |

## Security Notes
- Password Hashed
- JWT stored in HttpOnly cookies
- Refresh token stored in database 
- CSRF mitigated via SameSite cookie
