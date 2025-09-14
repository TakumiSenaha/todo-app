# Todo List Management Feature Implementation

## Overview
Implement the core Todo list management functionality for the Todo application. Build a complete Todo management system including CRUD operations, completion toggle, and sorting features.

## Feature List

### 1. Todo Creation Feature
- Todo title, due date, and priority registration
- Validation (required fields, character limits)
- Real-time form validation

### 2. Todo List Feature
- Display registered todos in a list
- Card-based display
- Responsive design

### 3. Todo Edit Feature
- Edit todo title, due date, and priority
- Inline editing or modal editing
- Save and cancel changes

### 4. Todo Delete Feature
- Delete todos
- Confirmation dialog
- Update list after deletion

### 5. Task Completion Toggle Feature
- Toggle completion with checkboxes
- Gray out completed tasks
- Persist completion state

### 6. Task Sorting Feature
- Sort by due date (ascending/descending)
- Sort by priority (high/medium/low)
- Sort by creation date
- Maintain sort state

## Technical Specifications

### Backend Implementation

#### 1. Todo Controller
```go
// File: backend/internal/interface/controller/todo_controller.go
type TodoController struct {
    todoUseCase usecase.TodoUseCase
}

// Methods
- CreateTodo(w http.ResponseWriter, r *http.Request)
- GetTodos(w http.ResponseWriter, r *http.Request)
- GetTodo(w http.ResponseWriter, r *http.Request)
- UpdateTodo(w http.ResponseWriter, r *http.Request)
- DeleteTodo(w http.ResponseWriter, r *http.Request)
- ToggleTodoComplete(w http.ResponseWriter, r *http.Request)
```

#### 2. Todo Use Case
```go
// File: backend/internal/usecase/todo_interactor.go
type TodoInteractor struct {
    todoRepo repository.TodoRepository
}

// Methods
- CreateTodo(ctx context.Context, userID int, todo *domain.Todo) error
- GetTodos(ctx context.Context, userID int, sortBy string) ([]*domain.Todo, error)
- GetTodo(ctx context.Context, userID int, todoID int) (*domain.Todo, error)
- UpdateTodo(ctx context.Context, userID int, todo *domain.Todo) error
- DeleteTodo(ctx context.Context, userID int, todoID int) error
- ToggleTodoComplete(ctx context.Context, userID int, todoID int) error
```

#### 3. Todo Repository (Using SQLC Auto-generated Code)
```go
// File: backend/internal/interface/repository/todo_repository.go
// Using SQLC auto-generated Queries struct
type TodoRepository struct {
    queries *Queries // SQLC auto-generated
}

// Implementation methods (wrapping SQLC auto-generated methods)
func (tr *TodoRepository) CreateTodo(ctx context.Context, userID int, todo *domain.Todo) error {
    // Using SQLC auto-generated CreateTodo method
    params := CreateTodoParams{
        UserID:      int32(userID),
        Title:       todo.Title,
        DueDate:     sql.NullTime{Time: *todo.DueDate, Valid: todo.DueDate != nil},
        Priority:    int32(todo.Priority),
        IsCompleted: todo.IsCompleted,
    }
    sqlcTodo, err := tr.queries.CreateTodo(ctx, params)
    if err != nil {
        return err
    }
    // Map to domain object
    todo.ID = int(sqlcTodo.ID)
    todo.CreatedAt = sqlcTodo.CreatedAt.Time
    todo.UpdatedAt = sqlcTodo.UpdatedAt.Time
    return nil
}

// Other methods also use SQLC auto-generated code similarly
```

#### 4. API Endpoints
```go
// Routing configuration
mux.HandleFunc("/api/v1/todos", todoController.CreateTodo).Methods("POST")
mux.HandleFunc("/api/v1/todos", todoController.GetTodos).Methods("GET")
mux.HandleFunc("/api/v1/todos/{id}", todoController.GetTodo).Methods("GET")
mux.HandleFunc("/api/v1/todos/{id}", todoController.UpdateTodo).Methods("PUT")
mux.HandleFunc("/api/v1/todos/{id}", todoController.DeleteTodo).Methods("DELETE")
mux.HandleFunc("/api/v1/todos/{id}/toggle", todoController.ToggleTodoComplete).Methods("PATCH")
```

### Frontend Implementation

#### 1. Todo List Screen
```typescript
// File: frontend/src/app/dashboard/page.tsx
// Features
- Todo list display
- New todo creation form (at the top)
- Sorting functionality
- Completion toggle
- Edit and delete buttons
```

#### 2. Todo Component
```typescript
// File: frontend/src/components/TodoCard.tsx
// Features
- Todo information display
- Completion toggle checkbox
- Edit and delete buttons
- Due date display (overdue highlighting)
- Priority display
```

