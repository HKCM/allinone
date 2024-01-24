获取MFA的token

```bash
function mfa_staging() {
  local STAGING_AWS_ID=0123456789012

  # MFAクレデンシャル情報取得
  local SESSION_JSON=$(aws --profile staging sts get-session-token --serial-number arn:aws:iam::$STAGING_AWS_ID:mfa/user --token-code $1 --output json)
  if [ $? -ne 0 ]; then
    echo "get-session-token failed, nothing to do"
    return
  fi

  # jqで必要な項目を変数にセットする
  local MFA_ACCESS_KEY=$(echo $SESSION_JSON | jq -r '.Credentials.AccessKeyId')
  local MFA_SECRET_ACCESS_KEY=$(echo $SESSION_JSON | jq -r '.Credentials.SecretAccessKey')
  local MFA_SESSION_TOKEN=$(echo $SESSION_JSON | jq -r '.Credentials.SessionToken')
  local MFA_EXPIRATION=$(echo $SESSION_JSON | jq -r '.Credentials.Expiration')

  # MFA用クレデンシャルをセットするプロファイル名
  local MFA_PROFILE_NAME=staging-mfa

  # MFA用プロファイルにクレデンシャルをセット
  aws --profile $MFA_PROFILE_NAME configure set aws_access_key_id $MFA_ACCESS_KEY
  aws --profile $MFA_PROFILE_NAME configure set aws_secret_access_key $MFA_SECRET_ACCESS_KEY
  aws --profile $MFA_PROFILE_NAME configure set aws_session_token $MFA_SESSION_TOKEN
  aws --profile $MFA_PROFILE_NAME configure set region ap-northeast-1
  aws --profile $MFA_PROFILE_NAME configure set output json

  echo "New credentials have been set successfully. (profile: $MFA_PROFILE_NAME, expiration: $MFA_EXPIRATION)"
}
```