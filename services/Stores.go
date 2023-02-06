package services

import (
	"encoding/json"
	"github.com/abhikeshri07/go-mux/constants"
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
	limit, limitError := strconv.Atoi(r.FormValue("limit"))
	start, startError := strconv.Atoi(r.FormValue("start"))
	id, err := strconv.Atoi(vars["id"])
	if limitError != nil && r.FormValue("limit") != "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Format of limit. Send an integer as the limit value")
		return

	}
	if startError != nil && r.FormValue("start") != "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Format of start. Send an integer as the start value")
		return
	}
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid format of Store ID. Send an integer as the Store ID")
		return
	}

	store := models.StoreModel{StoreId: int64(id)}
	checkStorePresent := store.CheckStoreId(s.conn)

	if checkStorePresent == constants.STORE_NOT_FOUND_ERROR {
		utils.RespondWithError(w, http.StatusBadRequest, constants.STORE_NOT_FOUND_ERROR)
		return
	}

	products := store.GetProductsInStore(s.conn, limit, start)

	if products == nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "No products found at the specified store")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, products)
}

func (s *Stores) AddProducts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid format of Store ID. Send an integer as the Store ID")
		return
	}

	store := models.StoreModel{StoreId: int64(id)}
	var productsIds []int64
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&productsIds); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	result := store.AddProducts(s.conn, productsIds)
	if result == constants.STORE_PRODUCT_ENTRY_FOREIGN_KEY_ERROR {
		utils.RespondWithError(w, http.StatusBadRequest, "The productIds are not present in product table")
		return
	}
	if result == constants.DB_TRANSACTION_ERROR {
		utils.RespondWithError(w, http.StatusInternalServerError, constants.DB_TRANSACTION_ERROR)
		return
	}
	if result == constants.STORE_PRODCUT_ENTRY_ERROR {
		utils.RespondWithError(w, http.StatusInternalServerError, constants.STORE_PRODCUT_ENTRY_ERROR)
		return
	}
	if result == constants.STORE_PRODCUT_ENTRY_SUCCESS {
		utils.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
		return
	}
}

func NewStore(conn *gorm.DB) *Stores {
	return &Stores{conn: conn}
}
