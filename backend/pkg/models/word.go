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

type ReadingMaterial struct {
    ID      uint   `json:"id" gorm:"primaryKey"`
    Title   string `json:"title" gorm:"type:text"`
    Content string `json:"content" gorm:"type:text"`
    Chinese string `json:"chinese" gorm:"type:text"`
}

type GrammarExample struct {
    ID        uint   `json:"-" gorm:"primaryKey"`
    GrammarID uint   `json:"-" gorm:"column:grammar_id"` // foreign key
    Example   string `json:"example" gorm:"type:text"`
    Chinese   string `json:"chinese" gorm:"type:text"`
}

type Grammar struct {
    ID          uint   `json:"id" gorm:"primaryKey"`
    Description string `json:"description" gorm:"type:text"`
    Examples    []GrammarExample `json:"example" gorm:"foreignKey:GrammarID"`
}

type SearchResult[T any] struct {
    Count    int64 `json:"count"`
    Page     int   `json:"page"`
    PageSize int   `json:"pageSize"`
    Results  []T   `json:"results"`
}
