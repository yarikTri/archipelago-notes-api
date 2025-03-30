import uuid
import pytest
from http.cookies import SimpleCookie
from conftest import create_note


def test_create_and_link_tag(api_client, test_note):
    """Test creating a new tag and linking it to a note"""
    # Test successful creation
    tag_data = {"name": "test-tag", "note_id": test_note}
    response = api_client.post(
        f"{api_client.base_url}/api/tags/create", 
        json=tag_data
    )
    assert response.status_code == 201, response.text
    assert "tag_id" in response.json()
    assert "name" in response.json()
    assert response.json()["name"] == "test-tag"
    tag_id = response.json()["tag_id"]
    
    # Verify tag is linked to the note
    response = api_client.get(f"{api_client.base_url}/api/tags/note/{test_note}")
    assert response.status_code == 200, response.text
    tags = response.json()
    assert any(tag["tag_id"] == tag_id for tag in tags), "Tag should be linked to the note"

    # Test creating a tag with same name for the same note (should fail)
    response = api_client.post(
        f"{api_client.base_url}/api/tags/create", 
        json=tag_data
    )
    assert response.status_code == 409, "Creating a tag with the same name for the same note should fail"
    
    # Test invalid request - empty name
    invalid_data = {"name": "", "note_id": test_note}
    response = api_client.post(
        f"{api_client.base_url}/api/tags/create", 
        json=invalid_data
    )
    assert response.status_code == 400, response.text


def test_create_and_link_tag_invalid_request(api_client, test_user):
    """Test creating a tag with invalid request data"""
    # Test with invalid note ID
    tag_data = {
        "name": "test-tag",
        "note_id": str(uuid.uuid4()),  # Non-existent note ID
    }
    response = api_client.post(f"{api_client.base_url}/api/tags/create", json=tag_data)
    assert response.status_code == 404, response.text


def test_unlink_tag_from_note(api_client, test_tag, test_note):
    """Test unlinking a tag from a note"""
    # Test successful unlinking
    unlink_data = {"tag_id": test_tag, "note_id": test_note}
    response = api_client.post(
        f"{api_client.base_url}/api/tags/unlink", json=unlink_data
    )
    assert response.status_code == 200

    # Verify note exists
    response = api_client.get(f"{api_client.base_url}/api/notes/{test_note}")
    assert response.status_code == 200, response.text

    # Verify tag is unlinked
    print(test_note)
    response = api_client.get(f"{api_client.base_url}/api/tags/note/{test_note}")
    assert response.status_code == 200, response.text
    print(response.text)
    tags = response.json()
    print(tags)
    # Handle both empty array and None responses
    assert not any(tag["tag_id"] == test_tag for tag in tags)


def test_unlink_tag_from_note_not_found(api_client, test_note):
    """Test unlinking a non-existent tag"""
    unlink_data = {
        "tag_id": str(uuid.uuid4()),  # Non-existent tag ID
        "note_id": test_note,
    }
    response = api_client.post(
        f"{api_client.base_url}/api/tags/unlink", json=unlink_data
    )
    assert response.status_code == 404


def test_get_notes_by_tag(api_client, test_tag, test_note):
    """Test getting notes by tag"""
    response = api_client.get(f"{api_client.base_url}/api/tags/{test_tag}/notes")
    assert response.status_code == 200
    notes = response.json()
    print(notes)
    assert len(notes) == 1
    assert notes[0]["id"] == test_note


def test_get_tags_by_note(api_client, test_tag, test_note):
    """Test getting tags by note"""
    response = api_client.get(f"{api_client.base_url}/api/tags/note/{test_note}")
    assert response.status_code == 200
    tags = response.json()
    assert len(tags) == 1
    assert tags[0]["tag_id"] == test_tag


