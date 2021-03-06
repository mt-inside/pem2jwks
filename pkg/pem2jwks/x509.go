package pem2jwks

import (
	"crypto"
	"crypto/x509"
	"fmt"
)

// ParsePublicKey extracts the public key from the given cryptographic object.
// That object is expected to be the DER encoding of:
// * an x509 encoding (ie ASN.1 serialization) of a public key
// * a PKCS#1 (ie ASN.1 encoding) containing a public key
// * an x509 encoding (ie ASN.1 serialization) of a certificate
// * a PKCS#1 (ie ASN.1 encoding) containing a private key
// * a PKCS#8 (ie ASN.1 encoding) containing a private key
// * a SEC 1 (ie ASN.1 encoding) containing an EC private key
func ParsePublicKey(der []byte) (crypto.PublicKey, error) {

	if pubKey, err := x509.ParsePKIXPublicKey(der); err == nil {
		return pubKey, nil
	} else if pubKey, err := x509.ParsePKCS1PublicKey(der); err == nil { // RSA only; type *rsa.PublicKey
		return pubKey, nil
	} else if cert, err := x509.ParseCertificate(der); err == nil {
		return cert.PublicKey, nil
	} else if privKey, err := x509.ParsePKCS1PrivateKey(der); err == nil { // RSA only; type *rsa.PrivateKey
		return privKey.Public(), nil
	} else if privKey, err := x509.ParsePKCS8PrivateKey(der); err == nil { // OpenSSL 3+ default. RSA, ECDSA, Ed25519; type any, however: https://pkg.go.dev/crypto#PrivateKey
		return privKey.(interface {
			Public() crypto.PublicKey
		}).Public(), nil
	} else if privKey, err := x509.ParseECPrivateKey(der); err == nil { // ECDSA only; type *ecdsa.PrivateKey
		return privKey.Public(), nil
	} else {
		return nil, fmt.Errorf("DER block does not encode a recognised cryptographic object")
	}
}

func ParsePrivateKey(der []byte) (crypto.PrivateKey /* alias: any */, error) {

	if _, err := x509.ParsePKIXPublicKey(der); err == nil {
		return nil, fmt.Errorf("need a private key; got a public")
	} else if _, err := x509.ParsePKCS1PublicKey(der); err == nil { // RSA only; type *rsa.PublicKey
		return nil, fmt.Errorf("need a private key; got a public")
	} else if _, err := x509.ParseCertificate(der); err == nil {
		return nil, fmt.Errorf("need a private key; got a cert")
	} else if privKey, err := x509.ParsePKCS1PrivateKey(der); err == nil { // RSA only; type *rsa.PrivateKey
		return privKey, nil
	} else if privKey, err := x509.ParsePKCS8PrivateKey(der); err == nil { // RSA, ECDSA, Ed25519; type any, however: https://pkg.go.dev/crypto#PrivateKey
		return privKey, nil
	} else if privKey, err := x509.ParseECPrivateKey(der); err == nil { // ECDSA only; type *ecdsa.PrivateKey
		return privKey, nil
	} else {
		return nil, fmt.Errorf("DER block does not encode a recognised cryptographic object")
	}
}
