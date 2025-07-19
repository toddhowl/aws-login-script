# Copilot Instructions for aws-login-script

## Project Overview
- This repo provides two main entry points for securely logging in to the AWS CLI using credentials stored in the macOS Keychain:
  - `aws-login.sh`: Zsh script for login, S3 usage, and (optionally) billing.
  - `aws-login.go`: Go CLI app with the same functionality.
- Both tools automate AWS credential setup, S3 usage reporting, and (optionally) AWS billing cost retrieval.
- Credentials are never echoed to the terminal and are only written to the standard AWS CLI config file.

## Key Workflows
- **Login and S3 Usage:**
  - `aws-login` or `./aws-login` logs in and prints S3 bucket usage.
- **Login, S3, and Billing:**
  - `aws-login $` or `./aws-login $` logs in, prints S3 usage, and queries current month AWS billing (may incur API costs).
- **Credential Source:**
  - Credentials are retrieved from a macOS Keychain entry named `aws-cli` (Service field), with the Access Key ID in the Account field and Secret Access Key in the Password field.
- **Shell Alias:**
  - Users are encouraged to alias the script or binary in their `~/.zshrc` for convenience.

## Build & Usage
- **Go App:**
  - Build: `go build -o aws-login aws-login.go`
  - Run: `./aws-login` or `./aws-login $`
- **Shell Script:**
  - Make executable: `chmod +x aws-login.sh`
  - Run: `zsh aws-login.sh` or `zsh aws-login.sh $`

## Patterns & Conventions
- **No secrets in logs:**
  - Only the Access Key ID is ever printed; the Secret Access Key is never echoed.
- **macOS Keychain integration:**
  - Uses the `security` CLI to fetch credentials.
- **S3 Usage:**
  - Uses `aws s3 ls ... --recursive --summarize` and parses `Total Size:` for each bucket.
- **Billing:**
  - Uses `aws ce get-cost-and-usage` and parses JSON output (with `jq` in shell, native in Go).
  - Billing is only checked if `$` is passed as the first argument.
- **.gitignore:**
  - Excludes `.aws/`, build artifacts, `.vscode/`, `.github/`, and other common sensitive or local files.

## Integration Points
- **AWS CLI:**
  - Assumes AWS CLI is installed and configured in the user's PATH.
- **Go:**
  - No external Go dependencies beyond the standard library.
- **Shell:**
  - Uses standard macOS tools (`security`, `jq` if available).

## Examples
- See `README.md` for setup, usage, and troubleshooting.
- See `aws-login.sh` and `aws-login.go` for implementation details and argument handling.

---

If you add new features, update this file to document new workflows, conventions, or integration points.
