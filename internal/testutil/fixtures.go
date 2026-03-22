package testutil

import "encoding/json"

func OpportunityFixture() json.RawMessage {
	return json.RawMessage(`{
		"data": {
			"id": "opp-123",
			"name": "Jane Doe",
			"headline": "Software Engineer",
			"stage": "stage-456",
			"origin": "sourced",
			"createdAt": 1640000000000,
			"updatedAt": 1640000000000
		}
	}`)
}

func OpportunityListFixture() json.RawMessage {
	return json.RawMessage(`{
		"data": [
			{"id": "opp-123", "name": "Jane Doe"},
			{"id": "opp-456", "name": "John Smith"}
		],
		"hasNext": false
	}`)
}

func NoteFixture() json.RawMessage {
	return json.RawMessage(`{
		"data": {
			"id": "note-123",
			"text": "Great candidate",
			"createdAt": 1640000000000
		}
	}`)
}

func NoteListFixture() json.RawMessage {
	return json.RawMessage(`{
		"data": [
			{"id": "note-123", "text": "Great candidate"},
			{"id": "note-456", "text": "Follow up scheduled"}
		],
		"hasNext": false
	}`)
}

func UserListFixture() json.RawMessage {
	return json.RawMessage(`{
		"data": [
			{"id": "user-123", "name": "Alice", "email": "alice@example.com"},
			{"id": "user-456", "name": "Bob", "email": "bob@example.com"}
		],
		"hasNext": false
	}`)
}

func StageListFixture() json.RawMessage {
	return json.RawMessage(`{
		"data": [
			{"id": "stage-123", "text": "New Lead"},
			{"id": "stage-456", "text": "Phone Screen"}
		],
		"hasNext": false
	}`)
}

func PostingFixture() json.RawMessage {
	return json.RawMessage(`{
		"data": {
			"id": "posting-123",
			"text": "Senior Engineer",
			"state": "published",
			"team": "Engineering"
		}
	}`)
}

func ArchiveReasonListFixture() json.RawMessage {
	return json.RawMessage(`{
		"data": [
			{"id": "reason-1", "text": "Hired", "type": "hired"},
			{"id": "reason-2", "text": "Not a fit", "type": "non-hired"}
		],
		"hasNext": false
	}`)
}

func EmptyListFixture() json.RawMessage {
	return json.RawMessage(`{"data": [], "hasNext": false}`)
}
