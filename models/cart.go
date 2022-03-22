package models

type Cart struct {
	Id     uint `json:"id"`
	UserId uint `json:"user_id" gorm:"unique"`

	Items    []CartItem `json:"items" gorm:"foreignKey:CartId;references:Id"`
	Reserved bool        `json:"reserved"`
}

type CartItem struct {
	Id        uint  `json:"id"`
	CartId    uint  `json:"cartId"`
	ProductId uint  `json:"productId"`
	Quantity  uint32 `json:"quantity"`
}
