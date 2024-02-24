package config

import (
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	type args struct {
		folder      string
		filename    string
		envfilename string
	}

	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "test config",
			args: args{
				folder:      "fixtures",
				filename:    "test",
				envfilename: ".test",
			},
			want: &Config{
				DB: DB{
					Host:     "localhost",
					Port:     "5432",
					Username: "postgres",
					Name:     "postgres",
					SSLMode:  "disable",
					Password: "qwerty",
				},
				Server: Server{
					Port: "8000",
				},
				GRPC: GRPC{
					Port: "9000",
				},
				JWT: JWT{
					AccessTTL:  15 * time.Minute,
					RefreshTTL: 1 * time.Hour,
				},
				Hash: Hash{
					Salt: "salt",
				},
				Auth: Auth{
					Secret: "secret",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.folder, tt.args.filename, tt.args.envfilename)
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Init() got = %v, want %v", got, tt.want)
			}
		})
	}
}
