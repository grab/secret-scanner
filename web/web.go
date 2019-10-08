package web

import (
	"fmt"

	"gitlab.myteksi.net/product-security/ssdlc/secret-scanner/scanner/session"
)

// InitRouter ...
func InitRouter(address, port string, sess *session.Session) {
	bind := fmt.Sprintf("%s:%s", address, port)
	r := NewRouter(sess)
	if err := r.Run(bind); err != nil {
		sess.Out.Fatal("Error when starting web server: %s\n", err)
	}
	sess.Out.Info("Web server started at %s", bind)
}
