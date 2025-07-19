# AWS CLI Login Script with macOS Keychain (Shell & Go)

This project provides two tools for securely logging in to the AWS CLI using credentials stored in the macOS Keychain:
- `aws-login.sh`: Zsh script for login, S3 usage, and (optionally) billing.
- `aws-login.go`: Go CLI app with the same functionality.

## Features

- Securely retrieves AWS credentials from the macOS Keychain (never echoes secrets).
- Sets up AWS CLI credentials for the default profile.
- Prints S3 bucket usage (total bytes per bucket).
- Optionally prints current month AWS billing (if you pass `$` as the first argument).

## Setup

### 1. Add AWS Credentials to Keychain

1. Open **Keychain Access** on your Mac.
2. Click the **+** button to add a new item.
3. Set:
   - **Keychain**: login
   - **Kind**: Internet password
   - **Service**: `aws-cli`
   - **Account**: Your AWS Access Key ID (e.g., `AKIA...`)
   - **Password**: Your AWS Secret Access Key
   - **Where**: (leave blank)
4. Save.

### 2. Use the Zsh Script

- **Prerequisite:** `jq` must be installed for billing parsing. Install with:
  ```sh
  brew install jq
  ```
- Make executable:
  ```sh
  chmod +x /<user>/Documents/Git/aws-login-script/aws-login.sh
  ```
- Add to your `~/.zshrc`:
  ```sh
  alias aws-login='zsh /<user>/Documents/Git/aws-login-script/aws-login.sh'
  source ~/.zshrc
  ```
- Usage:
  ```sh
  aws-login           # log in and check S3 only
  aws-login $         # log in, check S3, and check current month costs
  ```

### 3. Use the Go App

- Install Go if needed (`brew install go`).
- Build:
  ```sh
  cd /<user>/Documents/Git/aws-login-script
  go build -o aws-login aws-login.go
  ```
- Usage:
  ```sh
  ./aws-login         # log in and check S3 only
  ./aws-login $       # log in, check S3, and check current month costs
  ```

## Notes

- **Billing API calls may incur AWS charges.** Only run with `$` if you want to check costs.
- Credentials are only written to `~/.aws/credentials` (standard AWS CLI location).
- No secrets are ever printed to the terminal.

## Troubleshooting

- If you see a Keychain error, double-check the Service field is `aws-cli` and the entry is in the login keychain.
- If you see `InvalidAccessKeyId`, verify your credentials are correct and active in AWS.
- For billing info, your IAM user must have Cost Explorer permissions and billing access enabled.

## Security

- Credentials are stored securely in the macOS Keychain.
- The AWS CLI stores credentials in `~/.aws/credentials` in plaintext. Restrict file permissions or remove after use for extra security.

---

Feel free to extend for multiple profiles or other customizations as needed.
