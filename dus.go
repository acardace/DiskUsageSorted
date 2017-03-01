package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

type bySize struct{ file []os.FileInfo }

func (f *bySize) Less(i, j int) bool {
	return f.file[i].Size() < f.file[j].Size()
}

func (f *bySize) Len() int {
	return len(f.file)
}

func (f *bySize) Swap(i, j int) {
	f.file[i], f.file[j] = f.file[j], f.file[i]
}

func usagePrint(info os.FileInfo) (out string) {
	size := info.Size()
	if 0 < size && size < 1000 {
		out = fmt.Sprintf("%s %dB", info.Name(), size)
	} else if 1000 <= size && size < 1000000 {
		out = fmt.Sprintf("%s %.2fKB", info.Name(), float64(size)/1000.0)
	} else if 1000000 <= size && size < 1000000000 {
		out = fmt.Sprintf("%s %.2fMB", info.Name(), float64(size)/1000000.0)
	} else if 1000000000 <= size && size < 1000000000000 {
		out = fmt.Sprintf("%s %.2fGB", info.Name(), float64(size)/1000000000.0)
	}
	return
}

func usage() {
	usageString := `Usage: du [OPTION]... [FILE]...
Summarize disk usage of the set of FILEs.`
	fmt.Println(usageString)
}

func duDir(path string) {
	files, error := ioutil.ReadDir(path)
	if error != nil {
		log.Fatal(error)
	}
	sort.Sort(sort.Reverse(&bySize{files}))
	for _, f := range files {
		fmt.Println(usagePrint(f))
	}
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	fileInfo, error := os.Stat(os.Args[1])
	if error != nil {
		log.Fatal(error)
	}
	switch mode := fileInfo.Mode(); {
	case mode.IsDir():
		duDir(os.Args[1])
	case mode.IsRegular():
		fmt.Println(usagePrint(fileInfo))
	}

}