#### 3. Todo Creation Form
```typescript
// File: frontend/src/components/TodoForm.tsx
// Features
- Todo title input
- Due date selection
- Priority selection
- Validation
- Submit processing
```

#### 4. Todo Edit Modal
```typescript
// File: frontend/src/components/TodoEditModal.tsx
// Features
- Edit existing todo information
- Form validation
- Save and cancel processing
```

## Data Structure

### 1. Todo Domain
```go
// File: backend/internal/domain/todo.go
type Todo struct {
    ID          int       `json:"id"`
    UserID      int       `json:"user_id"`
    Title       string    `json:"title"`
    DueDate     *time.Time `json:"due_date"`
    Priority    int       `json:"priority"` // 0:Low, 1:Medium, 2:High
    IsCompleted bool      `json:"is_completed"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### 2. Request/Response
```typescript
// Todo creation request
interface CreateTodoRequest {
  title: string;
  due_date?: string; // ISO 8601 format
  priority: number; // 0:Low, 1:Medium, 2:High
}

// Todo update request
interface UpdateTodoRequest {
  title?: string;
  due_date?: string;
  priority?: number;
  is_completed?: boolean;
}

// Todo response
interface TodoResponse {
  id: number;
  user_id: number;
  title: string;
  due_date?: string;
  priority: number;
  is_completed: boolean;
  created_at: string;
  updated_at: string;
}
```

## Screen Design

### 1. Todo List Screen Layout
```typescript
// Layout structure
┌─────────────────────────────────────────────────────────┐
│ Todo App                           [Username] [☰]      │
├─────────────────────────────────────────────────────────┤
│                                                         │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ Add New Todo                                        │ │
│ │                                                     │ │
│ │ Title: [________________]                           │ │
│ │ Due Date: [2024-01-15]                             │ │
│ │ Priority: [Low ▼]                                  │ │
│ │ [Add] [Cancel]                                     │ │
│ └─────────────────────────────────────────────────────┘ │
│                                                         │
│ Todo List                                               │
│                                                         │
│ [Sort: Due Date ▼] [Priority ▼] [Created ▼]            │
│                                                         │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ ☐ Task 1                    [Edit] [Delete]        │ │
│ │   Due: 2024-01-15 Priority: High                   │ │
│ └─────────────────────────────────────────────────────┘ │
│                                                         │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ ☑ Task 2 (Completed)        [Edit] [Delete]        │ │
│ │   Due: 2024-01-10 Priority: Medium                 │ │
│ └─────────────────────────────────────────────────────┘ │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### 2. Todo Creation Form (Fixed at Top)
```typescript
// Form structure - Always visible at the top
┌─────────────────────────────────────────────────────────┐
│ Add New Todo                                            │
│                                                         │
│ Title: [________________]                               │
│                                                         │
│ Due Date: [2024-01-15]                                 │
│                                                         │
│ Priority: [Low ▼]                                      │
│                                                         │
│ [Add] [Cancel]                                         │
└─────────────────────────────────────────────────────────┘
```

## Validation Specifications

### 1. Backend Validation
```go
// Todo creation validation
- title: Required, 1-100 characters
- due_date: Date format, future date
- priority: Integer 0-2
- user_id: Authenticated user ID
```

### 2. Frontend Validation
```typescript
// Real-time validation
- title: Required, 1-100 characters
- due_date: Valid date format
- priority: Integer 0-2
- Final check on form submission
```

## Sorting Feature Specifications

### 1. Sort Options
```typescript
// Sort types
- Due date (ascending/descending)
- Priority (high/medium/low)
- Creation date (newest/oldest)
- Completion status (incomplete/completed)

// Default sort
- Prioritize incomplete tasks
- Due date (ascending)
- Creation date (newest first)
```

### 2. Sort Implementation (SQLC Auto-generated)
```sql
-- File: backend/internal/interface/repository/todo.sql
-- name: ListTodosWithSort :many
SELECT * FROM todos
WHERE user_id = $1
ORDER BY 
  CASE WHEN $2 = 'due_date_asc' THEN due_date END ASC,
  CASE WHEN $2 = 'due_date_desc' THEN due_date END DESC,
  CASE WHEN $2 = 'priority_desc' THEN priority END DESC,
  CASE WHEN $2 = 'created_desc' THEN created_at END DESC,
  is_completed ASC,
  created_at DESC;

