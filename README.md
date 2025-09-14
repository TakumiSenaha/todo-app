# Todo App - Full Stack Authentication System

A complete full-stack todo application with robust JWT-based authentication, built with Go backend and Next.js frontend.

## 🏗️ Architecture

### Backend (Go)
- **Clean Architecture** with dependency injection
- **JWT Authentication** with HttpOnly cookies
- **PostgreSQL** database with migrations
- **Middleware** for CORS, logging, and authentication
- **Validation** using go-playground/validator
- **Password hashing** with bcrypt

### Frontend (Next.js)
- **BFF (Backend for Frontend)** pattern with API routes
- **React Context** for global authentication state
- **Route protection** with Next.js middleware
- **TypeScript** for type safety
- **Tailwind CSS** for styling

## 🚀 Quick Start (Recommended)

### Prerequisites
- Docker & Docker Compose
- Git

### One-Command Setup

1. **Clone the repository:**
```bash
git clone git@github.com:TakumiSenaha/todo-app.git
cd todo-app
```

2. **Start the entire application:**
```bash
docker-compose up -d
```

3. **Access the application:**
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Database: PostgreSQL on localhost:5432

That's it! The application will automatically:
- Build and start the Go backend
- Build and start the Next.js frontend
- Set up PostgreSQL database with migrations
- Configure all necessary environment variables

### Stop the application:
```bash
docker-compose down
```

### View logs:
```bash
# All services
docker-compose logs

# Backend only
docker-compose logs backend

# Frontend only
docker-compose logs frontend
```

## 🛠️ Manual Setup (Development)

### Prerequisites
- Go 1.21+
- Node.js 18+
- PostgreSQL
- pnpm (recommended) or npm

### Database Setup

1. Create a PostgreSQL database:
```sql
CREATE DATABASE todo_db;
CREATE USER user WITH ENCRYPTED PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE todo_db TO user;
```

2. Run migrations:
```bash
cd backend
# Install golang-migrate if not already installed
# https://github.com/golang-migrate/migrate
migrate -path migrations -database "postgresql://user:password@localhost:5432/todo_db?sslmode=disable" up
```

### Backend Setup

1. Navigate to backend directory:
```bash
cd backend
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set environment variables (optional):
```bash
export DB_SOURCE="postgresql://user:password@localhost:5432/todo_db?sslmode=disable"
export JWT_SECRET="your-super-secret-jwt-key-change-in-production"
export PORT="8080"
```

4. Run the backend:
```bash
go run cmd/api/main.go
```

The backend will start on `http://localhost:8080`

### Frontend Setup

1. Navigate to frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
pnpm install
# or
npm install
```

3. Set environment variables (create `.env.local`):
```bash
BACKEND_URL=http://localhost:8080
```

4. Run the frontend:
```bash
pnpm dev
# or
npm run dev
```

The frontend will start on `http://localhost:3000`

## 🐳 Docker Configuration

### Services

The application consists of three Docker services:

1. **PostgreSQL Database** (`db`)
   - Port: 5432
   - Database: `todo_db`
   - User: `user`
   - Password: `password`

2. **Go Backend** (`backend`)
   - Port: 8080
   - Auto-rebuilds on code changes
   - Connects to PostgreSQL database

3. **Next.js Frontend** (`frontend`)
   - Port: 3000
   - Auto-rebuilds on code changes
   - Connects to backend via API routes

### Environment Variables

All environment variables are configured in `docker-compose.yml`:

```yaml
# Backend
DB_SOURCE: "postgresql://user:password@db:5432/todo_db?sslmode=disable"
JWT_SECRET: "your-super-secure-jwt-secret-key-here-change-this-in-production"
PORT: "8080"

# Frontend
BACKEND_URL: "http://backend:8080"
```

### Development with Docker

For development with hot reloading:

```bash
# Start in development mode
docker-compose up

# Rebuild specific service
docker-compose up --build backend

# Run in background
docker-compose up -d

# Stop all services
docker-compose down

# Remove volumes (clears database)
docker-compose down -v
```

## 🔐 Authentication Flow

