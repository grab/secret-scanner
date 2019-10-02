package web

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/static"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner/session"
)

const (
	GithubBaseUri   = "https://raw.githubusercontent.com"
	MaximumFileSize = 102400
	CspPolicy       = "default-src 'none'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'"
	ReferrerPolicy  = "no-referrer"
)

type binaryFileSystem struct {
	fs http.FileSystem
}

func NewRouter(s *session.Session) *gin.Engine {
	if *s.Options.Debug == true {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(static.Serve("/", BinaryFileSystem("static")))
	router.Use(secure.New(secure.Config{
		SSLRedirect:           false,
		IsDevelopment:         false,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: CspPolicy,
		ReferrerPolicy:        ReferrerPolicy,
	}))
	router.GET("/stats", func(c *gin.Context) {
		c.JSON(200, s.Stats)
	})
	router.GET("/findings", func(c *gin.Context) {
		c.JSON(200, s.Findings)
	})
	//router.GET("/targets", func(c *gin.Context) {
	//	c.JSON(200, s.Targets)
	//})
	router.GET("/repositories", func(c *gin.Context) {
		c.JSON(200, s.Repositories)
	})
	//router.GET("/files/:owner/:repo/:commit/*path", fetchFile)

	return router
}

func BinaryFileSystem(root string) *binaryFileSystem {
	fs := &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: root}
	return &binaryFileSystem{
		fs,
	}
}

func (b *binaryFileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

func (b *binaryFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := b.fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}
