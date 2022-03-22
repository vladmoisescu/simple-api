package controllers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"simple-api/database"
	"simple-api/models"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

func AddToCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// get token from cookie and check if it is valid
	c, err := r.Cookie(CookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	email := getEmail(c.Value)
	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// get user
	var user models.User
	database.DB.Where("email = ?", email).First(&user)
	if user.Id == 0 {
		w.WriteHeader(http.StatusInternalServerError)
	}

	var item models.CartItem
	err = json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// retrieve the product from database
	var product models.Product
	database.DB.Where("id = ?", item.ProductId).First(&product)
	if product.Id == 0 {
		http.Error(w, "product not found", http.StatusBadRequest)
		return
	}

	if !product.IsInStock() {
		http.Error(w, "product out of stock", http.StatusNotFound)
		return
	}

	var existingCart models.Cart
	database.DB.Where("user_id = ?", user.Id).First(&existingCart)
	// create new cart if current user has no cart
	if existingCart.Id == 0 {
		cart := models.Cart{
			UserId: user.Id,
			Items: []models.CartItem{{
				ProductId: product.Id,
				Quantity:  item.Quantity,
			}},
			Reserved: false,
		}
		database.DB.Create(&cart)
		return
	}

	// check if there is the same product already in cart. If it is, just increase the quantity.
	var cartItems []models.CartItem
	database.DB.Where("cart_id = ?", existingCart.Id).Find(&cartItems)
	for _, existingItem := range cartItems {
		if existingItem.ProductId == item.ProductId {
			existingItem.Quantity += item.Quantity
			database.DB.Updates(&existingItem)
			return
		}
	}

	// create new item in cart since this is a new product.
	existingCart.Items = append(existingCart.Items, models.CartItem{
		CartId:    existingCart.Id,
		ProductId: product.Id,
		Quantity: item.Quantity,
	})
	database.DB.Updates(existingCart)
}

func Checkout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// get token from cookie and check if it is valid
	c, err := r.Cookie(CookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	email := getEmail(c.Value)
	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// get user
	var user models.User
	database.DB.Where("email = ?", email).First(&user)
	if user.Id == 0 {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// get user cart
	var existingCart models.Cart
	database.DB.Where("user_id = ?", user.Id).First(&existingCart)
	if existingCart.Id == 0 {
		http.Error(w, "nothing to checkout", http.StatusNotFound)
		return
	}

	// get items from the cart
	var cartItems []models.CartItem
	database.DB.Where("cart_id = ?", existingCart.Id).Find(&cartItems)
	if len(cartItems) == 0 {
		http.Error(w, "nothing to checkout", http.StatusNotFound)
		return
	}

	// remove products from stock and mark the existing cart as reserved
	var products []models.Product
	var product models.Product
	for _, item := range cartItems {
		database.DB.Where("id = ?", item.ProductId).First(&product)
		if product.Stock < item.Quantity {
			http.Error(w, fmt.Sprintf("not enough products: %s in stock", product.Name), http.StatusBadRequest)
			return
		}
		product.Stock -= item.Quantity
		products = append(products, product)
	}
	existingCart.Reserved = true

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Updates(existingCart).Error; err != nil {
			return err
		}

		for _, product = range products {
			p, err := structToMap(product)
			if err != nil {
				return err
			}
			if err = tx.Model(models.Product{}).Where("name = ?", product.Name).Updates(&p).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// simulate payment
	success := sendPaymentRequest()
	// if success, all the items from cart are deleted
	if success {
		for _, item := range cartItems {
			if err = database.DB.Delete(item).Error; err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
		// else, stocks are restored
	} else {
		err = database.DB.Transaction(func(tx *gorm.DB) error {
			if err = tx.Updates(existingCart).Error; err != nil {
				return err
			}

			for _, item := range cartItems {
				for ii, p := range products {
					if item.ProductId == p.Id {
						products[ii].Stock += item.Quantity
						break
					}
				}
			}
			for _, product = range products {
				if err = tx.Updates(&product).Error; err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusInternalServerError)
	}
	notReserved := markCartAsNotReserved()
	if err = database.DB.Model(models.Cart{}).Where("user_id = ?", existingCart.UserId).Updates(&notReserved).Error; err != nil {
		fmt.Println(err)
	}
}

func getEmail(token string) (email string) {
	// Initialize a new instance of `Claims`
	claims := &jwt.StandardClaims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return
		}
		return
	}
	if !tkn.Valid {
		err = fmt.Errorf("invalid token")
		return
	}
	return claims.Issuer
}

// sendPaymentRequest is a dummy method that waits for 2 seconds and has a 50% chance of success
func sendPaymentRequest() (success bool) {
	randNum := rand.Int()
	fmt.Println(randNum)

	time.Sleep(2 * time.Second)

	if randNum%2 == 0 {
		success = true
		return
	}
	success = false
	return
}

func structToMap(obj interface{}) (objectToUpdate map[string]interface{}, err error) {
	data, err := json.Marshal(obj)

	if err != nil {
		return
	}

	err = json.Unmarshal(data, &objectToUpdate) // Convert to a map
	return
}

func markCartAsNotReserved() (cartData map[string]interface{}) {
	cartData = map[string]interface{}{
		"reserved": false,
	}
	return
}
