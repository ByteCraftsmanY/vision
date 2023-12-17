package forms

type UserSignUp struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Country  string `json:"country" binding:"required"`
	Password string `json:"password" binding:"required"`
}
