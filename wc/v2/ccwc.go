package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

/***

-c : byte count
-l : line count
-w : word count
-m : character count

***/

func inputProcess1(r io.Reader) (byteCount, line, word, char int) {
	reader := bufio.NewReader(r)
	inWord := false

	for {
		in, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if in == '\n' {
			line++
		}

		if unicode.IsSpace(in) {
			inWord = false
		} else if !inWord {
			inWord = true
			word++
		}

		char++
		byteCount += utf8.RuneLen(in)
	}

	return
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
	var filename string
	var reader io.Reader

	if len(args) == 0 { // pipeline
		reader = os.Stdin
	} else { // file
		filename = args[0]
		f, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		reader = f
		defer f.Close()
	}

	byteCount, lineCount, wordCount, charCount := inputProcess1(reader)

	if *c == false && *l == false && *w == false && *m == false {
		*c, *w, *l = true, true, true
	}

	if *l {
		output = append(output, strconv.Itoa(lineCount))
	}
	if *w {
		output = append(output, strconv.Itoa(wordCount))
	}
	if *c {
		output = append(output, strconv.Itoa(byteCount))
	}
	if *m {
		output = append(output, strconv.Itoa(charCount))
	}

	fmt.Println("  " + strings.Join(output, "   ") + "   " + filename)
}