### 1. User Registration
- **Frontend**: User fills registration form with username, email, password
- **Validation**: Client-side validation for password strength (8+ chars, letters + numbers)
- **Backend**: Server validates data, hashes password with bcrypt, stores in database
- **Response**: User data returned (excluding password hash)

### 2. User Login
- **Frontend**: User submits login credentials
- **Backend**: Validates credentials, generates JWT token with 24-hour expiry
- **Cookie**: JWT stored in HttpOnly, Secure, SameSite=Strict cookie
- **Response**: User data returned, client updates global state

### 3. Route Protection
- **Middleware**: Next.js middleware checks for auth_token cookie on protected routes
- **Redirect**: Unauthenticated users redirected to login page
- **Backend**: Go middleware validates JWT token and extracts user ID

### 4. User Session
- **Context**: React Context maintains authentication state across components
- **Persistence**: Auth state persists across browser refreshes via /api/users/me endpoint
- **Logout**: Clears HttpOnly cookie and resets client state

## 📡 API Endpoints

### Authentication
- `POST /api/v1/register` - User registration
- `POST /api/v1/login` - User login
- `POST /api/v1/logout` - User logout
- `GET /api/v1/me` - Get current user (protected)

### Health
- `GET /health` - Health check

## 🛡️ Security Features

### Password Security
- **Minimum 8 characters** with letters and numbers required
- **bcrypt hashing** with default cost (currently 10)
- **No plain text** passwords stored anywhere

### JWT Security
- **24-hour expiry** time
- **HttpOnly cookies** prevent XSS attacks
- **Secure flag** for HTTPS in production
- **SameSite=Strict** prevents CSRF attacks

### CORS Protection
- **Origin validation** for specific allowed domains
- **Credentials support** for cookie authentication
- **Pre-flight handling** for complex requests

### Input Validation
- **Server-side validation** using go-playground/validator
- **Client-side validation** for immediate feedback
- **SQL injection protection** with parameterized queries

## 🔧 Configuration

### Environment Variables

#### Backend
- `DB_SOURCE` - PostgreSQL connection string
- `JWT_SECRET` - Secret key for JWT signing
- `PORT` - Server port (default: 8080)

