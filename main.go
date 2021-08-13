package main

import (
	"flag"
	"fmt"
	"github.com/jonas-p/go-shp"
	"log"
	"os"
)

func shp2map(s *shp.Reader, key string) (int, map[string]string) {
	fields := s.Fields()
	m := make(map[string]string)
	var pos int
	for s.Next() {
		n, p := s.Shape()
		for k, f := range fields {
			if f.String() == key {
				pos = k
				val := s.ReadAttribute(n, k)
				bbox := fmt.Sprintf("%#v", p.BBox())
				m[val] = bbox
			}
		}
	}
	return pos, m
}

func main() {
	k := flag.String("k", "", "unique key to use as comparison (default: id)")
	h := flag.Bool("h", false, "help")
	verbose := flag.Bool("v", false, "verbose")
	flag.Parse()
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "Outputs a shapefile that contains what is in one shapefile and not the other or that is different based on a key\n")
		fmt.Fprint(os.Stderr, "Usage: shpdiff file1.shp file2.shp diff.shp \n")
		flag.PrintDefaults()
		fmt.Fprint(os.Stderr, "Example: shpdiff new.shp old.shp diff.shp\n")
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

	var diffshp string
	if flag.Arg(2) != "" {
		diffshp = flag.Arg(2)
	} else {
		diffshp = "diff.shp"
	}

	oldshp, err := shp.Open(flag.Arg(1))
	if err != nil {
		log.Fatal(err)
	}
	defer oldshp.Close()
	_, old := shp2map(oldshp, key)

	newshp, err := shp.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer newshp.Close()
	t := newshp.GeometryType
	np, new := shp2map(newshp, key)
	if len(old) == 0 {
		fmt.Println("no such field: ", key)
		os.Exit(1)
	}

	diffmap := make(map[string]string)
	for key, value := range new {
		//check if in old
		v, found := old[key]
		if found && v == value {
			if *verbose {
				fmt.Println("key", key, "matches.")
			}
			diffmap[key] = "m"
		} else if !found {
			if *verbose {
				fmt.Println("key", key, "not found.")
			}
			diffmap[key] = "w"
		} else {
			if *verbose {
				fmt.Println("key", key, "does not match.")
			}
			diffmap[key] = "w"
		}
	}
	//write diff
	d, err := shp.Create(diffshp, t)
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()

	newshp, err = shp.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer newshp.Close()
	fields := newshp.Fields()
	d.SetFields(fields)
	i := 0
	for newshp.Next() {
		n, p := newshp.Shape()
		attr := newshp.ReadAttribute(n, np)
		w, found := diffmap[attr]
		if found && w == "w" {

			d.Write(p)
			for k, _ := range fields {
				val := newshp.ReadAttribute(n, k)
				d.WriteAttribute(i, k, val)
			}
			i++
		}
	}
}
