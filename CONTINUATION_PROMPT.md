# AI 協作接續開發指令 (Continuation Prompt)

## 任務目標

接續 `taiwanstay-back` Go 後端專案的開發，目標是根據已制定的計畫，逐步完成所有 API 端點的遷移與實作。

## 核心原則

你必須嚴格遵守已經建立的架構、規範和開發模式。在進行任何開發之前，請務必先閱讀並理解以下關鍵文件：

1.  `api_migration_plan.md`: 這是我們遷移工作的藍圖，定義了所有需要實作的 API 端點。
2.  `user_stories.md`: 描述了每個功能的使用者情境，是功能驗收的標準。
3.  `architecture_diagram.md`: 視覺化地展示了系統的分層架構。
4.  `ARCHITECTURE_PHILOSOPHY.md` (如果存在): 解釋了我們選擇顯式分層架構的原因。
5.  `internal/domain/`: 包含了所有已定義的核心資料模型，**新的功能必須基於或擴充這些模型**。

## 開發流程

對於每一個新的 API 端點 (例如「使用者登入」或「建立工作機會」)，你都必須遵循以下**分層開發順序**：

1.  **Domain 層 (如果需要)**:
    -   檢查 `internal/domain/` 中是否已存在對應的模型。
    -   如果需要新的資料結構，請在此處新增，並確保其與前端模型 (`taiwanstay-front/models/`) 保持一致。

2.  **Repository 層**:
    -   在 `internal/repository/` 中對應的 `_repo.go` 檔案裡，先在 `interface` 中定義新的資料庫操作方法。
    -   然後，在 `struct` 中實作該方法。**在資料庫連線完成前，先使用模擬邏輯**。

3.  **Service 層**:
    -   在 `internal/service/` 中對應的 `_service.go` 檔案裡，於 `interface` 中定義新的業務邏輯方法。
    -   在 `struct` 中實作該方法。此處應包含所有核心業務邏輯，例如資料驗證、計算、呼叫 Repository 等。

4.  **Handler 層**:
    -   在 `internal/api/` 中對應的 `_handler.go` 檔案裡，建立新的 Handler 方法。
    -   此方法負責：
        a.  解析與驗證 HTTP 請求 (Request Body, URL Params)。
        b.  呼叫 `Service` 層的對應方法。
        c.  將 `Service` 的結果轉換為 JSON 響應回傳。

5.  **Router 層**:
    -   在 `internal/api/router.go` 的 `SetupRoutes` 函式中，將新的 Handler 方法註冊到對應的路由上。

## 程式碼風格與規範

-   **顯式優於隱式**: 嚴格遵守手動依賴注入的模式，所有依賴關係都必須在 `cmd/server/main.go` 中清晰地串連起來。
-   **介面導向**: 所有 `Repository` 和 `Service` 都必須先定義介面 (Interface)，然後再進行實作。
-   **錯誤處理**: 嚴格處理所有函式可能回傳的 `error`。
-   **安全性**: 密碼等敏感資訊必須加密處理，絕不能出現在日誌或 API 響應中。

## 你的下一個任務

現在，請你接續開發 **「使用者登入」** 功能。請依照上述的開發流程，從 `Repository` 層開始，逐步完成 `Service` 層、`Handler` 層，並最後將其註冊到路由中。
