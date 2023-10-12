package specialfolder

import "runtime"

type SpecialFolder interface {
	Home() string
	SSH() string
	ID() string
}

type specialFolder struct {
	home string
	ssh  string
}

func NewSpecialFolder() SpecialFolder {
	if runtime.GOOS == "windows" {
		return newWindowsSpecialFolder()
	}

	return newLinuxSpecialFolder()
}
