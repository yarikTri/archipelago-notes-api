import os

import psycopg2
import pytest
import requests
from dotenv import load_dotenv
from http.cookies import SimpleCookie

# Load environment variables
load_dotenv()


def create_note(api_client, title, automerge_url, root_dir):
    """Helper function to create a note"""
    print(f"TEST CREATE NOTE api_client.headers: {api_client.headers}")
    print(f"TEST CREATE NOTE cookies: {api_client.cookies}")
    note_data = {"title": title, "dir_id": root_dir, "automerge_url": automerge_url}
    response = api_client.post(f"{api_client.base_url}/api/notes", json=note_data)
    assert response.status_code == 201, response.text
    return response.json()["id"]


def print_tables(conn):
    with conn.cursor() as cur:
        cur.execute("""
            SELECT table_name
            FROM information_schema.tables
            WHERE table_schema = 'public'
            ORDER BY table_name;
        """)
        tables = cur.fetchall()
        print("\nDatabase tables:")
        for table in tables:
            print(f"- {table[0]}")
        print()


def print_users(conn):
    try:
        with conn.cursor() as cur:
            cur.execute('SELECT * FROM "user"')
            users = cur.fetchall()
            print("\nUsers in database:")
            for user in users:
                print(f"  {user}")
    except Exception as e:
        print(f"Error printing users: {e}")
        # Don't raise the exception, just log it and continue
        return
    finally:
        # Ensure we don't leave any uncommitted transactions
        conn.commit()


def clean_pg_database(conn):
    """Clean up all tables in the database"""
    try:
        with conn.cursor() as cur:
            # First truncate junction tables
            cur.execute("TRUNCATE TABLE tag_to_tag CASCADE")
            cur.execute("TRUNCATE TABLE tag_to_note CASCADE")
            # Then truncate the main tag table
            cur.execute("TRUNCATE TABLE tag CASCADE")
            # Also truncate related tables that might have references
            cur.execute("TRUNCATE TABLE note CASCADE")
            cur.execute("TRUNCATE TABLE dir CASCADE")
            cur.execute('TRUNCATE TABLE "user" CASCADE')
            cur.execute("TRUNCATE TABLE user_root_dir CASCADE")
            cur.execute("TRUNCATE TABLE note_access CASCADE")
            cur.execute("TRUNCATE TABLE summ_to_note CASCADE")
            cur.execute("TRUNCATE TABLE summ CASCADE")
    except Exception as e:
        print(f"Error cleaning database: {e}")
        conn.rollback()
        raise
    finally:
        conn.commit()

    print("PG DATABASE CLEANED")


def clean_sessions(auth_base_url):
    """Clean up all sessions from the Redis database"""
    # Clear all sessions using the admin endpoint
    header = {"X-Admin-Password": os.getenv("ADMIN_PASSWORD")}
    response = requests.post(
        f"{auth_base_url}/auth-service/clear-sessions",
        headers=header,
    )
    assert response.status_code == 200, response.text
    # if response.status_code != 200:
    #    print(f"Warning: Failed to clear sessions: {response.text}")
    print("SESSIONS CLEANED")


@pytest.fixture(scope="session")
def api_base_url():
    """Get the base URL for the API"""
    print(f"API_BASE_URL: {os.getenv('API_BASE_URL')}")
    return os.getenv("API_BASE_URL", "http://localhost:8080")


@pytest.fixture(scope="session")
def auth_base_url():
    """Get the base URL for the auth service"""
    print(f"AUTH_BASE_URL: {os.getenv('AUTH_BASE_URL')}")
    return os.getenv("AUTH_BASE_URL", "http://localhost:8082")


@pytest.fixture(scope="session")
def api_client(api_base_url):
    """Create a session for making API requests"""
    session = requests.Session()
    session.base_url = api_base_url
    print(f"API_CLIENT: {session}")
    return session


@pytest.fixture(scope="session")
def db_connection(auth_base_url):
    """Create a database connection"""
    conn = psycopg2.connect(
        dbname=os.getenv("DB_NAME", "archipelago_notes"),
        user=os.getenv("DB_USER", "postgres"),
        password=os.getenv("DB_PASSWORD", "postgres"),
        host=os.getenv("DB_HOST", "localhost"),
        port=os.getenv("DB_PORT", "5432"),
    )
    # Print all tables in the database
    print_tables(conn)

    # Print all users in the database
    print_users(conn)
    # Clean database before starting tests
    clean_pg_database(conn)
    print_users(conn)
    clean_sessions(auth_base_url)
    return conn


