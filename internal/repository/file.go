package repository

import (
	"certs-metrics/internal/model"
	"crypto/x509"
	"encoding/pem"
	"io"
	"os"
	"path/filepath"
)

type CertsLoader struct {
}

func (c *CertsLoader) Load(dir string) (cf model.CertificateFile, err error) {
	fullPath, err := filepath.Abs(dir)
	if err != nil {
		return model.CertificateFile{}, err
	}

	cf.Fullpath = fullPath
	f, err := os.Open(fullPath)
	if err != nil {
		return model.CertificateFile{}, err
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return model.CertificateFile{}, err
	}

	block, _ := pem.Decode(b)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return model.CertificateFile{}, err
	}

	cf.Cert = *cert
	return cf, nil
}
