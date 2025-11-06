package slack

import "os"

// 環境変数からwebhookを読み込む。設定されていない場合は空文字列（通知無効）
var webhook = os.Getenv("SLACK_WEBHOOK_URL")