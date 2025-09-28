# 📚 Specifications - Multi-user Personal Library API

## 🎯 Project Objective

Develop a RESTful API in Go allowing users to manage their personal book libraries with authentication and data isolation.

## 👥 System Actors

- **Unauthenticated user**: Can register and login
- **Authenticated user**: Can manage their personal library and wishlist

## 📋 Required Features

### 1. **User Management**

- [x] **Registration**: Account creation with email and password
- [x] **Login**: Authentication with JWT token generation
- [ ] **User profile**: View and edit profile
- [ ] **Password change**
- [ ] **Logout**: Token invalidation

### 2. **Book Management (Main Library)**

- [ ] **Add a book** to library
- [ ] **Edit information** of an existing book
- [ ] **Delete a book** from library
- [ ] **Mark as read/unread**
- [ ] **Add a rating** (1-5 stars) and personal comment
- [ ] **Set dates** (start and finish reading)
- [ ] **Reading tracking** (in progress, to read, finished)

### 3. **Wishlist Management**

- [ ] **Add a book** to wishlist
- [ ] **Remove a book** from wishlist
- [ ] **Move a book** from wishlist to main library
- [ ] **Prioritize wishes** with priority level (low, medium, high)
- [ ] **Add notes** on why this book is wanted
- [ ] **Mark as acquired** when added to library

### 4. **Search and Filtering**

- [ ] **List all books** with pagination
- [ ] **List wishlist** with pagination
- [ ] **Text search** by title and author (in library AND wishlist)
- [ ] **Filter by reading status** (to read, in progress, finished)
- [ ] **Filter by rating** (1 to 5 stars)
- [ ] **Filter by priority** in wishlist
- [ ] **Sort by**: date added, title, author, rating, priority
- [ ] **Advanced search** with filter combinations

### 5. **Advanced Features**

- [ ] **Personal statistics**:
  - Total number of books
  - Number of books read/in progress/to read
  - Number of books in wishlist
  - Average rating
  - Most read author
  - Average wishlist priority
- [ ] **Data export**: Export library and wishlist in JSON format
- [ ] **Book import** via ISBN (integration with external API)
- [ ] **Suggestions**: Popular books among other users' wishlists (anonymized)

## 🔐 Security and Authentication

- **JWT tokens** with short expiration (15-30 minutes)
- **Refresh tokens** to maintain session
- **Password hashing** with bcrypt
- **Authentication middleware** on all protected routes
- **Strict validation** of input data
- **Data isolation**: a user can only access their own books and wishes

## 🌐 Detailed API Endpoints

### Authentication

```
POST /api/auth/register     # Account creation
POST /api/auth/login        # Login
POST /api/auth/refresh      # Token renewal
POST /api/auth/logout       # Logout
GET  /api/auth/me           # User profile
PUT  /api/auth/password     # Change password
```

### Book Management (Library)

```
GET    /api/books           # List with filters/pagination
GET    /api/books/{id}      # Book details
POST   /api/books           # Add a book
PUT    /api/books/{id}      # Edit a book
DELETE /api/books/{id}      # Delete a book
GET    /api/books/stats     # Personal statistics
POST   /api/books/import    # Import by ISBN
GET    /api/books/export    # Export library
```

### Wishlist Management

```
GET    /api/wishlist        # List wishes with filters
GET    /api/wishlist/{id}   # Wish details
POST   /api/wishlist        # Add a book to wishlist
PUT    /api/wishlist/{id}   # Edit a wish (priority, notes)
DELETE /api/wishlist/{id}   # Remove a book from wishlist
POST   /api/wishlist/{id}/move-to-books  # Move to library
GET    /api/wishlist/stats  # Wishlist statistics
```

## 💾 Database

- **PostgreSQL** as main database
- **Main tables**: users, books, wishlist_items
- **Relations**: users → books (one-to-many), users → wishes (one-to-many)
- **Wishlist fields**: priority, wish notes, date added, acquisition date
- **Indexes** on frequently searched fields

## 🚀 Technical Constraints

- **Language**: Go (with web framework like Gin or Echo)
- **Database**: PostgreSQL
- **Authentication**: JWT
- **API**: RESTful JSON
- **Separation**: strict data separation by user
- **Error handling**: appropriate with meaningful HTTP codes

## 📦 Deliverables

- API Go source code
- PostgreSQL database configuration scripts
- API endpoints documentation
- Installation and deployment instructions

The wishlist adds an interesting dimension to the project while remaining manageable for learning Go and PostgreSQL.