#### Frontend
- `BACKEND_URL` - Backend API URL (default: http://localhost:8080)

## 📝 Testing the System

### Using Docker (Recommended)

1. **Start the application:**
```bash
docker-compose up
```

2. **Visit** `http://localhost:3000`

3. **Create an account** via the registration form:
   - Username: `testuser`
   - Email: `test@example.com`
   - Password: `password123`

4. **Sign in** with your credentials

5. **Access the dashboard** (protected route)

6. **Test profile editing** via the hamburger menu

7. **Sign out** and verify redirect to login

### Using Manual Setup

1. **Start PostgreSQL** database
2. **Start backend:** `cd backend && go run cmd/api/main.go`
3. **Start frontend:** `cd frontend && pnpm dev`
4. **Follow steps 2-7 above**

### Manual API Testing

```bash
# Register a new user
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}' \
  -c cookies.txt

# Access protected endpoint
curl -X GET http://localhost:8080/api/v1/me \
  -b cookies.txt

# Logout
curl -X POST http://localhost:8080/api/v1/logout \
  -b cookies.txt
```

## 🏢 Project Structure

```
todo-app/
├── backend/
│   ├── cmd/api/main.go              # Application entrypoint
│   ├── internal/
│   │   ├── domain/                  # Domain entities
│   │   │   ├── user.go
│   │   │   └── todo.go
│   │   ├── usecase/                 # Business logic
│   │   │   ├── user_repository.go
│   │   │   └── user_interactor.go
│   │   ├── interface/
│   │   │   ├── controller/          # HTTP handlers
│   │   │   │   └── user_controller.go
│   │   │   └── middleware/          # HTTP middleware
│   │   │       └── auth_middleware.go
│   │   └── infrastructure/
│   │       └── persistence/         # Database implementation
│   │           └── user_persistence.go
│   ├── migrations/                  # Database migrations
│   ├── go.mod
│   └── go.sum
└── frontend/
    ├── src/
    │   ├── app/
    │   │   ├── api/                 # BFF API routes
    │   │   │   ├── auth/
    │   │   │   └── users/
    │   │   ├── login/page.tsx       # Login page
    │   │   ├── register/page.tsx    # Registration page
    │   │   ├── dashboard/page.tsx   # Protected dashboard
    │   │   ├── layout.tsx           # Root layout with AuthProvider
    │   │   └── page.tsx             # Landing page
    │   └── contexts/
    │       └── AuthContext.tsx      # Global auth state
    ├── middleware.ts                # Route protection
    ├── package.json
    └── README.md
```

## 🎯 Next Steps

1. **Todo Management**: Implement CRUD operations for todos
2. **User Profile**: Add user profile editing functionality  
3. **Email Verification**: Add email verification for new accounts
4. **Password Reset**: Implement password reset functionality
5. **Rate Limiting**: Add rate limiting for auth endpoints
6. **Session Management**: Add active session management
7. **Audit Logging**: Add comprehensive audit logging
8. **Unit Tests**: Add comprehensive test coverage

## 🐛 Troubleshooting

### Docker Issues

**Container won't start:**
```bash
# Check logs
docker-compose logs

# Rebuild containers
docker-compose up --build

# Remove all containers and volumes
docker-compose down -v
docker-compose up
```

**Port conflicts:**
```bash
# Check what's using the ports
lsof -i :3000
lsof -i :8080
lsof -i :5432

# Stop conflicting services or change ports in docker-compose.yml
```

**Database connection issues:**
```bash
# Check database logs
docker-compose logs db

# Connect to database directly
docker-compose exec db psql -U user -d todo_db
```

### Backend Issues
- **Database connection**: Verify PostgreSQL is running and credentials are correct
- **Port conflicts**: Ensure port 8080 is available
- **JWT errors**: Check JWT_SECRET environment variable
- **Migration errors**: Check database logs and ensure migrations are applied

### Frontend Issues
- **API connection**: Verify backend is running on correct port
- **Cookie issues**: Check browser developer tools for cookie presence
- **CORS errors**: Verify CORS configuration in backend
- **Build errors**: Check Node.js version and dependencies

### Common Issues
- **"Method not allowed"**: Check HTTP method and endpoint spelling
- **"Unauthorized"**: Verify token is present and valid
- **"Validation failed"**: Check request payload format and requirements
- **"Internal Server Error"**: Check backend logs for detailed error messages

## 📋 Requirements Checklist

### User Management ✅
- [x] User registration with unique username
- [x] Password validation (8+ chars, letters + numbers)
- [x] Password encryption (bcrypt)
- [x] User login with credentials
- [x] Login error handling
- [x] User information editing capability  
- [x] Logout functionality

### Technical Requirements ✅
- [x] Clean Architecture implementation
- [x] JWT-based authentication
- [x] HttpOnly cookie security
- [x] CORS configuration
- [x] Input validation
- [x] Route protection
- [x] Error handling
- [x] Logging middleware

### Frontend Features ✅
- [x] BFF API integration
- [x] Global state management
- [x] Route protection middleware
- [x] Login/Register forms
- [x] Dashboard interface
- [x] Responsive design

This authentication system provides a solid foundation for a production-ready application with security best practices and modern architectural patterns.

## 🎯 Quick Reference

### Essential Commands

```bash
# Start everything
docker-compose up

# Stop everything
docker-compose down

# View logs
docker-compose logs -f

# Rebuild and start
docker-compose up --build

# Reset database
docker-compose down -v && docker-compose up
```

### Default Credentials

- **Database**: `user` / `password`
- **Database Name**: `todo_db`
- **JWT Secret**: `your-super-secure-jwt-secret-key-here-change-this-in-production`

### URLs

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Health Check**: http://localhost:8080/health

### Test User

Create a test account with:
- Username: `testuser`
- Email: `test@example.com`
- Password: `password123`

---

**Ready to start? Run `docker-compose up` and visit http://localhost:3000!**
