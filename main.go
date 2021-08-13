package main

import (
	"flag"
	"fmt"
	"github.com/jonas-p/go-shp"
	"log"
	"os"
)

func shp2map(s *shp.Reader, key string) map[string]string {
	fields := s.Fields()
	m := make(map[string]string)

	for s.Next() {
		n, p := s.Shape()
		for k, f := range fields {
			if f.String() == key {
				val := s.ReadAttribute(n, k)
				bbox := fmt.Sprintf("%#v", p.BBox())
				m[val] = bbox
			}
		}
	}
	return m
}

func main() {
	k := flag.String("k", "", "unique key to use as comparison")
	h := flag.Bool("h", false, "help")
	flag.Parse()
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "Usage: shpdiff file1.shp file2.shp \n")
		flag.PrintDefaults()
		fmt.Fprint(os.Stderr, "Example: shpdiff new.shp old.shp\n")
	}
	key := "id"
	if *h {
		flag.Usage()
		os.Exit(0)
	}

	if *k != "" {
		key = *k
	}
	fmt.Println("using key ", key)

	if flag.NArg() < 2 {
		fmt.Println("No input supplied")
		flag.Usage()
		os.Exit(1)
	}

	oldshp, err := shp.Open(flag.Arg(1))
	if err != nil {
		log.Fatal(err)
	}
	defer oldshp.Close()
	old := shp2map(oldshp, key)

	newshp, err := shp.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer newshp.Close()
	new := shp2map(newshp, key)

	if len(old) == 0 {
		fmt.Println("no such field: ", key)
		os.Exit(1)
	}
	for key, value := range old {
		fmt.Println(key, value)
	}

	for key, value := range new {
		fmt.Println(key, value)
	}
}
