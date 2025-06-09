package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/romsar/gonertia/v2"
)

func ProjectRoot() string {
	currentDir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		_, err := os.ReadFile(filepath.Join(currentDir, "go.mod"))
		if os.IsNotExist(err) {
			if currentDir == filepath.Dir(currentDir) {
				return ""
			}
			currentDir = filepath.Dir(currentDir)
			continue
		} else if err != nil {
			return ""
		}
		break
	}
	return currentDir
}

func NewInertia() (*gonertia.Inertia, error) {
	rootDir := ProjectRoot()
	viteHotFile := filepath.Join(rootDir, "public", "hot")
	rootViewFile := filepath.Join(rootDir, "ui", "index.html")
	manifestPath := filepath.Join(rootDir, "public", "build", "manifest.json")

	// check if laravel-vite-plugin is running in dev mode (it puts a "hot" file in the public folder)
	url, err := viteHotFileUrl(viteHotFile)
	if err != nil {
		return nil, err
	}

	if url != "" {
		return newInertiaDevMode(rootViewFile, url)
	}

	i, err := gonertia.NewFromFile(
		rootViewFile,
		gonertia.WithVersionFromFile(manifestPath),
	)
	if err != nil {
		return nil, err
	}

	viteFn, err := vite(manifestPath, "/public/build/")
	if err != nil {
		return nil, err
	}

	i.ShareTemplateFunc("vite", viteFn)
	i.ShareTemplateFunc("viteReactRefresh", func () (template.HTML, error) {
		return "", nil
	})

	return i, nil
}

func newInertiaDevMode(rootViewFile, url string) (*gonertia.Inertia, error) {
	i, err := gonertia.NewFromFile(
		rootViewFile,
	)

	if err != nil {
		return nil, err
	}

	i.ShareTemplateFunc("vite", func(entry string) (template.HTML, error) {
		if entry != "" && !strings.HasPrefix(entry, "/") {
			entry = "/" + entry
		}
		htmlTag := fmt.Sprintf(`<script type="module" src="%s%s"></script>`, url, entry)
		return template.HTML(htmlTag), nil
	})

	i.ShareTemplateFunc("viteReactRefresh", viteReactRefresh(url))

	return i, nil
}

// viteHotFileUrl Get the vite hot file url
func viteHotFileUrl(viteHotFile string) (string, error) {
	_, err := os.Stat(viteHotFile)
	if err != nil {
		return "", nil
	}
	content, err := os.ReadFile(viteHotFile)
	if err != nil {
		return "", err
	}
	url := strings.TrimSpace(string(content))
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		url = url[strings.Index(url, ":")+1:]
	} else {
		url = "//localhost:1323"
	}
	return url, nil
}

// viteReactRefresh Generate React refresh runtime script
func viteReactRefresh(url string) func() (template.HTML, error) {
	return func() (template.HTML, error) {
		if url == "" {
			return "", nil
		}

		script := fmt.Sprintf(`
<script type="module">
    import RefreshRuntime from '%s/@react-refresh'
    RefreshRuntime.injectIntoGlobalHook(window)
    window.$RefreshReg$ = () => {}
    window.$RefreshSig$ = () => (type) => type
    window.__vite_plugin_react_preamble_installed__ = true
</script>`, url)

		return template.HTML(script), nil
	}
}

func vite(manifestPath, buildDir string) (func(path string) (template.HTML, error), error) {
	f, err := os.Open(manifestPath)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	viteAssets := make(map[string]*struct {
		File   string `json:"file"`
		Source string `json:"src"`
	})

	err = json.NewDecoder(f).Decode(&viteAssets)
	if err != nil {
		return nil, err
	}

	return func(p string) (template.HTML, error) {
		htmlStr := `<script type="module" src="%s"></script>`
		
		if val, ok := viteAssets[p]; ok {
			tag := fmt.Sprintf(htmlStr, path.Join("/", buildDir, val.File))
			return template.HTML(tag), nil
		}
		return template.HTML(""), fmt.Errorf("asset %q not found", p)
	}, nil
}
