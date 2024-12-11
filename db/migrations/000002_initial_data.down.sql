-- 外部キー制約を無効化
SET FOREIGN_KEY_CHECKS = 0;

-- テーブルデータを削除 (すべての行を削除し、AUTO_INCREMENT をリセット)
TRUNCATE TABLE client_bank_account;
TRUNCATE TABLE client;
TRUNCATE TABLE user;
TRUNCATE TABLE organization;
TRUNCATE TABLE tax_rate;

-- 外部キー制約を有効化
SET FOREIGN_KEY_CHECKS = 1;