@pytest.fixture(autouse=True)
def cleanup_database(db_connection, auth_base_url):
    """Clean up the database after each test"""
    yield

    clean_sessions(auth_base_url)
    clean_pg_database(db_connection)


@pytest.fixture(scope="function", autouse=True)
def test_user(api_client, auth_base_url):
    """Create a test user and return its ID"""
    print("STARTING TEST USER")
    # Sign up a new user
    signup_data = {
        "email": "test@example.com",
        "name": "TestUser",
        "password": "testpassword123",
    }
    response = api_client.post(
        f"{auth_base_url}/auth-service/registration", json=signup_data
    )
    print(f"TEST USER registration response: {response.text}")
    assert response.status_code == 200, response.text
    user_id = response.json()["user_id"]

    print("TEST USER registration success")

    # Login to get session cookie
    login_data = {"email": "test@example.com", "password": "testpassword123"}
    response = api_client.post(f"{auth_base_url}/auth-service/login", json=login_data)
    assert response.status_code == 200, response.text
    assert response.json()["user_id"] == user_id

    print("TEST USER login success")

    # Set user ID in headers for subsequent requests
    api_client.headers["X-User-Id"] = user_id

    return user_id


@pytest.fixture(scope="function")
def root_dir(api_client, test_user):
    """Create a root directory for the test user and return its ID"""
    print("STARTING ROOT DIR")
    dir_data = {"name": "Root Directory", "parent_id": None}
    response = api_client.post(f"{api_client.base_url}/api/dirs", json=dir_data)
    print(f"ROOT DIR response: {response.text}")
    assert response.status_code == 200, response.text
    dir_id = response.json()["id"]
    print(f"ROOT DIR created with ID: {dir_id}")
    return dir_id


@pytest.fixture(scope="function")
def test_note(api_client, test_user, root_dir):
    """Create a test note and return its ID"""
    # Create a test note
    note_data = {
        "title": "Test Note",
        "dir_id": root_dir,  # Use the root directory ID
        "automerge_url": "test-automerge-url",  # Required by API
    }
    print(f"TEST NOTE api_client.headers: {api_client.headers}")
    print(f"TEST NOTE cookies: {api_client.cookies}")
    response = api_client.post(f"{api_client.base_url}/api/notes", json=note_data)
    assert response.status_code == 201, response.text
    return response.json()["id"]


@pytest.fixture(scope="function")
def test_tag(api_client, test_note):
    """Create a test tag and return its ID"""
    # Create and link a test tag
    tag_data = {"name": "test-tag", "note_id": test_note}
    response = api_client.post(f"{api_client.base_url}/api/tags/create", json=tag_data)
    assert response.status_code == 201
    tag_id = response.json()["tag_id"]
    print(f"TEST TAG response: {response.text}")
    print(f"TEST TAG tag_id: {tag_id}")
    return tag_id


@pytest.fixture(scope="function")
def second_user_client(api_base_url, auth_base_url):
    """Create a client for a second user"""
    # Create a new session
    client = requests.Session()
    client.base_url = api_base_url

    # Register second user
    registration_data = {
        "name": "Second User",
        "email": "second_user@example.com",
        "password": "password123"
    }
    response = client.post(f"{auth_base_url}/auth-service/registration", json=registration_data)
    assert response.status_code == 200
    user_id = response.json()["user_id"]

    # Login second user
    login_data = {
        "email": "second_user@example.com",
        "password": "password123"
    }
    response = client.post(f"{auth_base_url}/auth-service/login", json=login_data)
    assert response.status_code == 200
    assert response.json()["user_id"] == user_id

    # Set auth token cookie
    auth_token = response.cookies["auth_token"]
    client.cookies.set("auth_token", auth_token, domain="localhost.local")

    # Set user ID header
    client.headers["X-User-Id"] = user_id

    return client


@pytest.fixture(scope="function")
def second_user_note(second_user_client, root_dir):
    """Create a note for the second user"""
    note_id = create_note(
        second_user_client, "Second User Note", "second-user-automerge-url", root_dir
    )
    return note_id


@pytest.fixture(scope="function")
def second_user_tag(second_user_client, second_user_note):
    """Create a tag for the second user"""
    tag_data = {"name": "second-user-tag", "note_id": second_user_note}
    response = second_user_client.post(
        f"{second_user_client.base_url}/api/tags/create", json=tag_data
    )
    assert response.status_code == 201
    return response.json()["tag_id"]
