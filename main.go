/* projget - simple version control system
* copyright - (c) Abhigyan Ghosh 
* 2026 - Present
* Lisenced Under MIT Lisence
*/
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const MAGIC uint64 = 0x50524F4A474554 // PROJGET
const VERSION uint32 = 1

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: projget <command>")
		return
	}

	switch os.Args[1] {
	case "init":
		initRepo()

	case "bundle":
		bundle(".")

	case "push":
		if len(os.Args) < 4 {
			fmt.Println("Usage: projget push <dir> <url>")
			return
		}
		push(os.Args[2], os.Args[3])

	case "get":
		if len(os.Args) < 3 {
			fmt.Println("Usage: projget get <file-or-url>")
			return
		}
		get(os.Args[2])

	default:
		fmt.Println("Unknown command")
	}
}

func initRepo() {
	err := os.Mkdir(".Proj", 0755)
	if err != nil && !os.IsExist(err) {
		fmt.Println("Error creating .Proj:", err)
		return
	}

	content := fmt.Sprintf(
		"MAGIC=0x%X\nVERSION=%d\nCREATED=%s\n",
		MAGIC,
		VERSION,
		time.Now().Format(time.RFC3339),
	)

	err = os.WriteFile(".Proj/ident.txt", []byte(content), 0644)
	if err != nil {
		fmt.Println("Error writing ident:", err)
		return
	}

	fmt.Println("Initialized empty projget repository")
}

func bundle(dir string) {
	outFile := "repo.pf"

	file, err := os.Create(outFile)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	var files []string

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			if strings.Contains(path, ".Proj") {
				return filepath.SkipDir
			}
			return nil
		}

		cleanPath := strings.TrimPrefix(path, "./")
		files = append(files, cleanPath)
		return nil
	})

	// Write header
	binary.Write(file, binary.LittleEndian, MAGIC)
	binary.Write(file, binary.LittleEndian, VERSION)
	binary.Write(file, binary.LittleEndian, uint32(len(files)))

	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			continue
		}

		pathBytes := []byte(f)

		binary.Write(file, binary.LittleEndian, uint16(len(pathBytes)))
		file.Write(pathBytes)

		binary.Write(file, binary.LittleEndian, uint64(len(data)))
		file.Write(data)
	}

	fmt.Println("Bundled into", outFile)
}

func push(dir, url string) {
	bundle(dir)

	f, err := os.Open("repo.pf")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer f.Close()

	resp, err := http.Post(url+"/push", "application/octet-stream", f)
	if err != nil {
		fmt.Println("Push failed:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Push status:", resp.Status)
}

func get(source string) {
	var data []byte
	var err error

	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		resp, err2 := http.Get(source)
		if err2 != nil {
			fmt.Println("Download failed:", err2)
			return
		}
		defer resp.Body.Close()

		data, err = io.ReadAll(resp.Body)
	} else {
		data, err = os.ReadFile(source)
	}

	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	reader := bytes.NewReader(data)

	var magic uint64
	var version uint32
	var count uint32

	binary.Read(reader, binary.LittleEndian, &magic)
	binary.Read(reader, binary.LittleEndian, &version)
	binary.Read(reader, binary.LittleEndian, &count)

	if magic != MAGIC {
		fmt.Println("Invalid projget file")
		return
	}

	for i := 0; i < int(count); i++ {
		var pathLen uint16
		binary.Read(reader, binary.LittleEndian, &pathLen)

		pathBytes := make([]byte, pathLen)
		reader.Read(pathBytes)

		var fileSize uint64
		binary.Read(reader, binary.LittleEndian, &fileSize)

		fileData := make([]byte, fileSize)
		reader.Read(fileData)

		os.MkdirAll(filepath.Dir(string(pathBytes)), 0755)
		os.WriteFile(string(pathBytes), fileData, 0644)
	}

	fmt.Println("Project extracted successfully")
}
