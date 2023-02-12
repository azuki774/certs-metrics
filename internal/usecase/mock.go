package usecase

import (
	"certs-metrics/internal/model"
	"crypto/x509"
	"path/filepath"
	"time"
)

type mockCertsLoader struct {
	err error
}

func (m *mockCertsLoader) Load(dir string) (cf model.CertificateFile, err error) {
	if m.err != nil {
		return model.CertificateFile{}, m.err
	}

	cf = model.CertificateFile{
		Fullpath: filepath.Join("/root/" + dir),
		Cert: x509.Certificate{
			NotBefore: time.Date(2001, 1, 23, 12, 0, 0, 0, time.UTC),
			NotAfter:  time.Date(2001, 1, 24, 12, 0, 0, 0, time.UTC),
		},
	}
	return cf, nil
}
