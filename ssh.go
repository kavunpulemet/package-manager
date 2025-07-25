package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func UploadFileToServer(addr, user, pass, localPath, remotePath string) error {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return fmt.Errorf("SSH dial failed: %w", err)
	}
	defer conn.Close()

	client, err := sftp.NewClient(conn)
	if err != nil {
		return fmt.Errorf("SFTP client failed: %w", err)
	}
	defer client.Close()

	srcFile, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("cant open local file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := client.Create(remotePath)
	if err != nil {
		return fmt.Errorf("cant create remote file: %w", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("copy failed: %w", err)
	}

	fmt.Println("archive successfully sent to server:", remotePath)
	return nil
}

func ListFilesFromServer(addr, user, pass, remoteDir string) ([]string, error) {
	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(pass)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("ssh dial error: %w", err)
	}
	defer conn.Close()

	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, fmt.Errorf("sftp error: %w", err)
	}
	defer client.Close()

	files, err := client.ReadDir(remoteDir)
	if err != nil {
		return nil, fmt.Errorf("readdir error: %w", err)
	}

	var names []string
	for _, f := range files {
		if !f.IsDir() {
			names = append(names, f.Name())
		}
	}
	return names, nil
}

func DownloadFileFromServer(addr, user, pass, remotePath, localPath string) error {
	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(pass)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return fmt.Errorf("SSH dial failed: %w", err)
	}
	defer conn.Close()

	client, err := sftp.NewClient(conn)
	if err != nil {
		return fmt.Errorf("SFTP client failed: %w", err)
	}
	defer client.Close()

	srcFile, err := client.Open(remotePath)
	if err != nil {
		return fmt.Errorf("remote open failed: %w", err)
	}
	defer srcFile.Close()

	os.MkdirAll(filepath.Dir(localPath), os.ModePerm)
	dstFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("local create failed: %w", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
