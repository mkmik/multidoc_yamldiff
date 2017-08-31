package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/golang/glog"
	"github.com/juju/errors"
	"github.com/kylelemons/godebug/pretty"
)

var (
	context = flag.Int("C", 3, "Output N lines of copied context")
)

func init() {
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: multidoc_yamldiff file1.yaml file2.yaml\n")
	flag.PrintDefaults()
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

type queue []string

func (q *queue) append(s string, l int) {
	if len(*q) >= l {
		*q = (*q)[1:]
	}
	*q = append(*q, s)
}

// compact takes a raw diff as input and removes the unnecessay context,
// by keeping only at most `context` lines above and below.
func compact(src string, context int) string {
	var buf bytes.Buffer
	var before queue

	lines := strings.Split(src, "\n")

	hot := 0
	for _, l := range lines {

		if strings.HasPrefix(l, "+") || strings.HasPrefix(l, "-") {
			if hot == 0 {
				fmt.Fprintln(&buf, "@@ .. @@")
			}
			hot = context + 1

			for _, b := range before {
				fmt.Fprintln(&buf, b)
			}
			before = nil
		}

		if hot > 0 {
			hot--
			fmt.Fprintln(&buf, l)
		} else {
			before.append(l, context)
		}
	}

	return buf.String()
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

	fmt.Println(compact(pretty.Compare(v1, v2), *context))

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
