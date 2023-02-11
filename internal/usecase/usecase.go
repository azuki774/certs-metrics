package usecase

import (
	"certs-metrics/internal/model"
)

type CertsLoader interface {
	Load(dir string) (cf model.CertificateFile, err error)
}

// type Usecase struct {
// 	l *zap.Logger
// 	c CertsLoader
// }

// func (u *Usecase) LoadCertsInfo(dirs []string) (err error) {
// 	for _, dir := range dirs {
// 		ci, err := u.loadCertInfo(dir)
// 	}
// }

// func (u *Usecase) loadCertInfo(dir string) (ci []model.CertsMetricsInfo, err error) {

// }
