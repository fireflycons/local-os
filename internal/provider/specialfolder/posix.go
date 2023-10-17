package specialfolder

import (
	"fmt"
	"hash/fnv"
	"os"
	"os/user"
	"path/filepath"
)

// It is assumed that everything that is not windows follows the posix paradigm

type posixSpecialFolder struct {
	specialFolder
}

func newPosixSpecialFolder() *posixSpecialFolder {
	var home string
	u, err := user.Current()
	if err == nil {
		home = u.HomeDir
	} else {
		home = os.Getenv(("HOME"))
	}

	s := &posixSpecialFolder{
		specialFolder: specialFolder{
			home: home,
			ssh:  filepath.Join(home, ".ssh"),
		},
	}

	return s
}

func (f *posixSpecialFolder) Home() string {
	return f.home
}

func (f *posixSpecialFolder) SSH() string {
	return f.ssh
}

func (f *posixSpecialFolder) ID() string {
	h := fnv.New32a()
	h.Sum([]byte(os.Getenv("HOSTNAME")))
	return fmt.Sprintf("%08x", h.Sum32())
}
