package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

var refresh = flag.Bool("refresh", false, "force re-index the system")
var root = flag.String("root", "/", "base directory to copy porn to")
var every = flag.Int("every", 1000, "time between copying images IN MILLISECONDS")
var count = flag.Int("count", 5, "number of images to copy per <every>")
var debug = flag.Bool("debug", false, "print every instance of copying")

func main() {
	flag.Usage = Usage
	flag.Parse()
	pornPaths := append(flag.Args(), ".")
	rand.Seed(time.Now().Unix())

	cacheDir, err := os.UserCacheDir()
	Handle("getting cache directory", err)
	cacheOutPath := cacheDir + "/porncpy.json"

	if *root != "/" && !*refresh {
		log.Println("WARNING: --root specified without --refresh. If the cache was created for a",
			"different root, that will be used instead.")
		log.Println("WARNING: Waiting 5 seconds, just in case.")
		time.Sleep(5 * time.Second)
	}

	var allPaths []string
	if *refresh || !Exists(cacheOutPath) {
		log.Println("Refreshing. This may take a few minutes.")

		allPaths, err = AllFs(*root, false)
		Handle("getting directories", err)
		cacheFile, err := os.OpenFile(cacheOutPath, os.O_CREATE, 0777)
		Handle("opening cache file", err)
		jsonResult, err := json.Marshal(allPaths)
		Handle("encoding directory cache", err)
		_, err = cacheFile.Write(jsonResult)
		Handle("writing to cache file", err)
	} else {
		log.Println("Using directories cache")

		jsonStr, err := ioutil.ReadFile(cacheOutPath)
		Handle("reading cache file", err)
		err = json.Unmarshal(jsonStr, &allPaths)
		Handle("decoding cache file", err)
	}

	log.Println("Got directories")
	log.Println("Getting all images in", pornPaths)

	var allImages []string
	var visited []string
	for _, path := range pornPaths {
		// Just in case
		for _, i := range visited {
			if path == i {
				continue
			}
		}

		visited = append(visited, path)
		allT, err := AllFs(path, true)
		Handle("reading porn dir "+path, err)
		allImages = append(allImages, allT...)
	}

	log.Println("Got all porn")
	log.Println("Starting")

	for i := 0; i < *count; i++ {
		go PornCopy(allImages, allPaths, *every)
	}

	<-(chan int)(nil)
}

func AllFs(root string, files bool) ([]string, error) {
	var r []string

	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				if errors.Is(err, fs.ErrPermission) || errors.Is(err, fs.ErrNotExist) {
					return nil
				}

				return err
			}

			// Short circuiting my beloved
			if (files && info.IsDir()) || (!files && !info.IsDir()) {
				return nil
			}

			r = append(r, path)
			return nil
		})
	if err != nil {
		return []string{}, err
	}

	return r, nil
}

func PornCopy(images []string, paths []string, every int) {
	for {
		image := images[rand.Intn(len(images))]
		imageExtension := filepath.Ext(image)
		outDir := paths[rand.Intn(len(paths))]
		outPath := outDir + "/GOONER" + fmt.Sprint(rand.Intn(100000)) + imageExtension

		if *debug {
			log.Println("Writing", image, "to", outPath)
		}

		data, err := ioutil.ReadFile(image)
		Handle("reading porn image "+image, err)
		err = ioutil.WriteFile(outPath, data, 0644)

		if err != nil {
			if errors.Is(err, fs.ErrPermission) || errors.Is(err, fs.ErrNotExist) {
				continue
			}

			Handle("Writing porn image "+image+" to "+outPath, err)
		}

		time.Sleep(time.Millisecond * time.Duration(every))
	}
}

func Exists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

func Handle(spot string, err error) {
	if err != nil {
		log.Fatalln("failed", spot+".", "error:", err)
	}
}

func Usage() {
	fmt.Printf("Usage: %s [OPTIONS] PORN_PATHS...\n\n", os.Args[0])
	flag.PrintDefaults()
}
