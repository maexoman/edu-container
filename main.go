package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/maexoman/edu-container/cgroups"
	"github.com/maexoman/edu-container/containerfs"
	"github.com/maexoman/edu-container/utils"
)

var root = "/path/to/store/containers"
var layers = []string{"/path/to/layer2", "/path/to/layer1"}

func main() {
	switch os.Args[1] {
	case "run":
		forkSelfWithNamespaces(root)
	case "runCommand":
		containerId := os.Args[2]
		command := os.Args[3]
		arguments := os.Args[4:]
		runCommand(root, containerId, command, arguments)
	default:
		panic("invalid command")
	}
}

func forkSelfWithNamespaces(rootPath string) {
	var err error

	containerId := utils.NewContainerId()

	// create the rootfs for the container
	err = containerfs.Mount(rootPath, containerId, layers)
	utils.Must(err)

	// fork this process with new namespaces
	err = createProcessWithNewNamespaces(containerId)
	utils.Must(err)

	// cleanup after the container command has been run
	err = containerfs.Unmount(rootPath, containerId)
	utils.Must(err)
}

func createProcessWithNewNamespaces(containerId string) error {
	var err error

	// starts itself to be able to set the cgroups and container hostname
	cmd := exec.Command("/proc/self/exe", append([]string{"runCommand", containerId}, os.Args[2:]...)...)

	// add wanted namespaces
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}

	// default wireing
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// start the process
	err = cmd.Start()
	if err != nil {
		return err
	}

	// wait for the child to be able to cleanup once the "container" is done
	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

func runCommand(rootPath, containerId string, command string, arguments []string) {
	var err error

	// set container limit(s)
	err = cgroups.SetLimits(containerId, cgroups.Limits{
		MemoryLimitInBytes: "1000000",
	})
	utils.Must(err)

	// mount the root for the container
	err = mountRootFs(rootPath, containerId)
	utils.Must(err)

	// set the container name
	err = setHostname(containerId)
	utils.Must(err)

	// mount new "process hub" for the ps command to read available processes
	err = disconnectProcDirectory()
	utils.Must(err)

	// prepare the user provided command
	cmd := exec.Command(command, arguments...)

	// default wireing again
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// finally run the command :)
	fmt.Printf("Run: %s as %s(pid: %d)\n", command, containerId, os.Getpid())
	err = cmd.Run()
	utils.Must(err)

	// Cleanup the CGroups
	err = cgroups.RemoveLimits(containerId)
	utils.Must(err)

	// cleanup the "process hub"
	err = cleanupProcDirectory()
	utils.Must(err)
}

func mountRootFs(rootPath, containerId string) error {
	var err error

	rootFsPath := containerfs.GetRootFsPath(rootPath, containerId)

	// set the chroot - this is not safe (see man pages) but itll surfice for this educational project
	err = syscall.Chroot(rootFsPath)
	if err != nil {
		return err
	}

	// needed because otherwhiese / is "undefined"
	err = syscall.Chdir("/")
	if err != nil {
		return err
	}

	return nil
}

func setHostname(containerId string) error {
	return syscall.Sethostname([]byte(containerId))
}

func disconnectProcDirectory() error {
	// this mount is needed by ps to read running processes
	return syscall.Mount("proc", "proc", "proc", 0, "")
}

func cleanupProcDirectory() error {
	return syscall.Unmount("/proc", 0)
}
