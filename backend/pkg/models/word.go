package models

type ExampleSentence struct {
    ID             uint   `json:"-" gorm:"primaryKey"`
    JapaneseWordID uint   `json:"-" gorm:"column:japanese_word_id"` // foreign key
    Sentence       string `json:"example" gorm:"type:text"`
    Chinese        string `json:"chinese" gorm:"type:text"`
}

type JapaneseWord struct {
    ID        uint              `json:"id" gorm:"primaryKey;<-:create"`
    DictName  string            `json:"-" gorm:"index:idx_dict"`
    Kanji     string            `json:"kanji" gorm:"index:idx_dict,priority:1"`
    Chinese   string            `json:"chinese"`
    Katakana  string            `json:"katakana" gorm:"index:idx_dict,priority:2"` 
    Hiragana  string            `json:"hiragana"`
    Type      string            `json:"type"`
    Examples  []ExampleSentence `json:"example" gorm:"foreignKey:JapaneseWordID"`
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