package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
)

type position struct {
	x, y int
}

type pen struct {
	pos   position
	angle int
}

func getSource(path string) ([]byte, error) {
	input, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("Couldn't open source file for reading. Aborting.")
	}
	return input, nil
}

func getLoopEndTagIndex(splitSource [][]byte, currentIndex int, lr *regexp.Regexp, ler *regexp.Regexp) int {
	var i int
	for i = currentIndex + 1; i < len(splitSource); i++ {
		if lr.Match(splitSource[i]) {
			return -1
		}
		if ler.Match(splitSource[i]) {
			break
		}
	}
	return i
}

func generateOutput(source []byte) ([]byte, error) {
	fbRp := regexp.MustCompile("^(forward|backward) [1-9][0-9]*$")
	lrRp := regexp.MustCompile("^(left|right) -?[1-9][0-9]*$")
	loopRp := regexp.MustCompile("^loop [1-9][0-9]*$")
	loopEndRp := regexp.MustCompile("^loopend$")
	drawRp := regexp.MustCompile("^draw$")
	nodrawRp := regexp.MustCompile("^nodraw$")
	colorRp := regexp.MustCompile("^color (0|BLACK|1|BLUE|2|GREEN|3|CYAN|4|RED|5|MAGENTA|6|BROWN|7|LIGHTGRAY|8|DARKGRAY|9|LIGHTBLUE|10|LIGHTGREEN|11|LIGHTCYAN|12|LIGHTRED|13|LIGHTMAGENTA|14|YELLOW|15|WHITE)$")
	bgColorRp := regexp.MustCompile("^bgcolor (0|BLACK|1|BLUE|2|GREEN|3|CYAN|4|RED|5|MAGENTA|6|BROWN|7|LIGHTGRAY|8|DARKGRAY|9|LIGHTBLUE|10|LIGHTGREEN|11|LIGHTCYAN|12|LIGHTRED|13|LIGHTMAGENTA|14|YELLOW|15|WHITE)$")

	splitSource := bytes.Split(source, []byte("\n"))
	for i := 0; i < len(splitSource); i++ {
		splitSource[i] = bytes.TrimSpace(splitSource[i])
	}

	output := []byte("#include <graphics.h>\n#include <stdlib.h>\n#include <stdio.h>\n\nint main()\n{\n\tinitwindow(800, 800);\n\n")

	drawIndicator := true
	pen := pen{
		position{400, 400},
		90,
	}

	var loopBegin, loopEnd, loopCounter int

	for i := 0; i < len(splitSource); i++ {
		switch {
		case fbRp.Match(splitSource[i]):
			split := bytes.Split(splitSource[i], []byte(" "))
			steps, err := strconv.Atoi(string(split[1]))
			if err != nil {
				panic(err)
			}
			newPos := position{
				int((float64(steps) * math.Cos((float64(pen.angle)*math.Pi)/180.0)) + float64(pen.pos.x)),
				int((float64(steps) * math.Sin((float64(pen.angle)*math.Pi)/180.0)) + float64(pen.pos.y)),
			}

			if drawIndicator {
				output = append(output, []byte(fmt.Sprintf("\tline(%d, %d, %d, %d);\n", pen.pos.x, pen.pos.y, newPos.x, newPos.y))...)
				break
			}
			pen.pos = newPos
		case lrRp.Match(splitSource[i]):
			split := bytes.Split(splitSource[i], []byte(" "))
			rotationAngle, err := strconv.Atoi(string(split[1]))
			if err != nil {
				panic(err)
			}

			pen.angle = (pen.angle + rotationAngle) % 360
		case loopRp.Match(splitSource[i]):
			loopBegin = i + 1
			loopEnd = getLoopEndTagIndex(splitSource, i, loopRp, loopEndRp)
			if loopEnd == -1 {
				return nil, errors.New(fmt.Sprintf("Nested loops aren't supported. Error occured at line %d:\t\n%s", i+1, string(splitSource[i])))
			}
			if loopEnd >= len(splitSource) {
				return nil, errors.New(fmt.Sprintf("Expected loopend, found EOF. Error occured at line %d:\t\n%s", i+1, string(splitSource[i])))
			}

			splitLoopString := bytes.Split(splitSource[i], []byte(" "))
			var err error
			loopCounter, err = strconv.Atoi(string(splitLoopString[1]))
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Couldn't convert []byte to int. Error occured at line %d:\t\n%s", i+1, string(splitSource[i])))
			}
		case loopEndRp.Match(splitSource[i]):
			if loopCounter > 1 {
				loopCounter--
				i = loopBegin
			}
		case drawRp.Match(splitSource[i]):
			drawIndicator = true
		case nodrawRp.Match(splitSource[i]):
			drawIndicator = false
		case colorRp.Match(splitSource[i]):
			split := bytes.Split(splitSource[i], []byte(" "))
			color := split[1]

			toAppend := append([]byte("\tsetcolor("), color...)
			toAppend = append(toAppend, []byte(");\n")...)
			output = append(output, toAppend...)
		case bgColorRp.Match(splitSource[i]):
			split := bytes.Split(splitSource[i], []byte(" "))
			bgColor := split[1]

			toAppend := append([]byte("\tsetbkcolor("), bgColor...)
			toAppend = append(toAppend, []byte(");\n\tcleardevice();\n")...)
			output = append(output, toAppend...)
		case bytes.Equal([]byte(""), splitSource[i]):
			break
		default:
			return nil, errors.New(fmt.Sprintf("Unexpected identifier. Error occured at line %d:\t\n%s", i+1, string(splitSource[i])))
		}
	}

	output = append(output, []byte("\n\tgetch();\n\tclosegraph();\n\treturn 0;\n}\n")...)
	return output, nil
}

func writeOutput(filepath string, output []byte) error {
	err := ioutil.WriteFile(filepath, output, 0755)
	if err != nil {
		return errors.New("Couldn't open output file for writing. Aborting")
	}
	return nil
}

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Number of arguments isn't 2. First argument should be the source file, second should be the output file")
	}

	inputFilepath := os.Args[1]
	outputFilepath := os.Args[2]
	input, err := getSource(inputFilepath)
	if err != nil {
		log.Fatal(err)
	}

	output, err := generateOutput(input)
	if err != nil {
		log.Fatal(err)
	}

	if err := writeOutput(outputFilepath, output); err != nil {
		log.Fatal(err)
	}
}
