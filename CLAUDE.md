# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

D-Day management backend API built with Go and Fiber framework following a layered MVC architecture pattern. The project provides REST and API endpoints for managing D-Day events with MySQL/MariaDB database support.

## Development Commands

### Dependencies and Setup

```bash
# Install dependencies
go mod tidy

# Set up database (MariaDB/MySQL required)
mysql -u root -p < schema.sql

# Copy environment configuration
cp .env.example .env
# Edit .env with your database credentials
```

### Running the Applicationcd

```bash
# Run the server
go run main.go

# The server starts on port 8080 (configurable via PORT env var)
```

### Testing

```bash
# Test database connectivity
mysql -u root -p -e "USE dday; SELECT COUNT(*) FROM ddays;"

# Test API endpoints
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/ddays
```

## Architecture

### Directory Structure

```
dday/back/
├── main.go                    # Application entry point
├── go.mod                     # Go module dependencies
├── models/                    # Data layer
│   ├── db.go                 # Database connection and utilities
│   ├── dday.go               # DDay model and manager
│   └── dday/                 # DDay constants and helpers
│       └── dday.go
├── controllers/               # Request handling layer
│   ├── controllers.go        # Base controller with common functionality
│   ├── api/                  # Business logic controllers
│   │   └── dday.go           # DDay API controller with validation
│   └── rest/                 # Simple CRUD controllers
│       └── dday.go           # DDay REST controller
├── router/                    # Routing setup
│   └── router.go             # Route definitions and middleware
├── global/                    # Cross-cutting concerns
│   └── config/               # Configuration management
│       └── config.go
├── services/                  # Background services (future)
└── schema.sql                # Database schema
```

### Core Architecture Components

#### **1. Layered MVC Pattern**

- **Model Layer** (`models/`): Database entities, business logic, and data access
- **Controller Layer** (`controllers/`): Request handling and response formatting
- **Router Layer** (`router/`): HTTP routing and middleware configuration
- **Global Layer** (`global/`): Configuration and utilities

#### **2. Database Layer (`models/`)**

- **Connection Management**: Custom `Connection` struct wrapping `sql.DB`
- **Manager Pattern**: Each model has a dedicated manager (e.g., `DdayManager`)
- **Flexible Query Building**: Support for dynamic WHERE clauses, ordering, pagination
- **Transaction Support**: Built-in transaction handling with `Begin()`, `Commit()`, `Rollback()`
- **Dual Model Structure**:
  - `models/dday.go`: Main model and database operations
  - `models/dday/dday.go`: Constants, enums, and helper functions

#### **3. Controller Architecture**

- **Base Controller** (`controllers/Controller`): Common functionality shared across all controllers

  - Database connection management
  - Request parameter parsing (`Get()`, `Query()`, `Params()`)
  - Response handling (`Success()`, `Error()`, `Created()`)
  - Pagination utilities (`GetPagination()`)
  - Search and filtering helpers

- **API Controllers** (`controllers/api/`): Business logic with validation

  - Input validation and sanitization
  - Business rule enforcement
  - Error handling with appropriate HTTP status codes
  - Complex queries with filtering, searching, pagination

- **REST Controllers** (`controllers/rest/`): Simple CRUD operations
  - Direct model operations
  - Minimal business logic
  - Standardized REST endpoints

#### **4. Database Operations**

- **Query Building**: Dynamic query construction using interface slices
  - `Where`: Column-based conditions
  - `Custom`: Raw SQL fragments
  - `Paging`: Offset/limit pagination
  - `Ordering`: Sort specifications
- **Connection Pooling**: Configurable connection limits and lifetime
- **Transaction Isolation**: Support for transactional operations

### Key Dependencies

#### **Core Framework**

- **GoFiber v2**: High-performance HTTP framework
- **MySQL Driver**: `go-sql-driver/mysql` for database connectivity
- **UUID**: `google/uuid` for unique identifier generation

#### **Middleware**

- **CORS**: Cross-origin resource sharing support
- **Logger**: Request/response logging
- **Recover**: Panic recovery middleware

### API Endpoints

#### **API Routes** (`/api/v1`)

Business logic endpoints with validation:

