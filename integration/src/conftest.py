import pytest
import httpx
import os

BASE_URL = os.getenv("API_BASE_URL", "http://localhost:8080/v1")

@pytest.fixture(scope="session")
def client():
    with httpx.Client(base_url=BASE_URL) as client:
        yield client

@pytest.fixture
def cleanup(client):
    """Main cleanup fixture that tracks all created resources"""
    resources = {
        'authors': [],
        'borrowers': [],
        'books': []
    }
    yield resources
    # Cleanup in reverse order to respect dependencies
    for book_id in resources['books']:
        try:
            client.delete(f"/books/{book_id}")
        except httpx.HTTPError:
            pass
    for borrower_id in resources['borrowers']:
        try:
            client.delete(f"/borrowers/{borrower_id}")
        except httpx.HTTPError:
            pass
    for author_id in resources['authors']:
        try:
            client.delete(f"/authors/{author_id}")
        except httpx.HTTPError:
            pass