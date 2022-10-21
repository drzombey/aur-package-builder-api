package package_collector

import (
	"fmt"
	"os"
)

func (pc *PackageCollector) detectPackagesToUpload() ([]string, error) {
	files, err := os.ReadDir(pc.pkgDir)
	if err != nil {
		return nil, err
	}

	var fileList []string

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileList = append(fileList, file.Name())
	}
	return fileList, nil
}

func (pc *PackageCollector) StartTaskAddPackageToRepository() error {
	files, err := pc.detectPackagesToUpload()

	if err != nil {
		return err
	}

	fmt.Println(files)
	return nil
}
