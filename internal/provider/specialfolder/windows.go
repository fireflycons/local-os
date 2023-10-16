package specialfolder

import (
	"fmt"
	"hash/fnv"
	"os"
	"os/user"
	"path/filepath"
)

type windowsSpecialFolder struct {
	specialFolder
}

func newWindowsSpecialFolder() *windowsSpecialFolder {
	var home string
	u, err := user.Current()
	if err == nil {
		home = u.HomeDir
	} else {
		home = os.Getenv(("USERPROFILE"))
	}

	s := &windowsSpecialFolder{
		specialFolder: specialFolder{
			home: home,
			ssh:  filepath.Join(home, ".ssh"),
		},
	}

	return s
}

func (f *windowsSpecialFolder) Home() string {
	return f.home
}

func (f *windowsSpecialFolder) SSH() string {
	return f.ssh
}

func (f *windowsSpecialFolder) ID() string {
	h := fnv.New32a()
	h.Sum([]byte(os.Getenv("COMPUTERNAME")))
	return fmt.Sprintf("%08x", h.Sum32())
}
