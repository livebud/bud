package socket

import (
	"net"
	"os"
	"strconv"
	"syscall"
)

// Load the listener from a passed in file or start a new listener
func Load(path string) (net.Listener, error) {
	files := loadFiles()
	if len(files) == 0 {
		return listen(path)
	}
	file := files[0]
	ln, err := net.FileListener(file)
	if err != nil {
		return nil, err
	}
	file.Close()
	return ln, nil
}

const (
	// listenFdsStart corresponds to `SD_LISTEN_FDS_START`.
	listenFdsStart = 3
)

// Load files if the LISTEN_FDS environment variable is set.
// This is the same environment variable that systemd uses to support socket
// activation.
//
// See:
// - https://man.archlinux.org/man/sd_listen_fds.3.en
// - https://mgdm.net/weblog/systemd-socket-activation/
// - https://vincent.bernat.ch/en/blog/2018-systemd-golang-socket-activation
func loadFiles() (files []*os.File) {
	nfds, err := strconv.Atoi(os.Getenv("LISTEN_FDS"))
	if err != nil || nfds == 0 {
		return nil
	}
	files = make([]*os.File, 0, nfds)
	for fd := listenFdsStart; fd < listenFdsStart+nfds; fd++ {
		syscall.CloseOnExec(fd)
		name := "LISTEN_FD_" + strconv.Itoa(fd)
		files = append(files, os.NewFile(uintptr(fd), name))
	}
	return files
}
