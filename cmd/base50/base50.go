package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/james-antill/base50"
)

func main() {
	var (
		err error
		bin []byte

		help   = flag.Bool("h", false, "display this message")
		input  = flag.String("i", "", `input file (use: "-" for stdin, "" for arguments)`)
		output = flag.String("o", "-", `output file (use: "-" for stdout)`)
		base16 = flag.Bool("x", false, `treat input/output as base16`)
		decode = flag.Bool("d", false, `decode input`)
	)

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}
	if *input == "" && len(flag.Args()) < 1 {
		fmt.Fprintln(os.Stderr, "No arguments given for input.")
		flag.Usage()
		os.Exit(1)
	}

	var fin io.Reader
	var fout io.Writer
	fin, fout = os.Stdin, os.Stdout

	if *input != "-" && *input != "" {
		if fin, err = os.Open(*input); err != nil {
			fmt.Fprintln(os.Stderr, "input file err:", err)
			os.Exit(1)
		}
	}

	if *output != "-" {
		if fout, err = os.Create(*output); err != nil {
			fmt.Fprintln(os.Stderr, "output file err:", err)
			os.Exit(1)
		}
	}

	de16input := !*decode && *base16

	if de16input {
		fin = hex.NewDecoder(fin)
	}
	if *input == "" {
		done := false
		for _, b := range flag.Args() {
			if de16input && done {
				bin = append(bin, ' ')
			}
			done = true
			bin = append(bin, []byte(b)...)
		}
		if de16input {
			if len(bin)%2 != 0 { // Allow user the specify 0 instead of 00
				bin = append([]byte{'0'}, bin...)
			}
			dst := make([]byte, hex.DecodedLen(len(bin)))
			hex.Decode(dst, bin)
			bin = dst
		}
	} else if bin, err = ioutil.ReadAll(fin); err != nil {
		fmt.Fprintln(os.Stderr, "read input err:", err)
		os.Exit(1)
	}

	if *decode {
		lastBytePos := len(bin) - 1
		if bin[lastBytePos] == '\n' {
			bin = bin[:lastBytePos]
		}

		decoded := make([]byte, base50.DecodeLen(len(bin)))
		decoded, err := base50.Decode(decoded, bin)
		if err != nil {
			fmt.Fprintln(os.Stderr, "decode input err:", err)
			os.Exit(1)
		}

		if *base16 {
			fout = hex.NewEncoder(fout)
		}
		_, err = io.Copy(fout, bytes.NewReader(decoded))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		os.Exit(0)
	}

	fmt.Fprintln(fout, base50.EncodeToString(bin))
}
