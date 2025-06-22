package servers

import (
	"net"

	z "github.com/Oudwins/zog"
	"github.com/oklog/ulid/v2"
)

var createServerSchema = z.Struct(z.Shape{
	"name": z.String().
		Required(z.Message("Name is required")).
		Max(255, z.Message("Name max length is 255 characters")),
	"description": z.String().
		Optional().
		Max(255, z.Message("Description max length is 255 characters")),
	"ip": z.String().
		Required(z.Message("IP is required")).
		TestFunc(
			func(val *string, ctx z.Ctx) bool { return net.ParseIP(*val) != nil },
			z.Message("IP is not valid"),
		),
	"port": z.Int().Required(z.Message("Port is required")),
	"user": z.String().Required(z.Message("User is required")),
	"private_key_id": z.String().
		Required(z.Message("Private key ID is required")).
		TestFunc(
			func(val *string, ctx z.Ctx) bool {
				_, err := ulid.Parse(*val)
				return err == nil
			},
			z.Message("Unknown private key"),
		),
})