def test_unlink_tags(api_client, test_tag, test_note):
    """Test unlinking two tags"""
    # Create another tag
    tag_data = {"name": "second-tag", "note_id": test_note}
    response = api_client.post(f"{api_client.base_url}/api/tags/create", json=tag_data)
    assert response.status_code == 201
    second_tag_id = response.json()["tag_id"]

    # Link tags first
    link_data = {"tag1_id": test_tag, "tag2_id": second_tag_id}
    response = api_client.post(f"{api_client.base_url}/api/tags/link", json=link_data)
    assert response.status_code == 200

    # Test successful unlinking
    unlink_data = {"tag1_id": test_tag, "tag2_id": second_tag_id}
    response = api_client.post(
        f"{api_client.base_url}/api/tags/unlink-tags", json=unlink_data
    )
    assert response.status_code == 200

    # Verify tags are unlinked
    response = api_client.get(f"{api_client.base_url}/api/tags/{test_tag}/linked")
    assert response.status_code == 200
    linked_tags = response.json()
    assert len(linked_tags) == 0


def test_get_linked_tags(api_client, test_tag, test_note):
    """Test getting linked tags"""
    # Create another tag
    tag_data = {"name": "second-tag", "note_id": test_note}
    response = api_client.post(f"{api_client.base_url}/api/tags/create", json=tag_data)
    assert response.status_code == 201
    second_tag_id = response.json()["tag_id"]

    # Link tags
    link_data = {"tag1_id": test_tag, "tag2_id": second_tag_id}
    response = api_client.post(f"{api_client.base_url}/api/tags/link", json=link_data)
    assert response.status_code == 200

    # Test getting linked tags
    response = api_client.get(f"{api_client.base_url}/api/tags/{test_tag}/linked")
    assert response.status_code == 200
    linked_tags = response.json()
    assert len(linked_tags) == 1
    assert linked_tags[0]["tag_id"] == second_tag_id


def test_get_linked_tags_not_found(api_client, test_user):
    """Test getting linked tags for non-existent tag"""
    response = api_client.get(f"{api_client.base_url}/api/tags/{uuid.uuid4()}/linked")
    assert response.status_code == 404


def test_unlink_non_linked_tag(api_client, test_tag, root_dir):
    """Test unlinking a tag that is not linked to a note"""
    note1_id = create_note(api_client, "Note 1", "automerge-url-1", root_dir)

    unlink_data = {"tag_id": test_tag, "note_id": note1_id}
    response = api_client.post(
        f"{api_client.base_url}/api/tags/unlink", json=unlink_data
    )
    assert response.status_code == 404, response.text


def test_get_tags_for_note_with_no_tags(api_client, root_dir):
    """Test getting tags for a note with no tags"""
    note_id = create_note(api_client, "Note Without Tags", "automerge-url", root_dir)
    response = api_client.get(f"{api_client.base_url}/api/tags/note/{note_id}")
    assert response.status_code == 200, response.text
    assert response.json() == []


def test_delete_tag(api_client, test_tag, test_note):
    """Test deleting a tag completely"""
    # Create another tag and link it to test_tag
    tag_data = {"name": "second-tag", "note_id": test_note}
    response = api_client.post(f"{api_client.base_url}/api/tags/create", json=tag_data)
    assert response.status_code == 201
    second_tag_id = response.json()["tag_id"]

    # Link tags
    link_data = {"tag1_id": test_tag, "tag2_id": second_tag_id}
    response = api_client.post(f"{api_client.base_url}/api/tags/link", json=link_data)
    assert response.status_code == 200

    # Delete the tag
    response = api_client.post(f"{api_client.base_url}/api/tags/delete", json={"tag_id": test_tag})
    assert response.status_code == 200

    # Verify tag is deleted from note
    response = api_client.get(f"{api_client.base_url}/api/tags/note/{test_note}")
    assert response.status_code == 200
    tags = response.json()
    assert not any(tag["tag_id"] == test_tag for tag in tags)

    # Verify tag is deleted from linked tags
    response = api_client.get(f"{api_client.base_url}/api/tags/{second_tag_id}/linked")
    assert response.status_code == 200
    linked_tags = response.json()
    assert not any(tag["tag_id"] == test_tag for tag in linked_tags)

    # Verify tag doesn't exist anymore
    response = api_client.get(f"{api_client.base_url}/api/tags/{test_tag}/notes")
    assert response.status_code == 404


def test_delete_non_existent_tag(api_client, test_user):
    """Test deleting a non-existent tag"""
    response = api_client.post(
        f"{api_client.base_url}/api/tags/delete",
        json={"tag_id": str(uuid.uuid4())}
    )
    assert response.status_code == 404


