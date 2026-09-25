package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/stenic/go-rotate-backups/internal/drivers"
	"github.com/stenic/go-rotate-backups/internal/utils"
)

// Backing up files one run at a time with a shared --date must leave exactly
// the same tree as backing them all up in one run: same snapshots, same tier
// promotions, and rotation deleting the same (and only the old) snapshots.
func TestSequentialRunsMatchOneBatchRun(t *testing.T) {
	work := t.TempDir()
	files := []string{"./a.dump", "./b.dump", "./c.dump"}
	for _, file := range files {
		if err := os.WriteFile(filepath.Join(work, file), []byte(file), 0600); err != nil {
			t.Fatal(err)
		}
	}
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(work); err != nil {
		t.Fatal(err)
	}

	// Small retention so the run rotates on every tier within the window.
	saved := []int{keepDaily, keepWeekly, keepMonthly, keepYearly}
	keepDaily, keepWeekly, keepMonthly, keepYearly = 3, 2, 2, 1
	originalNow := now
	t.Cleanup(func() {
		keepDaily, keepWeekly, keepMonthly, keepYearly = saved[0], saved[1], saved[2], saved[3]
		now = originalNow
		os.Chdir(wd)
	})

	sequential, batch := newLocalUtils(t), newLocalUtils(t)

	// Unrelated content that rotation must never touch: a non-timestamp folder
	// inside a tier and a folder outside the tiers.
	unrelated := []string{"daily/keep-me/notes.txt", "other/data.bin"}
	for _, u := range []utils.Utils{sequential, batch} {
		for _, rel := range unrelated {
			p := filepath.Join(u.Driver.GetTargetPath(), rel)
			if err := os.MkdirAll(filepath.Dir(p), 0750); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(p, []byte("keep"), 0600); err != nil {
				t.Fatal(err)
			}
		}
	}
	run := func(u utils.Utils, files []string) {
		t.Helper()
		if err := addFunc(u, nil, files); err != nil {
			t.Fatal(err)
		}
		if err := rotateFunc(u, nil, files); err != nil {
			t.Fatal(err)
		}
	}

	// Nightly for 60 nights across a month and a year boundary.
	start := time.Date(2024, 12, 10, 23, 0, 0, 0, time.UTC)
	for day := 0; day < 60; day++ {
		now = start.AddDate(0, 0, day)
		for _, file := range files {
			run(sequential, []string{file})
		}
		run(batch, files)

		got, want := tree(t, sequential), tree(t, batch)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("%s: sequential runs left\n%v\nwant the batch run's\n%v", now.Format(DateFormat), got, want)
		}
	}

	// Guard against a vacuous pass: rotation really deleted snapshots and
	// every kept snapshot holds all files.
	final := tree(t, batch)
	for _, rel := range unrelated {
		if _, err := os.Stat(filepath.Join(sequential.Driver.GetTargetPath(), rel)); err != nil {
			t.Fatalf("rotation removed unrelated %s: %v", rel, err)
		}
	}
	if len(final) == len(unrelated) || (len(final)-len(unrelated))%len(files) != 0 {
		t.Fatalf("unexpected final tree %v", final)
	}
	var daily []string
	for _, entry := range snapshotFiles(t, filepath.Join(batch.Driver.GetTargetPath(), "daily")) {
		if _, err := time.Parse(DateFormat, entry); err == nil {
			daily = append(daily, entry)
		}
	}
	if len(daily) > keepDaily+1 {
		t.Fatalf("daily kept %d snapshots, rotation did not run: %v", len(daily), daily)
	}
}

func newLocalUtils(t *testing.T) utils.Utils {
	t.Helper()
	driver := &drivers.LocalDriver{}
	driver.SetTargetPath(t.TempDir())
	return utils.Utils{Driver: driver, DateFormat: DateFormat}
}

// tree lists every file under the target, relative to it.
func tree(t *testing.T, u utils.Utils) []string {
	t.Helper()
	root := u.Driver.GetTargetPath()
	var files []string
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			rel, _ := filepath.Rel(root, p)
			files = append(files, rel)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(files)
	return files
}
