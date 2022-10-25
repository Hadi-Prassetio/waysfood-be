package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
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

	Carts, err := h.CartRepository.FindCarts()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: 200, Data: Carts}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerCart) GetCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	transId := int(userInfo["id"].(float64))

	// id, _ := strconv.Atoi(mux.Vars(r)["id"])
	Cart, err := h.CartRepository.GetCart(transId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: 200, Data: Cart}
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

	var TransIdIsMatch = false
	var CartId int
	for !TransIdIsMatch {
		CartId = idUser + rand.Intn(10000) - rand.Intn(100)
		CartData, _ := h.CartRepository.GetCart(CartId)
		if CartData.ID == 0 {
			TransIdIsMatch = true
		}
	}

	Cart := models.Cart{
		ID:     CartId,
		UserID: idUser,
		Status: "active",
	}

	statusCheck, _ := h.CartRepository.FindbyIDCart(idUser, "active")
	if statusCheck.Status == "active" {
		response := dto.SuccessResult{Code: http.StatusOK, Data: Cart}
		json.NewEncoder(w).Encode(response)
	} else {
		data, _ := h.CartRepository.CreateCart(Cart)
		w.WriteHeader(http.StatusOK)
		response := dto.SuccessResult{Code: 200, Data: data}
		json.NewEncoder(w).Encode(response)
	}
}

func (h handlerCart) DeleteCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contetmt-type", "application/json")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	Cart, err := h.CartRepository.GetCart(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	data, err := h.CartRepository.DeleteCart(Cart)
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

	Cart, err := h.CartRepository.FindbyIDCart(idTrans, "active")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	if request.UserID != 0 {
		Cart.UserID = request.UserID
	}

	if request.SubTotal != 0 {
		Cart.SubTotal = request.SubTotal
	}

	if request.Status != "active" {
		Cart.Status = request.Status
	}

	_, err = h.CartRepository.UpdateCart(Cart)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	// dataCarts, err := h.CartRepository.FindbyIDCart(idTrans, request.Status)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	json.NewEncoder(w).Encode(err.Error())
	// 	return
}

// // // 1. Initiate Snap client
// var s = snap.Client{}
// s.New(os.Getenv("SERVER_KEY"), midtrans.Sandbox)
// // // Use to midtrans.Production if you want Production Environment (accept real Cart).

// // // 2. Initiate Snap request param
// req := &snap.Request{
// 	CartDetails: midtrans.CartDetails{
// 		OrderID:  strconv.Itoa(dataCarts.ID),
// 		GrossAmt: int64(dataCarts.Total),
// 	},
// 	CreditCard: &snap.CreditCardDetails{
// 		Secure: true,
// 	},
// 	CustomerDetail: &midtrans.CustomerDetails{
// 		FName: dataCarts.User.Fullname,
// 		Email: dataCarts.User.Email,
// 	},
// }

// // // 3. Execute request create Snap Cart to Midtrans Snap API
// snapResp, _ := s.CreateCart(req)

// w.WriteHeader(http.StatusOK)
// response := dto.SuccessResult{Code: 200, Data: snapResp}
// json.NewEncoder(w).Encode(response)
// }

func (h *handlerCart) FindbyIDCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))
	// id, _ := strconv.Atoi(mux.Vars(r)["id"])
	Cart, err := h.CartRepository.FindbyIDCart(userId, "active")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: 200, Data: Cart}
	json.NewEncoder(w).Encode(response)
}

// func (h *handlerCart) Notification(w http.ResponseWriter, r *http.Request) {
// 	var notificationPayload map[string]interface{}

// 	err := json.NewDecoder(r.Body).Decode(&notificationPayload)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
// 		json.NewEncoder(w).Encode(response)
// 		return
// 	}

// 	CartStatus := notificationPayload["Cart_status"].(string)
// 	fraudStatus := notificationPayload["fraud_status"].(string)
// 	orderId := notificationPayload["order_id"].(string)

// 	if CartStatus == "capture" {
// 		if fraudStatus == "challenge" {
// 			// TODO set Cart status on your database to 'challenge'
// 			// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
// 			h.CartRepository.UpdateCart("pending", orderId)
// 		} else if fraudStatus == "accept" {
// 			// TODO set Cart status on your database to 'success'
// 			h.CartRepository.UpdateCart("success", orderId)
// 		}
// 	} else if CartStatus == "settlement" {
// 		// TODO set Cart status on your databaase to 'success'
// 		h.CartRepository.UpdateCart("success", orderId)
// 	} else if CartStatus == "deny" {
// 		// TODO you can ignore 'deny', because most of the time it allows payment retries
// 		// and later can become success
// 		h.CartRepository.UpdateCart("failed", orderId)
// 	} else if CartStatus == "cancel" || CartStatus == "expire" {
// 		// TODO set Cart status on your databaase to 'failure'
// 		h.CartRepository.UpdateCart("failed", orderId)
// 	} else if CartStatus == "pending" {
// 		// TODO set Cart status on your databaase to 'pending' / waiting payment
// 		h.CartRepository.UpdateCart("pending", orderId)
// 	}

// 	w.WriteHeader(http.StatusOK)
// }

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
