package httpsserve

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"httpsproxy/proxy"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"
)

var logger = log.New(os.Stderr, "httpsproxy:", log.Llongfile|log.LstdFlags)

func Serve(listenAdress string){
	cert, err := genCertificate()
	if err != nil {
		logger.Fatal(err)
	}

	server := &http.Server{
		Addr: listenAdress,
		TLSConfig: 	&tls.Config{Certificates: []tls.Certificate{cert},},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u := r.URL
			if u.Host == "" {
				u.Host = r.Host
			}
			hostname := u.Hostname()
			switch {
			case strings.HasSuffix(hostname, "mobage.jp"):
				fallthrough
			case strings.HasSuffix(hostname, "mbga.jp"):
				fallthrough
			case strings.HasSuffix(hostname, "gree.net"):
				fallthrough
			case strings.HasSuffix(hostname, "granbluefantasy.jp"):
				fallthrough
			case strings.HasSuffix(hostname, "203.104.248.14"):
				proxy.Serve(w, r)
			default:
				w.WriteHeader(403)
				_, err := w.Write([]byte("Host not Allowed\r\n"))
				if err != nil {
					log.Fatal(err)
				}
			}
		}),
	}

	logger.Fatal(server.ListenAndServe())

}

func genCertificate() (cert tls.Certificate, err error){
	rawCert, rawKey, err := generateKeyPair()
	if err != nil {
		return
	}
	return tls.X509KeyPair(rawCert, rawKey)

}

func generateKeyPair() (rawCert, rawKey []byte, err error) {
	// Create private key and self-signed certificate
	// Adapted from https://golang.org/src/crypto/tls/generate_cert.go

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return
	}
	validFor := time.Hour * 24 * 365 * 10 // ten years
	notBefore := time.Now()
	notAfter := notBefore.Add(validFor)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Zarten"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return
	}

	rawCert = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	rawKey = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	return
}
