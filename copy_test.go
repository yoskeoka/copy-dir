package cpdir_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/yoskeoka/cpdir"
)

func setupSrcDir(t *testing.T) (tempdir string, teardown func()) {
	t.Helper()

	tempdir, err := ioutil.TempDir("", "cpdir_test")
	if err != nil {
		t.Fatal(err)
	}

	os.MkdirAll(filepath.Join(tempdir, "dir1/subdir1"), 0755)
	ioutil.WriteFile(filepath.Join(tempdir, "dir1/file1"), []byte("file1"), 0644)
	ioutil.WriteFile(filepath.Join(tempdir, "dir1/subdir1/file2"), []byte("file2"), 0644)
	ioutil.WriteFile(filepath.Join(tempdir, "dir1/subdir1/file3"), []byte("file3"), 0644)

	os.MkdirAll(filepath.Join(tempdir, "dir2"), 0755)
	ioutil.WriteFile(filepath.Join(tempdir, "dir2/.file"), []byte(".file"), 0644)
	ioutil.WriteFile(filepath.Join(tempdir, "dir2/file.go"), []byte("file.go"), 0644)

	teardown = func() {
		os.RemoveAll(tempdir)
	}
	return tempdir, teardown
}

func setupDestDir(t *testing.T) (tempdir string, teardown func()) {
	t.Helper()

	tempdir, err := ioutil.TempDir("", "cpdir_test")
	if err != nil {
		t.Fatal(err)
	}

	teardown = func() {
		os.RemoveAll(tempdir)
	}
	return tempdir, teardown
}

func equalFile(t *testing.T, srcFile, destFile string) {
	t.Helper()

	srcFileContent, err := ioutil.ReadFile(srcFile)
	if err != nil {
		t.Fatalf("could not read src file %v", srcFile)
	}

	destFileContent, err := ioutil.ReadFile(destFile)
	if err != nil {
		t.Fatalf("could not read dest file %v", destFile)
	}

	if !bytes.Equal(srcFileContent, destFileContent) {
		t.Errorf("src file and dest file are not equal, src = %+v, dest = %+v", srcFileContent, destFileContent)
	}
}

func TestCopyFile(t *testing.T) {
	type args struct {
		src  string
		dest string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"file1", args{"dir1/file1", "file1_copy"}, false},
		{".file", args{"dir2/.file", ".file_copy"}, false},
		{"file2", args{"dir1/subdir1/file2", "file2_copy"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srcdir, teardownSrc := setupSrcDir(t)
			defer teardownSrc()
			destdir, teardownDest := setupDestDir(t)
			defer teardownDest()

			srcFile := filepath.Join(srcdir, tt.args.src)
			destFile := filepath.Join(destdir, tt.args.dest)

			if err := cpdir.CopyFile(srcFile, destFile); (err != nil) != tt.wantErr {
				t.Errorf("CopyFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			equalFile(t, srcFile, destFile)
		})
	}
}

func TestCopyDirContents(t *testing.T) {
	type args struct {
		src  string
		dest string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"tmp", args{"", "tmp"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cwd, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}

			os.Setenv("TMPDIR", cwd)

			srcTemp, teardownSrc := setupSrcDir(t)
			defer teardownSrc()
			destTemp, teardownDest := setupDestDir(t)
			defer teardownDest()

			srcDir := filepath.Join(srcTemp, tt.args.src)
			t.Logf("src dir %s", srcDir)
			destDir := filepath.Join(destTemp, tt.args.dest)
			t.Logf("dest dir %s", destDir)
			os.MkdirAll(destDir, os.ModePerm)

			if err := cpdir.CopyDirContents(srcDir, destDir); (err != nil) != tt.wantErr {
				t.Errorf("CopyDirContents() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}
