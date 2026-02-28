# MS GA Identifier - Identity Service

A robust, enterprise-ready microservice for identity and authentication management built with Go.

## 🏗️ Architecture

The service follows **Clean Architecture** principles with clear separation of concerns:

- **API Layer**: HTTP handlers and routing using Gin framework
- **Service Layer**: Core business logic for authentication and authorization
- **Repository Layer**: Data persistence using GORM (PostgreSQL)
- **Domain Layer**: Core business entities and interfaces
- **Infrastructure Layer**: External service integrations (Redis, Kafka, Auth Service)

## 🚀 Features

- User registration and login with JWT authentication
- Token refresh mechanism
- Password reset flow with secure token generation
- Email verification support
- Session management with refresh tokens
- Integration with auth service for roles and permissions
- Event-driven architecture with Kafka integration
- Redis caching for improved performance

## 🛠️ Technology Stack

- **Language**: Go 1.23+
- **Framework**: Gin
- **Database**: PostgreSQL 15
- **Cache**: Redis
- **Message Queue**: Apache Kafka
- **ORM**: GORM

## 📋 Prerequisites

- Go 1.23 or higher
- Docker and Docker Compose
- PostgreSQL 15
- Redis 7

## 🏃 Getting Started

### 1. Start Infrastructure

```bash
docker-compose up -d
```

This will start:

- PostgreSQL (port 5432)
- Redis (port 6379)
- Kafka (port 9092)
- Zookeeper

### 2. Configure Environment

Copy `.env.example` to `.env` and update the values:

```bash
cp .env.example .env
```

### 3. Run Migrations

```bash
make migrate
```

### 4. Run the Service

```bash
make run
```

The server will start on port 8081 by default.

## 📖 API Documentation

### Public Endpoints

#### Register

```http
POST /identity/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "first_name": "John",
  "last_name": "Doe"
}
```

#### Login

```http
POST /identity/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

#### Refresh Token

```http
POST /identity/refresh
Content-Type: application/json

{
  "refresh_token": "your-refresh-token"
}
```

#### Forgot Password

```http
POST /identity/forgot-password
Content-Type: application/json

{
  "email": "user@example.com"
}
```

#### Reset Password

```http
POST /identity/reset-password
Content-Type: application/json

{
  "token": "reset-token",
  "new_password": "NewSecurePass123!"
}
```

### Protected Endpoints

All protected endpoints require a valid JWT token in the Authorization header:

```http
Authorization: Bearer <jwt-token>
```

#### Logout

```http
POST /identity/logout
```

#### Get Current User

```http
GET /identity/me
```

#### Change Password

```http
POST /identity/change-password
Content-Type: application/json

{
  "current_password": "OldPass123!",
  "new_password": "NewPass123!"
}
```

## 🧪 Testing

Run unit tests:

```bash
make test
```

## 📁 Project Structure

```
ms-ga-identifier/
├── api/                    # OpenAPI specifications
├── cmd/api/               # Application entry point
├── db/migrations/         # Database migrations
├── internal/
│   ├── api/
│   │   ├── handler/      # HTTP handlers
│   │   └── router/       # Route definitions
│   ├── domain/
│   │   ├── entity/       # Domain entities
│   │   └── repository/   # Repository interfaces
│   ├── infrastructure/
│   │   ├── external/     # External service clients
│   │   ├── messaging/    # Kafka producer
│   │   └── persistence/  # Database implementations
│   ├── middleware/       # HTTP middleware
│   └── service/          # Business logic
├── pkg/
│   ├── config/           # Configuration
│   ├── database/         # Database connection
│   ├── redis/           # Redis connection
│   └── utils/           # Utility functions
└── docker-compose.yml    # Docker orchestration
```

## 🔒 Security

- Passwords are hashed using bcrypt
- JWT tokens are signed using HMAC-SHA256
- Refresh tokens are hashed before storage
- Rate limiting via Redis for failed login attempts
- Token blacklisting support

## 📝 License

MIT
