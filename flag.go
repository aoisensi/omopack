package main

import "flag"

var (
	flagInclude = flag.String("include", "LICENSE,*.txt,*.md,_bundletool_mod_icon.png", "Include files that match the pattern")
)
