package main

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/stenic/go-rotate-backups/internal/drivers"
	"github.com/stenic/go-rotate-backups/internal/utils"
)

func TestAddFuncFillsEachTierOncePerPeriod(t *testing.T) {
	target := t.TempDir()
	source := filepath.Join(t.TempDir(), "backup.sql")
	if err := os.WriteFile(source, []byte("dump"), 0600); err != nil {
		t.Fatal(err)
	}

	driver := &drivers.LocalDriver{}
	driver.SetTargetPath(target)
	util := utils.Utils{Driver: driver, DateFormat: DateFormat}

	originalNow := now
	t.Cleanup(func() { now = originalNow })
	for _, at := range []time.Time{
		time.Date(2024, 4, 29, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
	} {
		now = at
		if err := addFunc(util, nil, []string{source}); err != nil {
			t.Fatal(err)
		}
	}

	want := map[string][]string{
		"daily":   {"2024-04-29_00-00-00", "2024-05-01_00-00-00"},
		"weekly":  {"2024-04-29_00-00-00"},
		"monthly": {"2024-04-29_00-00-00", "2024-05-01_00-00-00"},
		"yearly":  {"2024-04-29_00-00-00"},
	}
	for tier, wantEntries := range want {
		entries, err := os.ReadDir(filepath.Join(target, tier))
		if err != nil {
			t.Fatal(err)
		}
		got := make([]string, len(entries))
		for i, entry := range entries {
			got[i] = entry.Name()
		}
		if !reflect.DeepEqual(got, wantEntries) {
			t.Errorf("%s entries = %v, want %v", tier, got, wantEntries)
		}
	}
}

func TestAddFuncReturnsListingError(t *testing.T) {
	wantErr := errors.New("list failed")
	driver := &listErrorDriver{err: wantErr}
	driver.SetTargetPath(t.TempDir())
	util := utils.Utils{Driver: driver, DateFormat: DateFormat}

	if err := addFunc(util, nil, nil); !errors.Is(err, wantErr) {
		t.Fatalf("addFunc() error = %v, want %v", err, wantErr)
	}
	if driver.copyCalled {
		t.Fatal("addFunc() copied files after a listing error")
	}
}

func TestAddFuncRemovesPartialBackup(t *testing.T) {
	target := t.TempDir()
	driver := &drivers.LocalDriver{}
	driver.SetTargetPath(target)
	util := utils.Utils{Driver: driver, DateFormat: DateFormat}

	originalNow := now
	now = time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	t.Cleanup(func() { now = originalNow })

	err := addFunc(util, nil, []string{filepath.Join(t.TempDir(), "missing.dump")})
	if err == nil {
		t.Fatal("addFunc() succeeded with a missing source file")
	}

	entries, err := os.ReadDir(filepath.Join(target, "daily"))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("daily contains partial backup entries: %v", entries)
	}
}

type listErrorDriver struct {
	drivers.BaseDriver
	err        error
	copyCalled bool
}

func (d *listErrorDriver) Init() error                       { return nil }
func (d *listErrorDriver) ListDirs(string) ([]string, error) { return nil, d.err }
func (d *listErrorDriver) Mkdir(string) error                { return nil }
func (d *listErrorDriver) Delete(string) error               { return nil }
func (d *listErrorDriver) Copy(string, string) (int64, error) {
	d.copyCalled = true
	return 0, nil
}