def test_delete_tag_invalid_request(api_client):
    """Test deleting a tag with invalid request format"""
    # Test with missing tag_id
    response = api_client.post(
        f"{api_client.base_url}/api/tags/delete",
        json={}
    )
    assert response.status_code == 400

    # Test with invalid tag_id format
    response = api_client.post(
        f"{api_client.base_url}/api/tags/delete",
        json={"tag_id": "invalid-uuid"}
    )
    assert response.status_code == 400


def test_delete_tag_with_multiple_relations(api_client, root_dir):
    """Test deleting a tag that has multiple relations"""
    # Create multiple notes
    note1_id = create_note(api_client, "Note 1", "automerge-url-1", root_dir)
    note2_id = create_note(api_client, "Note 2", "automerge-url-2", root_dir)
    note3_id = create_note(api_client, "Note 3", "automerge-url-3", root_dir)

    # Create and link a tag to all notes
    tag_data = {"name": "multi-note-tag", "note_id": note1_id}
    response = api_client.post(f"{api_client.base_url}/api/tags/create", json=tag_data)
    assert response.status_code == 201
    tag_id = response.json()["tag_id"]

    # Link the same tag to other notes
    tag_data["note_id"] = note2_id
    response = api_client.post(f"{api_client.base_url}/api/tags/create", json=tag_data)
    assert response.status_code == 201

    tag_data["note_id"] = note3_id
    response = api_client.post(f"{api_client.base_url}/api/tags/create", json=tag_data)
    assert response.status_code == 201

    # Create and link another tag
    second_tag_data = {"name": "second-tag", "note_id": note1_id}
    response = api_client.post(f"{api_client.base_url}/api/tags/create", json=second_tag_data)
    assert response.status_code == 201
    second_tag_id = response.json()["tag_id"]

    # Link the tags together
    link_data = {"tag1_id": tag_id, "tag2_id": second_tag_id}
    response = api_client.post(f"{api_client.base_url}/api/tags/link", json=link_data)
    assert response.status_code == 200

    # Delete the first tag
    response = api_client.post(f"{api_client.base_url}/api/tags/delete", json={"tag_id": tag_id})
    assert response.status_code == 200

    # Verify tag is deleted from all notes
    for note_id in [note1_id, note2_id, note3_id]:
        response = api_client.get(f"{api_client.base_url}/api/tags/note/{note_id}")
        assert response.status_code == 200
        tags = response.json()
        assert not any(tag["tag_id"] == tag_id for tag in tags)

    # Verify tag is deleted from linked tag
    response = api_client.get(f"{api_client.base_url}/api/tags/{second_tag_id}/linked")
    assert response.status_code == 200
    linked_tags = response.json()
    assert not any(tag["tag_id"] == tag_id for tag in linked_tags)

    # Verify tag doesn't exist anymore
    response = api_client.get(f"{api_client.base_url}/api/tags/{tag_id}/notes")
    assert response.status_code == 404


def test_delete_tag_cascade_effects(api_client, root_dir):
    """Test that deleting a tag doesn't affect other tags or notes"""
    # Create two notes
    note1_id = create_note(api_client, "Note 1", "automerge-url-1", root_dir)
    note2_id = create_note(api_client, "Note 2", "automerge-url-2", root_dir)

    # Create and link two different tags to the first note
    tag1_data = {"name": "tag1", "note_id": note1_id}
    response = api_client.post(f"{api_client.base_url}/api/tags/create", json=tag1_data)
    assert response.status_code == 201
    tag1_id = response.json()["tag_id"]

    tag2_data = {"name": "tag2", "note_id": note1_id}
    response = api_client.post(f"{api_client.base_url}/api/tags/create", json=tag2_data)
    assert response.status_code == 201
    tag2_id = response.json()["tag_id"]

    # Create and link a tag to the second note
    tag3_data = {"name": "tag3", "note_id": note2_id}
    response = api_client.post(f"{api_client.base_url}/api/tags/create", json=tag3_data)
    assert response.status_code == 201
    tag3_id = response.json()["tag_id"]

    # Delete tag1
    response = api_client.post(f"{api_client.base_url}/api/tags/delete", json={"tag_id": tag1_id})
    assert response.status_code == 200

    # Verify tag1 is deleted from note1
    response = api_client.get(f"{api_client.base_url}/api/tags/note/{note1_id}")
    assert response.status_code == 200
    tags = response.json()
    assert not any(tag["tag_id"] == tag1_id for tag in tags)
    assert any(tag["tag_id"] == tag2_id for tag in tags)

    # Verify note2 and its tag are unaffected
    response = api_client.get(f"{api_client.base_url}/api/tags/note/{note2_id}")
    assert response.status_code == 200
    tags = response.json()
    assert len(tags) == 1
    assert tags[0]["tag_id"] == tag3_id

