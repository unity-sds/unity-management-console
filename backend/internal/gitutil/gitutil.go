package gitutil

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// Downloader handles HTTP downloads
type Downloader struct {
	client *http.Client
}

// NewDownloader creates a new downloader instance
func NewDownloader() *Downloader {
	return &Downloader{
		client: &http.Client{},
	}
}

// DownloadFile downloads a file from a URL to a local path
func (d *Downloader) DownloadFile(url, filepath string) error {
	// Get the data
	resp, err := d.client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// CloneRepo clones a git repository to a local directory
func CloneRepo(gitURL, destDir string) error {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git is not installed or not in PATH")
	}

	// Parse git URL format (git::https://...)
	url := gitURL
	if strings.HasPrefix(url, "git::") {
		url = strings.TrimPrefix(url, "git::")
	}

	// Run git clone
	cmd := exec.Command("git", "clone", "--depth", "1", url, destDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// CheckoutRef checks out a specific ref in a git repository
func CheckoutRef(repoDir, ref string) error {
	cmd := exec.Command("git", "checkout", ref)
	cmd.Dir = repoDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git checkout failed: %w\nOutput: %s", err, string(output))
	}
	return nil
}