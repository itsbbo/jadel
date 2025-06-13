package gonertia

import (
	"html/template"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"
)

var rootTemplate = `<html>
<head>{{ .inertiaHead }}</head>
<body>{{ .inertia }}</body>
</html>`

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		i, err := New(rootTemplate)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if i.rootTemplateHTML != rootTemplate {
			t.Fatalf("root template html=%s, want=%s", i.rootTemplateHTML, rootTemplate)
		}
	})

	t.Run("blank", func(t *testing.T) {
		t.Parallel()

		_, err := New("")
		if err == nil {
			t.Fatal("error expected")
		}
	})
}

func TestNewFromFile(t *testing.T) {
	t.Parallel()

	f := tmpFile(t, rootTemplate)

	i, err := NewFromFile(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if i.rootTemplateHTML != rootTemplate {
		t.Fatalf("root template html=%s, want=%s", i.rootTemplateHTML, rootTemplate)
	}
}

func TestNewFromFileFS(t *testing.T) {
	t.Parallel()

	testFS := fstest.MapFS{
		"root.html": {
			Data: []byte(rootTemplate),
		},
	}

	i, err := NewFromFileFS(testFS, "root.html")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if i.rootTemplateHTML != rootTemplate {
		t.Fatalf("root template html=%s, want=%s", i.rootTemplateHTML, rootTemplate)
	}
}

func TestNewFromReader(t *testing.T) {
	t.Parallel()

	i, err := NewFromReader(strings.NewReader(rootTemplate))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if i.rootTemplateHTML != rootTemplate {
		t.Fatalf("root template html=%s, want=%s", i.rootTemplateHTML, rootTemplate)
	}
}

func TestNewFromBytes(t *testing.T) {
	t.Parallel()

	i, err := NewFromBytes([]byte(rootTemplate))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if i.rootTemplateHTML != rootTemplate {
		t.Fatalf("root template html=%s, want=%s", i.rootTemplateHTML, rootTemplate)
	}
}

func TestNewFromTemplate(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		tmpl := template.Must(template.New("foo").Parse(`<div id="app"></div>`))
		i, err := NewFromTemplate(tmpl)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if i.rootTemplate == nil {
			t.Fatalf("missing root template")
		}
	})

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		i, err := NewFromTemplate(nil)
		if err == nil {
			t.Fatalf("expected error for passing a nil template")
		}
		if i != nil {
			t.Fatalf("expected Inertia instance to be nil, but got %v", i)
		}
	})
}

func TestInertia_SharedProps(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		props Props
	}{
		{
			"empty",
			Props{},
		},
		{
			"with values",
			Props{"foo": "bar"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			i := I(func(i *Inertia) {
				i.sharedProps = tt.props
			})

			got := i.SharedProps()

			if !reflect.DeepEqual(got, i.sharedProps) {
				t.Fatalf("sharedProps=%#v, want=%#v", got, i.sharedProps)
			}
		})
	}
}

func TestInertia_SharedProp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		props  Props
		key    string
		want   any
		wantOk bool
	}{
		{
			"empty props",
			Props{},
			"foo",
			nil,
			false,
		},
		{
			"not found",
			Props{"foo": 123},
			"bar",
			nil,
			false,
		},
		{
			"found",
			Props{"foo": 123},
			"foo",
			123,
			true,
		},
		{
			"found nil value",
			Props{"foo": nil},
			"foo",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			i := I(func(i *Inertia) {
				i.sharedProps = tt.props
			})

			got, ok := i.SharedProp(tt.key)
			if ok != tt.wantOk {
				t.Fatalf("ok=%t, want=%t", ok, tt.wantOk)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("SharedProp()=%#v, want=%#v", got, tt.want)
			}
		})
	}
}
