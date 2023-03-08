package social

type User struct {
	ID    interface{} `json:"id"`
	Name  string      `json:"name"`
	Key   string      `json:"key"`
	Token string      `json:"token"`
	Email string      `json:"email"`
}
