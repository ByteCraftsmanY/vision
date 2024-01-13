package dtos

type UserForm struct {
	Name     string `json:"name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required,e164"`
	Password string `json:"password" binding:"required,min=8,max=16"`
}

type UserDTO struct {
}

//func (u *UserDTO) ConvertToEntity(entity *entities.User) {
//
//}
