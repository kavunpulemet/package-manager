package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func UpdatePackages(jsonPath string) error {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("cant read file: %w", err)
	}

	var list UpdateList
	if err := json.Unmarshal(data, &list); err != nil {
		return fmt.Errorf("json parse error: %w", err)
	}

	fmt.Println("updating packages:")
	for _, pkg := range list.Packages {
		fmt.Printf(" - %s %s\n", pkg.Name, pkg.Ver)
	}

	const (
		serverAddr = "127.0.0.1:2222"
		serverUser = "user"
		serverPass = "qwerty"
		serverDir  = "/home/user/packages"
	)

	files, err := ListFilesFromServer(serverAddr, serverUser, serverPass, serverDir)
	if err != nil {
		return fmt.Errorf("failed to get file list: %w", err)
	}

	for _, pkg := range list.Packages {
		var (
			prefix   = pkg.Name + "_"
			bestVer  = ""
			bestFile = ""
		)

		for _, f := range files {
			if !strings.HasPrefix(f, prefix) || !strings.HasSuffix(f, ".tar.gz") {
				continue
			}

			ver := strings.TrimSuffix(strings.TrimPrefix(f, prefix), ".tar.gz")
			if pkg.Ver == "" || versionMatches(pkg.Ver, ver) {
				if bestVer == "" || compareVersions(ver, bestVer) > 0 {
					bestVer = ver
					bestFile = f
				}
			}
		}

		if bestFile == "" {
			fmt.Println("archive not found for:", pkg.Name, pkg.Ver)
		} else {
			fmt.Printf("archive selected: %s (version %s)\n", bestFile, bestVer)
		}

		localPath := filepath.Join("downloads", bestFile)
		remotePath := serverDir + "/" + bestFile
		extractPath := filepath.Join("output", pkg.Name)

		os.MkdirAll("downloads", os.ModePerm)
		os.MkdirAll("output", os.ModePerm)

		fmt.Println(remotePath)
		err := DownloadFileFromServer(serverAddr, serverUser, serverPass, remotePath, localPath)
		if err != nil {
			fmt.Println("downloading error:", err)
			continue
		}

		err = ExtractTarGz(localPath, extractPath)
		if err != nil {
			fmt.Println("extracting error:", err)
			continue
		}

		fmt.Printf("packet %s installed in: %s\n", pkg.Name, extractPath)
	}

	return nil
}
