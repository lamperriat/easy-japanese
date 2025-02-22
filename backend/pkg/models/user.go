package models

type User struct {
    ID          uint       `json:"-" gorm:"primaryKey"`
    Username    string     `json:"username" gorm:"uniqueIndex"`
    Keyhash     string     `json:"-" gorm:"uniqueIndex"`
    Learned     []UserWord `json:"-" gorm:"foreignKey:UserID"`
    ReviewCount int64      `json:"reviewCount"`
    // the number of words that have been reviewed. 
    // One word can be counted multiple times. 
    // It will be assigned to `lastSeen` when the user reviews a word.
}

const (
	DefaultFamiliarity = 50
	MinFamiliarity	  = 1
	MaxFamiliarity	  = 500
	ChangeRate    = 1
)

type UserWordExample struct {
    ID          uint   `json:"-" gorm:"primaryKey"`
    UserWordID  uint   `json:"-" gorm:"column:user_word_id"` // foreign key
    Example     string `json:"example" gorm:"type:text"`
    Chinese     string `json:"chinese" gorm:"type:text"`
}


type UserWord struct {
    ID          uint    `json:"id" gorm:"primaryKey"`
    Kanji       string  `json:"kanji"`
    Chinese     string  `json:"chinese"`
    Katakana    string  `json:"katakana"`
    Hiragana    string  `json:"hiragana"`
    Type        string  `json:"type"`
    UserID      string  `json:"-" gorm:"index:user_id"`
    Familiarity int     `json:"familiarity" gorm:"default:50"`
    LastSeen    int64   `json:"lastSeen" gorm:"column:last_seen"`
    Examples    []UserWordExample `json:"example" gorm:"foreignKey:UserWordID"`
}