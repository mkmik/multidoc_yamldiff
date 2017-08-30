package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"

	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/golang/glog"
	"github.com/juju/errors"
	"github.com/kylelemons/godebug/pretty"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: multidoc_yamldiff file1.yaml file2.yaml\n")
	os.Exit(1)
}

type documents []interface{}

func (d documents) Len() int           { return len(d) }
func (d documents) Less(i, j int) bool { return pretty.Sprint(d[i]) < pretty.Sprint(d[j]) }
func (d documents) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }

func unmarshal(fname string) (documents, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer f.Close()

	var res documents

	d := yaml.NewYAMLOrJSONDecoder(f, 1024)
	for {
		var v interface{}
		if err := d.Decode(&v); err != nil {
			if err == io.EOF {
				break
			}
			return nil, errors.Trace(err)
		}
		res = append(res, v)
	}
	return res, nil
}

func run(fname1, fname2 string) error {
	v1, err := unmarshal(fname1)
	if err != nil {
		return errors.Trace(err)
	}

	v2, err := unmarshal(fname2)
	if err != nil {
		return errors.Trace(err)
	}

	sort.Sort(v1)
	sort.Sort(v2)

	fmt.Println(pretty.Compare(v1, v2))

	return nil
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if flag.NArg() < 2 {
		usage()
	}

	if err := run(flag.Arg(0), flag.Arg(1)); err != nil {
		glog.Exitf("%+v", err)
	}
}
