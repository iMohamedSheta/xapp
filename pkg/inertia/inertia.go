package inertia

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	gonertia "github.com/romsar/gonertia/v2"
)

type Inertia struct {
	*gonertia.Inertia
	beforeRender func(c context.Context)
	mu           sync.RWMutex
}

type Props = gonertia.Props

var IsInertiaRequest = gonertia.IsInertiaRequest

func New(i *gonertia.Inertia) *Inertia {
	return &Inertia{Inertia: i}
}

func (i *Inertia) OnBeforeRender(f func(context.Context)) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.beforeRender = f
}

func WithProps(c *gin.Context, props Props) *http.Request {
	ctx := gonertia.SetProps(c.Request.Context(), props)
	return c.Request.WithContext(ctx)
}

func (i *Inertia) Render(c *gin.Context, component string, props Props) error {
	if i.beforeRender != nil {
		safeCall(func() {
			i.beforeRender(c)
		})
	}
	return i.Inertia.Render(c.Writer, c.Request, component, props)
}

func (i *Inertia) Back(c *gin.Context, status int) {
	i.Inertia.Back(c.Writer, c.Request, status)
}

func InitInertia(viewsFS embed.FS, rootViewPath string, opts ...gonertia.Option) *Inertia {
	viteHotFile := "./public/hot"

	// check if laravel-vite-plugin is running in dev mode
	_, err := os.Stat(viteHotFile)
	if err == nil {
		// Dev mode with hot reload
		opts = append(opts, gonertia.WithSSR())

		// Read from embedded FS
		rootViewContent, err := fs.ReadFile(viewsFS, rootViewPath)
		if err != nil {
			log.Fatal("failed to read root view from embedded FS:", err)
		}

		i, err := gonertia.New(
			string(rootViewContent),
			opts...,
		)
		if err != nil {
			log.Fatal(err)
		}

		_ = i.ShareTemplateFunc("vite", func(entry string) (template.HTML, error) {
			content, err := os.ReadFile(viteHotFile)
			if err != nil {
				fmt.Println(err)
				return "", err
			}
			url := strings.TrimSpace(string(content))
			if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
				url = url[strings.Index(url, ":")+1:]
			} else {
				url = "//localhost:8080"
			}
			if entry != "" && !strings.HasPrefix(entry, "/") {
				entry = "/" + entry
			}

			fullSrc := url + entry
			fmt.Println(fullSrc)

			// In dev mode, Vite handles CSS dynamically via the JS bundle
			tag := fmt.Sprintf(`<script type="module" src="http://localhost:5173/@vite/client"></script>`+"\n"+`<script type="module" src="http://localhost:5173%s"></script>`, entry)

			return template.HTML(tag), nil
		})
		return New(i)
	}

	// Production mode with manifest
	manifestPath := "./public/build/manifest.json"

	// check if the manifest file exists
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		err := os.Rename("./public/build/.vite/manifest.json", "./public/build/manifest.json")
		if err != nil {
			return nil
		}
	}

	opts = append(opts, gonertia.WithVersionFromFile(manifestPath), gonertia.WithSSR())

	// Read from embedded FS
	rootViewContent, err := fs.ReadFile(viewsFS, rootViewPath)
	if err != nil {
		log.Fatal("failed to read root view from embedded FS:", err)
	}

	i, err := gonertia.New(
		string(rootViewContent),
		opts...,
	)
	if err != nil {
		log.Fatal(err)
	}

	_ = i.ShareTemplateFunc("vite", vite(manifestPath, "/public/build/"))

	return New(i)
}

func vite(manifestPath, buildDir string) func(path string) (template.HTML, error) {
	// Check if manifest exists
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		log.Printf("[INERTIA] Manifest not found at: %s", manifestPath)
		return func(p string) (template.HTML, error) {
			return "", fmt.Errorf("manifest file not found: %s", manifestPath)
		}
	}

	f, err := os.Open(manifestPath)
	if err != nil {
		log.Printf("[INERTIA] Cannot open manifest file: %s", err)
		return func(p string) (template.HTML, error) {
			return "", fmt.Errorf("cannot open manifest: %w", err)
		}
	}
	defer func() { _ = f.Close() }()

	type ManifestEntry struct {
		File string   `json:"file"`
		Src  string   `json:"src"`
		CSS  []string `json:"css,omitempty"`
	}

	viteAssets := make(map[string]*ManifestEntry)
	if err := json.NewDecoder(f).Decode(&viteAssets); err != nil {
		log.Printf("[INERTIA] Cannot decode manifest: %s", err)
		return func(p string) (template.HTML, error) {
			return "", fmt.Errorf("cannot decode manifest: %w", err)
		}
	}

	// Debug: print all manifest entries
	log.Println("[INERTIA] ===== Vite Manifest Loaded =====")
	for k, v := range viteAssets {
		log.Printf("[INERTIA] Key: %q -> File: %q", k, v.File)
	}
	log.Println("[INERTIA] ==================================")

	return func(p string) (template.HTML, error) {
		log.Printf("[INERTIA] Looking up asset: %q", p)

		// Try exact match
		if val, ok := viteAssets[p]; ok {
			result := path.Join("/", buildDir, val.File)

			var tags []string
			for _, cssFile := range val.CSS {
				cssUrl := path.Join("/", buildDir, cssFile)
				tags = append(tags, fmt.Sprintf(`<link rel="stylesheet" href="%s">`, cssUrl))
			}
			tags = append(tags, fmt.Sprintf(`<script type="module" src="%s"></script>`, result))

			log.Printf("[INERTIA] ✓ Found: %q -> %q", p, result)
			return template.HTML(strings.Join(tags, "\n")), nil
		}

		// Try without leading slash
		pTrimmed := strings.TrimPrefix(p, "/")
		if val, ok := viteAssets[pTrimmed]; ok {
			result := path.Join("/", buildDir, val.File)

			var tags []string
			for _, cssFile := range val.CSS {
				cssUrl := path.Join("/", buildDir, cssFile)
				tags = append(tags, fmt.Sprintf(`<link rel="stylesheet" href="%s">`, cssUrl))
			}
			tags = append(tags, fmt.Sprintf(`<script type="module" src="%s"></script>`, result))

			log.Printf("[INERTIA] ✓ Found (trimmed): %q -> %q", p, result)
			return template.HTML(strings.Join(tags, "\n")), nil
		}

		// Asset not found - list available keys
		log.Printf("[INERTIA] ✗ Asset %q not found", p)
		log.Println("[INERTIA] Available manifest keys:")
		for k := range viteAssets {
			log.Printf("[INERTIA]   - %q", k)
		}

		return "", fmt.Errorf("asset %q not found in manifest", p)
	}
}

// Use safeCall to avoid panic if the callback function panics
func safeCall(f func()) {
	defer func() {
		_ = recover()
	}()
	f()
}