-- name: ToggleTodoComplete :one
UPDATE todos 
SET is_completed = NOT is_completed, updated_at = NOW()
WHERE id = $1 AND user_id = $2
RETURNING *;
```

```go
// Usage example after SQLC auto-generation
func (tr *TodoRepository) GetTodos(ctx context.Context, userID int, sortBy string) ([]*domain.Todo, error) {
    // Using SQLC auto-generated ListTodosWithSort method
    sqlcTodos, err := tr.queries.ListTodosWithSort(ctx, ListTodosWithSortParams{
        UserID: int32(userID),
        SortBy: sortBy,
    })
    if err != nil {
        return nil, err
    }
    
    // Map to domain objects
    todos := make([]*domain.Todo, len(sqlcTodos))
    for i, sqlcTodo := range sqlcTodos {
        todos[i] = &domain.Todo{
            ID:          int(sqlcTodo.ID),
            UserID:      int(sqlcTodo.UserID),
            Title:       sqlcTodo.Title,
            DueDate:     &sqlcTodo.DueDate.Time,
            Priority:    int(sqlcTodo.Priority),
            IsCompleted: sqlcTodo.IsCompleted,
            CreatedAt:   sqlcTodo.CreatedAt.Time,
            UpdatedAt:   sqlcTodo.UpdatedAt.Time,
        }
    }
    return todos, nil
}
```

## Error Handling

### 1. Backend Errors
```go
// Error response
type ErrorResponse struct {
    Status  string `json:"status"`
    Message string `json:"message"`
    Code    int    `json:"code"`
}

// Error codes
- 400: Validation error
- 401: Authentication error
- 403: Authorization error
- 404: Todo not found
- 500: Server error
```

### 2. Frontend Errors
```typescript
// Error display
- Form validation errors
- API communication errors
- Network errors
- User-friendly messages
```

## Implementation Order

### 1. Backend Implementation
1. Todo repository implementation (using SQLC auto-generated code)
2. Todo use case implementation
3. Todo controller implementation
4. API endpoint configuration
5. Validation addition

### 2. Frontend Implementation
1. Basic layout for todo list screen
2. Todo creation form (fixed at top)
3. Todo card component
4. Completion toggle functionality
5. Edit and delete functionality
6. Sorting functionality
7. Error handling

## Test Specifications

### 1. Backend Tests
```go
// Test cases
- Todo creation test
- Todo retrieval test
- Todo update test
- Todo deletion test
- Completion toggle test
- Sorting functionality test
- Validation test
- Authentication test
```

### 2. Frontend Tests
```typescript
// Test cases
- Form validation test
- API communication test
- Component rendering test
- User interaction test
```

## Completion Criteria

### Functional Requirements
- [ ] Can add todos
- [ ] Can display todo list
- [ ] Can edit todos
- [ ] Can delete todos
- [ ] Can toggle completion status
- [ ] Completed tasks are grayed out
- [ ] Can sort by due date
- [ ] Can sort by priority

### Technical Requirements
- [ ] Backend API works properly
- [ ] Frontend screen displays correctly
- [ ] Validation works appropriately
- [ ] Error handling is implemented
- [ ] Responsive design is applied
- [ ] Authentication works properly

### Quality Requirements
- [ ] Code follows clean architecture
- [ ] Error messages are user-friendly
- [ ] Performance is good
- [ ] Security is properly implemented

## SQLC Auto-generated Code Utilization

### 1. Auto-generated Files
```bash
# Files generated after SQLC execution
backend/internal/interface/repository/
├── todo.sql          # SQL query definitions
├── db.go             # Auto-generated Queries struct
└── models.go         # Auto-generated model structs
```

### 2. SQLC Benefits
- **Type Safety**: Detect SQL errors at compile time
- **Performance**: Optimized SQL queries
- **Maintainability**: Separation of SQL and Go code
- **Auto-generation**: No manual SQL mapping required

### 3. Implementation Notes
- Regeneration required with `go generate` after SQLC generation
- Manual implementation for domain object mapping
- Proper error handling implementation
- Add transaction processing as needed

## Important Notes

1. **Security**: Users can only operate on their own todos
2. **Performance**: Comfortable operation even with large numbers of todos
3. **UX**: Intuitive and user-friendly interface
4. **Responsive**: Mobile and tablet support
5. **Accessibility**: Keyboard operation and screen reader support
6. **SQLC Utilization**: Maximize development efficiency using auto-generated code

## Layout Design Notes

### Fixed Todo Creation Form
- **Position**: Always visible at the top of the screen
- **Purpose**: Easy access to create new todos without scrolling
- **Behavior**: Form remains accessible even when todo list is long
- **UX**: Users can always add new todos without losing their place in the list

### Scrollable Todo List
- **Behavior**: Todo list below the form is scrollable
- **Benefit**: Can handle large numbers of todos without UI issues
- **Design**: Form stays fixed while list scrolls independently

