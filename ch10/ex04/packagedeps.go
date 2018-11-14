package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type Package struct {
	ImportPath string   `json:",omitempty"`
	Deps       []string `json:",omitempty"`
}

var LeftBrace = byte(123)  // '{'
var RightBrace = byte(125) // '}'

func parseJsonOutput(out []byte) ([]*Package, error) {
	var pkgs []*Package
	var startIdx, depth int

	for i, b := range out {
		if b == LeftBrace {
			if depth == 0 {
				startIdx = i
			}
			depth++
		}
		if b == RightBrace {
			depth--
			if depth < 0 {
				return nil, fmt.Errorf("json parse: unexpected right brace")
			}
			if depth == 0 {
				var p Package
				if err := json.Unmarshal(out[startIdx:i+1], &p); err != nil {
					return nil, err
				}
				pkgs = append(pkgs, &p)
			}
		}
	}
	return pkgs, nil
}

func execGoList(targets []string) ([]*Package, error) {
	args := []string{"list", "-json"}
	for _, target := range targets {
		args = append(args, target)
	}

	out, err := exec.Command("go", args...).Output()
	if err != nil {
		return nil, fmt.Errorf("execGoList: %v", err)
	}

	pkgs, err := parseJsonOutput(out)
	if err != nil {
		return nil, err
	}
	if len(pkgs) == 0 {
		return nil, fmt.Errorf("execGoList: matched no packages")
	}
	return pkgs, err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s package [package ..]\n", os.Args[0])
		os.Exit(1)
	}

	machedPkgs, err := execGoList(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var pkgNames []string
	for _, pkg := range machedPkgs {
		pkgNames = append(pkgNames, pkg.ImportPath)
	}

	targetPkgs, err := execGoList(pkgNames)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, pkg := range targetPkgs {
		for _, pkgName := range pkg.Deps {
			fmt.Printf("%s -> %s\n", pkg.ImportPath, pkgName)
		}
	}
}
