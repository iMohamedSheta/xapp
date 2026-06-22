package tls

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"time"
)

// LoadOrCreateCA loads an existing CA from disk, or creates and saves a new one.
// Run this once — commit ca.crt to the repo so teammates can import it.
func LoadOrCreateCA(certFile, keyFile string) (tls.Certificate, *x509.Certificate, error) {
	// If CA already exists on disk, load it
	if _, err := os.Stat(certFile); err == nil {
		ca, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return tls.Certificate{}, nil, err
		}
		parsed, err := x509.ParseCertificate(ca.Certificate[0])
		if err != nil {
			return tls.Certificate{}, nil, err
		}
		return ca, parsed, nil
	}

	// Otherwise generate a new CA
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	template := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "My App Local CA"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(10 * 365 * 24 * time.Hour), // 10 years
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
	}

	certDER, _ := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)

	// Save CA cert and key to disk
	cf, _ := os.Create(certFile)
	pem.Encode(cf, &pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	cf.Close()

	kf, _ := os.Create(keyFile)
	keyDER, _ := x509.MarshalECPrivateKey(key)
	pem.Encode(kf, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})
	kf.Close()

	parsed, _ := x509.ParseCertificate(certDER)
	return tls.Certificate{Certificate: [][]byte{certDER}, PrivateKey: key}, parsed, nil
}

// NewServerCert generates a fresh server cert on every startup, signed by your CA.
func NewServerCert(ca tls.Certificate, caParsed *x509.Certificate, ips []string) (tls.Certificate, error) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	var ipAddrs []net.IP
	for _, ip := range ips {
		ipAddrs = append(ipAddrs, net.ParseIP(ip))
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()), // unique every time
		Subject:      pkix.Name{CommonName: "local-server"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(24 * time.Hour), // expires tomorrow — fresh each run
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  ipAddrs,
		DNSNames:     []string{"localhost"},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, caParsed, &key.PublicKey, ca.PrivateKey)
	if err != nil {
		return tls.Certificate{}, err
	}

	return tls.Certificate{
		Certificate: [][]byte{certDER},
		PrivateKey:  key,
	}, nil
}
