# Discord Notify MCP Server アーキテクチャ

## 概要

このプロジェクトは、Discord Webhookを使用してメッセージを送信するMCP（Model Context Protocol）サーバーです。MCPプロトコルを使用することで、AIアシスタントやその他のMCPクライアントから簡単にDiscordへの通知を送信できます。

## 技術スタック

- **言語**: Go 1.23+
- **MCPライブラリ**: github.com/mark3labs/mcp-go v0.31.0
- **環境変数管理**: github.com/joho/godotenv v1.5.1
- **通信プロトコル**: MCP (Model Context Protocol) over stdio

## アーキテクチャ

### コンポーネント構成

```
┌─────────────────────┐
│   MCP Client        │
│ (AI Assistant等)    │
└─────────┬───────────┘
          │ stdio
          ▼
┌─────────────────────┐
│   MCP Server        │
│ (main.go)           │
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│  DiscordNotifyServer│
│   (server.go)       │
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│  Discord Client     │
│   (client.go)       │
└─────────┬───────────┘
          │ HTTP POST
          ▼
┌─────────────────────┐
│  Discord Webhook    │
│      API            │
└─────────────────────┘
```

### モジュール詳細

#### 1. Main (cmd/mcp-server/main.go)
- エントリーポイント
- 環境変数の読み込み（.envファイル）
- MCPサーバーの初期化と起動
- stdioを使用した通信の開始

#### 2. DiscordNotifyServer (internal/server/server.go)
- MCPサーバーの実装
- ツール（discord_send_message、discord_send_embed）の登録
- リクエストハンドラーの実装
- Discord Clientへのブリッジ

#### 3. Discord Client (internal/discord/client.go)
- Discord Webhook APIへのHTTPクライアント
- メッセージとEmbedの送信機能
- エラーハンドリング

### データフロー

1. **MCPクライアント → MCPサーバー**
   - JSON-RPC形式でツール呼び出しリクエストを送信
   - stdioを通じて通信

2. **MCPサーバー → Discord Client**
   - ツールハンドラーがリクエストを解析
   - 適切なDiscord Clientメソッドを呼び出し

3. **Discord Client → Discord API**
   - Webhook URLにHTTP POSTリクエストを送信
   - JSON形式でメッセージデータを送信

### セキュリティ考慮事項

1. **Webhook URLの保護**
   - 環境変数を使用してWebhook URLを管理
   - .envファイルは.gitignoreに追加してバージョン管理から除外

2. **入力検証**
   - 必須パラメータのチェック
   - 型の検証

3. **エラーハンドリング**
   - 適切なエラーメッセージの返却
   - Discord APIのレスポンスコードチェック

### 拡張性

このアーキテクチャは以下の拡張が容易です：

1. **新しいツールの追加**
   - server.goの`registerTools()`メソッドに新しいツールを登録
   - 対応するハンドラーメソッドを実装

2. **Discord機能の拡張**
   - client.goに新しいメソッドを追加
   - より複雑なDiscord APIの機能（ファイル添付、リアクション等）をサポート可能

3. **認証機能の追加**
   - MCPサーバーレベルでの認証
   - Webhook URLの動的管理

### パフォーマンス考慮事項

- 同期的な処理（MCPプロトコルの制約）
- HTTPリクエストのタイムアウト設定が必要（現在は未実装）
- 大量のメッセージ送信時はレート制限に注意