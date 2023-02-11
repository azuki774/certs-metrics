package model

import (
	"certs-metrics/internal/util"
	"crypto/x509"
	"reflect"
	"testing"
	"time"
)

func TestNewCertsMetricsInfo(t *testing.T) {
	type args struct {
		cf CertificateFile
	}
	tests := []struct {
		name    string
		args    args
		want    CertsMetricsInfo
		nowTime time.Time
	}{
		{
			name: "normal",
			args: args{
				cf: CertificateFile{
					Fullpath: "/file/path/ca.crt",
					Cert: x509.Certificate{
						NotBefore: time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC),
						NotAfter:  time.Date(2011, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			want: CertsMetricsInfo{
				FullPath:             "/file/path/ca.crt",
				FileName:             "ca.crt",
				NotBefore:            time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC),
				NotAfter:             time.Date(2011, 1, 1, 0, 0, 0, 0, time.UTC),
				ValidPeriod:          time.Duration(3600 * 24 * 365 * time.Second),
				RemainingValidPeriod: time.Duration(3600 * 24 * 31 * time.Second),
			},
			nowTime: time.Date(2010, 12, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "expired",
			args: args{
				cf: CertificateFile{
					Fullpath: "/file/path/ca.crt",
					Cert: x509.Certificate{
						NotBefore: time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC),
						NotAfter:  time.Date(2011, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			want: CertsMetricsInfo{
				FullPath:             "/file/path/ca.crt",
				FileName:             "ca.crt",
				NotBefore:            time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC),
				NotAfter:             time.Date(2011, 1, 1, 0, 0, 0, 0, time.UTC),
				ValidPeriod:          time.Duration(3600 * 24 * 365 * time.Second),
				RemainingValidPeriod: time.Duration(0 * time.Second),
			},
			nowTime: time.Date(2012, 12, 1, 0, 0, 0, 0, time.UTC), // expired
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			util.TimeNow = tt.nowTime.UTC
			if got := NewCertsMetricsInfo(tt.args.cf); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCertsMetricsInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
