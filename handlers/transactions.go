package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	dto "waysfood/dto/result"
	transactiondto "waysfood/dto/transaction"
	"waysfood/models"
	"waysfood/repositories"

	"github.com/golang-jwt/jwt/v4"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

var c = coreapi.Client{
	ServerKey: os.Getenv("SERVER_KEY"),
	ClientKey: os.Getenv("CLIENT_KEY"),
}

type handlerTransaction struct {
	TransactionRepository repositories.TransactionRepository
}

func HandlerTransaction(TransactionRepository repositories.TransactionRepository) *handlerTransaction {
	return &handlerTransaction{TransactionRepository}
}

func (h *handlerTransaction) FindTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))

	transactions, err := h.TransactionRepository.FindTransactions(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	}

	// var responseTransaction []models.Transaction
	// for _, t := range transactions {
	// 	responseTransaction = append(responseTransaction, convertResponseTransaction(t))
	// }

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: transactions}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) FindIncomes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))

	incomes, err := h.TransactionRepository.FindIncomes(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	}

	// var responseTransaction []models.Transaction
	// for _, t := range incomes {
	// 	responseTransaction = append(responseTransaction, convertResponseTransaction(t))
	// }

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: incomes}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	cartId := int(userInfo["id"].(float64))

	var request transactiondto.RequestTransaction
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	cart, err := h.TransactionRepository.GetCartTransaction(cartId, "pending")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Create Unique Transaction Id here ...
	time := time.Now()
	miliTime := time.Unix()

	transaction := models.Transaction{
		ID:       int(miliTime),
		CartID:   cart.ID,
		BuyerID:  cartId,
		SellerID: request.SellerID,
		Total:    request.Total,
		Status:   "pending",
	}

	log.Print(transaction)

	newTransaction, err := h.TransactionRepository.CreateTransaction(transaction)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	dataTransactions, err := h.TransactionRepository.GetTransaction(newTransaction.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// Request payment token from midtrans here ...
	// 1. Initiate Snap client
	var s = snap.Client{}
	s.New(os.Getenv("SERVER_KEY"), midtrans.Sandbox)
	// Use to midtrans.Production if you want Production Environment (accept real transaction).

	// 2. Initiate Snap request param
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(dataTransactions.ID),
			GrossAmt: int64(dataTransactions.Total),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: dataTransactions.Buyer.Fullname,
			Email: dataTransactions.Buyer.Email,
		},
	}

	// 3. Execute request create Snap transaction to Midtrans Snap API
	snapResp, _ := s.CreateTransaction(req)

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: snapResp}
	json.NewEncoder(w).Encode(response)
}

// Notification method ...
func (h *handlerTransaction) Notification(w http.ResponseWriter, r *http.Request) {
	var notificationPayload map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&notificationPayload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	transactionStatus := notificationPayload["transaction_status"].(string)
	fraudStatus := notificationPayload["fraud_status"].(string)
	orderId := notificationPayload["order_id"].(string)

	if transactionStatus == "capture" {
		if fraudStatus == "challenge" {
			// TODO set transaction status on your database to 'challenge'
			// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
			h.TransactionRepository.UpdateTransaction("pending", orderId)
		} else if fraudStatus == "accept" {
			// TODO set transaction status on your database to 'success'
			h.TransactionRepository.UpdateTransaction("success", orderId)
		}
	} else if transactionStatus == "settlement" {
		// TODO set transaction status on your databaase to 'success'
		h.TransactionRepository.UpdateTransaction("success", orderId)
	} else if transactionStatus == "deny" {
		// TODO you can ignore 'deny', because most of the time it allows payment retries
		// and later can become success
		h.TransactionRepository.UpdateTransaction("failed", orderId)
	} else if transactionStatus == "cancel" || transactionStatus == "expire" {
		// TODO set transaction status on your databaase to 'failure'
		h.TransactionRepository.UpdateTransaction("failed", orderId)
	} else if transactionStatus == "pending" {
		// TODO set transaction status on your databaase to 'pending' / waiting payment
		h.TransactionRepository.UpdateTransaction("pending", orderId)
	}

	w.WriteHeader(http.StatusOK)
}

// func SendMail(status string, transaction models.Transaction) {

// 	if status != transaction.Status {
// 		var CONFIG_SMTP_HOST = "smtp.gmail.com"
// 		var CONFIG_SMTP_PORT = 587
// 		var CONFIG_SENDER_NAME = "DumbMerch <demo.dumbways@gmail.com>"
// 		var CONFIG_AUTH_EMAIL = os.Getenv("EMAIL_SYSTEM")
// 		var CONFIG_AUTH_PASSWORD = os.Getenv("PASSWORD_SYSTEM")

// 		var productName = transaction.Cart
// 		var price = strconv.Itoa(transaction.Total)

// 		mailer := gomail.NewMessage()
// 		mailer.SetHeader("From", CONFIG_SENDER_NAME)
// 		mailer.SetHeader("To", transaction.Buyer.Email)
// 		mailer.SetHeader("Subject", "Transaction Status")
// 		mailer.SetBody("text/html", fmt.Sprintf(`<!DOCTYPE html>
//     <html lang="en">
//       <head>
//       <meta charset="UTF-8" />
//       <meta http-equiv="X-UA-Compatible" content="IE=edge" />
//       <meta name="viewport" content="width=device-width, initial-scale=1.0" />
//       <title>Document</title>
//       <style>
//         h1 {
//         color: brown;
//         }
//       </style>
//       </head>
//       <body>
//       <h2>Product payment :</h2>
//       <ul style="list-style-type:none;">
//         <li>Name : %s</li>
//         <li>Total payment: Rp.%s</li>
//         <li>Status : <b>%s</b></li>
// 		<li>Note : Dika Uye</li>
//       </ul>
//       </body>
//     </html>`, productName, price, status))

// 		dialer := gomail.NewDialer(
// 			CONFIG_SMTP_HOST,
// 			CONFIG_SMTP_PORT,
// 			CONFIG_AUTH_EMAIL,
// 			CONFIG_AUTH_PASSWORD,
// 		)

// 		err := dialer.DialAndSend(mailer)
// 		if err != nil {
// 			log.Fatal(err.Error())
// 		}

// 		log.Println("Mail sent! to " + transaction.Buyer.Email)
// 	}
// }

// func convertResponseTransaction(t models.Transaction) transactiondto.ResponseTransaction {
// 	return transactiondto.ResponseTransaction{
// 		ID:     t.ID,
// 		Cart:   append(t.Cart.Order[], ),
// 		Buyer:  t.Buyer.Fullname,
// 		Seller: t.Seller.Fullname,
// 		Total:  t.Total,
// 		Status: t.Status,
// 	}
// }
