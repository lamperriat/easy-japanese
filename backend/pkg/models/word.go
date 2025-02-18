package models

type ExampleSentence struct {
    ID             uint   `gorm:"primaryKey"`
    JapaneseWordID uint    // foreign key
    Sentence       string `gorm:"type:text"`
    Chinese        string `gorm:"type:text"`
}

type JapaneseWord struct {
    ID        uint              `gorm:"primaryKey"`
    Kanji     string            `gorm:"index"`
    Chinese   string            `gorm:"index"`
    Katakana  string            `gorm:"type:varchar(255)"` 
    Hiragana  string            `gorm:"type:varchar(255)"`
    Type      string            `gorm:"type:varchar(255)"`
    Examples  []ExampleSentence `gorm:"foreignKey:JapaneseWordID"`
}

const (
	DefaultWeight = 50
	MinWeight	  = 1
	MaxWeight	  = 500
	ChangeRate    = 1
)

type UserWord struct {
    UserID    string `gorm:"primaryKey"`
    WordID    uint   `gorm:"primaryKey"`
    Weight    int    `gorm:"check:weight BETWEEN 1 AND 500"`
    UserNote  string `gorm:"type:text"`
}