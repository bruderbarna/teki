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

func checkIfLoopHasEndTag(splitSource [][]byte, currentIndex int, ler *regexp.Regexp) bool {
	var i int
	for i = currentIndex + 2; i < len(splitSource); i++ {
		if ler.Match(splitSource[i]) {
			break
		}
	}
	if i >= len(splitSource) {
		return false
	}
	return true
}

func getPrefix(indicator bool) []byte {
	if indicator == true {
		return []byte("\t")
	}
	return []byte("")
}

func generateOutput(source []byte) ([]byte, error) {
	fbRp := regexp.MustCompile("^(forward|backward) [1-9][0-9]*$")
	lrRp := regexp.MustCompile("^(left|right) -?[1-9][0-9]*$")
	loopRp := regexp.MustCompile("^loop [1-9][0-9]*$")
	loopEndRp := regexp.MustCompile("^loopend$")
	drawRp := regexp.MustCompile("^draw$")
	nodrawRp := regexp.MustCompile("^nodraw$")

	output := []byte("#include <graphics.h>\n#include <stdlib.h>\n#include <stdio.h>\n\nint main()\n{\n\tinitwindow(800, 800);\n\n")

	drawIndicator := true
	pen := pen{
		position{400, 400},
		90,
	}

	splitSource := bytes.Split(source, []byte("\n"))
	for i, v := range splitSource {
		v = bytes.TrimSpace(v)
		fmt.Printf("%q\n", v)

		switch {
		case fbRp.Match(v):
			split := bytes.Split(v, []byte(" "))
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
		case lrRp.Match(v):
			split := bytes.Split(v, []byte(" "))
			rotationAngle, err := strconv.Atoi(string(split[1]))
			if err != nil {
				panic(err)
			}

			pen.angle = (pen.angle + rotationAngle) % 360
		case loopRp.Match(v):
			if !checkIfLoopHasEndTag(splitSource, i, loopEndRp) {
				return nil, errors.New("expected loopend, found EOF")
			}

			splitLoopString := bytes.Split(v, []byte(" "))
			loopCounter := splitLoopString[1]
			toAppend := append([]byte("\tfor (int i = 0; i < "), loopCounter...)
			toAppend = append(toAppend, []byte("; i++) {\n")...)
			output = append(output, toAppend...)
		case loopEndRp.Match(v):
			output = append(output, []byte("\t}\n")...)
		case drawRp.Match(v):
			drawIndicator = true
		case nodrawRp.Match(v):
			drawIndicator = false
		}

		// TODO: handle error
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
