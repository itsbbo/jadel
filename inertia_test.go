package main

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProjectRoot(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "project-root-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	projectDir := filepath.Join(tmpDir, "myproject")
	subDir := filepath.Join(projectDir, "cmd", "app")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create test directories: %v", err)
	}

	goModPath := filepath.Join(projectDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte("module myproject"), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	tests := []struct {
		name     string
		startDir string
		want     string
	}{
		{
			name:     "From project root",
			startDir: projectDir,
			want:     projectDir,
		},
		{
			name:     "From subdirectory",
			startDir: subDir,
			want:     projectDir,
		},
		{
			name:     "From outside project",
			startDir: tmpDir,
			want:     "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			originalDir, _ := os.Getwd()
			if err := os.Chdir(test.startDir); err != nil {
				t.Fatalf("Failed to change directory: %v", err)
			}
			defer os.Chdir(originalDir)

			got := ProjectRoot()
			if got != test.want {
				t.Errorf("ProjectRoot() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestViteHotFileUrl(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vite-hot-file-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name    string
		content string
		want    string
		wantErr bool
	}{
		{
			name:    "File does not exist",
			content: "",
			want:    "",
			wantErr: false,
		},
		{
			name:    "Empty file",
			content: "",
			want:    "//localhost:1323",
			wantErr: false,
		},
		{
			name:    "HTTP URL",
			content: "http://localhost:3000",
			want:    "//localhost:3000",
			wantErr: false,
		},
		{
			name:    "HTTPS URL",
			content: "https://localhost:3000",
			want:    "//localhost:3000",
			wantErr: false,
		},
		{
			name:    "Non-URL content",
			content: "localhost:3000",
			want:    "//localhost:1323",
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, "hot")

			if test.name != "File does not exist" {
				err := os.WriteFile(testFile, []byte(test.content), 0644)
				if err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}
			}

			got, err := viteHotFileUrl(testFile)
			if (err != nil) != test.wantErr {
				t.Errorf("viteHotFileUrld() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("viteHotFileUrld() = %v, want %v", got, test.want)
			}

			if test.name != "File does not exist" {
				os.Remove(testFile)
			}
		})
	}
}

func TestViteReactRefresh(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    string
		wantErr bool
	}{
		{
			name:    "Empty URL",
			url:     "",
			want:    "",
			wantErr: false,
		},
		{
			name: "Valid URL",
			url:  "http://localhost:3000",
			want: `
<script type="module">
    import RefreshRuntime from 'http://localhost:3000/@react-refresh'
    RefreshRuntime.injectIntoGlobalHook(window)
    window.$RefreshReg$ = () => {}
    window.$RefreshSig$ = () => (type) => type
    window.__vite_plugin_react_preamble_installed__ = true
</script>`,
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			refreshFunc := viteReactRefresh(test.url)
			got, err := refreshFunc()

			if (err != nil) != test.wantErr {
				t.Errorf("viteReactRefresh() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			gotStr := string(got)
			if strings.TrimSpace(gotStr) != strings.TrimSpace(test.want) {
				t.Errorf("viteReactRefresh() = %v, want %v", gotStr, test.want)
			}
		})
	}
}

func TestVite(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vite-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	manifestContent := `{
		"src/main.tsx": {
			"file": "assets/main.1234.js",
			"src": "src/main.tsx"
		},
		"src/app.tsx": {
			"file": "assets/app.5678.js", 
			"src": "src/app.tsx"
		}
	}`

	manifestPath := filepath.Join(tmpDir, "manifest.json")
	if err = os.WriteFile(manifestPath, []byte(manifestContent), 0644); err != nil {
		t.Fatalf("Failed to write manifest file: %v", err)
	}

	tests := []struct {
		name      string
		buildDir  string
		assetPath string
		want      template.HTML
		wantErr   bool
	}{
		{
			name:      "Valid asset",
			buildDir:  "dist",
			assetPath: "src/main.tsx",
			want:      template.HTML(`<script type="module" src="/dist/assets/main.1234.js"></script>`),
			wantErr:   false,
		},
		{
			name:      "Non-existent asset",
			buildDir:  "dist",
			assetPath: "nonexistent.tsx",
			want:      template.HTML(""),
			wantErr:   true,
		},
		{
			name:      "Different build directory",
			buildDir:  "build",
			assetPath: "src/app.tsx",
			want:      template.HTML(`<script type="module" src="/build/assets/app.5678.js"></script>`),
			wantErr:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var (
				viteFunc func(string) (template.HTML, error)
				got      template.HTML
			)

			viteFunc, err = vite(manifestPath, test.buildDir)
			if err != nil {
				t.Fatalf("vite() error = %v", err)
			}

			got, err = viteFunc(test.assetPath)
			if (err != nil) != test.wantErr {
				t.Errorf("vite() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("vite() = %v, want %v", got, test.want)
			}
		})
	}

	invalidManifestPath := filepath.Join(tmpDir, "invalid.json")
	if err = os.WriteFile(invalidManifestPath, []byte("invalid json"), 0644); err != nil {
		t.Fatalf("Failed to write invalid manifest file: %v", err)
	}

	_, err = vite(invalidManifestPath, "dist")
	if err == nil {
		t.Error("vite() with invalid manifest should return error")
	}

	_, err = vite("nonexistent.json", "dist")
	if err == nil {
		t.Error("vite() with non-existent manifest should return error")
	}
}
