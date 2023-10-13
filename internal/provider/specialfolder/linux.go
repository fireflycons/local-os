package specialfolder

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fireflycons/local-os/internal/provider/hasher"
)

// It is assumed that everything that is not windows follows the linux paradigm

type linuxSpecialFolder struct {
	specialFolder
}

func newLinuxSpecialFolder() *linuxSpecialFolder {
	home := os.Getenv("HOME")

	s := &linuxSpecialFolder{
		specialFolder: specialFolder{
			home: home,
			ssh:  filepath.Join(home, ".ssh"),
		},
	}

	return s
}

func (f *linuxSpecialFolder) Home() string {
	return f.home
}

func (f *linuxSpecialFolder) SSH() string {
	return f.ssh
}

func (f *linuxSpecialFolder) ID() string {
	h := hasher.NewMarvin32(2163)
	h.Sum([]byte(os.Getenv("HOSTNAME")))
	return fmt.Sprintf("%08x", h.Sum32())
}
