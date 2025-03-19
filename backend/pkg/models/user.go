package models

type User struct {
    ID          uint       `json:"-" gorm:"primaryKey"`
    Username    string     `json:"username" gorm:"uniqueIndex"`
    Keyhash     string     `json:"-" gorm:"uniqueIndex"`
    Learned     []UserWord `json:"-" gorm:"foreignKey:user_id;references:ID"`
    Grammar     []UserGrammar `json:"-" gorm:"foreignKey:user_id;references:ID"`
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
    UserID      uint    `json:"-"` 
    User        User    `json:"-" gorm:"foreignKey:user_id;references:ID"`
    Familiarity int     `json:"familiarity" gorm:"default:50"`
    LastSeen    int64   `json:"lastSeen" gorm:"column:last_seen"`
    Examples    []UserWordExample `json:"example" gorm:"foreignKey:user_word_id"`
}

type UserGrammarExample struct {
    ID          uint   `json:"-" gorm:"primaryKey"`
    GrammarID   uint   `json:"-" gorm:"column:grammar_id"` // foreign key
    Example     string `json:"example" gorm:"type:text"`
    Chinese     string `json:"chinese" gorm:"type:text"`
}

type UserGrammar struct {
    ID          uint   `json:"id" gorm:"primaryKey"`
    UserID      uint   `json:"-"`
    User        User   `json:"-" gorm:"foreignKey:user_id;references:ID"`
    Description string `json:"description" gorm:"type:text"`
    Familiarity int     `json:"familiarity" gorm:"default:50"`
    LastSeen    int64   `json:"lastSeen" gorm:"column:last_seen"`
    Examples    []UserGrammarExample `json:"example" gorm:"foreignKey:grammar_id"`
}

type UserReadingMaterial struct {
    ID      uint   `json:"id" gorm:"primaryKey"`
    UserID  uint   `json:"-"`
    User    User   `json:"-" gorm:"foreignKey:user_id;references:ID"`
    Title   string `json:"title" gorm:"type:text"`
    Content string `json:"content" gorm:"type:text"`
    Chinese string `json:"chinese" gorm:"type:text"`
}

type Learnable interface {
    SetUserID(userID uint)
    GetUserID() uint
    GetID() uint
    GetFamiliarity() int
    GetLastSeen() int64
    GetName() string
}

func (u *UserWord) SetUserID(userID uint) { u.UserID = userID }
func (u *UserWord) GetUserID() uint { return u.UserID }
func (u *UserWord) GetID() uint { return u.ID }
func (u *UserWord) GetFamiliarity() int { return u.Familiarity }
func (u *UserWord) GetLastSeen() int64 { return u.LastSeen }
func (u *UserWord) GetName() string { return "word" }

func (u *UserGrammar) SetUserID(userID uint) { u.UserID = userID }
func (u *UserGrammar) GetUserID() uint { return u.UserID }
func (u *UserGrammar) GetID() uint { return u.ID }
func (u *UserGrammar) GetFamiliarity() int { return u.Familiarity }
func (u *UserGrammar) GetLastSeen() int64 { return u.LastSeen }
func (u *UserGrammar) GetName() string { return "grammar" }

