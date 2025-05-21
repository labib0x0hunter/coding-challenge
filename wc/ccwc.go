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
	"sync"
	"unicode"
)

/***

-c : byte count
-l : line count
-w : word count
-m : character count

***/

var tokens = make(chan struct{}, 10)

type Result struct {
	filename                                   string
	byteCount, lineCount, wordCount, charCount int
	err                                        error
}

type Flag struct {
	c, l, w, m bool
}

// Process rune by rune from input and count results
func inputProcess(r io.Reader) (result Result) {
	reader := bufio.NewReader(r)
	inWord := false

	for {
		in, size, err := reader.ReadRune()	// Reads a single rune
		if err == io.EOF {					// Check if for End Of File
			break
		}
		if err != nil {
			log.Println(err)
			result.err = err
			break
		}

		if in == '\n' {						// Newline count
			result.lineCount++
		}

		if unicode.IsSpace(in) {
			inWord = false
		} else if !inWord {
			inWord = true
			result.wordCount++				// Word count
		}

		result.charCount++					// Char count
		result.byteCount += size			// Byte count
	}
	return
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

	ch := make(chan Result) 			// Channel
	var wg sync.WaitGroup				// Waitgroup for goroutine

	if len(args) == 0 {					// Stdin input
		wg.Add(1)						// Wait for 1 goroutine
		go func(ch chan Result) {
			defer wg.Done()	 			// goroutine is finised
			r := inputProcess(os.Stdin)
			r.filename = "os.Stdin"
			ch <- r
		}(ch)
	} else {							 // File input
		for _, filename := range args {
			tokens <- struct{}{}		 // Acquire token
			f, err := os.Open(filename)
			if err != nil {
				log.Println(err)
				<-tokens				 // Release token
				continue
			}
			wg.Add(1)					 // Wait for 1 goroutine
			go func(f *os.File, filename string, ch chan Result) {
				defer f.Close()			 // Close file
				defer wg.Done()			 // goroutine is finised
				r := inputProcess(f)
				r.filename = filename
				ch <- r					 // Send data to channel
				<-tokens				 // Release token
			}(f, filename, ch)
		}
	}

	// No flags are provided, so default will be -l -w -c
	if !(flags.c) && !flags.l && !flags.w && !flags.m {
		flags.c, flags.w, flags.l = true, true, true
	}

	// Wait for all goroutine to finish and close the channel
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Receive data from channel
	for out := range ch {
		PrintOutput(flags, out)
	}

}
