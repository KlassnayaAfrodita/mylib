package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/KlassnayaAfrodita/mylib/storage"

	"github.com/gorilla/mux"
)

type Api struct {
	session        *storage.Session
	users          *storage.UserStorage
	productStorage *storage.ProductStore
}

func NewApi() *Api {
	return &Api{
		session:        storage.NewSession(),
		users:          storage.NewUserStorage(),
		productStorage: storage.NewProductStore(),
	}
}

func (api *Api) GetAllCart(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("hello"))
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, `{"error": "you dont auth"}`, 400)
		return
	}
	userId, err := api.session.GetSession(cookie.Value)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		return
	}
	user, err := api.users.GetUser(userId)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		return
	}

	products, err := user.Cart.GetProducts()
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	resp, err := json.Marshal(products)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(resp)
}

func (api *Api) GetProductCart(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, `{"error": "you dont auth"}`, 400)
		return
	}
	userId, err := api.session.GetSession(cookie.Value)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		return
	}
	user, err := api.users.GetUser(userId)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["product_id"])
	if err != nil {
		http.Error(w, `{"error": "bad id"}`, 400)
		return
	}

	product, err := user.Cart.GetProduct(id)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		return
	}

	resp, err := json.Marshal(product)
	if err != nil {
		http.Error(w, `"error": "server error"`, 500)
		return
	}

	w.Write(resp)
}

func (api *Api) AddProductCart(w http.ResponseWriter, r *http.Request) { //! ожидается json с новым пользователем

	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, `{"error": "you dont auth"}`, 400)
		return
	}
	userId, err := api.session.GetSession(cookie.Value)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		return
	}
	user, err := api.users.GetUser(userId)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var product storage.Product

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error": "server error"}`, 500)
	}
	defer r.Body.Close()

	newerr := json.Unmarshal(body, &product)
	if newerr != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = user.Cart.AddProduct(product)
	if err != nil {
		http.Error(w, `"error":"db error"`, 500)
	}

	//! меняем базу всех продуктов
	product.Quantity--
	storage.ProductStore.ChangeProduct(product)

	w.Write([]byte(body))
}

func (api *Api) ChangeProductCart(w http.ResponseWriter, r *http.Request) { //! получаем json полного пользователя

	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, `{"error": "you dont auth"}`, 400)
		return
	}
	userId, err := api.session.GetSession(cookie.Value)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		return
	}
	user, err := api.users.GetUser(userId)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		return
	}

	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var product storage.Product

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error": "server error"}`, 500)
	}
	defer r.Body.Close()

	newerr := json.Unmarshal(body, &product)
	if newerr != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = user.Cart.ChangeProduct(product)
	if err != nil {
		http.Error(w, `"error":"db error"`, 500)
		return
	}
	w.Write(body)
}

func (api *Api) DeleteProductCart(w http.ResponseWriter, r *http.Request) { //! получаем только id

	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, `{"error": "you dont auth"}`, 400)
		return
	}
	userId, err := api.session.GetSession(cookie.Value)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		return
	}
	user, err := api.users.GetUser(userId)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["product_id"])
	if err != nil {
		http.Error(w, `{"error":"bad id"}`, 400)
		return
	}

	product, err := user.Cart.GetProduct(id)
	if err != nil {
		http.Error(w, `"error":"db error"`, 500)
		return
	}

	product, err = user.Cart.DeleteProduct(product)
	if err != nil {
		http.Error(w, `"error":"db error"`, 500)
		return
	}

	resp, err := json.Marshal(product)
	if err != nil {
		http.Error(w, `"error":"server error"`, 500)
	}

	w.Write(resp)
}
