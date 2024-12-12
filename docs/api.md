# Invoice API ドキュメント

## 概要
この API は請求書管理システムの一部として機能します。主に請求書の作成と検索の機能を提供します。

---

## エンドポイント一覧

| メソッド | エンドポイント     | 説明                  |
|----------|--------------------|-----------------------|
| POST     | `/invoice`         | 請求書を新規作成する  |
| GET      | `/invoice`         | 請求書を検索する      |

---

## エンドポイント詳細

### 1. 請求書の作成

- **URL**: `/invoice`
- **HTTP メソッド**: POST
- **リクエストヘッダー**:
  - `Content-Type`: `application/json`

- **リクエストボディ**:

```json
{
  "userId": 1,
  "clientId": 1,
  "issueDate": "2023-12-01",
  "amount": 10000,
  "dueDate": "2023-12-15"
}
```

| フィールド | 型 | 必須 | 説明 |
|----------|----|-----|---------|
| userId	| uint	| 必須	| ユーザー ID |
| clientId	| uint	| 必須	| クライアント ID |
| issueDate	| string | 必須	| 請求書の発行日 (YYYY-MM-DD 形式) |
| amount	| int64	| 必須 | 	請求金額 |
| dueDate	| string | 必須| 	支払期日 (YYYY-MM-DD 形式) |

- **レスポンス**:
  - 成功時: 200 OK
  
  ```json
  {
  "id": 1,
  "organizationId": 1,
  "organizationName": "Test Organization",
  "clientId": 1,
  "clientName": "Test Client",
  "issueDate": "2023-12-01",
  "amount": 10000,
  "dueDate": "2023-12-15",
  "status": "pending"
  }
  ```

### 2. 請求書の検索

- **URL**: `/invoice`
- **HTTP メソッド**: POST
- **リクエストパラメータ**:

| フィールド | 型 | 必須 | 説明 |
|----------|----|-----|---------|
| startDate	|string|	必須|	検索開始日 (YYYY-MM-DD 形式)|
| endDate	|string|	必須|	検索終了日 (YYYY-MM-DD 形式)|

例: /invoice?startDate=2023-12-01&endDate=2023-12-31

- **レスポンス**:
  - 成功時: 200 OK
  
```json
{
  "invoices": [
    {
      "id": 1,
      "organizationId": 1,
      "organizationName": "Test Organization",
      "clientId": 1,
      "clientName": "Test Client",
      "issueDate": "2023-12-01",
      "amount": 10000,
      "dueDate": "2023-12-15",
      "status": "pending"
    }
  ]
}
```