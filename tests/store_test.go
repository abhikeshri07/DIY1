package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/abhikeshri07/go-mux/constants"
	"github.com/abhikeshri07/go-mux/models"
	"net/http"
	"regexp"
	"testing"
)

func TestGetProductsWithWrongID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "stores" WHERE store_id = $1`)).WithArgs(5).WillReturnRows(sqlmock.NewRows([]string{"store_id", "product_id", "is_available"}))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "product_id" FROM "stores" WHERE store_id = $1`)).WithArgs(5).WillReturnRows(sqlmock.NewRows([]string{"product_id"}))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE id IN (NULL)`)).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price"}))
	mock.ExpectCommit()
	req, _ := http.NewRequest("GET", "/stores/5/products", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
	var m []models.ProductModel
	err = json.Unmarshal(response.Body.Bytes(), &m)
	if err != nil {
		return
	}
	if len(m) != 0 {
		t.Errorf("Expected response size to be 0 but got %q", len(m))
	}
}

func TestGetProductsSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	testStoreRows := sqlmock.NewRows([]string{"store_id", "product_id", "is_available"}).
		AddRow(1, 2, true)
	testProductIdRows := sqlmock.NewRows([]string{"product_id"}).
		AddRow(2)
	testProductRow := sqlmock.NewRows([]string{"id", "name", "price"}).
		AddRow(2, "Coffee", 20)
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "stores" WHERE store_id = $1`)).WithArgs(1).WillReturnRows(testStoreRows)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "product_id" FROM "stores" WHERE store_id = $1`)).WithArgs(1).WillReturnRows(testProductIdRows)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE id IN (2)`)).WillReturnRows(testProductRow)
	mock.ExpectCommit()
	req, _ := http.NewRequest("GET", "/stores/1/products", nil)
	response := executeRequest(req)
	if response.Code == http.StatusBadRequest {
		var m map[string]string
		json.Unmarshal(response.Body.Bytes(), &m)
		if m["error"] == constants.STORE_NOT_FOUND_ERROR {
			t.Errorf("Expected response store exist bu store not found with given storeId")
		}
	}
	var m []models.ProductModel
	err = json.Unmarshal(response.Body.Bytes(), &m)
	if err != nil {
		return
	}
	if len(m) == 0 {
		t.Errorf("Expected response size to be greater than 0 but got %q", len(m))
	}
}

// todo mock
func TestAddProductsEmpty(t *testing.T) {
	query := "SELECT * FROM products"
	result := a.DB.Exec(query)
	initialCount := result.RowsAffected
	var m []models.ProductModel
	s, _ := json.Marshal(&m)
	req, _ := http.NewRequest("POST", "/stores/1", bytes.NewBuffer(s))
	response := executeRequest(req)
	result = a.DB.Exec(query)
	finalCount := result.RowsAffected
	fmt.Println(finalCount)
	checkResponseCode(t, http.StatusOK, response.Code)
	if finalCount != initialCount {
		t.Errorf("Unnecessary product added.")
	}
}

// todo use test equal
func TestAddProductsSuccess(t *testing.T) {

	query := "SELECT * FROM stores"
	result := a.DB.Exec(query)
	initialCount := result.RowsAffected
	var m []int64
	m = append(m, 5)
	s, _ := json.Marshal(&m)
	req, _ := http.NewRequest("POST", "/stores/1", bytes.NewBuffer(s))
	response := executeRequest(req)
	result = a.DB.Exec(query)
	finalCount := result.RowsAffected
	checkResponseCode(t, http.StatusOK, response.Code)
	if finalCount != initialCount+1 {
		t.Errorf("Product not added.")
	}
}
