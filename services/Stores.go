package services

import (
	"encoding/json"
	"github.com/abhikeshri07/go-mux/models"
	"github.com/abhikeshri07/go-mux/utils"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type Stores struct {
	products IProducts
	conn     *gorm.DB
}

type IStores interface {
	GetProducts(w http.ResponseWriter, r *http.Request)
	AddProducts(w http.ResponseWriter, r *http.Request)
}

func (s *Stores) GetProducts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	start, _ := strconv.Atoi(r.FormValue("start"))
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Store ID")
		return
	}
	store := models.StoreModel{StoreId: int64(id)}
	products := store.GetProductsInStore(s.conn, limit, start)

	if products == nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Some Error Occurred")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, products)
}

func (s *Stores) AddProducts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Store ID")
		return
	}

	store := models.StoreModel{StoreId: int64(id)}
	var products []models.ProductModel
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&products); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	result := store.AddProducts(s.conn, products)
	if !result {
		utils.RespondWithError(w, http.StatusInternalServerError, "Some Error Occurred")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func NewStore(conn *gorm.DB) *Stores {
	return &Stores{conn: conn}
}
