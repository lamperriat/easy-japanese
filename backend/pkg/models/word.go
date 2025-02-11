package models

type ExampleSentence struct {
    Sentence    string `json:"example"`
    Chinese     string `json:"chinese"`
}

// Empty string if a field is not available
type JapaneseWord struct {
	ID       int               `json:"id"` // unique identifier
	Kanji    string            `json:"kanji"`
	Chinese  string            `json:"chinese"`
	Katakana string            `json:"katakana"` 
	Hiragana string            `json:"hiragana"`
	Type     string            `json:"type"` // specific to verb and adj
	Example  []ExampleSentence `json:"example"`
}

const (
	DefaultWeight = 50
	MinWeight	  = 1
	MaxWeight	  = 500
	ChangeRate    = 1
)

type UserWord struct {
	ID       int    `json:"id"`
	Weight   int    `json:"weight"` // 1 to 500
	UserNote string `json:"note"`
}