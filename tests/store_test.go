package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/abhikeshri07/go-mux/models"
	"net/http"
	"testing"
)

func TestGetProductsWithWrongID(t *testing.T) {
	req, _ := http.NewRequest("GET", "/stores/5/products", nil)
	response := executeRequest(req)
	var m []models.ProductModel
	err := json.Unmarshal(response.Body.Bytes(), &m)
	if err != nil {
		return
	}
	if len(m) != 0 {
		t.Errorf("Expected response size to be 0 but got %q", len(m))
	}
}

func TestGetProductsSuccess(t *testing.T) {

	req, _ := http.NewRequest("GET", "/stores/1/products", nil)
	response := executeRequest(req)
	var m []models.ProductModel
	err := json.Unmarshal(response.Body.Bytes(), &m)
	if err != nil {
		return
	}
	if len(m) == 0 {
		t.Errorf("Expected response size to be greater than 0 but got %q", len(m))
	}
}

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

func TestAddProductsSuccess(t *testing.T) {

	query := "SELECT * FROM products"
	result := a.DB.Exec(query)
	initialCount := result.RowsAffected
	var m []models.ProductModel
	m = append(m, models.ProductModel{Name: "New Test Product", Price: 11.69})
	s, _ := json.Marshal(&m)
	req, _ := http.NewRequest("POST", "/stores/1", bytes.NewBuffer(s))
	response := executeRequest(req)
	result = a.DB.Exec(query)
	finalCount := result.RowsAffected
	checkResponseCode(t, http.StatusOK, response.Code)
	if finalCount != initialCount+1 {
		t.Errorf("product not added added.")
	}

}
