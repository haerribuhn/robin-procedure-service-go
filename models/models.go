package models

import "time"

// Procedure schema of the procedure table
type Procedure struct {
	ID               int64     `json:"id"`
	LastModifiedOn   time.Time `json:"lastModifiedOn"`
	StructureID      int64     `json:"structureId"`
	StructureVersion int64     `json:"structureVersion"`
	Name             string    `json:"name"`
	Commodity        string    `json:"commodity"`
	ConsultantID     int64     `json:"consultantId"`
	DeadLine         time.Time `json:"deadLine"`
}
