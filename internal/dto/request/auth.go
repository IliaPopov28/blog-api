package request

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50,alphanum"`
	Password string `json:"password" binding:"required,min=8,max=100"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
