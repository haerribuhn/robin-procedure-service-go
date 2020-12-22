package middleware

import (
	"database/sql"
	"encoding/json" // package to encode and decode the json into struct and vice versa
	"fmt"
	"log"
	"net/http"                          // used to access the request and response object of the api
	"os"                                // used to read the environment variable
	"robin-procedure-service-go/models" // models package where Procedure schema is defined
	"strconv"                           // package used to covert string into int type

	"github.com/gorilla/mux" // used to get the params from the route

	"github.com/joho/godotenv" // package used to read the .env file
	_ "github.com/lib/pq"      // postgres golang driver
)

// response format
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// create connection with postgres db
func createConnection() *sql.DB {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Open the connection
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	// check the connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	// return the connection
	return db
}

func CreateProcedure(w http.ResponseWriter, r *http.Request) {

	// create an empty procedure of type models.Procedure
	var procedure models.Procedure

	// decode the json request to procedure
	err := json.NewDecoder(r.Body).Decode(&procedure)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// call insert procedure function and pass the procedure
	insertID := insertProcedure(procedure)

	// format a response object
	res := response{
		ID:      insertID,
		Message: "Procedure created successfully",
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

func UpdateProcedure(w http.ResponseWriter, r *http.Request) {

	// get the procedure id from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// create an empty procedure of type models.Procedure
	var procedure models.Procedure

	// decode the json request to procedure
	err = json.NewDecoder(r.Body).Decode(&procedure)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// call update procedure to update the procedure
	updatedRows := updateProcedure(int64(id), procedure)

	// format the message string
	msg := fmt.Sprintf("Procedure updated successfully. Total rows/record affected %v", updatedRows)

	// format the response message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

func DeleteProcedure(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	// convert the id in string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	deletedRows := deleteProcedure(int64(id))

	// format the message string
	msg := fmt.Sprintf("Procedure deleted successfully. Total rows/record affected %v", deletedRows)

	// format the response message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

func GetAllProcedures(w http.ResponseWriter, r *http.Request) {

	// get all the procedures in the db
	procedures, err := getAllProcedures()

	if err != nil {
		log.Fatalf("Unable to get all procedures. %v", err)
	}

	// send all the procedures as response
	json.NewEncoder(w).Encode(procedures)
}

func GetProcedure(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	procedure, err := getProcedure(int64(id))

	if err != nil {
		log.Fatalf("Unable to get procedure. %v", err)
	}

	// send the response
	json.NewEncoder(w).Encode(procedure)
}

func insertProcedure(procedure models.Procedure) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the insert sql query
	// returning procedure id will return the id of the inserted procedure
	sqlStatement := `INSERT INTO procedure (last_modified_on, structure_id, structure_version, name, commodity, consultant_id, deadline) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	// the inserted id will store in this id
	var id int64

	// execute the sql statement
	// Scan function will save the insert id in the id
	err := db.QueryRow(sqlStatement, procedure.LastModifiedOn, procedure.StructureID, procedure.StructureVersion, procedure.Name, procedure.Commodity, procedure.ConsultantID, procedure.DeadLine).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)

	// return the inserted id
	return id
}

func updateProcedure(id int64, procedure models.Procedure) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the update sql query
	sqlStatement := `UPDATE procedure SET last_modified_on=$2, structure_id=$3, structure_version=$4, name=$5, commodity=$6, consultant_id=$7, deadline=$8 WHERE id=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id, procedure.LastModifiedOn, procedure.StructureID, procedure.StructureVersion, procedure.Name, procedure.Commodity, procedure.ConsultantID, procedure.DeadLine)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}

func deleteProcedure(id int64) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the delete sql query
	sqlStatement := `DELETE FROM procedure WHERE id=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}

func getAllProcedures() ([]models.Procedure, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	var procedures []models.Procedure

	// create the select sql query
	sqlStatement := `SELECT p.id, p.last_modified_on, p.structure_id, p.structure_version, p.name, p.commodity, p.consultant_id, p.deadline FROM procedure p`

	// execute the sql statement
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var procedure models.Procedure

		// unmarshal the row object to procedure
		err = rows.Scan(&procedure.ID, &procedure.LastModifiedOn, &procedure.StructureID, &procedure.StructureVersion, &procedure.Name, &procedure.Commodity, &procedure.ConsultantID, &procedure.DeadLine)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		procedures = append(procedures, procedure)

	}

	return procedures, err
}

func getProcedure(id int64) (models.Procedure, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create a procedure of models.Procedure type
	var procedure models.Procedure

	// create the select sql query
	sqlStatement := `SELECT p.id, p.last_modified_on, p.structure_id, p.structure_version, p.name, p.commodity, p.consultant_id, p.deadline FROM procedure p WHERE p.id=$1`

	// execute the sql statement
	row := db.QueryRow(sqlStatement, id)

	// unmarshal the row object to procedure
	err := row.Scan(&procedure.ID, &procedure.LastModifiedOn, &procedure.StructureID, &procedure.StructureVersion, &procedure.Name, &procedure.Commodity, &procedure.ConsultantID, &procedure.DeadLine)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return procedure, nil
	case nil:
		return procedure, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	// return empty procedure on error
	return procedure, err
}
