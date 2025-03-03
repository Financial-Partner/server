# Financial Partner Backend

理財小夥伴後端服務，提供使用者記帳、預算管理等功能的 API 服務。

[API Document](https://financial-partner.github.io/server)

## 前置準備

1. 放置必要檔案：
   ```bash
   # 在 config/ 目錄下放置以下檔案
   config/
   ├── config.yaml        # 本地配置文件（從 config.example.yaml 複製修改）
   └── firebase_credential.json # Firebase Admin SDK 憑證
   ```

2. 安裝開發工具：
   ```bash
   # 安裝 air（用於本地開發熱重載，可選）
   go install github.com/cosmtrek/air@latest

   # 執行設定腳本（安裝 linter 和設定 Git Hooks）
   ./scripts/setup.sh
   ```

## 開發規範

專案使用 Git Hooks 進行程式碼品質控管：

- **Pre-commit**: 執行 lint 檢查，包含：
  - 程式碼格式檢查 (gofmt)
  - 靜態程式碼分析
  - 安全性檢查

- **Pre-push**: 執行所有檢查，包含：
  - 所有 pre-commit 檢查
  - 單元測試

## 啟動服務

### 只啟動依賴服務（MongoDB、Redis）

```bash
# 啟動 MongoDB 和 Redis
docker-compose up -d

# 只啟動特定服務
docker-compose up -d mongodb
docker-compose up -d redis
```

### 本地開發（使用 air）

```bash
air
```

### 使用 Docker 啟動完整環境

```bash
# 啟動所有服務（API、MongoDB、Redis）
docker-compose -f docker-compose.dev.yml up --build

# 背景執行
docker-compose -f docker-compose.dev.yml up -d --build
```

服務啟動後，API 可通過 http://localhost:8080 訪問

## API 文件

使用 Swagger 撰寫 API 文件和提供互動式 UI。

### 更新文件

當 API 有變動時，執行：

```bash
./scripts/swagger.sh
```

### 查看文件

http://localhost:8080/swagger/index.html
