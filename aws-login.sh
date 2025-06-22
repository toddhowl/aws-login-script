#!/bin/zsh

# Retrieve AWS Access Key ID (stored as "User Name" in the "aws-cli" entry)
AWS_ACCESS_KEY_ID=$(security find-internet-password -s aws-cli -g 2>&1 | grep '"acct"' | awk -F'"' '{print $4}')
echo "AWS Access Key ID: $AWS_ACCESS_KEY_ID"
# Retrieve AWS Secret Access Key (stored as "Password" in the "aws-cli" entry)
AWS_SECRET_ACCESS_KEY=$(security find-internet-password -s aws-cli -w)
# (Do not echo the secret key)
# Check for missing keychain entry
if [[ "$AWS_ACCESS_KEY_ID" == *"could not be found"* || "$AWS_SECRET_ACCESS_KEY" == *"could not be found"* ]]; then
  echo "Error: The 'aws-cli' entry was not found in your Keychain."
  echo "Please add an Internet Password entry named 'aws-cli' with your AWS Access Key ID as the User Name and your AWS Secret Access Key as the Password."
  exit 2
fi

# Check for empty credentials
if [[ -z "$AWS_ACCESS_KEY_ID" || -z "$AWS_SECRET_ACCESS_KEY" ]]; then
  echo "Failed to retrieve AWS credentials from Passwords app (Keychain)."
  exit 1
fi

# Set credentials in AWS CLI config
aws configure set aws_access_key_id "$AWS_ACCESS_KEY_ID"
aws configure set aws_secret_access_key "$AWS_SECRET_ACCESS_KEY"

echo "AWS CLI credentials set from Passwords app."

echo "\n==== AWS Billing for Current Month ===="
# Get the first day of the current month and today
START_DATE=$(date +%Y-%m-01)
END_DATE=$(date -v+1d +%Y-%m-%d)

# Query AWS Cost Explorer for current month's billing (requires permissions)
BILLING_OUTPUT=$(aws ce get-cost-and-usage \
  --time-period Start=$START_DATE,End=$END_DATE \
  --granularity MONTHLY \
  --metrics "UnblendedCost" \
  --region us-east-1 \
  --output json 2>/dev/null)

if [[ -z "$BILLING_OUTPUT" ]]; then
  echo "(Billing info unavailable: check permissions or Cost Explorer access)"
else
  if command -v jq >/dev/null 2>&1; then
    AMOUNT=$(echo "$BILLING_OUTPUT" | jq -r '.ResultsByTime[0].Total.UnblendedCost.Amount')
    CURRENCY=$(echo "$BILLING_OUTPUT" | jq -r '.ResultsByTime[0].Total.UnblendedCost.Unit')
    if [[ -n "$AMOUNT" && -n "$CURRENCY" && "$AMOUNT" != "null" && "$CURRENCY" != "null" ]]; then
      echo "Current month AWS cost: $AMOUNT $CURRENCY"
    else
      echo "(Could not parse billing info)"
    fi
  else
    echo "(jq not found, cannot parse billing info cleanly)"
  fi
fi

echo "\n==== S3 Bucket Usage (bytes) ===="
# List all buckets and show total size for each
for bucket in $(aws s3api list-buckets --query 'Buckets[].Name' --output text); do
  size=$(aws s3 ls s3://$bucket --recursive --summarize 2>/dev/null | grep 'Total Size' | awk '{print $3}')
  echo "$bucket: ${size:-0} bytes"
done
