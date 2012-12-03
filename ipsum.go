package main

import (
	"flag"
	"fmt"
	"github.com/daaku/go.flagbytes"
	"github.com/drhodes/golorem"
	"github.com/dustin/go-humanize"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
)

type ipsumReaderT struct{}

var ipsumReader *ipsumReaderT

func (r *ipsumReaderT) Read(p []byte) (n int, err error) {
	t := lorem.Paragraph(5, 10)
	tl := len(t)
	pl := len(p)
	for tl < pl {
		t = t + "\n\n" + lorem.Paragraph(5, 10)
		tl = len(t)
	}
	b := []byte(t)
	for i := range p {
		p[i] = b[i]
	}
	return pl, nil
}

func main() {
	dirCount := flag.Int("dirs", 10, "number of directories to generate")
	byteCount := flagbytes.Bytes("bytes", "5mb", "approximate amount of data to generate")
	baseDir := flag.String("root", "", "root output directory")
	flag.Parse()

	var err error
	if *baseDir == "" {
		*baseDir, err = ioutil.TempDir("", "ipsum")
		if err != nil {
			panic(err)
		}
		fmt.Println("Output in", *baseDir)
	} else {
		err := os.MkdirAll(*baseDir, os.FileMode(0700))
		if err != nil {
			panic(err)
		}
	}

	dirs, err := makeDirs(*baseDir, *dirCount)
	if err != nil {
		panic(err)
	}

	err = makeFiles(dirs, *byteCount)
	if err != nil {
		panic(err)
	}
}

func makeDirs(base string, count int) (dirs []string, err error) {
	dirs = make([]string, count)
	current := base
	for i := 0; i < count; i++ {
		dirs[i], err = ioutil.TempDir(current, "")
		if err != nil {
			return
		}
		if i != 0 {
			current = dirs[rand.Intn(i)]
		}
	}
	return
}

func makeFiles(dirs []string, size uint64) error {
	dirsl := len(dirs)
	var i uint64
	for i < size {
		dir := dirs[rand.Intn(dirsl)]
		file, err := ioutil.TempFile(dir, "")
		if err != nil {
			return err
		}
		filel := randSize()
		io.CopyN(file, ipsumReader, filel)
		i += uint64(filel)
		file.Close()
	}
	return nil
}

func h(f string) int64 {
	r, err := humanize.ParseBytes(f)
	if err != nil {
		panic(err)
	}
	return int64(r)
}

func randSize() int64 {
	r := rand.Float32()
	switch {
	case r < 0.05:
		return rand.Int63n(h("3mb")) + h("1mb")
	case r < 0.2:
		return rand.Int63n(h("200k")) + h("100k")
	case r < 0.5:
		return rand.Int63n(h("50k")) + h("50k")
	case r < 0.6:
		return rand.Int63n(h("50k")) + h("10k")
	}
	return rand.Int63n(h("10k"))
}

func min(x, y int) int {
	if x <= y {
		return x
	}
	return y
}
