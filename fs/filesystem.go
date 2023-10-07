package fs

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type Config []Database

var C Config

type Database struct {
	Backup string
	Output string
}

func ZipBackups(backupDir string, outputDir string, ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Printf("Backing up database dumps from %s to %s", backupDir, outputDir)

	files := getFilesInDir(backupDir)
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	zipName := "backup-" + time.Now().Format("20060102T150405") + ".zip"
	a, err := createZip(zipName, outputDir)

	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()

	z := zip.NewWriter(a)
	defer z.Close()

	for _, f := range files {
		addFileToZip(f.Name(), backupDir, z)
	}

	ch <- fmt.Sprintf("Compressing database dumps from %s to %s was successful", backupDir, outputDir)
}

func createZip(n string, p string) (*os.File, error) {
	path := fmt.Sprintf("%s/%s", p, n)
	f, err := os.Create(path)

	if err != nil {
		log.Fatal(err)
	}

	return f, err
}

func getFilesInDir(d string) []os.DirEntry {
	files, err := os.ReadDir(d)
	if err != nil {
		log.Fatal(err)
	}

	return files
}

func addFileToZip(f string, d string, w *zip.Writer) {
	fn := fmt.Sprintf("%s/%s", d, f)

	file, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	zf, err := w.Create(f)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := io.Copy(zf, file); err != nil {
		log.Fatal(err)
	}
}

func DeleteFilesInDir(d string) {
	err := os.RemoveAll(d)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Directory %s deleted successfully", d)
}