- `GET /api/v1/ddays` - Get D-Days with filtering, searching, pagination
- `POST /api/v1/ddays` - Create D-Day with validation
- `GET /api/v1/ddays/:id` - Get specific D-Day
- `PUT /api/v1/ddays/:id` - Update D-Day with validation
- `DELETE /api/v1/ddays/:id` - Delete D-Day

#### **REST Routes** (`/rest`)

Simple CRUD operations:

- `GET /rest/ddays` - List D-Days with basic pagination
- `POST /rest/ddays` - Create D-Day (minimal validation)
- `GET /rest/ddays/:id` - Get D-Day by ID
- `PUT /rest/ddays/:id` - Update D-Day
- `DELETE /rest/ddays/:id` - Delete D-Day

### Environment Configuration

All configuration via environment variables with sensible defaults:

**Database Settings:**

- `DB_HOST` (default: localhost)
- `DB_PORT` (default: 3306)
- `DB_USER` (default: root)
- `DB_PASSWORD` (default: empty)
- `DB_NAME` (default: dday)
- `DB_MAX_OPEN_CONNS` (default: 25)
- `DB_MAX_IDLE_CONNS` (default: 25)
- `DB_CONN_MAX_LIFETIME` (default: 300)

**Server Settings:**

- `PORT` (default: 8080)
- `ENV` (default: development)

### Data Model

**DDay Struct:**

```go
type DDay struct {
    ID          string    `json:"id"`           // UUID primary key
    Title       string    `json:"title"`        // Event title
    TargetDate  string    `json:"target_date"`  // Target date (YYYY-MM-DD)
    Category    string    `json:"category"`     // Event category
    Memo        string    `json:"memo"`         // Optional memo
    IsImportant bool      `json:"is_important"` // Importance flag
    CreatedAt   time.Time `json:"created_at"`   // Creation timestamp
    UpdatedAt   time.Time `json:"updated_at"`   // Last update timestamp
}
```

**Categories:**

- 개인 (Personal) - Default
- 학업 (Study)
- 업무 (Work)
- 기타 (Other)

### Database Schema

- **Primary Table**: `ddays_tb` with optimized indexes
- **Column Naming**: All columns prefixed with `d_` (e.g., `d_id`, `d_title`, `d_target_date`)
- **Table Naming**: All tables suffixed with `_tb` (e.g., `ddays_tb`, `categories_tb`)
- **Indexes**: d_target_date, d_category, d_is_important, d_created_at
- **Character Set**: UTF8MB4 for full Unicode support
- **Auto-timestamps**: d_created_at, d_updated_at with automatic management

### Naming Conventions

#### **Database Naming Standards**

- **Tables**: All table names end with `_tb` suffix (e.g., `ddays_tb`, `users_tb`, `categories_tb`)
- **Columns**: All columns are prefixed with abbreviated table name (e.g., `d_id`, `d_title` for ddays_tb)
- **Indexes**: Named with `idx_` prefix followed by column name (e.g., `idx_d_target_date`)

#### **Column Prefix Rules**

- Single letter prefix for short table names (e.g., `d_` for ddays_tb)
- Abbreviated prefix for longer table names to avoid conflicts (e.g., `cat_` for categories_tb)
- If conflicts arise, use meaningful abbreviations (e.g., `usr_` for users*tb, `usr*` for user_settings_tb)

### Development Patterns

#### **Adding New Models**

1. Create `models/newmodel.go` with struct and manager
2. Create `models/newmodel/` directory for constants/enums
3. Add database operations to manager
4. Create API and REST controllers
5. Add routes to router

#### **Query Building Pattern**

```go
args := []interface{}{
    models.NewWhere("d_category", "개인", "="),
    models.NewCustom("d_title LIKE ?", "%search%"),
    models.NewOrdering("d_target_date ASC"),
    models.NewPaging(1, 10),
}
results, err := manager.GetAll(args...)
```

#### **Controller Pattern**

```go
func (ctrl *Controller) HandlerMethod(c *fiber.Ctx) error {
    ctrl.Controller = controllers.NewController(c)

    // Parse parameters
    id := ctrl.Params("id")
    page, pageSize := ctrl.GetPagination()

    // Business logic
    result, err := ctrl.manager.GetByID(id)
    if err != nil {
        return ctrl.NotFound("Resource not found")
    }

    return ctrl.Success(result)
}
```
