package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/RangelReale/fproto-doc/gen-html-default"
	"github.com/RangelReale/fproto/fdep"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "array flags"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var (
	incPaths   = arrayFlags{}
	protoPaths = arrayFlags{}
	outputPath = flag.String("output_path", "", "Output root path")
)

func main() {
	// parse flags
	flag.Var(&incPaths, "inc_path", "Include paths (can be set multiple times)")
	flag.Var(&protoPaths, "proto_path", "Application proto files root paths (can be set multiple times)")
	flag.Parse()

	if *outputPath == "" {
		log.Fatal("The output path is required")
	}

	// create dependency parser
	parsedep := fdep.NewDep()
	// Accept not found dependencies, will appear on the output without links
	parsedep.IgnoreNotFoundDependencies = true

	// add include paths
	parsedep.IncludeDirs = append(parsedep.IncludeDirs, incPaths...)

	// add application proto files
	for _, pp := range protoPaths {
		if s, err := os.Stat(pp); err != nil {
			log.Fatalf("Error reading proto_path: %v", err)
		} else if !s.IsDir() {
			log.Fatalf("proto_path isn't a directory: %s", pp)
		}

		err := parsedep.AddPath(pp, fdep.DepType_Own)
		if err != nil {
			log.Fatal(err)
		}
	}

	// create output directory
	if err := os.MkdirAll(*outputPath, os.ModePerm); err != nil {
		log.Fatalf("Error creating output_path '%s': %v", *outputPath, err)
	}

	// create output file
	outfile, err := os.Create(filepath.Join(*outputPath, "index.html"))
	if err != nil {
		log.Fatal("Error creating html file: %v", err)
	}

	defer outfile.Close()

	// creates the HTML generator
	gen := fproto_doc_html_default.NewGenerator()

	// generate the files
	err = gen.Generate(parsedep, outfile)
	if err != nil {
		log.Fatal(err)
	}
}
