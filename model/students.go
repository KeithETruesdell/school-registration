package model

// Student represents a student record.
type Student struct {
	ID         int64  `json:"id"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Nickname   string `json:"nickname"`
	ParentName string `json:"parentName"`
	Address1   string `json:"address1"`
	Address2   string `json:"address2"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Email      string `json:"email"`
	Grade      string `json:"grade"`
}
