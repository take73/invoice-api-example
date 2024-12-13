# invoice-api-exsample

## APIドキュメント

API仕様は[こちら参照](./docs/api.md)

## 環境セットアップ

以下の手順に従って、ローカル環境でプロジェクトをセットアップしてください。


## 必須要件
- MySQL 8.0 以上がインストールされている（dockerでもOK）
- Golang 1.23 以上がインストールされている
- [golang-migrate](https://github.com/golang-migrate/migrate) 4.15.1 以上がインストールされている
- VSCodeのRestClientがインストールされている

## 1. MySQLにdatabaseを作成してください
```
CREATE DATABASE exampledb CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
```

## 2. .env.local ファイルをコピーして、内容を環境に合わせて修正してください
```
cp .env.local .env
```

## 3. migrateを実行してテーブルを作成してください（SQLを手動で流してもok）
```
make migrate-up
```

## 4. サーバーを起動します
```
% make run-local 
go run cmd/server/main.go
2024/12/13 02:30:56 Starting server on :1323

   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v4.13.0
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
⇨ http server started on [::]:1323
```

## 5. test.httpを実行して動作確認します



# 補足
## Authorization
今回の実装では、単純にクライアント単位でのscopeを設定しています。そのため、ユーザーごとに細かなアクセス制御が必要な場合は、ロールベースのアクセスコントロール (RBAC) を検討します

- [Auth0による認可](https://auth0.com/docs/quickstart/backend/golang/interactive)を行う
- [go-jwt-middleware](https://github.com/auth0/go-jwt-middleware)

