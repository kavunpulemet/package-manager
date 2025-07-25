package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type ResolvedTarget struct {
	Included []string
	Skipped  []string
}

func ResolveTargets(rawTargets []interface{}) (files []string, err error) {
	for _, entry := range rawTargets {
		switch v := entry.(type) {
		case string:
			matches, err := filepath.Glob(v)
			if err != nil {
				return nil, fmt.Errorf("glob error: %w", err)
			}
			files = append(files, matches...)
		case map[string]interface{}:
			path, ok := v["path"].(string)
			if !ok {
				return nil, fmt.Errorf("missing 'path' in target")
			}
			matches, err := filepath.Glob(path)
			if err != nil {
				return nil, fmt.Errorf("glob error: %w", err)
			}

			exclude := ""
			if ex, ok := v["exclude"].(string); ok {
				exclude = ex
			}

			for _, file := range matches {
				if exclude != "" {
					exMatches, _ := filepath.Match(exclude, filepath.Base(file))
					if exMatches {
						continue
					}
				}
				files = append(files, file)
			}
		default:
			tmp, _ := json.Marshal(v)
			var target Target
			json.Unmarshal(tmp, &target)

			matches, err := filepath.Glob(target.Path)
			if err != nil {
				return nil, fmt.Errorf("glob error: %w", err)
			}
			for _, file := range matches {
				if target.Exclude != "" {
					exMatches, _ := filepath.Match(target.Exclude, filepath.Base(file))
					if exMatches {
						continue
					}
				}
				files = append(files, file)
			}
		}
	}

	return removeDuplicates(files), nil
}

func removeDuplicates(input []string) []string {
	seen := map[string]bool{}
	var result []string
	for _, f := range input {
		if !seen[f] {
			result = append(result, f)
			seen[f] = true
		}
	}
	return result
}

func CreateTarGz(outputPath string, files []string) error {
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	gzw := gzip.NewWriter(outFile)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	for _, file := range files {
		err = addFileToTar(tw, file)
		if err != nil {
			return fmt.Errorf("error while adding %s: %w", file, err)
		}
	}

	return nil
}

func addFileToTar(tw *tar.Writer, filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	headerName := filepath.Base(filePath)

	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}
	header.Name = headerName

	if err := tw.WriteHeader(header); err != nil {
		return err
	}
	_, err = io.Copy(tw, file)
	return err
}

func ExtractTarGz(archivePath, outputDir string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(outputDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			os.MkdirAll(target, os.FileMode(header.Mode))
		case tar.TypeReg:
			os.MkdirAll(filepath.Dir(target), os.ModePerm)
			outFile, err := os.Create(target)
			if err != nil {
				return err
			}
			_, err = io.Copy(outFile, tr)
			outFile.Close()
			if err != nil {
				return err
			}
		default:
		}
	}
	return nil
}
