package backend_test

import (
	"bytes"
	"encoding/json"
	"example.com/backend"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a backend.App

const tableProductCreationQuery = `create table if not exists products 
(
	id INT NOT NULL PRIMARY KEY AUTOINCREMENT,
	prodcutCode varchar(25) not null,
	name varchar(256) not null,
	inventory int not null,
	price int not null,
	status varchar(64) no null)`

func TestMain(m *testing.M) {
	a = backend.App{}
	a.Initialize()
	ensureTableExists()
	code := m.Run()
	// clearProductTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableProductCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearProductTable() {
	a.DB.Exec("delete products")
	a.DB.Exec("delete from sqlite_sequence where name='products'")
}

func testGetNonExistentProduct(t *testing.T) {
	clearProductTable()
	req, _ := http.NewRequest("GET", "/product/11", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusInternalServerError, response.Code)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "sql: no rows in result set" {
		t.Errorf(
			"Expected the error key of the response to be set to 'sql: no rows in result set' got '%s'",
			m["error"],
		)
	}
}

func TestCreateProduct(t *testing.T) {
	clearProductTable()
	payload := []byte(
		`{"productCode":"TEST12345","name":"ProductTest","inventory":1,"price":1,"status":"testing"}`,
	)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(payload))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["productCode"] != "TEST12345" {
		t.Errorf("exepcted prod code TEST12345 got '%v'", m["productCode"])
	}
	if m["name"] != "ProductTest" {
		t.Errorf("exepcted name ProductTest got '%v'", m["name"])
	}
	if m["inventory"] != 1.0 {
		t.Errorf("exepcted inv 1 got '%v'", m["inventory"])
	}
	if m["price"] != 1.0 {
		t.Errorf("exepcted price 1 got '%v'", m["price"])
	}
	if m["status"] != "testing" {
		t.Errorf("exepcted status testing got '%v'", m["status"])
	}
	if m["id"] != 1.0 {
		t.Errorf("exepcted id 1 got '%v'", m["id"])
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected int, actual int) {
	if expected != actual {
		t.Errorf("Expect response code %d. got %d", expected, actual)
	}
}
