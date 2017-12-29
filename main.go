package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	errNoArgs = 1
	errRead   = 2
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %v <left file> <right file>\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	w, _, err := terminal.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot get terminal size width: %v; fallback to 80 characters default\n", err)
		w = 80
	}

	width := flag.Int("width", w, "Terminal width. You can use $(tput cols) to get the current size")
	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(errNoArgs)
	}

	ll, err := readLines(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(errRead)
	}

	lr, err := readLines(flag.Arg(1))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(errRead)
	}

	half := *width / 2

	for i, j := 0, max(len(ll), len(lr)); i < j; i++ {
		var ls, rs string
		if i < len(ll) {
			ls = strings.TrimSuffix(ll[i], "\n")
			ls = ls[:min(half, len(ls))]
		}

		if i < len(lr) {
			rs = strings.TrimSuffix(lr[i], "\n")
			rs = rs[:min(half, len(rs))]
		}

		if ls != "" || rs != "" {
			fmt.Printf("%s%s%s\n", ls, strings.Repeat(" ", half-len(ls)), rs)
		}
	}
}

func readLines(file string) ([]string, error) {
	fp, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	s := bufio.NewScanner(fp)

	var ss []string

	for s.Scan() {
		ss = append(ss, s.Text())
	}

	return ss, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
