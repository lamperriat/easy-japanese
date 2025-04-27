package models

type AdminAccount struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique;not null"`
	PasswordHash string `json:"passwordHash" gorm:"not null"`
}

type ApiKey struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	KeyHash  string `json:"keyHash" gorm:"unique;not null"`
}

