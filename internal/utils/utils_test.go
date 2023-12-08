package utils

import (
	"reflect"
	"testing"
	"time"
)

func TestUtils_getDeleteDirs(t *testing.T) {
	df := "2006-01-02_15-04-05"
	now := time.Now()

	type args struct {
		dirs   []string
		cutoff time.Time
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "empty",
			args: args{
				dirs:   []string{},
				cutoff: now,
			},
			want: []string{},
		},
		{
			name: "one-passed",
			args: args{
				dirs:   []string{"2020-01-01_00-00-00"},
				cutoff: now,
			},
			want: []string{"2020-01-01_00-00-00"},
		},
		{
			name: "one-future",
			args: args{
				dirs:   []string{now.AddDate(1, 0, 0).Format(df)},
				cutoff: now,
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Utils{
				Driver:     nil,
				DateFormat: df,
			}
			if got := u.getDeleteDirs(tt.args.dirs, tt.args.cutoff); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Utils.getDeleteDirs() = %v, want %v", got, tt.want)
			}
		})
	}
}
