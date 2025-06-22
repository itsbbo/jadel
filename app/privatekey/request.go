package privatekey

import (
	z "github.com/Oudwins/zog"
)

var createPrivateKeySchema = z.Struct(z.Shape{
	"name": z.String().
		Required(z.Message("Name is required")).
		Max(255, z.Message("Name must be less than 255 characters")),
	"description": z.String().
		Optional().
		Max(255, z.Message("Description must be less than 255 characters")),
	"publicKey": z.String().
		Required(z.Message("Public key is required")),
	"privateKey": z.String().
		Required(z.Message("Private key is required")),
})
