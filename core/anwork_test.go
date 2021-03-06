package core

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

const (
	testZipPath      = "data/test.zip"
	otherTestZipPath = "data/other.zip"

	defaultVersion = 2
)

func TestMakeAnwork(t *testing.T) {
	t.Parallel()

	for version := 1; true; version++ {
		if !fileExists(makeAnworkZipPath(version)) {
			break
		}

		t.Run(fmt.Sprintf("Version%d", version), func(t *testing.T) {
			anwork, err := MakeAnwork(version)
			if err != nil {
				t.Fatal("Failed to make anwork struct: ", err)
			}
			defer anwork.Close()

			out, err := anwork.Run("version")
			if err != nil {
				t.Fatal("Failed to successfully run anwork command:", err)
			} else if len(out) == 0 {
				t.Fatal("Failed to properly get the output from the version command:", out)
			}
		})
	}
}

func TestCloseAnwork(t *testing.T) {
	t.Parallel()

	anwork, err := MakeAnwork(defaultVersion)
	if err != nil {
		t.Fatal("Failed to make anwork struct:", err)
	}

	anwork.Close()

	_, err = anwork.Run("version")
	if err == nil {
		t.Fatal("We should have returned an error for a closed Anwork struct!")
	}
}

func TestParallelAnworkCreation(t *testing.T) {
	const anworksCount = 4
	anworkChan := make(chan *Anwork, anworksCount)
	defer close(anworkChan)

	for i := 0; i < cap(anworkChan); i++ {
		go func(i int) {
			anwork, err := MakeAnwork(defaultVersion)
			if err != nil {
				t.Errorf("Failed to make %dth anwork struct: %s", i, err)
			} else {
				t.Logf("Created anwork struct: %s", anwork)
			}
			anworkChan <- anwork
		}(i)
	}

	for i := 0; i < cap(anworkChan); i++ {
		anwork := <-anworkChan
		if anwork == nil {
			// We can get here if an above call to MakeAnwork fails.
			continue
		}

		output, err := anwork.Run("version")
		if err != nil {
			t.Errorf("Failed to run anwork struct: %s. Output: %s", err, output)
		} else if len(output) == 0 {
			t.Errorf("Did not get any output from anwork struct")
		} else {
			t.Logf("Ran Anwork struct: %s", anwork)
		}

		if err := anwork.Close(); err != nil {
			t.Errorf("Could not close anwork struct: %s", anwork)
		}
	}
}

func TestParallelAnworkRunning(t *testing.T) {
	const anworksCount = 4
	anworkChan := make(chan *Anwork, anworksCount)
	defer close(anworkChan)
	ranChan := make(chan bool, anworksCount)
	defer close(ranChan)

	for i := 0; i < cap(anworkChan); i++ {
		anwork, err := MakeAnwork(defaultVersion)
		if err != nil {
			t.Fatalf("Failed to make %dth anwork struct: %s", i, err)
		}
		defer anwork.Close()
		anworkChan <- anwork
		t.Logf("Created anwork struct: %s", anwork)
	}

	for i := 0; i < cap(anworkChan); i++ {
		go func() {
			anwork := <-anworkChan
			output, err := anwork.Run("version")
			if err != nil {
				t.Errorf("Failed to run anwork struct: %s. Output: %s", err, output)
			} else if len(output) == 0 {
				t.Errorf("Did not get any output from anwork struct")
			} else {
				t.Logf("Ran Anwork struct: %s", anwork)
			}
			ranChan <- true
		}()
	}

	// Sync here so that we ensure that the above go functions all returned.
	for i := 0; i < cap(anworkChan); i++ {
		_ = <-ranChan
	}
}

func TestParallelAnworkCreationAndRunning(t *testing.T) {
	const anworksCount = 4
	createdChan := make(chan *Anwork, anworksCount)
	defer close(createdChan)
	ranChan := make(chan *Anwork, anworksCount)
	defer close(ranChan)

	for i := 0; i < cap(createdChan); i++ {
		go func(i int) {
			anwork, err := MakeAnwork(defaultVersion)
			if err != nil {
				t.Errorf("Failed to make %dth anwork struct: %s", i, err)
			} else {
				t.Logf("Created anwork struct: %s", anwork)
			}
			createdChan <- anwork
		}(i)
	}

	for i := 0; i < cap(ranChan); i++ {
		go func(i int) {
			anwork := <-createdChan
			output, err := anwork.Run("version")
			if err != nil {
				t.Errorf("Failed to run anwork struct: %s. Output: %s", err, output)
			} else if len(output) == 0 {
				t.Errorf("Did not get any output from anwork struct")
			} else {
				t.Logf("Ran Anwork struct: %s", anwork)
			}
			ranChan <- anwork
		}(i)
	}

	// Sync here so that we ensure that the above go functions all returned.
	for i := 0; i < cap(ranChan); i++ {
		anwork := <-ranChan
		anwork.Close()
	}
}

