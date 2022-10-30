package package_collector

type PackageCollector struct {
	pkgDir string
}

func NewPackageCollector(pkgDir string) *PackageCollector {
	return &PackageCollector{
		pkgDir: pkgDir,
	}
}
