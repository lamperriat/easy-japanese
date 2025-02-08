package models

type User struct {
    ID       string     `json:"id"`
    Username string     `json:"username"`
    Learned  []UserWord `json:"learned"`
}

