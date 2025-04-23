package update

import (
	"archive/zip"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-management-console/backend/internal/aws"
)

// Asset represents a release asset in GitHub
type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// Release represents a GitHub release
type Release struct {
	TagName string  `json:"tag_name"`
	Name    string  `json:"name"`
	Assets  []Asset `json:"assets"`
}

func getLatestMCRelease(accessToken string) (*Release, error) {
	owner := "unity-sds"
	repo := "unity-management-console"

	latestReleaseURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	req, err := http.NewRequest("GET", latestReleaseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set the User-Agent header which is required by GitHub
	req.Header.Set("User-Agent", "Unity-Management-Console-Updater")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	// Only set auth header if we have a token
	if accessToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	} else {
		log.Info("Making unauthenticated GitHub API request")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest release: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get latest release: status %d, body: %s", resp.StatusCode, string(body))
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode release: %v", err)
	}

	log.Infof("Latest release is: %s (%s)", release.Name, release.TagName)

	return &release, nil
}

func getMCAssetByName(release Release, assetName string) (string, error) {
	var assetURL string
	var downloadAssetName string

	if assetName == "" && len(release.Assets) > 0 {
		// If no specific asset is requested, download the first one
		assetURL = release.Assets[0].BrowserDownloadURL
		downloadAssetName = release.Assets[0].Name
		fmt.Printf("No specific asset requested. Downloading: %s\n", downloadAssetName)
	} else {
		// Find the requested asset
		for _, asset := range release.Assets {
			if asset.Name == assetName {
				assetURL = asset.BrowserDownloadURL
				downloadAssetName = assetName
				break
			}
		}
	}

	if assetURL == "" {
		if assetName != "" {
			return "", fmt.Errorf("asset %s not found in latest release", assetName)
		} else {
			return "", fmt.Errorf("no assets found in latest release")
		}
	}

	return assetURL, nil
}

func downloadMCAsset(accessToken string, assetName string, assetURL string) (string, error) {
	req, err := http.NewRequest("GET", assetURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "unity-release-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %w", err)
	}

	outfilePath := filepath.Join(tempDir, assetName)
	outfile, err := os.Create(outfilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer outfile.Close()

	// Download asset
	// Set required headers for GitHub API
	req.Header.Set("User-Agent", "Unity-Management-Console-Updater")
	req.Header.Set("Accept", "application/octet-stream")

	// Only set auth header if we have a token
	if accessToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download asset: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to download asset: status %d, body: %s", resp.StatusCode, string(body))
	}

	_, err = io.Copy(outfile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	log.Infof("Release downloaded to %s", outfilePath)
	return outfilePath, nil
}

func UpdateManagementConsoleInPlace() error {
	assetName := "managementconsole.zip"

	// Get Github Application credentials from SSM
	appID, privateKey, err := getGitHubCredsFromSSM()
	if err != nil {
		return err
	}

	// Get Github installation token from SSM
	installationID, err := getGitHubAppInstallationIDFromSSM()
	if err != nil {
		return err
	}

	// Get JWT for authentication
	appJWT, err := createJWTForGitHubApp(appID, privateKey)
	if err != nil {
		return err
	}

	// Exchange JWT, installation ID for an installation token
	accessToken, err := getInstallationToken(appJWT, installationID)
	if err != nil {
		return err
	}

	log.Infof("Got installation token...")
	latestRelease, err := getLatestMCRelease(accessToken)
	if err != nil {
		return err
	}

	assetURL, err := getMCAssetByName(*latestRelease, assetName)
	if err != nil {
		return err
	}

	downloadedZipfilePath, err := downloadMCAsset(accessToken, assetName, assetURL)
	if err != nil {
		return err
	}

	// Unzip into parent dir
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Failed to determine current working directory")
	}

	destDir := filepath.Join(currentDir, "../updated-management-console")
	err = unzipFile(downloadedZipfilePath, destDir)
	if err != nil {
		log.Errorf("Error unzipping. Unzipped %s to %s, err: %v", downloadedZipfilePath, destDir, err)
		return err
	}
	log.Infof("Release unzipped to: %s", destDir)

	// Delete ZIP now that we've successfully unzipped it.
	err = os.Remove(downloadedZipfilePath)
	if err != nil {
		return err
	}

	// If we get here, the new release has been successfully downloaded and unzipped.
	// Kick off the updater service. We use a service because we want a process that
	// can run independently of this one.

	// Run service command to start the update service
	cmd := exec.Command("sudo", "service", "managementconsole-update", "start")
	// Run the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.WithError(err).Error("Failed to start update service")
		log.Errorf("Command output: %s", string(output))
		return err
	}

	// Nothing to return, the updater will kill this process and start a new one.
	return nil
}