@pytest.mark.xfail
def test_tag_isolation_between_users(api_client, second_user_client, test_note, second_user_note):
    """Test that users cannot see or modify other users' tags"""
    # First user creates a tag
    tag_data = {"name": "isolated-tag", "note_id": test_note}
    response = api_client.post(
        f"{api_client.base_url}/api/tags/create", 
        json=tag_data
    )
    assert response.status_code == 201
    first_user_tag_id = response.json()["tag_id"]
    
    # Second user cannot see the first user's tag
    response = second_user_client.get(f"{second_user_client.base_url}/api/tags/{first_user_tag_id}/notes")
    assert response.status_code == 404, "Second user should not be able to see first user's tag"
    
    # Second user creates a tag with the same name
    tag_data = {"name": "isolated-tag", "note_id": second_user_note}
    response = second_user_client.post(
        f"{second_user_client.base_url}/api/tags/create", 
        json=tag_data
    )
    assert response.status_code == 201
    second_user_tag_id = response.json()["tag_id"]
    assert first_user_tag_id != second_user_tag_id, "Tags with same name for different users should have different IDs"
    
    # Second user cannot link first user's tag to their note
    response = second_user_client.post(
        f"{second_user_client.base_url}/api/tags/{first_user_tag_id}/link/{second_user_note}"
    )
    assert response.status_code == 404, "Second user should not be able to link first user's tag"
    
    # First user cannot link their tag to second user's note
    response = api_client.post(
        f"{api_client.base_url}/api/tags/{first_user_tag_id}/link/{second_user_note}"
    )
    assert response.status_code == 403, "First user should not be able to link to second user's note"

@pytest.mark.xfail
def test_tag_deletion_isolation(api_client, second_user_client, test_note, second_user_note):
    """Test that users cannot delete other users' tags"""
    # First user creates a tag
    tag_data = {"name": "delete-isolation-tag", "note_id": test_note}
    response = api_client.post(
        f"{api_client.base_url}/api/tags/create",
        json=tag_data
    )
    assert response.status_code == 201
    first_user_tag_id = response.json()["tag_id"]

    # Second user cannot delete first user's tag
    response = second_user_client.post(f"{second_user_client.base_url}/api/tags/delete", json={"tag_id": first_user_tag_id})
    print(response.text)
    assert response.status_code == 404, "Second user should not be able to delete first user's tag"

    # Verify first user's tag still exists
    response = api_client.get(f"{api_client.base_url}/api/tags/note/{test_note}")
    assert response.status_code == 200
    tags = response.json()
    assert any(tag["tag_id"] == first_user_tag_id for tag in tags), "First user's tag should still exist"


def test_user_cannot_link_tag_to_other_users_note(api_client, second_user_client, test_note, second_user_note):
    """Test that users cannot link tags to other users' notes"""
    # First user creates a tag
    tag_data = {"name": "cross-linking-tag", "note_id": test_note}
    response = api_client.post(
        f"{api_client.base_url}/api/tags/create", 
        json=tag_data
    )
    assert response.status_code == 201
    tag_id = response.json()["tag_id"]
    
    # First user cannot link their tag to second user's note
    response = api_client.post(f"{api_client.base_url}/api/tags/{tag_id}/link/{second_user_note}")
    assert response.status_code == 403, "First user should not be able to link their tag to second user's note"


