package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"text/tabwriter"
)

type flagArgs struct {
	ascendingOrder bool
	help           bool
}

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

func infoPrint(info os.FileInfo, w io.Writer) {
	var err error
	size := info.Size()
	if 0 < size && size < 1000 {
		_, err = fmt.Fprintf(w, "%s\t%dB\t\n", info.Name(), size)
	} else if 1000 <= size && size < 1000000 {
		_, err = fmt.Fprintf(w, "%s\t%.2fKB\t\n", info.Name(), float64(size)/1000.0)
	} else if 1000000 <= size && size < 1000000000 {
		_, err = fmt.Fprintf(w, "%s\t%.2fMB\t\n", info.Name(), float64(size)/1000000.0)
	} else if 1000000000 <= size && size < 1000000000000 {
		_, err = fmt.Fprintf(w, "%s\t%.2fGB\t\n", info.Name(), float64(size)/1000000000.0)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func usage() {
	usageString := `Usage: du [OPTION]... [FILE]...
Summarize disk usage of the set of FILEs.

Arguments:`
	fmt.Println(usageString)
	flag.PrintDefaults()
}

func duDir(path string) {
	files, error := ioutil.ReadDir(path)
	if error != nil {
		log.Fatal(error)
	}
	if args.ascendingOrder {
		sort.Sort(&bySize{files})
	} else {
		sort.Sort(sort.Reverse(&bySize{files}))
	}
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	for _, f := range files {
		infoPrint(f, w)
	}
	if error = w.Flush(); error != nil {
		log.Fatal(error)
	}
	fmt.Print(buf.String())
}

// cmdline arguments
var args flagArgs

func init() {
	flag.BoolVar(&args.ascendingOrder, "asc", false, "sort entries in ascending order")
	flag.BoolVar(&args.help, "help", false, "display this help and exit")
}

func main() {
	flag.Parse()
	if args.help || len(flag.Args()) < 1 {
		usage()
		os.Exit(1)
	}
	fileInfo, error := os.Stat(flag.Arg(0))
	if error != nil {
		log.Fatal(error)
	}
	switch mode := fileInfo.Mode(); {
	case mode.IsDir():
		duDir(flag.Arg(0))
	case mode.IsRegular():
		infoPrint(fileInfo, os.Stdout)
	}

}
