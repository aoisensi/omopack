package main

import "flag"

var (
	flagInclude = flag.String("include", "*.txt,*.md,_bundletool_mod_icon.png", "Include files that match the pattern")
)
