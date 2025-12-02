package converters

import (
	"github.com/timurzdev/mentorship-test-task/internal/entity"
	"github.com/timurzdev/mentorship-test-task/internal/generated"
)

func UserFromGenRegister(reqGen generated.PostRegisterJSONBody) entity.User {
	return entity.User{
		Email:        string(*reqGen.Email),
		PasswordHash: *reqGen.Password,
		UserType:     string(*reqGen.UserType),
	}
}

func UserToGenRegister(user entity.User) generated.PostRegisterJSONBody {
	return generated.PostRegisterJSONBody{
		Email:    (*generated.Email)(&user.Email),
		Password: &user.PasswordHash,
		UserType: (*generated.UserType)(&user.UserType),
	}
}

func UserLoginFromGen(reqGen generated.PostLoginJSONBody) (userId string, password string) {
	if reqGen.Id != nil {
		userId = reqGen.Id.String()
	}
	if reqGen.Password != nil {
		password = *reqGen.Password
	}
	return userId, password
}
