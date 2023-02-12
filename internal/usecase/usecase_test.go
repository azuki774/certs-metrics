package usecase

import (
	"certs-metrics/internal/model"
	"certs-metrics/internal/util"
	"fmt"
	"reflect"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var l *zap.Logger

func init() {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	config.EncoderConfig.EncodeTime = jSTTimeEncoder
	l, _ = config.Build()

	l.WithOptions(zap.AddStacktrace(zap.ErrorLevel))
}

func jSTTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	const layout = "2006-01-02T15:04:05+09:00"
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	enc.AppendString(t.In(jst).Format(layout))
}

func TestUsecase_LoadCertsInfo(t *testing.T) {
	type fields struct {
		l *zap.Logger
		c CertsLoader
	}
	type args struct {
		dirs []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantCis []model.CertsMetricsInfo
		wantErr bool
		nowTime time.Time
	}{
		{
			name: "load ok",
			fields: fields{
				l: l,
				c: &mockCertsLoader{},
			},
			args: args{
				dirs: []string{"test1.crt", "test2.crt"},
			},
			wantCis: []model.CertsMetricsInfo{
				{
					FullPath:             "/root/test1.crt",
					FileName:             "test1.crt",
					NotBefore:            time.Date(2001, 1, 23, 12, 0, 0, 0, time.UTC),
					NotAfter:             time.Date(2001, 1, 24, 12, 0, 0, 0, time.UTC),
					ValidPeriod:          time.Duration(3600 * 24 * time.Second),
					RemainingValidPeriod: time.Duration(3600 * 12 * time.Second),
				},
				{
					FullPath:             "/root/test2.crt",
					FileName:             "test2.crt",
					NotBefore:            time.Date(2001, 1, 23, 12, 0, 0, 0, time.UTC),
					NotAfter:             time.Date(2001, 1, 24, 12, 0, 0, 0, time.UTC),
					ValidPeriod:          time.Duration(3600 * 24 * time.Second),
					RemainingValidPeriod: time.Duration(3600 * 12 * time.Second),
				},
			},
			wantErr: false,
			nowTime: time.Date(2001, 1, 24, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "load ng",
			fields: fields{
				l: l,
				c: &mockCertsLoader{err: fmt.Errorf("anything error")},
			},
			args: args{
				dirs: []string{"test1.crt", "test2.crt"},
			},
			wantCis: []model.CertsMetricsInfo{},
			wantErr: true,
			nowTime: time.Date(2001, 1, 24, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				l: tt.fields.l,
				c: tt.fields.c,
			}
			util.TimeNow = tt.nowTime.UTC
			gotCis, err := u.LoadCertsInfo(tt.args.dirs)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.LoadCertsInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCis, tt.wantCis) {
				t.Errorf("Usecase.LoadCertsInfo() = %v, want %v", gotCis, tt.wantCis)
			}
		})
	}
}
