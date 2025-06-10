package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Oudwins/zog"
	"github.com/labstack/echo/v4"
	"github.com/romsar/gonertia/v2"
)

func TestSetInertiaValidationError(t *testing.T) {
	tests := []struct {
		name     string
		errInput zog.ZogIssueMap
		want     gonertia.ValidationErrors
	}{
		{
			name: "single error",
			errInput: zog.ZogIssueMap{
				"email": zog.ZogIssueList{
					{Message: "Invalid email format"},
				},
			},
			want: gonertia.ValidationErrors{
				"email": "Invalid email format",
			},
		},
		{
			name: "multiple fields with errors",
			errInput: zog.ZogIssueMap{
				"email": zog.ZogIssueList{
					{Message: "Invalid email format"},
				},
				"password": zog.ZogIssueList{
					{Message: "Password too short"},
				},
			},
			want: gonertia.ValidationErrors{
				"email":    "Invalid email format",
				"password": "Password too short",
			},
		},
		{
			name:     "empty error map",
			errInput: zog.ZogIssueMap{},
			want:     gonertia.ValidationErrors{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			SetInertiaValidationErrorsZog(c, tt.errInput)

			ctx := c.Request().Context()

			got := gonertia.ValidationErrorsFromContext(ctx)

			if len(got) != len(tt.want) {
				t.Errorf("Expected %d validation errors length, got %d", len(tt.want), len(got))
			}
		})
	}
}
