-- 企業テーブルのテストデータ
INSERT INTO organization (name, representative_name, phone_number, postal_code, address)
VALUES
    ('株式会社サンプル', '山田 太郎', '03-1234-5678', '100-0001', '東京都千代田区丸の内1-1-1'),
    ('有限会社テスト', '鈴木 花子', '03-8765-4321', '150-0002', '東京都渋谷区渋谷2-2-2');

-- ユーザーテーブルのテストデータ
INSERT INTO user (organization_id, name, email, password)
VALUES
    (1, '佐藤 一郎', 'ichiro.sato@example.com', 'password123'),
    (1, '田中 二郎', 'jiro.tanaka@example.com', 'password123'),
    (2, '高橋 三郎', 'saburo.takahashi@example.com', 'password123');

-- 取引先テーブルのテストデータ
INSERT INTO client (organization_id, name, representative_name, phone_number, postal_code, address)
VALUES
    (1, '取引先A', '取引先担当者A', '03-1234-0001', '100-0010', '東京都港区芝公園1-1-1'),
    (1, '取引先B', '取引先担当者B', '03-1234-0002', '100-0020', '東京都新宿区新宿2-2-2'),
    (2, '取引先C', '取引先担当者C', '03-8765-0003', '150-0030', '東京都目黒区目黒3-3-3');

-- 取引先銀行口座テーブルのテストデータ
INSERT INTO client_bank_account (client_id, bank_name, branch_name, account_number, account_name)
VALUES
    (1, 'みずほ銀行', '本店', '1234567', '取引先A口座名義'),
    (2, '三菱UFJ銀行', '新宿支店', '2345678', '取引先B口座名義'),
    (3, 'りそな銀行', '目黒支店', '3456789', '取引先C口座名義');

-- 消費税
INSERT INTO tax_rate (start_date, end_date, rate)
VALUES
    ('2014-04-01', '2019-09-30', 8.00),  -- 消費税8%期間
    ('2019-10-01', NULL, 10.00);         -- 消費税10%（現在も有効）