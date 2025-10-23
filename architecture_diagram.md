# TaiwanStay 後端架構圖 (Mermaid)

本文檔使用 Mermaid 語法繪製了 `taiwanstay-back` 專案的系統架構圖與請求流程圖，以視覺化方式呈現系統設計。

## 1. 系統分層架構圖 (Layered Architecture)

此圖展示了系統的核心分層結構，以及各層之間的單向依賴關係。

```mermaid
graph TD
    subgraph "外部請求 (External Requests)"
        Client[用戶端 / Client]
    end

    subgraph "Go 後端應用 (Go Backend Application)"
        direction LR
        subgraph "接入層 (API Layer)"
            Router["路由 (Gin Router)"]
            Handlers[API Handlers]
            Middleware[中介軟體]
        end

        subgraph "業務邏輯層 (Service Layer)"
            Services[業務邏輯服務]
        end

        subgraph "資料存取層 (Repository Layer)"
            Repositories[資料倉儲]
        end

        subgraph "領域模型 (Domain Layer)"
            Models["資料模型 (Structs)"]
        end
    end

    subgraph "外部服務 (External Services)"
        MongoDB[(MongoDB)]
        GoogleVision[Google Vision API]
        Cloudinary[Cloudinary API]
    end

    Client --> Router
    Router --> Middleware
    Middleware --> Handlers
    Handlers --> Services
    Services --> Repositories
    Repositories --> Models
    Repositories --> MongoDB
    Services --> GoogleVision
    Services --> Cloudinary
```

## 2. 使用者註冊請求流程圖 (Sequence Diagram)

此圖以「使用者註冊」(故事 1.1) 為例，展示了一個典型請求在系統各層之間的詳細處理流程。

```mermaid
sequenceDiagram
    participant Client as 用戶端
    participant Router as Gin 路由
    participant AuthHandler as 認證 Handler
    participant AuthService as 認證 Service
    participant UserRepo as 用戶 Repository
    participant MongoDB as 資料庫

    Client->>+Router: POST /api/v1/auth/register (含註冊資料)
    Router->>+AuthHandler: 呼叫 Register 函式
    AuthHandler->>+AuthService: 呼叫 RegisterUser 服務 (傳入 DTO)
    AuthService->>AuthService: 驗證資料、加密密碼
    AuthService->>+UserRepo: 呼叫 CreateUser (傳入 User 模型)
    UserRepo->>+MongoDB: 插入新的使用者文件
    MongoDB-->>-UserRepo: 返回成功訊息
    UserRepo-->>-AuthService: 返回建立的使用者物件
    AuthService->>AuthService: 生成 JWT (JSON Web Token)
    AuthService-->>-AuthHandler: 返回使用者物件與 Token
    AuthHandler-->>-Router: 回傳 201 Created JSON 響應
    Router-->>-Client: 回傳使用者資料與 Token
```

## 3. 總結

-   **分層架構圖** 宏觀地展示了系統的模組劃分和依賴方向，確保了程式碼的低耦合與高內聚。
-   **序列圖** 微觀地展示了單一功能（如註冊）的完整生命週期，有助於理解程式碼的執行細節，並可作為撰寫整合測試的參考。
