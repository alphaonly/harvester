package mainanalyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const (
	mainPackageName    = "main" // Name of package in which it needs to find a call
	firstFunctionName  = "main" // Name of the function in which it needs to find a call
	osPackageName      = "os"   // Name of the package that uses to call of searched function
	secondFunctionName = "Exit" // Name of the function that is sought
)

// run - a function that describes how mainOsExit analyzes a project
func run(pass *analysis.Pass) (interface{}, error) {

	for _, file := range pass.Files {

		//Check the exact file package name is "main"
		if file.Name.String() != mainPackageName {
			continue
		}

		// A traverse on AST tree for search for a node
		ast.Inspect(file, func(n ast.Node) bool {

			//Check if a node is a function declaration
			if functionDeclaration, ok := n.(*ast.FuncDecl); ok {
				//Check if the function name is "main" and there is some body inside
				if functionDeclaration.Name.Name == firstFunctionName &&
					functionDeclaration.Body != nil {
					//Make a traverse on the function statements
					for _, statement := range functionDeclaration.Body.List {
						//Check the statement is an expression
						if expressionStatement, ok := statement.(*ast.ExprStmt); ok {
							//Check the expression statement is a call expression with params
							if callExpression, ok := expressionStatement.X.(*ast.CallExpr); ok {
								//Check whether the expression statement is Selector expression like a.b()
								if selector, ok := callExpression.Fun.(*ast.SelectorExpr); ok {
									//Check package identifier in expression name is "os"
									if identifier, ok := selector.X.(*ast.Ident); ok && identifier.Name == osPackageName {
										//Check function name is "Exit"
										if selector.Sel.Name == secondFunctionName {
											pass.Reportf(selector.Sel.NamePos, "os.Exit is used")
										}
									}
								}
							}
						}
					}
				}
			}

			return true

		})

	}

	return nil, nil
}

// MainOsExit - an Analyzer structure that describes additional Analyzer for os.Exit call occurrence check in main function of a project
var MainOsExit = &analysis.Analyzer{
	Name: "OsExitChecker",
	Doc:  "Checks whether os.Exit function call occurred in main function in file with package main",
	Run:  run,
}
