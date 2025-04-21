import uuid

def test_create_author(client, cleanup):
    """Test creating a new author"""
    data = {"name": "J.R.R. Tolkien"}
    response = client.post("/authors", json=data)
    
    assert response.status_code == 201
    author = response.json()
    cleanup['authors'].append(author['id'])  # Register for cleanup
    assert author["name"] == data["name"]
    assert "id" in author
    assert "books" in author
    assert isinstance(author["books"], list)

def test_get_authors(client, cleanup):
    """Test getting all authors"""
    # Create test author
    create_response = client.post("/authors", json={"name": "Test Author"})
    author_id = create_response.json()["id"]
    cleanup['authors'].append(author_id)
    
    response = client.get("/authors")
    assert response.status_code == 200
    authors = response.json()
    assert isinstance(authors, list)
    assert any(a["id"] == author_id for a in authors)

def test_get_single_author(client, cleanup):
    """Test getting a specific author by ID"""
    create_data = {"name": "George R.R. Martin"}
    create_response = client.post("/authors", json=create_data)
    author_id = create_response.json()["id"]
    cleanup['authors'].append(author_id)
    
    response = client.get(f"/authors/{author_id}")
    assert response.status_code == 200
    author = response.json()
    assert author["id"] == author_id
    assert author["name"] == create_data["name"]

def test_get_nonexistent_author(client):
    """Test getting an author that doesn't exist"""
    fake_id = str(uuid.uuid4())
    response = client.get(f"/authors/{fake_id}")
    assert response.status_code == 404
    error = response.json()
    assert "error" in error

def test_create_author_invalid_data(client):
    """Test creating an author with invalid data"""
    response = client.post("/authors", json={})
    assert response.status_code == 400
    error = response.json()
    assert "error" in error
