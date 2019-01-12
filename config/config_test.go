package config

import (
	"testing"
)

func TestLoadFile(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "ok",
			args:    args{filename: "testdata/conf.good.yml"},
			wantErr: false,
		},
		{
			name:    "empty_shell_path",
			args:    args{filename: "testdata/empty_shell_path.yml"},
			wantErr: true,
		},
		{
			name:    "empty_services",
			args:    args{filename: "testdata/empty_services.yml"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LoadFile(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
