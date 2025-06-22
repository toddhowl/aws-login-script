# AWS CLI Login Script with macOS Keychain (Shell & Go)

This repo allows you to securely log in to the AWS CLI using credentials stored in the macOS Keychain. Your AWS Access Key ID and Secret Access Key are never stored in plaintext files by the login process.

You can use either the provided Zsh script or the Go app to automate login, check your current AWS billing, and see S3 bucket usage.

## Setup Instructions

### 1. Add AWS Credentials to Keychain

1. Open **Keychain Access** on your Mac.
2. Click the **+** button to add a new item.
3. Set the following fields:
   - **Keychain**: login
   - **Kind**: Internet password
   - **Service**: `aws-cli` (this is the most important field; it must match exactly)
   - **Account**: Your AWS Access Key ID (e.g., `AKIA...`)
   - **Password**: Your AWS Secret Access Key
   - **Where**: (leave blank)
4. Save the entry.

---

## Option 1: Zsh Script

- Save the `aws-login.sh` script in this directory.
- Make it executable:
  ```sh
  chmod +x /<user>/Documents/Git/aws-login-script/aws-login.sh
  ```
- Add a shortcut to your `~/.zshrc`:
  ```sh
  alias aws-login='zsh /<user>/Documents/Git/aws-login-script/aws-login.sh'
  source ~/.zshrc
  ```
- Run:
  ```sh
  aws-login
  ```

---

## Option 2: Go App

- Make sure you have Go installed (`brew install go` if needed).
- Build the app:
  ```sh
  cd /<user>/Documents/Git/aws-login-script
  go build -o aws-login aws-login.go
  ```
- Run:
  ```sh
  ./aws-login
  ```

---

## What Happens
- Retrieves your AWS Access Key ID and Secret Access Key from the Keychain.
- Sets them in your AWS CLI configuration for the default profile.
- Shows your current month's AWS billing (if you have Cost Explorer permissions).
- Lists each S3 bucket and displays its total space used in bytes.

## Troubleshooting
- If you see an error about the Keychain entry not being found, double-check that the **Service** field is exactly `aws-cli` and that the entry is in the **login** keychain.
- If you see `InvalidAccessKeyId`, verify that your credentials are correct and active in AWS.
- For billing info, your IAM user must have Cost Explorer permissions and billing access enabled.

## Security
- Your credentials are stored securely in the macOS Keychain and never written to disk in plaintext by this script or app.
- The AWS CLI itself does store credentials in `~/.aws/credentials` in plaintext. Restrict file permissions or remove after use for extra security.

---

Feel free to modify the script or Go app for multiple profiles or other customizations as needed.
