package lib

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// DataDir returns the wcup data directory (~/.wcup).
func DataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, wcupDirName), nil
}

// ExportData writes all wcup data to a gzip-compressed tar archive.
func ExportData(destPath string) error {
	if err := EnsureResultsFile(); err != nil {
		return err
	}
	if err := EnsureParticipantsDir(); err != nil {
		return err
	}

	file, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	results, err := resultsPath()
	if err != nil {
		return err
	}
	if err := addFileToTar(tarWriter, results, resultsFileName); err != nil {
		return fmt.Errorf("export results: %w", err)
	}

	participants, err := participantsDir()
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(participants)
	if err != nil {
		return fmt.Errorf("export participants: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := filepath.Join(participants, entry.Name())
		archiveName := filepath.ToSlash(filepath.Join(participantsDirName, entry.Name()))
		if err := addFileToTar(tarWriter, path, archiveName); err != nil {
			return fmt.Errorf("export participant %q: %w", entry.Name(), err)
		}
	}

	return nil
}

// ImportData replaces all wcup data with the contents of a gzip-compressed tar archive.
func ImportData(srcPath string) error {
	if err := clearDataDir(); err != nil {
		return err
	}
	if err := extractArchive(srcPath); err != nil {
		return err
	}
	if err := LoadResults(); err != nil {
		return err
	}
	return LoadParticipants()
}

func addFileToTar(tw *tar.Writer, path, archiveName string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if info.IsDir() {
		return nil
	}

	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}
	header.Name = archiveName

	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	content, err := os.Open(path)
	if err != nil {
		return err
	}
	defer content.Close()

	_, err = io.Copy(tw, content)
	return err
}

func clearDataDir() error {
	dataDir, err := DataDir()
	if err != nil {
		return err
	}

	results, err := resultsPath()
	if err != nil {
		return err
	}
	if err := os.Remove(results); err != nil && !os.IsNotExist(err) {
		return err
	}

	participants, err := participantsDir()
	if err != nil {
		return err
	}
	if entries, err := os.ReadDir(participants); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			path := filepath.Join(participants, entry.Name())
			if err := os.Remove(path); err != nil {
				return err
			}
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(participants, 0o755); err != nil {
		return err
	}

	return nil
}

func extractArchive(srcPath string) error {
	file, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("invalid archive: %w", err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	dataDir, err := DataDir()
	if err != nil {
		return err
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read archive: %w", err)
		}
		if header.Typeflag != tar.TypeReg {
			continue
		}

		relPath, err := validateArchiveEntry(header.Name)
		if err != nil {
			return err
		}

		destPath := filepath.Join(dataDir, relPath)
		if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
			return err
		}

		out, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
		if err != nil {
			return err
		}

		if _, err := io.Copy(out, tarReader); err != nil {
			out.Close()
			return err
		}
		if err := out.Close(); err != nil {
			return err
		}
	}

	if err := EnsureResultsFile(); err != nil {
		return err
	}
	return EnsureParticipantsDir()
}

func validateArchiveEntry(name string) (string, error) {
	name = filepath.ToSlash(filepath.Clean(name))
	if name == "." || strings.HasPrefix(name, "../") || strings.Contains(name, "/../") {
		return "", fmt.Errorf("invalid path in archive: %s", name)
	}

	switch {
	case name == resultsFileName:
		return name, nil
	case strings.HasPrefix(name, participantsDirName+"/"):
		participant := strings.TrimPrefix(name, participantsDirName+"/")
		if participant == "" || strings.Contains(participant, "/") {
			return "", fmt.Errorf("invalid participant path in archive: %s", name)
		}
		return name, nil
	default:
		return "", fmt.Errorf("unexpected path in archive: %s", name)
	}
}
