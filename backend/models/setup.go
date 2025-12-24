package models

import (
	"time"
)

// 1. Users
// ID, Username, Password, Role
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"` 
	Password  string    `gorm:"not null" json:"-"`               
	Role      string    `gorm:"type:varchar(10);default:'staff'" json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 2. Suppliers
// ID, Name, Email, Address
type Supplier struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"unique" json:"email"`
	Address   string    `gorm:"type:text" json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 3. Items
// ID, Name, Stock, Price
type Item struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Stock     int       `gorm:"not null;check:stock >= 0" json:"stock"` 
	Price     int64     `gorm:"not null" json:"price"`                 
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 4. Purchasings (Header Transaksi)
// ID, Date, SupplierID, UserID, Grand Total
type Purchasing struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Date       time.Time `json:"date"`
	GrandTotal int64     `json:"grand_total"`

	// Foreign Keys
	SupplierID uint `json:"supplier_id"`
	UserID     uint `json:"user_id"`

	Supplier Supplier `gorm:"foreignKey:SupplierID" json:"supplier"`
	User     User     `gorm:"foreignKey:UserID" json:"user"`

	Details []PurchasingDetail `gorm:"foreignKey:PurchasingID" json:"details"`
}

// 5. Purchasing Details (Detail Barang per Transaksi)
// ID, PurchasingID, ItemID, Qty, SubTotal
type PurchasingDetail struct {
	ID           uint  `gorm:"primaryKey" json:"id"`
	PurchasingID uint  `json:"purchasing_id"`
	ItemID       uint  `json:"item_id"`
	Qty          int   `json:"qty"`
	SubTotal     int64 `json:"sub_total"`

	// Relasi
	Item Item `gorm:"foreignKey:ItemID" json:"item"`
}