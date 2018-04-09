// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package deplist

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"cmd/go/internal/base"
	"cmd/go/internal/cmdflag"
	"cmd/go/internal/work"
)

const cmd = "deplist"

// deplistFlagDefn is the set of flags we process.
var deplistFlagDefn = []*cmdflag.Defn{
	{Name: "test", BoolVar: new(bool)},
	{Name: "build", BoolVar: new(bool)},
}

var deplistTool string

// add build flags to deplistFlagDefn.
func init() {
	var cmd base.Command
	work.AddBuildFlags(&cmd)
	cmd.Flag.VisitAll(func(f *flag.Flag) {
		deplistFlagDefn = append(deplistFlagDefn, &cmdflag.Defn{
			Name:  f.Name,
			Value: f.Value,
		})
	})
}

// deplistFlags processes the command line, splitting it at the first non-flag
// into the list of flags and list of packages.
func deplistFlags(args []string) (passToDeplist, packageNames []string) {
	for i := 0; i < len(args); i++ {
		if !strings.HasPrefix(args[i], "-") {
			return args[:i], args[i:]
		}

		f, value, extraWord := cmdflag.Parse(cmd, deplistFlagDefn, args, i)
		if f == nil {
			fmt.Fprintf(os.Stderr, "deplist: flag %q not defined\n", args[i])
			fmt.Fprintf(os.Stderr, "Run \"go help deplist\" for more information\n")
			os.Exit(2)
		}
		if f.Value != nil {
			if err := f.Value.Set(value); err != nil {
				base.Fatalf("invalid flag argument for -%s: %v", f.Name, err)
			}
		}
		if extraWord {
			i++
		}
	}
	return args, nil
}
