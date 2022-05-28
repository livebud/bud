package extrafile

import (
	"os"
	"strconv"
	"syscall"
)

type File interface {
	File() (*os.File, error)
}

// Prepare is a low-level function for preparing files to be passed through a
// subprocess via "os/exec".*Cmd.ExtraFiles. The prefix must be the same on both
// sides. Use the offset to prepare multiple extra files to be passed through.
// The offset should start at 0.
func Prepare(prefix string, offset int, files ...File) ([]*os.File, []string, error) {
	lenFiles := len(files)
	osFiles := make([]*os.File, lenFiles)
	for i, file := range files {
		osFile, err := file.File()
		if err != nil {
			return nil, nil, err
		}
		osFiles[i] = osFile
	}
	env := prepareEnv(prefix, offset, osFiles)
	return osFiles, env, nil
}

// Half-hearted attempt to support Systemd out of the box.
// https://github.com/coreos/go-systemd/blob/main/activation/files_unix.go
// TODO: test this assumption
func prepareEnv(prefix string, offset int, files []*os.File) []string {
	if len(files) == 0 {
		return nil
	}
	return []string{
		prefix + "_FDS_START=" + strconv.Itoa(offset),
		prefix + "_FDS=" + strconv.Itoa(len(files)),
	}
}

// Loading extra file descriptors should start at 3 because the first 3 are:
// stdin (0), stdout (1) and stderr (2).
const startAt = 3

// Load the passed in files using the prefix from within the subprocess. If
// there are no files, it will return an empty list of files. This process is
// also known as socket activation.
//
// See the following references for more details:
// - https://man.archlinux.org/man/sd_listen_fds.3.en
// - https://mgdm.net/weblog/systemd-socket-activation/
// - https://vincent.bernat.ch/en/blog/2018-systemd-golang-socket-activation
func Load(prefix string) []*os.File {
	len, err := strconv.Atoi(os.Getenv(prefix + "_FDS"))
	if err != nil || len == 0 {
		return nil
	}
	offset, err := strconv.Atoi(os.Getenv(prefix + "_FDS_START"))
	if err != nil {
		return nil
	}
	var files []*os.File
	for fd := startAt + offset; fd < startAt+offset+len; fd++ {
		syscall.CloseOnExec(fd)
		name := prefix + "_FD_" + strconv.Itoa(fd)
		files = append(files, os.NewFile(uintptr(fd), name))
	}
	return files
}
