import uuid

def test_create_borrower(client, cleanup):
    """Test creating a new borrower"""
    data = {"name": "John Doe"}
    response = client.post("/borrowers", json=data)
    
    assert response.status_code == 201
    borrower = response.json()
    cleanup['borrowers'].append(borrower['id'])  # Register for cleanup
    assert borrower["name"] == data["name"]
    assert "id" in borrower
    assert "books" in borrower
    assert isinstance(borrower["books"], list)

def test_get_borrowers(client, cleanup):
    """Test getting all borrowers"""
    # Create test borrower
    create_response = client.post("/borrowers", json={"name": "Test Borrower"})
    borrower_id = create_response.json()["id"]
    cleanup['borrowers'].append(borrower_id)
    
    response = client.get("/borrowers")
    assert response.status_code == 200
    borrowers = response.json()
    assert isinstance(borrowers, list)
    assert any(b["id"] == borrower_id for b in borrowers)

def test_get_single_borrower(client, cleanup):
    """Test getting a specific borrower by ID"""
    create_data = {"name": "Jane Smith"}
    create_response = client.post("/borrowers", json=create_data)
    borrower_id = create_response.json()["id"]
    cleanup['borrowers'].append(borrower_id)
    
    response = client.get(f"/borrowers/{borrower_id}")
    assert response.status_code == 200
    borrower = response.json()
    assert borrower["id"] == borrower_id
    assert borrower["name"] == create_data["name"]

def test_get_borrower_books(client, cleanup):
    """Test getting books borrowed by a borrower"""
    # Create borrower
    borrower_data = {"name": "Book Reader"}
    borrower_response = client.post("/borrowers", json=borrower_data)
    borrower_id = borrower_response.json()["id"]
    cleanup['borrowers'].append(borrower_id)
    
    # Create book
    author_response = client.post("/authors", json={"name": "Book Author"})
    author_id = author_response.json()["id"]
    cleanup['authors'].append(author_id)
    
    book_data = {
        "title": "Sample Book",
        "authorID": author_id
    }
    book_response = client.post("/books", json=book_data)
    book_id = book_response.json()["id"]
    cleanup['books'].append(book_id)
    
    # Borrow book
    client.post(f"/books/{book_id}/borrow", json={"borrowerID": borrower_id})
    
    # Test borrower books endpoint
    response = client.get(f"/borrowers/{borrower_id}/books")
    assert response.status_code == 200
    books = response.json()
    assert isinstance(books, list)
    assert any(b["id"] == book_id for b in books)

def test_get_nonexistent_borrower(client):
    """Test getting a borrower that doesn't exist"""
    fake_id = str(uuid.uuid4())
    response = client.get(f"/borrowers/{fake_id}")
    assert response.status_code == 404
    error = response.json()
    assert "error" in error
