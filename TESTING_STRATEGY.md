# 後端測試策略

本文檔定義了 `taiwanstay-back` 專案的測試策略，涵蓋手動測試、自動化測試以及持續整合 (CI) 流程。

## 1. 測試層次

我們將採用多層次的測試策略，確保從單一函式到整個系統的穩定性。

-   **單元測試 (Unit Tests)**: 針對單一函式或組件（特別是 `Service` 層和 `Repository` 層的 mock 測試）進行的白箱測試。
-   **整合測試 (Integration Tests)**: 測試多個組件協同工作的正確性，例如測試從 `Handler` 到 `Service` 再到真實資料庫的完整流程。
-   **API 端點測試 (E2E / Manual)**: 驗證 API 端點是否符合預期規格的黑箱測試。

## 2. 手動 API 測試流程 (以 `curl` 為例)

在開發初期或需要快速驗證時，我們使用 `curl` 進行手動端點測試。

### 範例：測試使用者註冊

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
-H "Content-Type: application/json" \
-d '{
  "name": "testuser",
  "email": "test@example.com",
  "password": "password123"
}'
```

**驗證標準**:
1.  HTTP 狀態碼應為 `201 Created`。
2.  回應的 JSON Body 應包含新建立的使用者資料。
3.  回應的 JSON Body **絕不能**包含 `password` 欄位。

所有新的 API 端點在開發完成後，都應提供類似的 `curl` 範例，並記錄在對應的功能文件中。

## 3. 自動化測試

### 單元測試

-   **目標**: 測試 `Service` 層的業務邏輯和 `Repository` 層的資料庫操作邏輯。
-   **工具**: Go 內建的 `testing` 套件。
-   **實踐**:
    -   測試檔案命名為 `*_test.go`。
    -   大量使用 **Mocking**。例如，在測試 `UserService` 時，我們會傳入一個模擬的 `MockUserRepository`，而不是連線到真實資料庫。

### 整合測試

-   **目標**: 測試 API 端點的完整流程，包含資料庫互動。
-   **工具**: `testing` 套件, `net/http/httptest`。
-   **實踐**:
    -   我們會啟動一個測試專用的 Gin 伺服器實例。
    -   使用 `httptest` 來模擬 HTTP 請求，並驗證回應的狀態碼和內容。
    -   測試應在一個獨立的測試資料庫上執行，避免污染開發資料。

## 4. 持續整合 (Continuous Integration - CI)

我們將使用 **GitHub Actions** 來自動化測試與品質檢查流程。設定檔位於 `.github/workflows/ci.yml`。

CI 流程將包含以下步驟：

1.  **程式碼取出 (Checkout)**: 每次 `push` 或 `pull_request` 到 `main` 分支時觸發。
2.  **設定 Go 環境 (Setup Go)**: 設定指定的 Go 版本。
3.  **程式碼風格檢查 (Lint)**: 使用 `golangci-lint` 檢查程式碼風格是否一致。
4.  **執行單元測試 (Unit Tests)**: 執行所有 `*_test.go` 檔案。
5.  **建置應用程式 (Build)**: 執行 `go build` 確保專案可以成功編譯。

此流程確保了所有提交到主要分支的程式碼都符合品質標準。
