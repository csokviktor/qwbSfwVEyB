import pytest
import uuid

@pytest.fixture
def sample_author(client, cleanup):
    """Fixture to create a sample author for book tests"""
    response = client.post("/authors", json={"name": "Test Author"})
    author_id = response.json()["id"]
    cleanup['authors'].append(author_id)
    return response.json()

@pytest.fixture
def sample_borrower(client, cleanup):
    """Fixture to create a sample borrower for book tests"""
    response = client.post("/borrowers", json={"name": "Test Borrower"})
    borrower_id = response.json()["id"]
    cleanup['borrowers'].append(borrower_id)
    return response.json()

def test_create_book(client, sample_author, cleanup):
    """Test creating a new book"""
    data = {
        "title": "The Hobbit",
        "authorID": sample_author["id"]
    }
    response = client.post("/books", json=data)
    
    assert response.status_code == 201
    book = response.json()
    cleanup['books'].append(book["id"])  # Register for cleanup
    assert book["title"] == data["title"]
    assert book["authorID"] == data["authorID"]
    assert "id" in book

def test_get_books(client, sample_author, cleanup):
    """Test getting all books"""
    book_data = {
        "title": "Sample Book",
        "authorID": sample_author["id"]
    }
    book_response = client.post("/books", json=book_data)
    book_id = book_response.json()["id"]
    cleanup['books'].append(book_id)
    
    response = client.get("/books")
    assert response.status_code == 200
    books = response.json()
    assert isinstance(books, list)
    assert any(b["id"] == book_id for b in books)

def test_borrow_book(client, sample_author, sample_borrower, cleanup):
    """Test borrowing a book"""
    # Create book
    book_data = {
        "title": "Book to Borrow",
        "authorID": sample_author["id"]
    }
    book_response = client.post("/books", json=book_data)
    book_id = book_response.json()["id"]
    cleanup['books'].append(book_id)
    
    # Borrow book
    borrow_data = {"borrowerID": sample_borrower["id"]}
    response = client.post(f"/books/{book_id}/borrow", json=borrow_data)
    assert response.status_code == 201
    
    # Verify book is borrowed
    book_response = client.get("/books")
    books = book_response.json()
    borrowed_book = next(b for b in books if b["id"] == book_id)
    assert borrowed_book["borrowerID"] == sample_borrower["id"]

def test_create_book_invalid_author(client, cleanup):
    """Test creating a book with invalid author ID"""
    fake_author_id = str(uuid.uuid4())
    data = {
        "title": "Invalid Book",
        "authorID": fake_author_id
    }
    response = client.post("/books", json=data)
    assert response.status_code == 500
    error = response.json()
    assert "error" in error

def test_borrow_nonexistent_book(client, sample_borrower):
    """Test borrowing a book that doesn't exist"""
    fake_book_id = str(uuid.uuid4())
    borrow_data = {"borrowerID": sample_borrower["id"]}
    response = client.post(f"/books/{fake_book_id}/borrow", json=borrow_data)
    assert response.status_code == 404
    error = response.json()
    assert "error" in error