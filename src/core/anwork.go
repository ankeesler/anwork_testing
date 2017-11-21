// This package contains utilities for running Anwork tests.
package core

import (
	"archive/zip"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"os"
	"os/exec"
	"path"
)

// This is the path to where the Anwork release zip files are kept.
const ReleasePath string = "../../release"

// Anwork represents an Anwork program that can be executed.
type Anwork struct {
	// This is the path to the expanded package. If this path is the empty string, then this Anwork
	// struct has been Close'd and is no longer usable.
	packagePath string

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

	reader, err := makeAnworkZipReader(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	destinationPath, err := makeDestinationDirectory()
	if err != nil {
		return nil, err
	}

	err = unzip(reader, destinationPath)
	if err != nil {
		return nil, err
	}

	binary, exists := findBinary(version, destinationPath)
	if !exists {
		return nil, errors.New("Cannot find anwork binary at destinationPath: " + binary)
	}

	err = os.Chmod(binary, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return &Anwork{packagePath: destinationPath, binaryPath: binary}, nil
}

// Run a command with an instance of an anwork package. This function will return whatever the
// command printed to stdout, or a non-nil error is something failed.
func (anwork *Anwork) Run(command ...string) (string, error) {
	arguments := make([]string, 0, 2+len(command))
	arguments = append(arguments, "-o", anwork.packagePath)
	arguments = append(arguments, command...)
	cmd := exec.Command(anwork.binaryPath, arguments...)
	output, err := cmd.Output()
	return string(output), err
}

// Close an Anwork instance, i.e., delete the unexpanded package associated with it. This means the
// Anwork object will no longer be usable.
func (anwork *Anwork) Close() error {
	err := os.RemoveAll(anwork.packagePath)
	anwork.packagePath = ""
	return err
}

func makeAnworkZipPath(version int) string {
	return fmt.Sprintf("%s/v%d/anwork-%d.zip", ReleasePath, version, version)
}

func makeAnworkZipReader(path string) (*zip.ReadCloser, error) {
	reader, err := zip.OpenReader(path)
	return reader, err
}

func makeDestinationDirectory() (string, error) {
	maxRandom := big.NewInt(math.MaxUint16)
	random, err := rand.Int(rand.Reader, maxRandom)
	if err != nil {
		return "", err
	}
	path := fmt.Sprintf("tmp_%04x", random.Int64())

	// If the destination directory exists already, let's fail. Fail fast is good!
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return "", errors.New("Anwork destination directory is already in use: " + path)
	}

	if err := os.Mkdir(path, os.ModeDir|os.ModePerm); err != nil {
		return "", err
	}

	return path, nil
}

func unzip(reader *zip.ReadCloser, destinationPath string) error {
	for _, file := range reader.File {
		info := file.FileInfo()
		var err error
		if info.IsDir() {
			err = handleDir(destinationPath, file, &info)
		} else {
			err = handleFile(destinationPath, file, &info)
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

	buffer := make([]byte, 64) // read size
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
