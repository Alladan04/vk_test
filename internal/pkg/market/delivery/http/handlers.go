package http

import (
	"net/http"
	"strconv"

	"github.com/Alladan04/vk_test/internal/models"
	"github.com/Alladan04/vk_test/internal/pkg/market"
	"github.com/Alladan04/vk_test/internal/pkg/utils"
	"github.com/satori/uuid"
)

type MarketHandler struct {
	uc market.MarketUsecase
}

func NewMarketHandler(uc market.MarketUsecase) *MarketHandler {
	return &MarketHandler{
		uc: uc,
	}
}

func (h *MarketHandler) AddItem(w http.ResponseWriter, r *http.Request) {

	jwtPayload, ok := r.Context().Value(models.PayloadContextKey).(models.JwtPayload)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	item := models.ItemForm{}
	err := utils.GetRequestData(r, &item)
	if err != nil {
		utils.WriteErrorMessage(w, http.StatusBadRequest, "error unmarshalling")
		return
	}

	res, err := h.uc.AddItem(r.Context(), item, jwtPayload.Id)
	if err != nil {
		utils.WriteErrorMessage(w, http.StatusBadRequest, "error Adding data")
		return

	}

	utils.WriteResponseData(w, res, http.StatusOK)

}

// TODO: move all the logic from handler to usecase.
func (h *MarketHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	sortingOptions := []string{"price", "date"}

	//resolve sorting direction (ASC by default)
	sortOrder := r.URL.Query().Get("sort")
	if sortOrder == "up" {
		sortOrder = "ASC"
	} else if sortOrder == "down" || sortOrder == "" {
		sortOrder = "DESC"
	} else {
		utils.WriteErrorMessage(w, http.StatusBadRequest, "wrong sort parameter")
		return
	}
	//check if sorting option is a valid value
	sortField := r.URL.Query().Get("field")
	flag := false
	for _, i := range sortingOptions {
		if sortField == i {
			flag = true
		}
	}
	if !flag && sortField != "" {
		utils.WriteErrorMessage(w, http.StatusBadRequest, "wrong field parameter")
		return
	}
	if sortField != "price" {
		sortField = "create_time"
	}

	//process offset and count params
	countParam := r.URL.Query().Get("count")
	offsetParam := r.URL.Query().Get("offset")
	count, err := strconv.ParseInt(countParam, 10, 64)
	if err != nil && countParam != "" {
		utils.WriteErrorMessage(w, http.StatusBadRequest, "wrong count param")
		return
	}
	offset, err := strconv.ParseInt(offsetParam, 10, 64)
	if err != nil && offsetParam != "" {
		utils.WriteErrorMessage(w, http.StatusBadRequest, "wrong offset param")
		return
	}
	if countParam == "" {
		count = 10
	}
	if offsetParam == "" {
		offset = 0
	}

	//check min and max  price params in usecase
	minParam := r.URL.Query().Get("min")
	maxParam := r.URL.Query().Get("max")

	//get userId from context if user is authorized or set userId to zero value if not authorized
	userId := uuid.UUID{}
	jwtPayload, ok := r.Context().Value(models.PayloadContextKey).(models.JwtPayload)
	if ok {
		userId = jwtPayload.Id
	}

	//get result from usecase
	result, err := h.uc.GetAll(r.Context(), count, offset, sortOrder, sortField, minParam, maxParam, userId)
	if err != nil {
		utils.WriteErrorMessage(w, http.StatusBadRequest, err.Error())
		return

	}
	utils.WriteResponseData(w, result, http.StatusOK)

}
