package containerfs

import (
	"fmt"
	"os"
	"path"
	"strings"
	"syscall"
)

type containerPaths struct {
	rootPath string

	viewPath string
	workPath string
	dataPath string
}

func GetRootFsPath(rootPath, containerId string) string {
	return path.Join(rootPath, containerId, "view")
}

func Mount(rootPath, containerId string, layers []string) error {
	paths := createContainerPaths(rootPath, containerId)

	var err error

	fmt.Printf("removing potential leftovers..\n")
	err = removeDirectories(paths)
	if err != nil {
		return err
	}

	fmt.Printf("create new folder structure..\n")
	err = createDirectories(paths)
	if err != nil {
		return err
	}

	fmt.Printf("mount the overlay fs to %s..\n", paths.viewPath)
	err = mount(paths, layers)
	if err != nil {
		return err
	}

	return nil
}

func Unmount(rootPath, containerId string) error {
	paths := createContainerPaths(rootPath, containerId)

	var err error

	fmt.Printf("mount the overlay fs at %s..\n", paths.viewPath)
	err = syscall.Unmount(paths.viewPath, 0)
	if err != nil {
		return err
	}

	fmt.Printf("removing folder structure..\n")
	err = removeDirectories(paths)
	if err != nil {
		return err
	}

	return nil
}

func createContainerPaths(rootPath, containerId string) containerPaths {
	containerRoot := path.Join(rootPath, containerId)

	return containerPaths{
		rootPath: containerRoot,

		viewPath: path.Join(containerRoot, "view"),
		workPath: path.Join(containerRoot, "work"),
		dataPath: path.Join(containerRoot, "data"),
	}
}

func removeDirectories(paths containerPaths) error {
	return os.RemoveAll(paths.rootPath)
}

func createDirectories(paths containerPaths) error {
	var err error

	err = os.Mkdir(paths.rootPath, 0755)
	if err != nil {
		return err
	}

	err = os.Mkdir(paths.viewPath, 0755)
	if err != nil {
		return err
	}

	err = os.Mkdir(paths.dataPath, 0755)
	if err != nil {
		return err
	}

	err = os.Mkdir(paths.workPath, 0755)
	if err != nil {
		return err
	}

	return nil
}

func mount(paths containerPaths, layers []string) error {
	layersArgument := fmt.Sprintf(
		"lowerdir=%s,upperdir=%s,workdir=%s",
		strings.Join(layers, ":"),
		paths.dataPath,
		paths.workPath,
	)
	return syscall.Mount("overlay", paths.viewPath, "overlay", 0, layersArgument)
}
