package auth

type SessionCreateInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
