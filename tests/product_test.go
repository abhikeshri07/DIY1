package tests

import (
	"bytes"
	"encoding/json"
	"github.com/abhikeshri07/go-mux/models"
	"net/http"
	"strconv"
	"testing"
)

func TestGetNonExistentProduct(t *testing.T) {
	//clearTable()

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	err := json.Unmarshal(response.Body.Bytes(), &m)
	if err != nil {
		return
	}
	if m["error"] != "Product Not Found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Product not found'. Got '%s'", m["error"])
	}
}

// tom: rewritten function
func TestCreateProduct(t *testing.T) {

	var jsonStr = []byte(`{"name":"test product", "price": 11.22}`)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &m)
	if err != nil {
		return
	}

	if m["name"] != "test product" {
		t.Errorf("Expected product name to be 'test product'. Got '%v'", m["name"])
	}

	if m["price"] != 11.22 {
		t.Errorf("Expected product price to be '11.22'. Got '%v'", m["price"])
	}

	if m["id"] == nil {
		t.Errorf("Expected product ID to be '%v', got nil", m["id"])
	}
}

func TestGetProduct(t *testing.T) {

	req, _ := http.NewRequest("GET", "/product/3", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateProduct(t *testing.T) {
	var jsonStr = []byte(`{"name":"test product", "price": 11.22}`)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	var m models.ProductModel
	err := json.Unmarshal(response.Body.Bytes(), &m)
	if err != nil {
		return
	}

	id := m.ID

	req, _ = http.NewRequest("GET", "/product/"+strconv.FormatInt(id, 10), nil)
	response = executeRequest(req)
	var originalProduct map[string]interface{}
	err = json.Unmarshal(response.Body.Bytes(), &originalProduct)
	if err != nil {
		return
	}
	//
	jsonStr = []byte(`{"name":"test product - updated name", "price": 11.69}`)

	req, _ = http.NewRequest("PUT", "/product/"+strconv.FormatInt(id, 10), bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var updatedProduct map[string]interface{}
	err = json.Unmarshal(response.Body.Bytes(), &updatedProduct)
	if err != nil {
		return
	}

	if updatedProduct["id"] != originalProduct["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalProduct["id"], updatedProduct["id"])
	}

	if updatedProduct["name"] == originalProduct["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalProduct["name"], updatedProduct["name"], updatedProduct["name"])
	}

	if updatedProduct["price"] == originalProduct["price"] {
		t.Errorf("Expected the price to change from '%v' to '%v'. Got '%v'", originalProduct["price"], updatedProduct["price"], updatedProduct["price"])
	}
}

func TestDeleteProduct(t *testing.T) {

	var jsonStr = []byte(`{"name":"test product", "price": 11.22}`)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	var m models.ProductModel
	err := json.Unmarshal(response.Body.Bytes(), &m)
	if err != nil {
		return
	}

	id := m.ID

	req, _ = http.NewRequest("GET", "/product/"+strconv.FormatInt(id, 10), nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/product/"+strconv.FormatInt(id, 10), nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/product/"+strconv.FormatInt(id, 10), nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
