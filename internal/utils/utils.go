package utils

import (
	"path"
	"sort"

	"github.com/sirupsen/logrus"
	"github.com/stenic/go-rotate-backups/internal/drivers"
)

type Utils struct {
	Driver drivers.Driver
}

func (u *Utils) CleanFolder(dirPath string, keepCount int) error {
	logrus.Debugf("Listing %s", dirPath)
	dirs, err := u.Driver.ListDirs(dirPath)
	if err != nil {
		logrus.Error(err.Error())
	}
	if len(dirs) > keepCount {
		for _, old := range u.GetOldestN(dirs, len(dirs)-keepCount) {
			path := path.Join(dirPath, old)
			logrus.Debugf("Cleaning up %s", path)
			if err := u.Driver.Delete(path); err != nil {
				return err
			}
		}
	}

	return nil
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
		if _, err := u.Driver.Copy(file, path.Join(target, path.Base(file))); err != nil {
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
