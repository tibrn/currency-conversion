package config

import (
	"testing"

	"github.com/gobuffalo/envy"
)

func Test_mustGet(t *testing.T) {
	type args struct {
		param string
	}

	envy.Set("test", "test")

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "No error",
			args: args{
				param: "test",
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mustGet(tt.args.param); got != tt.want {
				t.Errorf("mustGet() = %v, want %v", got, tt.want)
			}
		})
	}
}
