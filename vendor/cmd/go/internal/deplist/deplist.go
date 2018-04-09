// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package deplist implements the ``go deplist'' command.
package deplist

import (
	"cmd/go/internal/base"
	"cmd/go/internal/load"
	"cmd/go/internal/work"
)

var CmdDeplist = &base.Command{
	Run:         runDeplist,
	CustomFlags: true,
	UsageLine:   "deplist [-n] [-x] [build flags] [deplist flags] [packages]",
	Short:       "provide JSON package build information for dependencies of packages",
	Long: `
Deplist provide JSON package build information for dependencies of packages
	`,
}

func runDeplist(cmd *base.Command, args []string) {
	deplistFlags, pkgArgs := deplistFlags(args)

	work.BuildInit()

	test := false
	build := false

	for _, v := range deplistFlags {
		switch v {
		case "-test":
			test = true
		case "-build":
			build = true
		}
	}

	if !build {
		base.Fatalf("don't yet know what to do without -build flag")
	}

	deps := make(map[string]bool)
	var testDeps []string

	pkgs := load.PackagesAndErrors(pkgArgs)
	if len(pkgs) == 0 {
		base.Fatalf("no packages to deplist")
	}

	for _, p := range pkgs {
		for _, d := range p.Deps {
			deps[d] = true
		}

		if test {
			for _, d := range p.TestImports {
				testDeps = append(testDeps, d)
			}
			for _, d := range p.XTestImports {
				testDeps = append(testDeps, d)
			}
		}
	}

	var testPkgArgs []string

	for _, d := range testDeps {
		if !deps[d] {
			testPkgArgs = append(testPkgArgs, d)
		}
	}

	if len(testPkgArgs) > 0 {
		for _, p := range load.PackagesAndErrors(testPkgArgs) {
			deps[p.ImportPath] = true
			for _, d := range p.Deps {
				deps[d] = true
			}
		}
	}

	var uniqDeps []string
	for p := range deps {
		uniqDeps = append(uniqDeps, p)
	}

	depPkgs := load.PackagesForBuild(uniqDeps)

	var b work.Builder
	b.Init()

	root := &work.Action{Mode: "go deplist"}
	for _, p := range depPkgs {
		root.Deps = append(root.Deps, b.DeplistAction(work.ModeBuild, work.ModeBuild, p))
	}
	b.Do(root)
}
