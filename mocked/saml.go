package mocked

import (
	"auth-guardian/config"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/crewjam/saml/samlidp"
	"golang.org/x/crypto/bcrypt"
)

// RunMockSAMLIDP runs a mocked SAML IDP
func RunMockSAMLIDP() *http.Server {
	// Overrite config
	config.SAMLCrt = "saml_mock.crt"
	config.SAMLKey = "saml_mock.key"
	config.IdpMetadataURL = "http://localhost:3002/metadata"
	config.IdpRegisterURL = "http://localhost:3002/services/sp"
	config.SelfRootURL = "http://localhost:3000"

	keyRaw, err := ioutil.ReadFile(config.SAMLKey)
	if err != nil {
		log.Fatal(err)
	}

	var key = func() crypto.PrivateKey {
		b, _ := pem.Decode([]byte(keyRaw))
		k, _ := x509.ParsePKCS1PrivateKey(b.Bytes)
		return k
	}()

	crtRaw, err := ioutil.ReadFile(config.SAMLCrt)
	if err != nil {
		log.Fatal(err)
	}
	var cert = func() *x509.Certificate {
		b, _ := pem.Decode([]byte(crtRaw))
		c, _ := x509.ParseCertificate(b.Bytes)
		return c
	}()

	baseURL, err := url.Parse("http://localhost:3002")
	idpServer, err := samlidp.New(samlidp.Options{
		URL:         *baseURL,
		Key:         key,
		Certificate: cert,
		Store:       &samlidp.MemoryStore{},
	})
	if err != nil {
		log.Fatalf("%s", err)
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("SpiritedAway"), bcrypt.DefaultCost)
	err = idpServer.Store.Put("/users/noface", samlidp.User{
		Name:           "noface",
		HashedPassword: hashedPassword,
		Groups:         []string{"Spirit", "Antagonist"},
		Email:          "no.face@local.com",
		CommonName:     "No Face",
		Surname:        "Face",
		GivenName:      "No",
	})
	if err != nil {
		log.Fatalf("%s", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", idpServer)
	server := &http.Server{Addr: ":3002", Handler: mux}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Panic(err)
		}
	}()

	return server
}
