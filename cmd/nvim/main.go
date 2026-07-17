package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const nvimUrl = "https://github.com/neovim/neovim/releases/latest/download/nvim-win64.zip"

func main() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	binDir := filepath.Dir(exePath)
	nvimDir := filepath.Join(binDir, "nvim-win64")
	nvimExe := filepath.Join(nvimDir, "bin", "nvim.exe")

	if _, err := os.Stat(nvimExe); os.IsNotExist(err) {
		fmt.Println("Neovim is not installed. Downloading the official binary...")
		if err := downloadAndExtract(binDir); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to install Neovim: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Neovim installed successfully!")
	}

	cmd := exec.Command(nvimExe, os.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		os.Exit(cmd.ProcessState.ExitCode())
	}
}

func downloadAndExtract(destDir string) error {
	zipPath := filepath.Join(destDir, "nvim-temp.zip")
	
	// Download
	resp, err := http.Get(nvimUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	outFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return err
	}
	outFile.Close()
	
	// Extract
	defer os.Remove(zipPath)
	return unzip(zipPath, destDir)
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		
		// Check for ZipSlip
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			continue
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}
