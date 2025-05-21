package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

/***

-c : byte count
-l : line count
-w : word count
-m : character count

***/

func byteCount(data *[]byte) int {
	return len(*data)
}

func lineCount(data *[]byte) int {
	return bytes.Count(*data, []byte("\n"))
}

func wordCount(data *[]byte) int {
	return len(bytes.Fields(*data))
}

func charCount(data *[]byte) int {
	return utf8.RuneCount(*data)
}

func fileInput(filename string) []byte {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	return input(f)
}

func input(r io.Reader) []byte {
	data, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func main() {

	// *bool
	c := flag.Bool("c", false, "byte count")
	l := flag.Bool("l", false, "line count")
	w := flag.Bool("w", false, "word count")
	m := flag.Bool("m", false, "character count")

	flag.Parse()

	args := flag.Args()
	output := make([]string, 0)
	data := make([]byte, 0)
	var filename string

	if len(args) == 0 {
		data = input(os.Stdin)
	} else {
		filename = args[0]
		data = fileInput(filename)
	}

	if *c == false && *l == false && *w == false && *m == false {
		*c, *w, *l = true, true, true
	}

	if *l {
		output = append(output, strconv.Itoa(lineCount(&data)))
	}
	if *w {
		output = append(output, strconv.Itoa(wordCount(&data)))
	}
	if *c {
		output = append(output, strconv.Itoa(byteCount(&data)))
	}
	if *m {
		output = append(output, strconv.Itoa(charCount(&data)))
	}
	fmt.Println("  " + strings.Join(output, " ") + " " + filename)
}
