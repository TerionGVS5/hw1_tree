package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

const basicSymbolRight string = "├───"
const basicSymbolRightLast string = "└───"
const basicSymbolVertical string = "│"

type node struct {
	name     string
	children []*node
	level    int
}

func processNode(currentNode *node, printFiles bool, level int, prefixPath string) error {
	dirsAndFiles, err := ioutil.ReadDir(prefixPath + currentNode.name)
	if err != nil {
		return err
	}
	for _, dirOrFile := range dirsAndFiles {
		isFile := !dirOrFile.IsDir()
		nodeName := dirOrFile.Name()
		if isFile {
			if !printFiles {
				continue
			}
			sizeOfBytesStr := fmt.Sprintf("%sb", strconv.Itoa(int(dirOrFile.Size())))
			if sizeOfBytesStr == "0b" {
				sizeOfBytesStr = "empty"
			}
			nodeName = fmt.Sprintf("%s (%s)", nodeName, sizeOfBytesStr)
		}
		newNode := node{
			name:  nodeName,
			level: level,
		}
		currentNode.children = append(currentNode.children, &newNode)
		if !isFile {
			err = processNode(&newNode, printFiles, level+1, fmt.Sprintf("%s%s/", prefixPath, currentNode.name))
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func printNode(out io.StringWriter, currentNode *node, parentsLast []bool) error {
	countChildren := len(currentNode.children)
	for index, nodeElement := range currentNode.children {
		var mainText string
		if index == countChildren-1 {
			mainText = basicSymbolRightLast
		} else {
			mainText = basicSymbolRight
		}
		var indentText string
		for _, parentLastFlag := range parentsLast {
			if parentLastFlag {
				indentText = indentText + "\t"
			} else {
				indentText = indentText + basicSymbolVertical + "\t"
			}
		}
		_, err := out.WriteString(fmt.Sprintf("%s%s%s\n", indentText, mainText, nodeElement.name))
		if err != nil {
			return err
		}
		newParentsLast := make([]bool, len(parentsLast))
		copy(newParentsLast, parentsLast)
		newParentsLast = append(newParentsLast, index == countChildren-1)
		err = printNode(out, nodeElement, newParentsLast)
		if err != nil {
			return err
		}
	}
	return nil
}

func dirTree(out io.StringWriter, path string, printFiles bool) error {
	rootNode := node{
		name: path,
	}
	err := processNode(&rootNode, printFiles, 1, "")
	if err != nil {
		return err
	}
	err = printNode(out, &rootNode, []bool{})
	if err != nil {
		return err
	}
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