// getGitHubCredsFromSSM retrieves GitHub client ID and secret from SSM
func getGitHubCredsFromSSM() (int64, *rsa.PrivateKey, error) {
	// Get Github Application ID
	paramName := "/unity/mc/github-app-id"
	withDecryption := true

	log.Infof("Retrieving GitHub Application ID from SSM parameter: %s", paramName)
	paramOutput, err := aws.ReadSSMParameterWithDecryption(paramName, withDecryption)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to retrieve SSM parameter: %w", err)
	}

	if paramOutput.Parameter == nil || paramOutput.Parameter.Value == nil {
		return 0, nil, fmt.Errorf("empty value returned for SSM parameter")
	}
	appIDString := *paramOutput.Parameter.Value

	// Convert app ID string to int64
	var appID int64
	_, err = fmt.Sscanf(appIDString, "%d", &appID)
	if err != nil {
		return 0, nil, fmt.Errorf("Invalid GITHUB_APP_ID: %v\n", err)

	}

	// Get Gethub Application private key
	paramName = "/unity/mc/github-app-private-key"

	log.Infof("Retrieving GitHub Application Private Key from SSM parameter: %s", paramName)
	paramOutput, err = aws.ReadSSMParameterWithDecryption(paramName, withDecryption)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to retrieve SSM parameter: %w", err)
	}

	if paramOutput.Parameter == nil || paramOutput.Parameter.Value == nil {
		return 0, nil, fmt.Errorf("empty value returned for SSM parameter")
	}
	privateKeyString := *paramOutput.Parameter.Value
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyString))
	if err != nil {
		return 0, nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	return appID, privateKey, nil
}

func getGitHubAppInstallationIDFromSSM() (int64, error) {
	// Get Github Application ID
	paramName := "/unity/mc/github-app-installation-id"
	withDecryption := true

	log.Infof("Retrieving GitHub Installation ID from SSM parameter: %s", paramName)
	paramOutput, err := aws.ReadSSMParameterWithDecryption(paramName, withDecryption)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve SSM parameter: %w", err)
	}

	if paramOutput.Parameter == nil || paramOutput.Parameter.Value == nil {
		return 0, fmt.Errorf("empty value returned for SSM parameter")
	}
	idString := *paramOutput.Parameter.Value

	// Convert app ID string to int64
	var installationID int64
	_, err = fmt.Sscanf(idString, "%d", &installationID)
	if err != nil {
		return 0, fmt.Errorf("Invalid GITHUB_INSTALLATION_ID: %v\n", err)

	}
	return installationID, nil
}

// createJWTForGitHubApp creates a JWT for GitHub App authentication
func createJWTForGitHubApp(appID int64, privateKey *rsa.PrivateKey) (string, error) {
	// Token expires in 10 minutes (GitHub's maximum is 10 minutes)
	now := time.Now()
	expiresAt := now.Add(10 * time.Minute)

	// Create the JWT claims
	claims := jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		Issuer:    strconv.FormatInt(appID, 10),
	}

	// Create the token with claims and sign it
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return signedToken, nil
}

// GetInstallationToken exchanges a GitHub App JWT for an installation access token
func getInstallationToken(jwt string, installationID int64) (string, error) {
	// Create the HTTP client
	client := &http.Client{Timeout: 10 * time.Second}

	// Build the URL for requesting an installation token
	url := fmt.Sprintf("https://api.github.com/app/installations/%d/access_tokens", installationID)

	// Create a new POST request
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Add required headers
	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("GitHub API error: status %d, body: %s",
			resp.StatusCode, string(bodyBytes))
	}

	// Parse the response
	var tokenResponse struct {
		Token     string    `json:"token"`
		ExpiresAt time.Time `json:"expires_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return tokenResponse.Token, nil
}

// unzipFile extracts a zip file to a destination directory
func unzipFile(zipPath, destDir string) error {
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	// Extract each file
	for _, file := range zipReader.File {
		err := extractZipFile(file, destDir)
		if err != nil {
			return err
		}
	}

	return nil
}

// extractZipFile extracts a single file from a zip archive
func extractZipFile(file *zip.File, destDir string) error {
	// Clean the file path to prevent zip slip vulnerability
	filePath := filepath.Join(destDir, file.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(destDir)+string(os.PathSeparator)) {
		return fmt.Errorf("illegal file path: %s", filePath)
	}

	// Handle directories
	if file.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, file.Mode()); err != nil {
			return err
		}
		return nil
	}

	// Make sure the parent directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	// Open the file
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// Create the file
	outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Copy the contents
	_, err = io.Copy(outFile, rc)
	return err
}
