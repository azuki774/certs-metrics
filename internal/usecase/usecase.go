package usecase

import (
	"certs-metrics/internal/model"

	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type CertsLoader interface {
	Load(dir string) (cf model.CertificateFile, err error)
}

type Usecase struct {
	l *zap.Logger
	c CertsLoader
}

// LoadCertsInfo uses 'check' commands
func (u *Usecase) LoadCertsInfo(dirs []string) (cis []model.CertsMetricsInfo, err error) {
	cis = []model.CertsMetricsInfo{}
	for _, dir := range dirs {
		ld := u.l.With(zap.String("dir", dir))
		ci, nerr := u.loadCertInfo(dir)
		if nerr != nil {
			err = multierr.Append(err, nerr)
			ld.Error("failed to load certs information", zap.Error(err))
		} else {
			// no error
			cis = append(cis, ci)
			ld.Info("load cert", zap.Time("notBefore", ci.NotBefore))
			ld.Info("load cert", zap.Time("notAfter", ci.NotAfter))
			ld.Info("load cert", zap.Duration("validPeriod", ci.ValidPeriod))
			ld.Info("load cert", zap.Duration("remainingValidPeriod", ci.RemainingValidPeriod))
		}
	}

	return cis, err
}

func (u *Usecase) loadCertInfo(dir string) (ci model.CertsMetricsInfo, err error) {
	cf, err := u.c.Load(dir)
	if err != nil {
		return model.CertsMetricsInfo{}, err
	}

	ci = model.NewCertsMetricsInfo(cf)
	return ci, nil
}
