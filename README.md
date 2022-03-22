# Simple Market API

## Auth

The auth API consists of 3 endpoints:
* /register
* /login
* /logout

### Register 
I implemented the register endpoint because I want to have a hash in users table, not the actual password.

example of json to send:
```json
{
  "name": "vlad",
  "email": "vlad@gmail.com",
  "password": "password"
}
```

### Login
When a user wants to log in, the credentials are checked and if everything is alright, a token is returned on a cookie.

example of json to send:
```json
{
  "email": "vlad@gmail.com",
  "password": "password"
}
```

### Logout
Simply sets the expiration date somewhere in the past 

## Market
The market API consists of 2 endpoints:
* /add
* /checkout

### Add
A new product is added in cart if the desired quantity is in stock.

example of json to send:
```json
{
  "productId": "1",
  "quantity": 7
}
```

### Checkout
It will try to checkout. It checks if all the products in cart are still in stock. If they are, the products are removed 
from stock and it will continue with payment. There is a chance of 50% for payment to succeed
If the payment succeeded, all the products in the cart are removed, if not, the stocks are restored.

## Database:
There are 4 tables required:
* users
* products
* carts
* cart_items

Users table is independent of the other 3.
The definition of the 4 database can be found under models.


## Improvements to be done:
* add validations fore received objects
* add tests

## Ideal implementation (future):
* The auth API and the market API are two different microservices (possible because of JWT)
* The two microservices run in a kubernetes cluster. Each of these microservices has its own database
    