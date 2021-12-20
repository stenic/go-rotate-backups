package drivers

import "github.com/sirupsen/logrus"

type DryRunDriver struct {
	BaseDriver
	Wrapped Driver
}

func (d *DryRunDriver) Init() error {
	return d.Wrapped.Init()
}
func (d *DryRunDriver) ListDirs(path string) ([]string, error) {
	return d.Wrapped.ListDirs(path)
}
func (d *DryRunDriver) Mkdir(path string) error {
	logrus.Warnf("No Mkdir in dry-run - %v", path)
	return nil
}
func (d *DryRunDriver) Delete(src string) error {
	logrus.Warnf("No Mkdir in dry-run - %v", src)
	return nil
}
func (d *DryRunDriver) Copy(src, dst string) (int64, error) {
	logrus.Warnf("No Mkdir in dry-run - %v to %v", src, dst)
	return 0, nil
}
func (d *DryRunDriver) SetTargetPath(path string) {
	d.Wrapped.SetTargetPath(path)
}
func (d *DryRunDriver) GetTargetPath() string {
	return d.Wrapped.GetTargetPath()
}
