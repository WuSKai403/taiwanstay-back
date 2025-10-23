# API 遷移計畫：從 Next.js 到 Go (已根據程式碼校正)

本文檔根據 `taiwanstay-front/pages/api` 目錄下的實際程式碼進行規劃，旨在將現有 API 端點準確地遷移至 `taiwanstay-back` (Go) 後端。

## 1. 待遷移的 API 端點 (Code-Verified)

以下是根據程式碼分析後，確認需要遷移的 API 資源列表。

### 核心資源 (Core Resources)

-   **/api/auth**: 處理用戶註冊、登入等身份驗證。
-   **/api/users**: 管理員對用戶的 CRUD 操作。
-   **/api/user**: 處理當前登入用戶的資料查詢與更新。
-   **/api/opportunities**: 工作機會的 CRUD 操作與搜尋。
-   **/api/organizations**: 組織的 CRUD 操作。
-   **/api/applications**: 申請的 CRUD 操作。
-   **/api/hosts**: 主人 (Host) 的 CRUD 操作。
-   **/api/bookmarks**: 書籤功能。
-   **/api/notifications**: 通知功能。
-   **/api/admin**: 管理員專用端點。

### 功能性端點 (Functional Endpoints)

-   **/api/check-image**: 圖片安全審核 (需整合 Google Vision API)。
-   **/api/upload** & **/api/cloudinary**: 檔案（主要是圖片）上傳功能。

### 開發用腳本 (Development Scripts)

-   **/api/seed**: 資料庫填充腳本。此功能建議在 Go 後端改為一個獨立的 CLI 命令 (`cmd/seeder/main.go`)，而非一個公開的 API 端點，以提高安全性。

## 2. Go 後端對應結構規劃

我們將為每個資源建立對應的 `handler`、`service` 和 `repository`。

| 資源 (Domain) | 路由前綴 (Route Prefix) | Handler | Service | Repository | Domain |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **認證** | `/api/v1/auth` | `auth_handler.go` | `auth_service.go` | `user_repo.go` | `user.go`, `token.go` |
| **用戶** | `/api/v1/users` | `user_handler.go` | `user_service.go` | `user_repo.go` | `user.go` |
| **工作機會** | `/api/v1/opportunities` | `opportunity_handler.go` | `opportunity_service.go` | `opportunity_repo.go` | `opportunity.go` |
| **組織** | `/api/v1/organizations` | `organization_handler.go` | `organization_service.go` | `organization_repo.go` | `organization.go` |
| **申請** | `/api/v1/applications` | `application_handler.go` | `application_service.go` | `application_repo.go` | `application.go` |
| **主人** | `/api/v1/hosts` | `host_handler.go` | `host_service.go` | `host_repo.go` | `host.go` |
| **書籤** | `/api/v1/bookmarks` | `bookmark_handler.go` | `bookmark_service.go` | `bookmark_repo.go` | `bookmark.go` |
| **通知** | `/api/v1/notifications` | `notification_handler.go` | `notification_service.go` | `notification_repo.go` | `notification.go` |
| **圖片上傳** | `/api/v1/upload` | `upload_handler.go` | `upload_service.go` | `(Cloudinary)` | `media.go` |
| **圖片審核** | `/api/v1/images/check` | `image_check_handler.go` | `image_check_service.go` | `(Google Vision)` | `media.go` |

## 3. 遷移步驟

1.  **建立專案結構**: 建立 `cmd`, `internal`, `pkg` 等目錄。
2.  **定義 Domain 模型**: 在 `internal/domain/` 目錄下，根據前端的 `models` 和 `types`，建立所有必要的 Go struct。
3.  **實現 Repository 層**: 為每個 Domain 建立 `repository`，負責與 MongoDB 互動。
4.  **實現 Service 層**: 建立 `service`，處理業務邏輯，並與外部 API (如 Google Vision, Cloudinary) 互動。
5.  **實現 Handler 層**: 建立 `api` handler，處理 HTTP 請求、驗證並呼叫 service。
6.  **設定路由**: 在 `main.go` 中將所有路由與 handler 綁定。
7.  **遷移認證機制**: 使用 Go 的函式庫（如 `jwt-go`）取代 NextAuth.js。
8.  **建立 Seeder 命令**: 建立一個 `cmd/seeder/main.go` 來處理資料庫填充。
9.  **撰寫測試**: 為各層級撰寫單元測試與整合測試。
