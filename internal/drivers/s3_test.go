package drivers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_s3ItemsToFolders(t *testing.T) {
	type args struct {
		path  string
		items []string
	}
	tests := []struct {
		name    string
		args    args
		wantRes []string
	}{
		{
			name: "should return empty list",
			args: args{
				path:  "test",
				items: []string{},
			},
			wantRes: []string{},
		},
		{
			name: "should return list with one item",
			args: args{
				path:  "test",
				items: []string{"test/1"},
			},
			wantRes: []string{"1"},
		},
		{
			name: "should return list with two items",
			args: args{
				path:  "test",
				items: []string{"test/1", "test/2"},
			},
			wantRes: []string{"1", "2"},
		},
		{
			name: "trim subdirs",
			args: args{
				path:  "test",
				items: []string{"test/1/2/3"},
			},
			wantRes: []string{"1"},
		},
		{
			name: "remove duplicate subdirs",
			args: args{
				path:  "test",
				items: []string{"test/1/2/3", "test/1/2/4"},
			},
			wantRes: []string{"1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes := s3ItemsToFolders(tt.args.path, tt.args.items)
			assert.Equal(t, tt.wantRes, gotRes)
		})
	}
}
