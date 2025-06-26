package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

// AWSCredentials holds the AWS access key ID and secret access key
type AWSCredentials struct {
	AccessKeyID     string
	SecretAccessKey string
}

// getAWSCredentialsFromKeychain retrieves AWS credentials from the macOS Keychain
func getAWSCredentialsFromKeychain() (AWSCredentials, error) {
	var creds AWSCredentials

	// Get Access Key ID from Keychain
	cmd := exec.Command("security", "find-internet-password", "-s", "aws-cli", "-g")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return creds, fmt.Errorf("failed to find aws-cli keychain entry: %v", err)
	}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "\"acct\"") {
			parts := strings.Split(line, "\"")
			if len(parts) > 3 {
				creds.AccessKeyID = parts[3]
			}
		}
	}
	if creds.AccessKeyID == "" {
		return creds, fmt.Errorf("could not extract AWS Access Key ID from keychain")
	}

	// Get Secret Access Key from Keychain
	cmd = exec.Command("security", "find-internet-password", "-s", "aws-cli", "-w")
	out, err = cmd.CombinedOutput()
	if err != nil {
		return creds, fmt.Errorf("failed to get secret access key: %v", err)
	}
	creds.SecretAccessKey = strings.TrimSpace(string(out))
	if creds.SecretAccessKey == "" {
		return creds, fmt.Errorf("could not extract AWS Secret Access Key from keychain")
	}

	return creds, nil
}

// setAWSCredentials writes the credentials to the AWS CLI credentials file
func setAWSCredentials(creds AWSCredentials) error {
	usr, err := user.Current()
	if err != nil {
		return err
	}
	awsDir := filepath.Join(usr.HomeDir, ".aws")
	os.MkdirAll(awsDir, 0700)
	credFile := filepath.Join(awsDir, "credentials")
	profile := "default"
	content := fmt.Sprintf("[%s]\naws_access_key_id = %s\naws_secret_access_key = %s\n", profile, creds.AccessKeyID, creds.SecretAccessKey)
	return os.WriteFile(credFile, []byte(content), 0600)
}

// getBilling queries AWS Cost Explorer for the current month's billing
func getBilling() {
	// Get the first day of the current month and today
	start := time.Now().Format("2006-01-02")[:8] + "01"
	end := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	cmd := exec.Command("aws", "ce", "get-cost-and-usage",
		"--time-period", fmt.Sprintf("Start=%s,End=%s", start, end),
		"--granularity", "MONTHLY",
		"--metrics", "UnblendedCost",
		"--region", "us-east-1",
		"--output", "json",
	)
	out, err := cmd.CombinedOutput()
	if err != nil || len(out) == 0 {
		fmt.Println("(Billing info unavailable: check permissions or Cost Explorer access)")
		return
	}
	// Parse JSON output for amount and currency
	type costResult struct {
		ResultsByTime []struct {
			Total struct {
				UnblendedCost struct {
					Amount string `json:"Amount"`
					Unit   string `json:"Unit"`
				} `json:"UnblendedCost"`
			} `json:"Total"`
		} `json:"ResultsByTime"`
	}
	var result costResult
	if err := json.Unmarshal(out, &result); err == nil && len(result.ResultsByTime) > 0 {
		amount := result.ResultsByTime[0].Total.UnblendedCost.Amount
		currency := result.ResultsByTime[0].Total.UnblendedCost.Unit
		if amount != "" && currency != "" {
			fmt.Printf("Current month AWS cost: %s %s\n", amount, currency)
			return
		}
	}
	fmt.Println("(Could not parse billing info)")
}

// getS3Usage lists all S3 buckets and prints the total size used in each
func getS3Usage() {
	cmd := exec.Command("aws", "s3api", "list-buckets", "--query", "Buckets[].Name", "--output", "text")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("(Could not list S3 buckets: %v)\n", err)
		return
	}
	buckets := strings.Fields(string(out))
	for _, bucket := range buckets {
		cmd := exec.Command("aws", "s3", "ls", fmt.Sprintf("s3://%s", bucket), "--recursive", "--summarize")
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("  %s: error retrieving bucket size: %v\n", bucket, err)
			continue
		}
		size := "0"
		for _, line := range strings.Split(string(out), "\n") {
			if strings.Contains(line, "Total Size:") {
				parts := strings.Fields(line)
				if len(parts) >= 3 {
					size = parts[2]
				}
			}
		}
		fmt.Printf("%s: %s bytes\n", bucket, size)
	}
}

// main is the entry point of the application
func main() {
	fmt.Println("Retrieving AWS credentials from Keychain...")
	creds, err := getAWSCredentialsFromKeychain()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Printf("AWS Access Key ID: %s\n", creds.AccessKeyID)
	// Do not print the secret key
	fmt.Println("Setting AWS CLI credentials...")
	if err := setAWSCredentials(creds); err != nil {
		fmt.Println("Error writing credentials:", err)
		os.Exit(1)
	}
	fmt.Println("AWS CLI credentials set from KeyChain.")

	// Check for billing flag: if first argument is "$", check billing
	checkCost := false
	if len(os.Args) > 1 && os.Args[1] == "$" {
		checkCost = true
	}

	if checkCost {
		fmt.Println("\n==== AWS Billing for Current Month ====")
		getBilling()
	}

	fmt.Println("\n==== S3 Bucket Usage (bytes) ====")
	getS3Usage()
}
