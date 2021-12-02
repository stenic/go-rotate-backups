package drivers

import "fmt"

type Driver interface {
	// Implementation
	Init() error
	ListDirs(path string) ([]string, error)
	Mkdir(path string) error
	Delete(src string) error
	Copy(src, dst string) (int64, error)

	// baseDriver
	SetTargetPath(path string)
	GetTargetPath() string
}

type BaseDriver struct {
	TargetPath string
}

func (d *BaseDriver) SetTargetPath(targetPath string) {
	d.TargetPath = targetPath
}
func (d *BaseDriver) GetTargetPath() string {
	return d.TargetPath
}

var driverRegistry = map[string]Driver{}

func GetDriver(name string) (Driver, error) {
	if val, ok := driverRegistry[name]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("Driver %s could not be found", name)
}

func AddDriver(name string, driver Driver) {
	driverRegistry[name] = driver
}
