package utils

import (
	"path"
	"path/filepath"
	"sort"
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
		logrus.Error(err.Error())
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
// same period as at. A listing error is logged and treated as "not covered", so
// a missing or unreadable tier still receives a backup.
func (u *Utils) CoversPeriod(dirPath string, at time.Time, period Period) bool {
	dirs, err := u.Driver.ListDirs(dirPath)
	if err != nil {
		logrus.Error(err.Error())
	}

	return u.coversPeriod(dirs, at, period)
}

func (u *Utils) coversPeriod(dirs []string, at time.Time, period Period) bool {
	for _, dir := range dirs {
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

func (u *Utils) GetPaths(targetPath string) (string, string, string, string) {
	daily := path.Join(targetPath, "daily")
	u.Driver.Mkdir(daily)
	weekly := path.Join(targetPath, "weekly")
	u.Driver.Mkdir(weekly)
	monthly := path.Join(targetPath, "monthly")
	u.Driver.Mkdir(monthly)
	yearly := path.Join(targetPath, "yearly")
	u.Driver.Mkdir(yearly)

	return daily, weekly, monthly, yearly
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

func (u *Utils) GetOldestN(list []string, cnt int) []string {
	if len(list) < cnt {
		return []string{}
	}

	sort.Strings(list)
	return list[0:cnt]
}
