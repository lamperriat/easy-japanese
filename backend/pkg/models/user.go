package models

type User struct {
    ID       string     `gorm:"primaryKey"`
    Username string     `gorm:"uniqueIndex"`
    Learned  []UserWord `gorm:"foreignKey:UserID"`
}
