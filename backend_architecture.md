# Go 後端架構推薦

本文檔為 `taiwanstay-back` 專案提供建議的技術架構、專案結構和函式庫選擇，旨在建立一個現代化、可維護且高效能的 Go 後端服務。

## 1. 核心技術棧

| 類別 | 推薦函式庫 | 說明 |
| :--- | :--- | :--- |
| **Web 框架** | [Gin](https://github.com/gin-gonic/gin) | 一個高效能、API 精簡的 Web 框架。社群龐大，文件豐富，對新手友好。 |
| **資料庫驅動** | [Official MongoDB Driver](https://github.com/mongodb/mongo-go-driver) | 官方支援的 MongoDB 驅動，穩定可靠。 |
| **環境變數管理** | [Viper](https://github.com/spf13/viper) | 強大的設定檔管理工具，能輕鬆讀取 `.env` 檔案和 YAML。 |
| **請求驗證** | [Validator v10](https://github.com/go-playground/validator) | 基於 struct tag 的驗證函式庫，可以與 Gin 完美整合。 |
| **相依性管理** | Go Modules (內建) | Go 官方的相依性管理工具。 |

## 2. 專案結構 (Project Layout)

我們將採用一個清晰的、分層的專案結構，以實現關注點分離。

```
/taiwanstay-back
├── cmd/
│   └── server/
│       └── main.go            # 程式進入點，初始化與啟動伺服器
├── internal/
│   ├── api/                   # HTTP Handlers (控制器)
│   │   ├── handler.go         # Handler 的基礎結構與路由設定
│   │   ├── middleware.go      # 中介軟體 (如：認證、日誌)
│   │   └── bookmark_handler.go # 處理書籤相關的 API
│   ├── service/               # 業務邏輯層
│   │   └── bookmark_service.go
│   ├── repository/            # 資料存取層
│   │   └── bookmark_repo.go
│   └── domain/                # 核心領域模型 (資料結構)
│       └── bookmark.go
├── pkg/
│   ├── config/                # 設定檔讀取 (Viper)
│   └── database/              # 資料庫連線
├── .env.example               # 環境變數範例
├── .gitignore
├── go.mod                     # Go Modules 檔案
└── go.sum
```

### 目錄職責說明

*   **`cmd/server/main.go`**: 應用程式的唯一入口。負責：
    1.  初始化設定 (Viper)。
    2.  建立資料庫連線。
    3.  注入依賴 (將 repository 注入 service，再將 service 注入 handler)。
    4.  設定 Gin 路由。
    5.  啟動 HTTP 伺服器。

*   **`internal/`**: 存放所有核心應用程式程式碼。`internal` 是 Go 的一個特殊目錄，意味著這裡的程式碼不能被其他專案匯入，保證了其內部性。
    *   **`domain/`**: 定義應用程式的核心資料結構 (structs)，例如 `User`, `Opportunity`, `Image`。這些是純粹的資料結構，不包含業務邏輯。
    *   **`repository/`**: 負責與資料庫進行所有互動。每個 `domain` 都會有一個對應的 repository 檔案。例如，`bookmark_repo.go` 將包含 `CreateBookmark`, `GetBookmarkByID` 等函式。
    *   **`service/`**: 包含所有的業務邏輯。它會呼叫 `repository` 來存取資料，並執行計算、驗證等操作。Service 層不應該知道任何有關 HTTP 的事情。
    *   **`api/`**: 處理所有 HTTP 相關的邏輯。Handler 函式負責：
        1.  解析和驗證傳入的請求 (JSON body, URL 參數)。
        2.  呼叫 `service` 層來執行業務邏輯。
        3.  將 `service` 返回的結果格式化為 JSON 並回傳給客戶端。

*   **`pkg/`**: 存放可以被外部應用程式安全使用的共享程式碼（"public" code）。例如，一個通用的資料庫連線模組或設定檔模組。

## 3. 架構模式：分層架構 (Layered Architecture)

我們將遵循經典的分層架構，確保程式碼的低耦合與高內聚。

**一個請求的生命週期如下：**

1.  **`main.go`** 註冊路由，將請求導向到 **`api/` (Handler)**。
2.  **Handler** 解析 HTTP 請求，並呼叫 **`service/` (Service)** 中的對應函式。
3.  **Service** 執行業務邏輯，並透過 **`repository/` (Repository)** 介面來操作資料庫。
4.  **Repository** 實作與 MongoDB 的具體互動，並將從 **`domain/` (Domain)** 取得或轉換的資料模型返回給 Service。
5.  資料沿著相同的路徑返回，最終由 **Handler** 將其格式化為 JSON 回應。

**依賴關係流向：** `api` -> `service` -> `repository` -> `domain`

這種單向依賴確保了核心業務邏輯 (service, domain) 的獨立性，使其易於測試和維護。

## 4. 下一步

1.  **環境設定**: 安裝 Go 語言環境。
2.  **初始化專案**: 執行 `go mod init github.com/your-username/taiwanstay-back`。
3.  **安裝依賴**: `go get` 上述推薦的函式庫。
4.  **建立專案結構**: 按照上面的結構建立目錄和檔案。
5.  **從 `main.go` 開始**: 撰寫第一個 "Hello World" 伺服器。
