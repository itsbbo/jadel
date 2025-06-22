package projects

import z "github.com/Oudwins/zog"

var createProjectSchema = z.Struct(z.Shape{
	"Name": z.String().
		Required(z.Message("Name is required")).
		Max(255, z.Message("Name max length is 255 characters")),
	"Description": z.String().
		Optional().
		Max(255, z.Message("Description max length is 255 characters")),
})
