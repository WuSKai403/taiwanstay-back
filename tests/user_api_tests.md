# 使用者 API 整合測試規格

本文檔定義了 `taiwanstay-back` 使用者相關 API 的自動化整合測試案例，遵循 `TESTING_STRATEGY.md` 中定義的整合測試流程。

測試將使用 Go 的 `testing` 套件與 `net/http/httptest` 來模擬 API 請求並驗證回應。

## Epic 1: 使用者認證 (User Authentication)

### 故事 1.1: 使用者註冊

**使用者故事**:
> **身為** 一位新訪客，
> **我想要** 註冊一個新的帳號，
> **以便** 我可以申請工作機會並儲存我的最愛。

**端點**: `POST /api/v1/auth/register`

---

#### 測試案例 1: `TestRegister_Success`

**目的**: 驗證一位新使用者可以使用有效的資料成功註冊。

**測試步驟**:
1.  設定一個測試用的 Gin 引擎與 `UserHandler`。
2.  建立一個模擬的 HTTP `POST` 請求到 `/api/v1/auth/register`，Body 包含有效的使用者資料 (name, email, password)。
3.  使用 `httptest.ResponseRecorder` 記錄回應。
4.  斷言 (Assert) HTTP 狀態碼為 `201 Created`。
5.  斷言回應的 JSON Body 包含 `id`, `name`, `email` 欄位。
6.  斷言回應的 JSON Body **不包含** `password` 欄位。

---

#### 測試案例 2: `TestRegister_EmailAlreadyExists`

**目的**: 驗證系統不允許使用重複的 Email 註冊。

**測試步驟**:
1.  先在測試資料庫中建立一個使用者。
2.  建立一個模擬的 HTTP `POST` 請求，Body 使用與步驟 1 中相同的 `email`。
3.  使用 `httptest.ResponseRecorder` 記錄回應。
4.  斷言 HTTP 狀態碼為 `409 Conflict`。
5.  斷言回應的 JSON Body 包含指定的錯誤訊息。

---

#### 測試案例 3: `TestRegister_MissingFields`

**目的**: 驗證當請求中缺少必要欄位 (如 `password`) 時，系統會回傳錯誤。

**測試步驟**:
1.  建立一個模擬的 HTTP `POST` 請求，Body 中故意缺少 `password` 欄位。
2.  使用 `httptest.ResponseRecorder` 記錄回應。
3.  斷言 HTTP 狀態碼為 `400 Bad Request`。
4.  斷言回應的 JSON Body 包含關於欄位驗證的錯誤訊息。
