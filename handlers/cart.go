package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	cartdto "waysfood/dto/cart"
	dto "waysfood/dto/result"
	"waysfood/models"
	"waysfood/repositories"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type handlerCart struct {
	CartRepository repositories.CartRepository
}

func HandlerCart(CartRepository repositories.CartRepository) *handlerCart {
	return &handlerCart{CartRepository}
}

func (h *handlerCart) FindCarts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	carts, err := h.CartRepository.FindCarts()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: 200, Data: carts}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerCart) GetCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	transId := int(userInfo["id"].(float64))

	// id, _ := strconv.Atoi(mux.Vars(r)["id"])
	cart, err := h.CartRepository.GetCart(transId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: 200, Data: cart}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerCart) CreateCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	idUser := int(userInfo["id"].(float64))

	request := new(cartdto.CreateCart)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	validate := validator.New()
	err := validate.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	time := time.Now()
	miliTime := time.Unix()

	cart := models.Cart{
		ID:     int(miliTime),
		UserID: idUser,
		Status: "pending",
	}

	statusCheck, _ := h.CartRepository.FindbyIDCart(idUser, "pending")
	if statusCheck.Status == "pending" {
		response := dto.SuccessResult{Code: http.StatusOK, Data: cart}
		json.NewEncoder(w).Encode(response)
	} else {
		data, _ := h.CartRepository.CreateCart(cart)
		w.WriteHeader(http.StatusOK)
		response := dto.SuccessResult{Code: 200, Data: data}
		json.NewEncoder(w).Encode(response)
	}
}

func (h handlerCart) DeleteCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contetmt-type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	cart, err := h.CartRepository.GetCart(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	data, err := h.CartRepository.DeleteCart(cart)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: 200, Data: data}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerCart) UpdateCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	idTrans := int(userInfo["id"].(float64))

	request := new(cartdto.UpdateCart)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	cart, err := h.CartRepository.FindbyIDCart(idTrans, "pending")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	if request.UserID != 0 {
		cart.UserID = request.UserID
	}

	if request.SubTotal != 0 {
		cart.SubTotal = request.SubTotal
	}

	if request.Status != "pending" {
		cart.Status = request.Status
	}

	_, err = h.CartRepository.UpdateCart(cart)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}
}

func (h *handlerCart) FindbyIDCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))
	// id, _ := strconv.Atoi(mux.Vars(r)["id"])
	cart, err := h.CartRepository.FindbyIDCart(userId, "pending")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: 200, Data: cart}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerCart) AllProductById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))

	Carts, err := h.CartRepository.AllProductById(userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: 200, Data: Carts}
	json.NewEncoder(w).Encode(response)
}
