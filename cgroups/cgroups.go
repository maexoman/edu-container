package cgroups

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

const memoryCGroupRoot = "/sys/fs/cgroup/memory"

type Limits struct {
	MemoryLimitInBytes string
}

func SetLimits(containerId string, limits Limits) error {
	var err error = nil

	if limits.MemoryLimitInBytes != "" {
		err = setMemoryLimit(containerId, limits.MemoryLimitInBytes)
	}

	return err
}

func RemoveLimits(containerId string) error {
	return removeMemoryLimit(containerId)
}

func setMemoryLimit(containerId string, limitInBytes string) error {
	var err error

	cgroupName := createCGroupName(containerId)
	cgroup := path.Join(memoryCGroupRoot, cgroupName)

	// create the new cgroup folder in the synthetic fs
	err = os.MkdirAll(cgroup, 0755)
	if err != nil {
		return err
	}

	// set the actual limit
	err = ioutil.WriteFile(path.Join(cgroup, "memory.limit_in_bytes"), []byte(limitInBytes), 0755)
	if err != nil {
		return err
	}

	// add this process to cgroup
	err = ioutil.WriteFile(path.Join(cgroup, "cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0644)
	if err != nil {
		return err
	}

	fmt.Printf("setting memory limit to %s for %s\n", limitInBytes, containerId)
	return nil
}

func removeMemoryLimit(containerId string) error {
	cgroupName := createCGroupName(containerId)
	cgroup := path.Join(memoryCGroupRoot, cgroupName)

	fmt.Printf("removing memory limit for %s\n", containerId)
	return os.RemoveAll(cgroup)
}

func createCGroupName(containerId string) string {
	return fmt.Sprintf("container-%s", containerId)
}
