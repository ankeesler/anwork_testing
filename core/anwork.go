// This package contains utilities for running Anwork tests.
package core

import (
	"archive/zip"
	"crypto/md5"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/big"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
)

// This is the path to where the Anwork release zip files are kept.
const ReleasePath string = "../release"

// This is the lock that guards the unzipping procedure.
var unzipMutex sync.Mutex

// Anwork represents an Anwork program that can be executed.
type Anwork struct {
	// This is the path to the context directory for the anwork executable to use.
	contextPath string

	// This is the path to the actual executable.
	binaryPath string
}

// Make an Anwork struct for the provided version. This function will look in the correct version
// directory in the release path directory (see ReleasePath).
func MakeAnwork(version int) (*Anwork, error) {
	path := makeAnworkZipPath(version)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	hash, err := getAnworkZipHash(path)
	if err != nil {
		return nil, err
	}

	reader, err := makeAnworkZipReader(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	unzipPath := makeAnworkDestinationPath(hash)
	err = unzip(reader, unzipPath)
	if err != nil {
		return nil, err
	}

	binary, exists := findBinary(version, unzipPath)
	if !exists {
		return nil, errors.New("Cannot find anwork binary at destinationPath: " + binary)
	}

	err = os.Chmod(binary, os.ModePerm)
	if err != nil {
		return nil, err
	}

	contextPath, err := makeContextPath()
	if err != nil {
		return nil, err
	}

	return &Anwork{contextPath: contextPath, binaryPath: binary}, nil
}

// Run a command with an instance of an anwork package. This function will return whatever the
// command printed to stdout, or a non-nil error is something failed.
func (anwork *Anwork) Run(command ...string) (string, error) {
	arguments := make([]string, 0, 2+len(command))
	arguments = append(arguments, "-o", anwork.contextPath)
	arguments = append(arguments, command...)
	cmd := exec.Command(anwork.binaryPath, arguments...)
	output, err := cmd.Output()
	return string(output), err
}

// Close an Anwork instance, i.e., delete the context directory for this Anwork instance. This
// Anwork instance will not be able to be used after this method is called.
func (anwork *Anwork) Close() error {
	err := os.RemoveAll(anwork.contextPath)
	anwork.binaryPath = ""
	return err
}

func makeAnworkZipPath(version int) string {
	return fmt.Sprintf("%s/v%d/anwork-%d.zip", ReleasePath, version, version)
}

func getAnworkZipHash(path string) (string, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	hash := md5.Sum(contents)
	bytes := make([]string, len(hash))
	for i, bite := range hash {
		bytes[i] = fmt.Sprintf("%x", bite)
	}

	return strings.Join(bytes, ""), nil
}

func makeAnworkDestinationPath(hash string) string {
	// This destination path is per executable so that we can run multiple test packages at the same
	// time and the two packages won't stomp on each other.
	executableName := path.Base(os.Args[0])
	return fmt.Sprintf("../.anwork-%s-%s", executableName, hash)
}

func makeAnworkZipReader(path string) (*zip.ReadCloser, error) {
	reader, err := zip.OpenReader(path)
	return reader, err
}

func unzip(reader *zip.ReadCloser, destinationPath string) error {
	var err error = nil

	unzipMutex.Lock()
	if _, err = os.Stat(destinationPath); os.IsNotExist(err) {
		err = reallyUnzip(reader, destinationPath)
	}
	unzipMutex.Unlock()

	return err
}

func reallyUnzip(reader *zip.ReadCloser, path string) error {
	// If the destination directory exists, then let's delete it.
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		if err = os.RemoveAll(path); err != nil {
			return errors.New("Anwork destination directory cannot be deleted: " + err.Error())
		}
	}

	if err := os.Mkdir(path, os.ModeDir|os.ModePerm); err != nil {
		return errors.New("Anwork destination directory cannot be created: " + err.Error())
	}

	for _, file := range reader.File {
		info := file.FileInfo()
		var err error = nil
		if info.IsDir() {
			err = handleDir(path, file, &info)
		} else {
			err = handleFile(path, file, &info)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func handleDir(destinationPath string, file *zip.File, info *os.FileInfo) error {
	dirPath := path.Join(destinationPath, file.Name)
	return os.Mkdir(dirPath, os.ModeDir|os.ModePerm)
}

func handleFile(destinationPath string, file *zip.File, info *os.FileInfo) error {
	filePath := path.Join(destinationPath, file.Name)
	osFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer osFile.Close()

	ioReader, err := file.Open()
	if err != nil {
		return err
	}
	defer ioReader.Close()

	buffer := make([]byte, 1024) // read size
	for {
		readCount, err := ioReader.Read(buffer)
		if readCount > 0 {
			_, err := osFile.Write(buffer[:readCount])
			if err != nil {
				return nil
			}
		}

		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
	}

	return nil
}

func findBinary(version int, packageRoot string) (string, bool) {
	packageName := fmt.Sprintf("anwork-%d", version)
	binaryPath := path.Join(packageRoot, packageName, "bin", "anwork")
	_, err := os.Stat(binaryPath)
	return binaryPath, !os.IsNotExist(err)
}

func makeContextPath() (string, error) {
	maxRandom := big.NewInt(math.MaxUint16)
	random, err := rand.Int(rand.Reader, maxRandom)
	if err != nil {
		return "", err
	}
	path := fmt.Sprintf("tmp_%04x", random.Int64())

	// If the destination path exists, let's blow up!
	if _, err = os.Stat(path); !os.IsNotExist(err) {
		return "", errors.New("Context directory already exists!")
	}

	return path, nil
}
