package structs

type PublicUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type Message struct {
	Id     int        `json:"id"`
	UserId int        `json:"userId"`
	Value  string     `json:"value"`
	User   PublicUser `json:"user"`
}
