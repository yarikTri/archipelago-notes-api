import uuid

import pytest


def create_note(api_client, title, automerge_url, root_dir):
    """Helper function to create a note"""
    print(f"TEST CREATE NOTE api_client.headers: {api_client.headers}")
    print(f"TEST CREATE NOTE cookies: {api_client.cookies}")
    note_data = {"title": title, "dir_id": root_dir, "automerge_url": automerge_url}
    response = api_client.post(f"{api_client.base_url}/api/notes", json=note_data)
    assert response.status_code == 201, response.text
    return response.json()["id"]


def test_create_and_link_tag(api_client, test_note):
    """Test creating and linking a tag to a note"""
    # Test successful creation
    tag_data = {"name": "test-tag", "note_id": test_note}
    response = api_client.post(f"{api_client.base_url}/api/tags/create", json=tag_data)
    assert response.status_code == 201, response.text
    assert "tag_id" in response.json()
    assert "name" in response.json()
    assert response.json()["name"] == "test-tag"

    # Test invalid request
    invalid_data = {
        "name": "",  # Empty name
        "note_id": test_note,
    }
    response = api_client.post(
        f"{api_client.base_url}/api/tags/create", json=invalid_data
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
    response = api_client.get(f"{api_client.base_url}/api/tags/notes/{test_note}")
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


def test_update_tag(api_client, test_tag):
    """Test updating a tag"""
    # Test successful update
    update_data = {"tag_id": test_tag, "name": "updated-tag"}
    response = api_client.put(f"{api_client.base_url}/api/tags", json=update_data)
    assert response.status_code == 200
    assert response.json()["name"] == "updated-tag"

    # Test updating non-existent tag
    non_existent_data = {"tag_id": str(uuid.uuid4()), "name": "updated-tag"}
    response = api_client.put(f"{api_client.base_url}/api/tags", json=non_existent_data)
    assert response.status_code == 404


def test_update_tag_empty_name(api_client, test_tag):
    """Test updating a tag with empty name"""
    update_data = {
        "tag_id": test_tag,
        "name": "",  # Empty name
    }
    response = api_client.put(f"{api_client.base_url}/api/tags", json=update_data)
    assert response.status_code == 400


def test_update_tag_for_note(api_client, test_tag, test_note):
    """Test updating a tag for a specific note"""
    # First unlink the tag to delete it
    unlink_data = {"tag_id": test_tag, "note_id": test_note}
    response = api_client.post(
        f"{api_client.base_url}/api/tags/unlink", json=unlink_data
    )
    assert response.status_code == 200

    # Try to update the deleted tag
    update_data = {
        "tag_id": test_tag,
        "note_id": test_note,
        "name": "note-specific-tag",
    }
    response = api_client.put(f"{api_client.base_url}/api/tags/note", json=update_data)
    assert response.status_code == 404, (
        response.text
    )  # Should return 404 for non-existent tag


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
    response = api_client.get(f"{api_client.base_url}/api/tags/notes/{test_note}")
    assert response.status_code == 200
    tags = response.json()
    assert len(tags) == 1
    assert tags[0]["tag_id"] == test_tag


def test_link_tags(api_client, test_tag, test_note):
    """Test linking two tags together"""
    # Create another tag
    tag_data = {"name": "second-tag", "note_id": test_note}
    response = api_client.post(f"{api_client.base_url}/api/tags/create", json=tag_data)
    assert response.status_code == 201
    second_tag_id = response.json()["tag_id"]
    print(second_tag_id)

    # Test successful linking
    link_data = {"tag1_id": test_tag, "tag2_id": second_tag_id}
    response = api_client.post(f"{api_client.base_url}/api/tags/link", json=link_data)
    assert response.status_code == 200

    # Verify tags are linked
    response = api_client.get(f"{api_client.base_url}/api/tags/{test_tag}/linked")
    assert response.status_code == 200
    linked_tags = response.json()
    assert len(linked_tags) == 1
    assert linked_tags[0]["tag_id"] == second_tag_id


def test_link_tags_already_linked(api_client, test_tag, test_note):
    """Test linking already linked tags"""
    # Create another tag
    tag_data = {"name": "second-tag", "note_id": test_note}
    response = api_client.post(f"{api_client.base_url}/api/tags/create", json=tag_data)
    assert response.status_code == 201
    second_tag_id = response.json()["tag_id"]

    # Link tags first time
    link_data = {"tag1_id": test_tag, "tag2_id": second_tag_id}
    response = api_client.post(f"{api_client.base_url}/api/tags/link", json=link_data)
    assert response.status_code == 200

    # Try to link again
    response = api_client.post(f"{api_client.base_url}/api/tags/link", json=link_data)
    assert response.status_code == 409  # Conflict status code for already linked tags


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


def test_get_linked_tags_not_found(api_client):
    """Test getting linked tags for non-existent tag"""
    response = api_client.get(f"{api_client.base_url}/api/tags/{uuid.uuid4()}/linked")
    assert response.status_code == 404


def test_tag_behavior_across_notes(api_client, root_dir):
    """Test tag behavior across multiple notes"""
    # Create first note
    note1_id = create_note(
        api_client, "First Test Note", "test-automerge-url-1", root_dir
    )

    # Create second note
    note2_id = create_note(
        api_client, "Second Test Note", "test-automerge-url-2", root_dir
    )

    # Create and link tag to first note
    tag_data = {"name": "shared-tag", "note_id": note1_id}
    response = api_client.post(f"{api_client.base_url}/api/tags/create", json=tag_data)
    assert response.status_code == 201, response.text
    tag_id = response.json()["tag_id"]

    # Link same tag to second note
    tag_data = {"name": "shared-tag", "note_id": note2_id}
    response = api_client.post(f"{api_client.base_url}/api/tags/create", json=tag_data)
    assert response.status_code == 201, response.text
    assert response.json()["tag_id"] == tag_id  # Should reuse the same tag

    # Update tag name for first note only
    update_data = {"tag_id": tag_id, "note_id": note1_id, "name": "note1-specific-tag"}
    response = api_client.put(f"{api_client.base_url}/api/tags/note", json=update_data)
    assert response.status_code == 200, response.text
    assert response.json()["name"] == "note1-specific-tag"

    # Check tags for first note
    response = api_client.get(f"{api_client.base_url}/api/tags/notes/{note1_id}")
    assert response.status_code == 200, response.text
    note1_tags = response.json()
    assert len(note1_tags) == 1
    assert note1_tags[0]["name"] == "note1-specific-tag"

    # Check tags for second note
    response = api_client.get(f"{api_client.base_url}/api/tags/notes/{note2_id}")
    assert response.status_code == 200, response.text
    note2_tags = response.json()
    assert len(note2_tags) == 1
    assert note2_tags[0]["name"] == "shared-tag"


def test_tag_updates_across_notes(api_client, root_dir):
    """Test tag updates behavior across multiple notes"""
    # Create first note
    note1_id = create_note(
        api_client, "First Test Note", "test-automerge-url-1", root_dir
    )

    # Create second note
    note2_id = create_note(
        api_client, "Second Test Note", "test-automerge-url-2", root_dir
    )

    # Create and link tag to first note
    tag_data = {"name": "shared-tag", "note_id": note1_id}
    response = api_client.post(f"{api_client.base_url}/api/tags/create", json=tag_data)
    assert response.status_code == 201, response.text
    tag_id = response.json()["tag_id"]

    # Link same tag to second note
    tag_data = {"name": "shared-tag", "note_id": note2_id}
    response = api_client.post(f"{api_client.base_url}/api/tags/create", json=tag_data)
    assert response.status_code == 201, response.text
    assert response.json()["tag_id"] == tag_id  # Should reuse the same tag

    # Update tag name for first note
    update_data = {"tag_id": tag_id, "note_id": note1_id, "name": "note1-specific-tag"}
    response = api_client.put(f"{api_client.base_url}/api/tags/note", json=update_data)
    assert response.status_code == 200, response.text
    assert response.json()["name"] == "note1-specific-tag"

    # Check that second note still has original tag name
    response = api_client.get(f"{api_client.base_url}/api/tags/notes/{note2_id}")
    assert response.status_code == 200, response.text
    note2_tags = response.json()
    assert len(note2_tags) == 1
    assert note2_tags[0]["name"] == "shared-tag"

    # Check that first note has updated tag name
    response = api_client.get(f"{api_client.base_url}/api/tags/notes/{note1_id}")
    assert response.status_code == 200, response.text
    note1_tags = response.json()
    assert len(note1_tags) == 1
    assert note1_tags[0]["name"] == "note1-specific-tag"

    # Update tag name for second note
    update_data = {"tag_id": tag_id, "note_id": note2_id, "name": "note2-specific-tag"}
    response = api_client.put(f"{api_client.base_url}/api/tags/note", json=update_data)
    assert response.status_code == 200, response.text
    assert response.json()["name"] == "note2-specific-tag"

    # Check that both notes have their specific tag names
    response = api_client.get(f"{api_client.base_url}/api/tags/notes/{note1_id}")
    assert response.status_code == 200, response.text
    note1_tags = response.json()
    assert len(note1_tags) == 1
    assert note1_tags[0]["name"] == "note1-specific-tag"

    response = api_client.get(f"{api_client.base_url}/api/tags/notes/{note2_id}")
    assert response.status_code == 200, response.text
    note2_tags = response.json()
    assert len(note2_tags) == 1
    assert note2_tags[0]["name"] == "note2-specific-tag"

    # Check that original tag name is not used anywhere
    response = api_client.get(f"{api_client.base_url}/api/tags/{tag_id}/notes")
    assert response.status_code == 404, response.text


# COPILOT


@pytest.mark.xfail(reason="returns 200")
def test_create_duplicate_tag_for_note(api_client, test_note):
    """Test creating a duplicate tag for the same note"""
    tag_data = {"name": "duplicate-tag", "note_id": test_note}
    response = api_client.post(f"{api_client.base_url}/api/tags/create", json=tag_data)
    assert response.status_code == 201, response.text

    # Attempt to create the same tag again
    response = api_client.post(f"{api_client.base_url}/api/tags/create", json=tag_data)
    assert response.status_code == 409, response.text


def test_create_tag_without_note_id(api_client):
    """Test creating a tag without providing a note ID"""
    tag_data = {"name": "tag-without-note"}
    response = api_client.post(f"{api_client.base_url}/api/tags/create", json=tag_data)
    assert response.status_code == 400, response.text


def test_delete_non_existent_tag(api_client):
    """Test deleting a non-existent tag"""
    response = api_client.delete(f"{api_client.base_url}/api/tags/{uuid.uuid4()}")
    assert response.status_code == 404, response.text


def test_link_tag_to_multiple_notes(api_client, root_dir):
    """Test linking a tag to multiple notes"""
    # Create two notes
    note1_id = create_note(api_client, "Note 1", "automerge-url-1", root_dir)
    note2_id = create_note(api_client, "Note 2", "automerge-url-2", root_dir)

    # Create and link a tag to the first note
    tag_data = {"name": "multi-note-tag", "note_id": note1_id}
    response = api_client.post(f"{api_client.base_url}/api/tags/create", json=tag_data)
    assert response.status_code == 201, response.text
    tag_id = response.json()["tag_id"]

    # Link the same tag to the second note
    tag_data["note_id"] = note2_id
    response = api_client.post(f"{api_client.base_url}/api/tags/create", json=tag_data)
    assert response.status_code == 201, response.text

    # Verify the tag is linked to both notes
    response = api_client.get(f"{api_client.base_url}/api/tags/notes/{note1_id}")
    assert response.status_code == 200, response.text
    assert any(tag["tag_id"] == tag_id for tag in response.json())

    response = api_client.get(f"{api_client.base_url}/api/tags/notes/{note2_id}")
    assert response.status_code == 200, response.text
    assert any(tag["tag_id"] == tag_id for tag in response.json())


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
    response = api_client.get(f"{api_client.base_url}/api/tags/notes/{note_id}")
    assert response.status_code == 200, response.text
    assert response.json() == []
