package core

import (
	"archive/zip"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestMakeAnwork(t *testing.T) {
	anwork, err := MakeAnwork(1) // version 1
	if err != nil {
		t.Fatal("Failed to make anwork struct: ", err)
	}

	out, err := anwork.Run("version")
	if err != nil {
		t.Fatal("Failed to successfully run anwork command:", err)
	} else if len(out) == 0 {
		t.Fatal("Failed to properly get the output from the version command:", out)
	}
}

func TestCloseAnwork(t *testing.T) {
	anwork, err := MakeAnwork(1) // version 1
	if err != nil {
		t.Fatal("Failed to make anwork struct: ", err)
	}

	anwork.Close()

	_, err = anwork.Run("version")
	if err == nil {
		t.Fatal("We should have returned an error for a closed Anwork struct!")
	}
}

func TestAnworkZipPath(t *testing.T) {
	path := makeAnworkZipPath(1) // version 1
	if !fileExists(path) {
		t.Fatal("Zip path (", path, ") does not exist")
	}
}

func TestAnworkZipReaderCreation(t *testing.T) {
	path := makeAnworkZipPath(1) // version 1
	_, err := makeAnworkZipReader(path)
	if err != nil {
		t.Fatal("Zip reader cannot be created from path (", path, "):", err)
	}
}

func TestCreatingTempDirectory(t *testing.T) {
	path, err := makeDestinationDirectory()
	if err != nil {
		t.Fatal("Did not successfully create destination directory", err)
	}

	if !fileExists(path) {
		t.Fatal("Did not actually create the destination directory")
	}

	os.RemoveAll(path)
}

func TestUnzip(t *testing.T) {
	testZipPath := "test.zip"
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

	err = unzip(reader, tmpDirPath)
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
