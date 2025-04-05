package auth

type SessionCreateInput struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}
