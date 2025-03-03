# Financial Partner Backend

理財小夥伴後端服務，提供使用者記帳、預算管理等功能的 API 服務。

## 前置準備

1. 放置必要檔案：
   ```bash
   # 在 config/ 目錄下放置以下檔案
   config/
   ├── config.yaml        # 本地配置文件（從 config.example.yaml 複製修改）
   └── firebase_credential.json # Firebase Admin SDK 憑證
   ```

2. 安裝 air（可選，用於本地開發熱重載）：
   ```bash
   # 使用 go install 安裝
   go install github.com/cosmtrek/air@latest

   # 確認 air 已安裝
   air -v
   ```

## 啟動服務

### 只啟動依賴服務（MongoDB、Redis）

```bash
# 啟動所有依賴服務
docker-compose up -d

# 只啟動特定服務
docker-compose up -d mongodb
docker-compose up -d redis
```

### 本地開發

使用 air
```bash
air
```

直接啟動
```bash
go run cmd/server/main.go
```

### 使用 Docker 啟動完整環境（包含 API 服務）

```bash
docker-compose -f docker-compose.dev.yml up -d --build
```

服務啟動後，API 可通過 http://localhost:8080