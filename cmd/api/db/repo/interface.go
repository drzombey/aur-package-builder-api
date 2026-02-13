package repository

import "github.com/drzombey/aur-package-builder-api/pkg/aur"

// PackageRepository defines the interface for package repository operations
type PackageRepository interface {
	GetAurPackageInfoByName(name string) (*aur.Package, error)
	GetPackageFromAur(name string) ([]aur.Package, error)
	GetAlreadyBuildPackageByAurIdAndVersion(id int64, version string) (*aur.Package, error)
	GetAlreadyBuildPackages() (*[]aur.Package, error)
	AddAurPackage(model aur.Package) error
	DeleteAurPackage(model aur.Package) error
}