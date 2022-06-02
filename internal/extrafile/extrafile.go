package extrafile

import (
	"os"
	"strconv"
	"syscall"
)

// prepareEnv prepares the environment variables so the file descriptors can be
// recovered in the subprocess.
//
// This is also a half-hearted attempt to support systemd out of the box.
// https://github.com/coreos/go-systemd/blob/main/activation/files_unix.go
// TODO: test systemd support
func prepareEnv(prefix string, offset int, files ...*os.File) []string {
	if len(files) == 0 {
		return nil
	}
	return []string{
		prefix + "_FDS_START=" + strconv.Itoa(offset),
		prefix + "_FDS=" + strconv.Itoa(len(files)),
	}
}

// Inject files and environment for a subprocess
func Inject(extras *[]*os.File, env *[]string, prefix string, files ...*os.File) {
	if len(files) == 0 {
		return
	}
	offset := len(*extras)
	environ := prepareEnv(prefix, offset, files...)
	*extras = append(*extras, files...)
	*env = append(*env, environ...)
}

// Forward an existing prefix into a subprocesses ExtraFiles and Env
// parameters.
func Forward(extras *[]*os.File, env *[]string, prefix string) {
	files := Load(prefix)
	if len(files) == 0 {
		return
	}
	Inject(extras, env, prefix, files...)
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
		name := prefix + "_FD_" + strconv.Itoa(fd-startAt)
		files = append(files, os.NewFile(uintptr(fd), name))
	}
	return files
}
