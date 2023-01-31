package services

import (
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

// AddProducts todo
func (s *Stores) AddProducts(w http.ResponseWriter, r *http.Request) {

}
func NewStore(conn *gorm.DB) *Stores {
	return &Stores{conn: conn}
}
