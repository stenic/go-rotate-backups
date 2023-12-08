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
