package routes

import (
	"net/http"
	"simple-api/controllers"
)

func Setup() {
	http.HandleFunc("/register", controllers.Register)
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/logout", controllers.Logout)
	http.HandleFunc("/add", controllers.AddToCart)
	http.HandleFunc("/checkout", controllers.Checkout)
}
