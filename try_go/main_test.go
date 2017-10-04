package main_test

import (
	"os"
	"log"
	"testing"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"bytes"
	"."
)

var a main.App


func TestMain(m *testing.M) {
	a = main.App{}

	a.Initialize("user", "pswd", "db")

	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/persons", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentPerson(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/persons/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Person not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Person not found'. Got '%s'", m["error"])
	}
}

func TestCreateProduct(t *testing.T) {
	clearTable()

	payload := []byte(`{"name":"Jon Doe","data":"text"}`)

	req, _ := http.NewRequest("POST", "/persons", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "Jon Doe" {
		t.Errorf("Expected product name to be 'test product'. Got '%v'", m["name"])
	}

	if m["data"] != "text" {
		t.Errorf("Expected product price to be '11.22'. Got '%v'", m["price"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected person ID to be '1'. Got '%v'", m["id"])
	}
}

//TODO related to TestCreateProduct
func TestGetProduct(t *testing.T) {

	req, _ := http.NewRequest("GET", "/persons/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}



func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM person")
	a.DB.Exec("ALTER SEQUENCE person_id_seq RESTART WITH 1")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS person
(
id SERIAL PRIMARY KEY,
name varchar(100) NOT NULL,
data varchar(50) NOT NULL DEFAULT 'text'
)`