func TestAnworkZipPath(t *testing.T) {
	t.Parallel()

	path := makeAnworkZipPath(defaultVersion)
	if !fileExists(path) {
		t.Fatal("Zip path (", path, ") does not exist")
	}
}

func TestAnworkZipHash(t *testing.T) {
	t.Parallel()

	hash, err := getAnworkZipHash(testZipPath)
	t.Run("Single", func(t *testing.T) {
		if err != nil {
			t.Fatalf("Error when calculating zip file hash: %s", err)
		}
	})
	t.Run("Double", func(t *testing.T) {
		hashAgain, errAgain := getAnworkZipHash(testZipPath)
		if errAgain != nil {
			t.Fatalf("Error when calculcating zip file hash for the second time: %s", err)
		} else if hashAgain != hash {
			t.Fatalf("Expected for two hashes to be the same value: %s vs %s", hash, hashAgain)
		}
	})
	t.Run("BadFile", func(t *testing.T) {
		_, err = getAnworkZipHash("this/path/does/not/exist.zip")
		if err == nil {
			t.Fatal("Expected error from hasing a non-existent file!")
		}
	})
	t.Run("OtherFile", func(t *testing.T) {
		otherHash, err := getAnworkZipHash(otherTestZipPath)
		if err != nil {
			t.Fatal("Got an error from unzipping other test data!")
		} else if hash == otherHash {
			t.Fatalf("Hmmm...I don't think these two hashes should be equal...%s vs %s", hash, otherHash)
		} else {
			t.Logf("Hashes look like '%s' and '%s'", hash, otherHash)
		}
	})
}

func TestAnworkZipReaderCreation(t *testing.T) {
	t.Parallel()

	path := makeAnworkZipPath(defaultVersion)
	_, err := makeAnworkZipReader(path)
	if err != nil {
		t.Fatal("Zip reader cannot be created from path (", path, "):", err)
	}
}

func TestUnzip(t *testing.T) {
	t.Parallel()

	reader, err := zip.OpenReader(testZipPath)
	if err != nil {
		t.Fatal("Cannot create reader for zipfile:", err)
	}
	defer reader.Close()

	tmpDirPath := "tmp"
	if err = os.Mkdir(tmpDirPath, os.ModeDir|os.ModePerm); err != nil {
		t.Fatal("Could not create tmp directory:", err)
	}
	defer os.RemoveAll(tmpDirPath)

	err = reallyUnzip(reader, tmpDirPath)
	if err != nil {
		t.Fatal("Did not unzip file successfully:", err)
	}
	filepath.Walk(tmpDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			t.Logf("Failed to visit file %s: %s", path, err)
		} else {
			t.Logf("Visited file %s", path)
		}
		return nil
	})

	// This map mimics the file tree in the test.zip file.
	fileTree := map[string][]string{
		"":      []string{"file-1", "file-2"},
		"dir-a": []string{"file-a-1", "file-a-2"},
		"dir-b": []string{"file-b-1"},
	}
	checkFileTree(t, path.Join(tmpDirPath, "test"), fileTree)
}

func TestNonExistentAnworkVersion(t *testing.T) {
	t.Parallel()

	_, err := MakeAnwork(65535) // I sure hope this version never exists...
	if err == nil {
		t.Fatal("Should have received an error from bad anwork version!")
	}
}

func checkFileTree(t *testing.T, root string, fileTree map[string][]string) {
	for directory, files := range fileTree {
		if path := path.Join(root, directory); !fileExists(path) {
			t.Error("Did not find unzipped directory:", path)
		}
		for _, file := range files {
			filepath := path.Join(root, directory, file)
			if !fileExists(filepath) {
				t.Error("Did not find unzipped file:", filepath)
			} else {
				contents, err := ioutil.ReadFile(filepath)
				stringContents := string(contents)
				if err != nil {
					t.Error("Could not read contents of file:", err)
				} else if basename := path.Base(filepath); stringContents != basename {
					t.Error("The contents of the file (", stringContents,
						") did not match basename of file:", basename)
				}
			}
		}
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
