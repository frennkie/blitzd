// Copyright (c) 2013-2017 The btcsuite developers
// Copyright (c) 2015-2016 The Decred developers
// Copyright (C) 2015-2017 The Lightning Network Developers

package util

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"time"
)

const (
	// Make certificate valid for 14 months.
	autogenCertValidity = 14 /*months*/ * 30 /*days*/ * 24 * time.Hour
)

var (
	// End of ASN.1 time.
	endOfTime = time.Date(2049, 12, 31, 23, 59, 59, 0, time.UTC)

	// Max serial number.
	serialNumberLimit = new(big.Int).Lsh(big.NewInt(1), 128)
)

func GenRootCaSignedClientServerCert(alias, rootCaCertPath, serverCertPath, serverKeyPath, clientCertPath, clientKeyPath string) error {
	// generate Root CA and save as Public Key as PEM. Private Key is never saved!
	rootCaCert, rootCaPriv := GenRootCa(alias)
	rootCaCertPem := pem.EncodeToMemory(pemBlockForCert(rootCaCert))
	_ = exportCert(rootCaCert, rootCaCertPath)
	rootCaPrivPem := pem.EncodeToMemory(pemBlockForKey(rootCaPriv))

	// create signed certificate for "server" and "client" with RootCA public/private key from memory
	signedServerCert, signedServerCertPriv := GenCaSignedServerCert(alias, rootCaCertPem, rootCaPrivPem)
	_ = exportCert(signedServerCert, serverCertPath)
	_ = exportKey(signedServerCertPriv, serverKeyPath)

	signedClientCert, signedClientCertPriv := GenCaSignedClientCert(alias, rootCaCertPem, rootCaPrivPem)
	_ = exportCert(signedClientCert, clientCertPath)
	_ = exportKey(signedClientCertPriv, clientKeyPath)

	return nil
}

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForCert(cert []byte) *pem.Block {
	return &pem.Block{Type: "CERTIFICATE", Bytes: cert}
}

func pemBlockForKey(priv *ecdsa.PrivateKey) *pem.Block {
	b, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
		os.Exit(2)
	}
	return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
}

func exportCert(cert []byte, path string) error {

	// Public key
	certOut, err := os.Create(path)
	if err != nil {
		log.Printf("failed to open "+path+" for writing: %s", err)
		return err
	}

	_ = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: cert})
	_ = certOut.Close()
	log.Print("written " + path + "\n")

	return nil
}

func exportKey(key *ecdsa.PrivateKey, path string) error {

	// Private key
	keyOut, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Printf("failed to open "+path+" for writing: %s", err)
		return err
	}

	_ = pem.Encode(keyOut, pemBlockForKey(key))
	_ = keyOut.Close()
	log.Print("written " + path + "\n")

	return nil

}

func genPrivKey(ecdsaCurve string) *ecdsa.PrivateKey {

	var priv *ecdsa.PrivateKey
	var err error
	switch ecdsaCurve {
	case "":
		priv, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case "P256":
		priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "P384":
		priv, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case "P521":
		priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		_, _ = fmt.Fprintf(os.Stderr, "Unrecognized elliptic curve: %q", ecdsaCurve)
		os.Exit(1)
	}
	if err != nil {
		log.Printf("failed to generate private key: %s", err)
		panic(err)
	}

	return priv
}

func GenRootCa(alias string) ([]byte, *ecdsa.PrivateKey) {

	now := time.Now()
	validUntil := now.Add(autogenCertValidity)

	// Check that the certificate validity isn't past the ASN.1 end of time.
	if validUntil.After(endOfTime) {
		validUntil = endOfTime
	}

	// Generate a serial number that's below the serialNumberLimit.
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		panic(err)
	}

	ca := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Blitzd Ephemeral CA (" + alias + ")"},
			Country:      []string{"DE"},
			CommonName:   "Blitzd Ephemeral CA (" + alias + ")",
		},
		NotBefore:             now.AddDate(0, 0, -1),
		NotAfter:              validUntil,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		MaxPathLen:            0,
		MaxPathLenZero:        true,
	}

	priv := genPrivKey("P256")

	signedRootCaCert, err := x509.CreateCertificate(rand.Reader, ca, ca, publicKey(priv), priv)
	if err != nil {
		log.Printf("Failed to create Root CA certificate: %s", err)
		panic(err)
	}

	return signedRootCaCert, priv

}

