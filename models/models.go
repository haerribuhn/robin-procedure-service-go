package models

// User schema of the user table
type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Age      int64  `json:"age"`
}

// Procedure schema of the procedure table
type Procedure struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Commodity    string `json:"commodity"`
	ConsultantID int64  `json:"consultantId"`
}
