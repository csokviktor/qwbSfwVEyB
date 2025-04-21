# Introduction and Goals

The goal of the project is to showcase some knowledge regarding micro-service development.

## Requirements

### Functional Requirements

1. Book Endpoints:
    - List Books: Retrieve a list of all available books.
    - Add Book: Add a new book to the library.
    - Borrow Book: Borrow a book by specifying the book ID and borrower ID.

2. Author Endpoints:
    - Create Author: Add a new author.
    - Get Author: Retrieve details of a specific author.

3. Borrower Endpoints:
    - Create Borrower: Register a new borrower.
    - Get Borrower: Retrieve details of a specific borrower.
    - Borrowed Books: Retrieve the list of books borrowed by a specific borrower.

### Non-Functional Requirements

1. Database:
    - Use any in-memory database (e.g., SQLite, H2) or external database (e.g.,MySQL, PostgreSQL) of your choice.
    - Ensure the application can run locally with the database configured appropriately.

### Testing Requirements

1. Unit Tests:
    - Write unit tests for the main components of your application.
    - Write end-to-end integration tests to verify the complete workflow.
    - Use any testing framework of your choice.
2. Integration Tests:
    - Verify the complete workflow of user creation, adding a book, borrowing a book, and retrieving borrowed books.
    - Ensure that the application works end-to-end.