package repository

import (
	"certs-metrics/internal/model"
	"crypto/x509"
	"reflect"
	"testing"
	"time"
)

func TestCertsLoader_Load(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		c       *CertsLoader
		args    args
		wantCf  model.CertificateFile
		wantErr bool
	}{
		{
			name: "ok",
			c:    &CertsLoader{},
			args: args{dir: "../../test/testca.crt"},
			wantCf: model.CertificateFile{
				Fullpath: "DUMMY", // not tested
				Cert: x509.Certificate{
					// notBefore=Feb  4 11:53:30 2023 GMT
					// notAfter=Feb  1 11:53:30 2033 GMT
					NotBefore: time.Date(2023, 2, 4, 11, 53, 30, 0, time.UTC),
					NotAfter:  time.Date(2033, 2, 1, 11, 53, 30, 0, time.UTC),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CertsLoader{}
			gotCf, err := c.Load(tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("CertsLoader.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCf.Cert.NotBefore, tt.wantCf.Cert.NotBefore) {
				t.Errorf("CertsLoader.Load() NotBefore = %#v, want %v", gotCf.Cert.NotBefore, tt.wantCf.Cert.NotBefore)
			}

			if !reflect.DeepEqual(gotCf.Cert.NotAfter, tt.wantCf.Cert.NotAfter) {
				t.Errorf("CertsLoader.Load() NotAfter = %#v, want %v", gotCf.Cert.NotAfter, tt.wantCf.Cert.NotAfter)
			}
		})
	}
}
