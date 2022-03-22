package models

type Product struct {
	Id    uint   `json:"id" gorm:"foreignKey:ProductId;references:Id"`
	Name  string `json:"name" gorm:"unique"`
	Price uint32 `json:"price"`
	Stock uint32  `json:"stock"`
}

func (p Product) IsInStock() bool {
	return p.Stock > 0
}
