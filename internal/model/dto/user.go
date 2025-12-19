package dto

type UserDto struct {
	Login    string `json:"login"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Password string `json:"password"`
}

func NewUserDto() *UserDto {
	return &UserDto{}
}

func NewUserDtoWithParams(login string, name string, surname string) *UserDto {
	return &UserDto{Login: login, Name: name, Surname: surname}
}
