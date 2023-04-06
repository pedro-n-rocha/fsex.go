package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	arg "github.com/alexflint/go-arg"
	log "github.com/sirupsen/logrus"
)

type Manifest []ManifestEntry

type ManifestEntry struct {
	Config   string
	RepoTags []string
	Layers   []string
}

func main() {

	var args struct {
		Input  string `arg:"required,positional" help:"input image tar file"`
		Output string `arg:"required,-o"  help:"output directory"`
	}

	arg.MustParse(&args)

	var inp = args.Input
	var out = args.Output

	var path string = ".tmp"
	err := Untar(inp, path)
	if err != nil {
		log.Error(err)
		return
	}

	fmanifest, err := os.Open(path + "/manifest.json")
	if err != nil {
		log.Error(err)
		return
	}
	defer fmanifest.Close()
	byteValue, _ := ioutil.ReadAll(fmanifest)

	var manifest Manifest
	err = json.Unmarshal(byteValue, &manifest)
	if err != nil {
		log.Error(err)
		return
	}

	if len(manifest) > 1 {
		log.Error("multi-image archives cannot be extracted: contains %d images", len(manifest))
		return
	}
	if len(manifest) < 1 {
		log.Error("invalid archive")
		return
	}

	for _, v := range manifest {
		for _, l := range v.Layers {
			err = UnTarGzip(path+"/"+l, out)
			if err != nil {
				log.Error(err)
			}
		}
	}
	os.RemoveAll(".tmp")
}

func Untar(tarball, target string) error {
	reader, err := os.Open(tarball)
	if err != nil {
		return err
	}
	defer reader.Close()
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := filepath.Join(target, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
	}
	return nil
}

func UnTarGzip(source, target string) error {
	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(archive)

	defer archive.Close()
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := filepath.Join(target, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
	}
	return err
}
