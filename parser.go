package mane

import (
	"fmt"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/packages"
)

type Parser struct {
	PackageLoader *PackageLoader
}

func NewParser() *Parser {
	packageLoader := NewPackageLoader()
	return &Parser{
		PackageLoader: &packageLoader,
	}
}

func (p *Parser) ParseFile(filePath string) ([]Interface, error) {
	loaded, err := p.PackageLoader.LoadFile(filePath)
	if err != nil {
		return nil, err
	}
	return p.parsePackage(loaded)
}

func (p *Parser) ParsePackage(packagePath string) ([]Interface, error) {
	loaded, err := p.PackageLoader.Load(packagePath)
	if err != nil {
		return nil, err
	}
	return p.parsePackage(loaded)
}

func (p *Parser) parsePackage(pkg *packages.Package) ([]Interface, error) {
	var neededInterfaces []Interface
	for _, fileDecl := range pkg.Syntax {
		for _, decl := range fileDecl.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
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
				iface, innerOk := typeSpec.Type.(*ast.InterfaceType)
				if !innerOk {
					continue
				}
				methods, err := p.iface(pkg, iface)
				if err != nil {
					return nil, fmt.Errorf("error parsing interface %s: %w", typeSpec.Name.Name, err)
				}
				neededInterfaces = append(neededInterfaces, Interface{
					Name:       typeSpec.Name.Name,
					Package:    pkg,
					TypeParams: typeSpec.TypeParams,
					Methods:    methods,
				})
			}
		}
	}
	return neededInterfaces, nil
}

func (p *Parser) iface(pkg *packages.Package, decl *ast.InterfaceType) ([]Method, error) {
	var neededMethods []Method
	for _, method := range decl.Methods.List {
		var foundMethods []Method
		var err error
		switch methodType := method.Type.(type) {
		case *ast.FuncType:
			foundMethods = []Method{p.funcDel(pkg, method.Names[0].Name, methodType)}
		case *ast.SelectorExpr:
			// Embedded interface referencing other package
			foundMethods, err = p.foreignPackageIface(methodType)
		case *ast.Ident:
			// Embedded interface referencing same package
			foundMethods, err = p.localPackageIface(pkg, methodType)
		}
		if err != nil {
			return nil, err
		}
		neededMethods = append(neededMethods, foundMethods...)
	}
	return neededMethods, nil
}

func (p *Parser) localPackageIface(pkg *packages.Package, decl *ast.Ident) ([]Method, error) {
	depIface := FindInterface(pkg, decl.Name)
	if depIface == nil {
		return nil, fmt.Errorf("interface %s not found in package %s", decl.Name, pkg.Name)
	}
	depMethods, err := p.iface(pkg, depIface.Type.(*ast.InterfaceType))
	return depMethods, err
}

func (p *Parser) foreignPackageIface(decl *ast.SelectorExpr) ([]Method, error) {
	convX, convXOk := decl.X.(*ast.Ident)
	if !convXOk {
		return nil, fmt.Errorf("expected ident, got %T", decl.X)
	}
	dependentPkg, err := p.PackageLoader.Load(convX.Name)
	if err != nil {
		return nil, fmt.Errorf("error loading package %s: %w", convX.Name, err)
	}

	depIface := FindInterface(dependentPkg, decl.Sel.Name)
	if depIface == nil {
		return nil, fmt.Errorf("interface %s not found in package %s", decl.Sel.Name, convX.Name)
	}
	depMethods, err := p.iface(dependentPkg, depIface.Type.(*ast.InterfaceType))
	return depMethods, err
}

func (p *Parser) funcDel(pkg *packages.Package, name string, decl *ast.FuncType) Method {
	return Method{
		Name:       name,
		Package:    pkg,
		TypeParams: decl.TypeParams,
		Params:     decl.Params,
		Results:    decl.Results,
	}
}
