package model

import (
	"encoding/json"
	"testing"
)

func TestStudent_JSONMarshaling(t *testing.T) {
	student := Student{
		ID:         1,
		FirstName:  "John",
		LastName:   "Doe",
		Nickname:   "JD",
		ParentName: "Jane Doe",
		Address1:   "123 Main St",
		Address2:   "Apt 4B",
		City:       "Springfield",
		State:      "IL",
		PostalCode: "62701",
		Email:      "john.doe@example.com",
		Grade:      "10",
	}

	data, err := json.Marshal(student)
	if err != nil {
		t.Fatalf("unexpected error marshaling student: %v", err)
	}

	expectedJSON := `{"id":1,"firstName":"John","lastName":"Doe","nickname":"JD","parentName":"Jane Doe","address1":"123 Main St","address2":"Apt 4B","city":"Springfield","state":"IL","postalCode":"62701","email":"john.doe@example.com","grade":"10"}`
	if string(data) != expectedJSON {
		t.Fatalf("expected JSON %s, got %s", expectedJSON, string(data))
	}

	var decoded Student
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unexpected error unmarshaling student: %v", err)
	}

	if decoded != student {
		t.Fatalf("expected decoded student %+v, got %+v", student, decoded)
	}
}
