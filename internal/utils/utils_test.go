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

func TestUtils_coversPeriod(t *testing.T) {
	df := "2006-01-02_15-04-05"

	type args struct {
		dirs   []string
		at     time.Time
		period Period
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty tier is never covered",
			args: args{
				dirs:   []string{},
				at:     time.Date(2024, 3, 14, 9, 0, 0, 0, time.UTC),
				period: PeriodMonth,
			},
			want: false,
		},
		{
			name: "same month, different day",
			args: args{
				dirs:   []string{"2024-03-01_02-00-00"},
				at:     time.Date(2024, 3, 14, 9, 0, 0, 0, time.UTC),
				period: PeriodMonth,
			},
			want: true,
		},
		{
			name: "previous month does not cover",
			args: args{
				dirs:   []string{"2024-02-28_02-00-00"},
				at:     time.Date(2024, 3, 14, 9, 0, 0, 0, time.UTC),
				period: PeriodMonth,
			},
			want: false,
		},
		{
			name: "same month one year earlier does not cover",
			args: args{
				dirs:   []string{"2023-03-14_02-00-00"},
				at:     time.Date(2024, 3, 14, 9, 0, 0, 0, time.UTC),
				period: PeriodMonth,
			},
			want: false,
		},
		{
			name: "same ISO week across a month boundary",
			args: args{
				// Both fall in ISO week 2024-W18.
				dirs:   []string{"2024-04-29_02-00-00"},
				at:     time.Date(2024, 5, 3, 9, 0, 0, 0, time.UTC),
				period: PeriodWeek,
			},
			want: true,
		},
		{
			name: "previous ISO week does not cover",
			args: args{
				dirs:   []string{"2024-04-28_02-00-00"},
				at:     time.Date(2024, 5, 3, 9, 0, 0, 0, time.UTC),
				period: PeriodWeek,
			},
			want: false,
		},
		{
			name: "same year",
			args: args{
				dirs:   []string{"2024-01-01_02-00-00"},
				at:     time.Date(2024, 12, 31, 9, 0, 0, 0, time.UTC),
				period: PeriodYear,
			},
			want: true,
		},
		{
			name: "previous year does not cover",
			args: args{
				dirs:   []string{"2023-12-31_02-00-00"},
				at:     time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
				period: PeriodYear,
			},
			want: false,
		},
		{
			name: "unparsable entries are skipped",
			args: args{
				dirs:   []string{"not-a-date", "2024-03-02_02-00-00"},
				at:     time.Date(2024, 3, 14, 9, 0, 0, 0, time.UTC),
				period: PeriodMonth,
			},
			want: true,
		},
		{
			name: "only unparsable entries",
			args: args{
				dirs:   []string{"not-a-date"},
				at:     time.Date(2024, 3, 14, 9, 0, 0, 0, time.UTC),
				period: PeriodMonth,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Utils{
				Driver:     nil,
				DateFormat: df,
			}
			if got := u.coversPeriod(tt.args.dirs, tt.args.at, tt.args.period); got != tt.want {
				t.Errorf("Utils.coversPeriod() = %v, want %v", got, tt.want)
			}
		})
	}
}
