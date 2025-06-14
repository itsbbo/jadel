package settings

import z "github.com/Oudwins/zog"

var changeProfileSchema = z.Struct(z.Shape{
	"name": z.String().
		Required(z.Message("Name is required")).
		Max(255, z.Message("Name max length is 255 characters")),
	"email": z.String().
		Required(z.Message("Email is required")).
		Max(255, z.Message("Email max length is 255 characters")).
		Email(z.Message("Email is not valid")),
})

var changePasswordSchema = z.Struct(z.Shape{
	"currentPassword": z.String().
		Required(z.Message("Current password is required")).
		Max(50, z.Message("Current password max length is 50 characters")),
	"password": z.String().
		Required(z.Message("New password is required")).
		Min(8, z.Message("New password min length is 8 characters")).
		Max(50, z.Message("New password max length is 50 characters")),
	"passwordConfirmation": z.String().
		Required(z.Message("Password confirmation is required")).
		Min(8, z.Message("Password confirmation min length is 8 characters")).
		Max(50, z.Message("Password confirmation max length is 50 characters")),
})

var destroyAccountSchema = z.Struct(z.Shape{
	"password": z.String().
		Required(z.Message("Password is required")).
		Min(8, z.Message("Password min length is 8 characters")).
		Max(50, z.Message("Password max length is 50 characters")),
})
