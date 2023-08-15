package mane

import (
	"go/ast"
	"golang.org/x/tools/go/packages"
)

type Method struct {
	Name       string            // Method name
	Package    *packages.Package // Package context of where this method is defined
	TypeParams *ast.FieldList    // type parameters; or nil
	Params     *ast.FieldList    // (incoming) parameters; non-nil
	Results    *ast.FieldList    // (outgoing) results; or nil
}

type Interface struct {
	Name       string            // Interface name
	Package    *packages.Package // Package context of where this interface is defined
	TypeParams *ast.FieldList    // type parameters; or nil
	Methods    []Method          // Interface methods
}
