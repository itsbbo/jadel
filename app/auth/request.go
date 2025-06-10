package auth

import z "github.com/Oudwins/zog"

var loginSchema = z.Struct(z.Shape{
	"email": z.String().
		Required(z.Message("Email is required")).
		Max(255, z.Message("Email max length is 255 characters")).
		Email(z.Message("Email is not valid")),
	"password": z.String().
		Required(z.Message("Password is required")).
		Max(255, z.Message("Password max length is 255 characters")),
})

var registerSchema = z.Struct(z.Shape{
	"name": z.String().
		Required(z.Message("Name is required")).
		Max(255, z.Message("Name max length is 255 characters")),
	"email": z.String().
		Required(z.Message("Email is required")).
		Max(255, z.Message("Email max length is 255 characters")).
		Email(z.Message("Email is not valid")),
	"password": z.String().
		Required(z.Message("Password is required")).
		Min(8, z.Message("Password min length is 8 characters")).
		Max(255, z.Message("Password max length is 255 characters")),
	"passwordConfirmation": z.String().
		Required(z.Message("Password confirmation is required")).
		Min(8, z.Message("Password confirmation min length is 8 characters")).
		Max(255, z.Message("Password confirmation max length is 255 characters")),
})
