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
)

/***

-c : byte count
-l : line count
-w : word count
-m : character count

***/

type Result struct {
	filename                                   string
	byteCount, lineCount, wordCount, charCount int
	err                                        error
}

type Flag struct {
	c, l, w, m bool
}

// Process rune by rune from input and count results
func inputProcess(r io.Reader, filename string, ch chan Result) {
	result := Result{filename: filename}
	reader := bufio.NewReader(r)
	inWord := false
	result.err = nil

	for {
		in, size, err := reader.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			result.err = err
			break
		}

		if in == '\n' {
			result.lineCount++
		}

		if unicode.IsSpace(in) {
			inWord = false
		} else if !inWord {
			inWord = true
			result.wordCount++
		}

		result.charCount++
		result.byteCount += size
	}
	if f, ok := r.(*os.File); ok{
		f.Close()
	}
	ch <- result
}

// Print output
func PrintOutput(flag Flag, result Result) {
	output := make([]string, 0)
	if flag.l {
		output = append(output, strconv.Itoa(result.lineCount))
	}
	if flag.w {
		output = append(output, strconv.Itoa(result.wordCount))
	}
	if flag.c {
		output = append(output, strconv.Itoa(result.byteCount))
	}
	if flag.m {
		output = append(output, strconv.Itoa(result.charCount))
	}

	if result.err == nil {
		fmt.Println("  " + strings.Join(output, "   ") + "   " + result.filename)
	}
}

func main() {

	// Flag input
	var flags Flag
	flag.BoolVar(&flags.c, "c", false, "byte count")
	flag.BoolVar(&flags.l, "l", false, "line count")
	flag.BoolVar(&flags.w, "w", false, "word count")
	flag.BoolVar(&flags.m, "m", false, "character count")

	flag.Parse()
	args := flag.Args()

	// Channel
	ch := make(chan Result)
	var channel int = 0

	// Pipeline and file input
	if len(args) == 0 {
		channel = 1
		go inputProcess(os.Stdin, "", ch)
	} else {
		for _, filename := range args {
			f, err := os.Open(filename)
			if err != nil {
				// log.Fatal(err)
				log.Println(err)
				continue
			}
			channel++
			go inputProcess(f, filename, ch)
		}
	}

	// No flags are provided, so default will be -l -w -c
	if !(flags.c) && !flags.l && !flags.w && !flags.m {
		flags.c, flags.w, flags.l = true, true, true
	}

	// Wait for goroutine to send output and print result
	for i := 0; i < channel; i++ {
		output := <-ch
		PrintOutput(flags, output)
	}
}
