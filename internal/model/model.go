package model

import (
	"certs-metrics/internal/util"
	"crypto/x509"
	"path/filepath"
	"time"
)

type CertificateFile struct {
	Fullpath string
	Cert    x509.Certificate
}

type CertsMetricsInfo struct {
	FullPath             string
	FileName             string
	NotBefore            time.Time
	NotAfter             time.Time
	ValidPeriod          time.Duration
	RemainingValidPeriod time.Duration
}

func NewCertsMetricsInfo(cf CertificateFile) CertsMetricsInfo {
	ci := CertsMetricsInfo{}
	ci.FullPath = cf.Fullpath
	ci.FileName = filepath.Base(cf.Fullpath)
	ci.NotBefore = cf.Cert.NotBefore
	ci.NotAfter = cf.Cert.NotAfter
	ci.ValidPeriod = ci.NotAfter.Sub(ci.NotBefore)
	t := util.TimeNow()
	ci.RemainingValidPeriod = ci.NotAfter.Sub(t)
	if ci.RemainingValidPeriod < 0 {
		ci.RemainingValidPeriod = 0
	}
	return ci
}