def test_link_existing_tag(api_client, test_tag, root_dir):
    """Test linking an existing tag to a different note"""
    # Create a second note
    note2_id = create_note(api_client, "Second Note", "automerge-url-2", root_dir)
    
    # Link the existing tag to the second note
    response = api_client.post(f"{api_client.base_url}/api/tags/{test_tag}/link/{note2_id}")
    assert response.status_code == 201, response.text
    
    # Verify tag is linked to both notes
    response = api_client.get(f"{api_client.base_url}/api/tags/{test_tag}/notes")
    assert response.status_code == 200, response.text
    notes = response.json()
    assert len(notes) == 2, "Tag should be linked to two notes"
    note_ids = [note["id"] for note in notes]
    assert note2_id in note_ids, "Tag should be linked to the second note"
    
    # Try linking the same tag to the same note again (should fail)
    response = api_client.post(f"{api_client.base_url}/api/tags/{test_tag}/link/{note2_id}")
    assert response.status_code == 409, "Linking the same tag to the same note again should fail"
    
    # Test with non-existent tag ID
    non_existent_tag_id = str(uuid.uuid4())
    response = api_client.post(f"{api_client.base_url}/api/tags/{non_existent_tag_id}/link/{note2_id}")
    assert response.status_code == 404, "Linking a non-existent tag should fail"


def test_tag_with_same_name_for_different_users(api_client, second_user_client, test_note, second_user_note):
    """Test that different users can create tags with the same name"""
    # Create a tag with the same name for both users
    tag_name = "shared-tag-name"
    
    # First user creates a tag
    tag_data = {"name": tag_name, "note_id": test_note}
    response = api_client.post(f"{api_client.base_url}/api/tags/create", json=tag_data)
    assert response.status_code == 201
    first_user_tag_id = response.json()["tag_id"]
    
    # Second user creates a tag with the same name
    tag_data = {"name": tag_name, "note_id": second_user_note}
    response = second_user_client.post(f"{second_user_client.base_url}/api/tags/create", json=tag_data)
    assert response.status_code == 201
    second_user_tag_id = response.json()["tag_id"]
    
    # Verify the tags have different IDs
    assert first_user_tag_id != second_user_tag_id, "Tags with same name for different users should have different IDs"
    
    # Verify first user can see their tag
    response = api_client.get(f"{api_client.base_url}/api/tags/note/{test_note}")
    assert response.status_code == 200
    tags = response.json()
    assert any(tag["name"] == tag_name for tag in tags), "First user should see their tag"
    
    # Verify second user can see their tag
    response = second_user_client.get(f"{second_user_client.base_url}/api/tags/note/{second_user_note}")
    assert response.status_code == 200
    tags = response.json()
    assert any(tag["name"] == tag_name for tag in tags), "Second user should see their tag"

@pytest.mark.xfail
def test_user_cannot_link_tag_to_other_users_note(api_client, second_user_client, test_tag, second_user_note):
    """Test that a user cannot link their tag to another user's note"""
    # First user tries to link their tag to second user's note
    tag_data = {"name": "first-user-tag", "note_id": second_user_note}
    response = api_client.post(f"{api_client.base_url}/api/tags/create", json=tag_data)
    assert response.status_code == 404, "User should not be able to link tag to another user's note"


def test_tag_operations_after_user_switch(api_client, second_user_client, test_tag, second_user_tag, test_note, second_user_note):
    """Test tag operations after switching between users"""
    # First user creates a tag
    tag_data = {"name": "temp-tag", "note_id": test_note}
    response = api_client.post(f"{api_client.base_url}/api/tags/create", json=tag_data)
    assert response.status_code == 201
    temp_tag_id = response.json()["tag_id"]
    
    # Second user creates a tag
    tag_data = {"name": "temp-tag-2", "note_id": second_user_note}
    response = second_user_client.post(f"{second_user_client.base_url}/api/tags/create", json=tag_data)
    assert response.status_code == 201
    temp_tag_id_2 = response.json()["tag_id"]
    
    # First user deletes their temporary tag
    response = api_client.post(f"{api_client.base_url}/api/tags/delete", json={"tag_id": temp_tag_id})
    assert response.status_code == 200
    
    # Second user verifies their tag still exists
    response = second_user_client.get(f"{second_user_client.base_url}/api/tags/note/{second_user_note}")
    assert response.status_code == 200
    tags = response.json()
    assert any(tag["tag_id"] == temp_tag_id_2 for tag in tags), "Second user's tag should still exist"
    
    # Second user deletes their temporary tag
    response = second_user_client.post(f"{second_user_client.base_url}/api/tags/delete", json={"tag_id": temp_tag_id_2})
    assert response.status_code == 200
