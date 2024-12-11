-- 外部キー制約を無効化
SET FOREIGN_KEY_CHECKS = 0;

-- テーブルの削除
DROP TABLE IF EXISTS invoice;
DROP TABLE IF EXISTS client_bank_account;
DROP TABLE IF EXISTS client;
DROP TABLE IF EXISTS user;
DROP TABLE IF EXISTS organization;
DROP TABLE IF EXISTS tax_rate;

-- 外部キー制約を有効化
SET FOREIGN_KEY_CHECKS = 1;

