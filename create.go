package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func CreatePacket(jsonPath string) error {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("cant read file: %w", err)
	}

	var packet PacketConfig
	if err := json.Unmarshal(data, &packet); err != nil {
		return fmt.Errorf("json parse error: %w", err)
	}

	files, err := ResolveTargets(packet.Targets)
	if err != nil {
		return fmt.Errorf("resolve targets failed: %w", err)
	}

	fmt.Println("files for packaging:")
	for _, f := range files {
		fmt.Println("  -", f)
	}

	outputPath := fmt.Sprintf("archive/%s_%s.tar.gz", packet.Name, packet.Version)
	os.MkdirAll("archive", os.ModePerm)

	err = CreateTarGz(outputPath, files)
	if err != nil {
		return fmt.Errorf("archiving error: %w", err)
	}

	remotePath := fmt.Sprintf("/home/user/packages/%s_%s.tar.gz", packet.Name, packet.Version)

	err = UploadFileToServer("127.0.0.1:2222", "user", "qwerty", outputPath, remotePath)
	if err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}

	return nil
}
