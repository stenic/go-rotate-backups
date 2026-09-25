package utils

import (
	"path"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stenic/go-rotate-backups/internal/drivers"
)

type Utils struct {
	Driver     drivers.Driver
	DateFormat string
}

func (u *Utils) CleanFolder(dirPath string, cutoff time.Time) error {
	dirs, err := u.Driver.ListDirs(dirPath)
	if err != nil {
		return err
	}
	logrus.Infof("Listing %s: %v", dirPath, dirs)
	for _, dir := range u.getDeleteDirs(dirs, cutoff) {
		logrus.Debugf("Cleaning up %s", dir)
		if err := u.Driver.Delete(filepath.Join(dirPath, dir)); err != nil {
			return err
		}
	}
	return nil
}

func (u *Utils) getDeleteDirs(dirs []string, cutoff time.Time) []string {
	deleteList := []string{}
	for _, dir := range dirs {
		dirTime, err := time.Parse(u.DateFormat, dir)
		if err != nil {
			logrus.Warnf("Could not parse %s as date: %v", dir, err)
			continue
		}
		if dirTime.Before(cutoff) {
			deleteList = append(deleteList, dir)
		}
	}
	return deleteList

}

// Period is the span of time a backup tier covers. A tier holds at most one
// backup per period, no matter how often the command runs.
type Period int

const (
	PeriodWeek Period = iota
	PeriodMonth
	PeriodYear
)

// CoversPeriod reports whether dirPath already holds a backup that falls in the
// same period as at.
func (u *Utils) CoversPeriod(dirPath string, at time.Time, period Period) (bool, error) {
	dirs, err := u.Driver.ListDirs(dirPath)
	if err != nil {
		return false, err
	}

	return u.coversPeriod(dirs, at, period), nil
}

func (u *Utils) coversPeriod(dirs []string, at time.Time, period Period) bool {
	current := at.Format(u.DateFormat)
	for _, dir := range dirs {
		// The snapshot being written right now does not cover its own period:
		// a later run with the same timestamp still has to add its files.
		if dir == current {
			continue
		}
		dirTime, err := time.Parse(u.DateFormat, dir)
		if err != nil {
			logrus.Warnf("Could not parse %s as date: %v", dir, err)
			continue
		}
		if samePeriod(dirTime, at, period) {
			return true
		}
	}

	return false
}

func samePeriod(a, b time.Time, period Period) bool {
	switch period {
	case PeriodWeek:
		aYear, aWeek := a.ISOWeek()
		bYear, bWeek := b.ISOWeek()
		return aYear == bYear && aWeek == bWeek
	case PeriodMonth:
		return a.Year() == b.Year() && a.Month() == b.Month()
	case PeriodYear:
		return a.Year() == b.Year()
	}

	return false
}

func (u *Utils) GetPaths(targetPath string) (string, string, string, string, error) {
	daily := path.Join(targetPath, "daily")
	weekly := path.Join(targetPath, "weekly")
	monthly := path.Join(targetPath, "monthly")
	yearly := path.Join(targetPath, "yearly")
	for _, dir := range []string{daily, weekly, monthly, yearly} {
		if err := u.Driver.Mkdir(dir); err != nil {
			return "", "", "", "", err
		}
	}

	return daily, weekly, monthly, yearly, nil
}

// HasEntry reports whether dirPath already holds a snapshot for exactly at.
func (u *Utils) HasEntry(dirPath string, at time.Time) (bool, error) {
	dirs, err := u.Driver.ListDirs(dirPath)
	if err != nil {
		return false, err
	}
	current := at.Format(u.DateFormat)
	for _, dir := range dirs {
		if dir == current {
			return true, nil
		}
	}

	return false, nil
}

// RemovePartial undoes a failed CopyFiles. In a snapshot shared with other runs
// only this run's files are removed; otherwise the whole target goes.
func (u *Utils) RemovePartial(files []string, target string, shared bool) error {
	if !shared {
		return u.Driver.Delete(target)
	}
	for _, file := range files {
		if err := u.Driver.Delete(path.Join(target, file)); err != nil {
			return err
		}
	}

	return nil
}

func (u *Utils) CopyFiles(files []string, target string) error {
	if err := u.Driver.Mkdir(target); err != nil {
		return err
	}
	for _, file := range files {
		if _, err := u.Driver.Copy(file, path.Join(target, file)); err != nil {
			return err
		}
	}

	return nil
}
