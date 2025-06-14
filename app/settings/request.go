package settings

import z "github.com/Oudwins/zog"

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