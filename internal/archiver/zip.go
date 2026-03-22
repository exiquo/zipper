package archiver

import (
	"archive/zip"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// CreateArchive creates a zip archive from the given source directory
// and writes it to the specified target file.
// The archive will contain the source directory as its root.
// An error is returned if validation fails or any file operation fails.
func CreateArchive(srcDirectory string, targetFile string) error {
	srcDirectory = filepath.Clean(srcDirectory)
	targetFile = filepath.Clean(targetFile)

	if err := validate(srcDirectory, targetFile); err != nil {
		return err
	}

	archive, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	files, err := listFiles(srcDirectory)
	if err != nil {
		return err
	}

	baseDir := filepath.Base(srcDirectory)

	for _, filePath := range files {
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(srcDirectory, filePath)
		if err != nil {
			file.Close()
			return err
		}

		zipPath := filepath.ToSlash(filepath.Join(baseDir, relativePath))

		writer, err := zipWriter.Create(zipPath)
		if err != nil {
			file.Close()
			return err
		}

		if _, err = io.Copy(writer, file); err != nil {
			file.Close()
			return err
		}

		if err := file.Close(); err != nil {
			return err
		}
	}

	return nil
}

func listFiles(directory string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

func validate(srcDirectory string, targetFile string) error {
	if err := validateSourceDirectory(srcDirectory); err != nil {
		return err
	}

	if err := validateTargetFile(targetFile); err != nil {
		return err
	}

	return validateTargetOutsideSource(srcDirectory, targetFile)
}

func validateSourceDirectory(srcDirectory string) error {
	if srcDirectory == "" {
		return errors.New("source directory is required")
	}

	info, err := os.Stat(srcDirectory)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("source directory does not exist")
		}
		return err
	}

	if !info.IsDir() {
		return errors.New("source path is not a directory")
	}

	return nil
}

func validateTargetFile(targetFile string) error {
	if targetFile == "" {
		return errors.New("target file is required")
	}

	if filepath.Ext(targetFile) != ".zip" {
		return errors.New("target file must have .zip extension")
	}

	targetDir := filepath.Dir(targetFile)

	info, err := os.Stat(targetDir)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("target directory does not exist")
		}
		return err
	}

	if !info.IsDir() {
		return errors.New("target directory is not a directory")
	}

	return nil
}

// validateTargetOutsideSource ensures that the target file
// is not located inside the source directory.
func validateTargetOutsideSource(srcDirectory string, targetFile string) error {
	absSrc, err := filepath.Abs(srcDirectory)
	if err != nil {
		return err
	}

	absTarget, err := filepath.Abs(targetFile)
	if err != nil {
		return err
	}

	rel, err := filepath.Rel(absSrc, absTarget)
	if err != nil {
		return err
	}

	// If the relative path doesn't start with "..", the target is located within the source directory
	if !strings.HasPrefix(rel, "..") {
		return errors.New("target file must not be inside source directory")
	}

	return nil
}
