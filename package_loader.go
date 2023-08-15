package mane

import (
	"errors"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/packages"
)

type PackageLoader struct {
	locks        LockMap
	cache        map[string]*packages.Package
	PackagesConf packages.Config
}

func NewPackageLoader() PackageLoader {
	return PackageLoader{
		locks: NewLockMap(),
		cache: make(map[string]*packages.Package),
		PackagesConf: packages.Config{
			Mode: packages.NeedTypes |
				packages.NeedTypesSizes |
				packages.NeedSyntax |
				packages.NeedTypesInfo |
				packages.NeedImports |
				packages.NeedName |
				packages.NeedFiles |
				packages.NeedCompiledGoFiles,
		},
	}
}

func (pc *PackageLoader) LoadFile(filePath string) (*packages.Package, error) {
	return pc.Load("file=" + filePath)
}

func (pc *PackageLoader) Load(packagePath string) (*packages.Package, error) {
	pc.locks.Lock(packagePath)
	defer pc.locks.Unlock(packagePath)

	if existing, ok := pc.cache[packagePath]; ok {
		return existing, nil
	}
	found, err := packages.Load(&pc.PackagesConf, packagePath)
	if err != nil {
		return nil, err
	}
	if len(found) == 0 {
		return nil, errors.New("no packages found")
	}
	pc.cache[packagePath] = found[0]
	return found[0], nil
}

func (pc *PackageLoader) Get(packagePath string) (*packages.Package, error) {
	pc.locks.Lock(packagePath)
	defer pc.locks.Unlock(packagePath)

	if pkg, ok := pc.cache[packagePath]; ok {
		return pkg, nil
	}
	return nil, errors.New("package not found")
}

func FindInterface(pkg *packages.Package, interfaceName string) *ast.TypeSpec {
	for _, fileDecl := range pkg.Syntax {
		for _, innerDecl := range fileDecl.Decls {
			genDecl, ok := innerDecl.(*ast.GenDecl)
			if !ok {
				continue
			}
			if genDecl.Tok != token.TYPE {
				continue
			}
			for _, spec := range genDecl.Specs {
				typeSpec, innerOk := spec.(*ast.TypeSpec)
				if !innerOk {
					continue
				}
				if typeSpec.Name.Name == interfaceName {
					return typeSpec
				}
			}
		}
	}
	return nil
}
