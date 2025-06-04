# Discord Notify MCP Server

Discord通知を送信するためのMCP (Model Context Protocol) サーバーです。

## 機能

- シンプルなテキストメッセージの送信
- リッチな埋め込みメッセージ（Embed）の送信

## セットアップ

### Discord Webhook URLの取得

1. Discordサーバーの設定を開く
2. 「連携サービス」→「ウェブフック」を選択
3. 「新しいウェブフック」をクリック
4. 名前とチャンネルを設定し、「ウェブフックURLをコピー」

### プロジェクトのセットアップ

1. 依存関係のインストール:
```bash
go mod download
```

2. `.env`ファイルを作成し、Discord Webhook URLを設定:
```
DISCORD_WEBHOOK_URL=your_webhook_url_here
```

3. ビルド:
```bash
go build -o bin/discord-notify-mcp ./cmd/mcp-server
```

## 使用方法

### Claude Codeでの使用

1. MCPサーバーをClaude Codeに登録:
```bash
claude mcp add discord-notify -e DISCORD_WEBHOOK_URL="your_webhook_url_here" -- /path/to/bin/discord-notify-mcp
```

または、`.env`ファイルがある場合:
```bash
claude mcp add discord-notify -e DISCORD_WEBHOOK_URL="$(cat .env | grep DISCORD_WEBHOOK_URL | cut -d= -f2-)" -- $(pwd)/bin/discord-notify-mcp
```

2. 登録を確認:
```bash
claude mcp list
```

3. Claude Code内で`discord_send_message`や`discord_send_embed`ツールが使用可能になります。

### スタンドアロンでの使用

MCPサーバーを直接起動:
```bash
./bin/discord-notify-mcp
```

## 利用可能なツール

### discord_send_message
シンプルなテキストメッセージをDiscordに送信します。

パラメータ:
- `content` (string, required): 送信するメッセージ内容

### discord_send_embed
リッチな埋め込みメッセージをDiscordに送信します。

パラメータ:
- `title` (string): 埋め込みのタイトル
- `description` (string): 埋め込みの説明
- `color` (number): 埋め込みの色（10進数）
- `fields_json` (string): フィールドのJSON配列。各フィールドは`name`、`value`、オプションの`inline`プロパティを持ちます

## 使用例

### テキストメッセージの送信
```
discord_send_message(content="Hello from MCP Server!")
```

### 埋め込みメッセージの送信
```
discord_send_embed(
    title="通知タイトル",
    description="これは埋め込みメッセージです",
    color=16711680,  # 赤色
    fields_json='[{"name":"フィールド1","value":"値1","inline":true}]'
)
```

## プロジェクト構造

```
discord-notify-mcp/
├── cmd/
│   └── mcp-server/
│       └── main.go          # メインエントリーポイント
├── internal/
│   ├── discord/
│   │   └── client.go        # Discord Webhookクライアント
│   └── server/
│       └── server.go        # MCPサーバー実装
├── blueprint/
│   └── architecture.md      # アーキテクチャドキュメント
├── .env                     # 環境変数設定
├── .gitignore              # Git除外ファイル
├── go.mod                  # Goモジュール定義
├── go.sum                  # 依存関係のチェックサム
└── README.md               # このファイル
```