package dtos

type OTPCreateRequest struct {
	Phone string `json:"phone,omitempty" binding:"required,e164"`
}

type OTPVerifyRequest struct {
	//ID    *types.ID `json:"id,omitempty" binding:"required,uuid"`
	Phone string `json:"phone,omitempty" binding:"required,e164"`
	Code  string `json:"code,omitempty" binding:"required,len=4"`
}

type OTPVerifyResponse struct {
	Status string `json:"status,omitempty"`
	Token  string `json:"token,omitempty"`
}

type Login struct {
	Email    string `json:"email,omitempty" binding:"email"`
	Password string `json:"password,omitempty" binding:"required,min=3,max=16"`
}