func GenCaSignedServerCert(alias string, certPEMBlock, keyPEMBlock []byte) ([]byte, *ecdsa.PrivateKey) {

	// Load CA
	catls, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		panic(err)
	}
	ca, err := x509.ParseCertificate(catls.Certificate[0])
	if err != nil {
		panic(err)
	}

	// ToDO check (currently doing the same things as LND for now)

	now := time.Now()
	validUntil := now.Add(autogenCertValidity)

	// Check that the certificate validity isn't past the ASN.1 end of time.
	if validUntil.After(endOfTime) {
		validUntil = endOfTime
	}

	// Generate a serial number that's below the serialNumberLimit.
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		panic(err)
	}

	// Collect the host's IP addresses, including loopback, in a slice.
	ipAddresses := []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")}

	// addIP appends an IP address only if it isn't already in the slice.
	addIP := func(ipAddr net.IP) {
		for _, ip := range ipAddresses {
			if bytes.Equal(ip, ipAddr) {
				return
			}
		}
		ipAddresses = append(ipAddresses, ipAddr)
	}

	// Add all the interface IPs that aren't already in the slice.
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, a := range addrs {
		ipAddr, _, err := net.ParseCIDR(a.String())
		if err == nil {
			addIP(ipAddr)
		}
	}

	// Add extra IPs to the slice.
	//for _, ip := range cfg.TLSExtraIPs {
	//	ipAddr := net.ParseIP(ip)
	//	if ipAddr != nil {
	//		addIP(ipAddr)
	//	}
	//}

	// Collect the host's names into a slice.
	host, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	dnsNames := []string{host}
	if host != "localhost" {
		dnsNames = append(dnsNames, "localhost")
	}
	//dnsNames = append(dnsNames, cfg.TLSExtraDomains...)

	// Also add fake hostnames for unix sockets, otherwise hostname
	// verification will fail in the client.
	dnsNames = append(dnsNames, "unix", "unixpacket")

	// Prepare certificate
	cert := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Blitzd Certificates (" + alias + ")"},
			Country:      []string{"DE"},
			CommonName:   "Blitzd (" + alias + ")",
		},
		NotBefore:    now.AddDate(0, 0, -1),
		NotAfter:     validUntil,
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
		DNSNames:     dnsNames,
		IPAddresses:  ipAddresses,
	}

	priv := genPrivKey("P256")

	signedCert, err := x509.CreateCertificate(rand.Reader, cert, ca, publicKey(priv), catls.PrivateKey)
	if err != nil {
		panic(err)
	}

	return signedCert, priv
}

func GenCaSignedClientCert(alias string, certPEMBlock, keyPEMBlock []byte) ([]byte, *ecdsa.PrivateKey) {

	// Load CA
	catls, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		panic(err)
	}
	ca, err := x509.ParseCertificate(catls.Certificate[0])
	if err != nil {
		panic(err)
	}

	now := time.Now()
	validUntil := now.Add(autogenCertValidity)

	// Check that the certificate validity isn't past the ASN.1 end of time.
	if validUntil.After(endOfTime) {
		validUntil = endOfTime
	}

	// Generate a serial number that's below the serialNumberLimit.
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		panic(err)
	}

	// Prepare certificate
	cert := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Blitzd Certificates (" + alias + ")"},
			Country:      []string{"DE"},
			CommonName:   "Blitzd Client",
		},
		NotBefore:    now.AddDate(0, 0, -1),
		NotAfter:     validUntil,
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	priv := genPrivKey("P256")

	signedCert, err := x509.CreateCertificate(rand.Reader, cert, ca, publicKey(priv), catls.PrivateKey)
	if err != nil {
		panic(err)
	}

	return signedCert, priv
}
