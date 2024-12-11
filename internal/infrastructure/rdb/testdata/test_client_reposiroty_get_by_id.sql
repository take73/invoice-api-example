INSERT INTO client (organization_id, name, representative_name, phone_number, postal_code, address)
VALUES
    (1, '取引先A', '取引先担当者A', '03-1234-0001', '100-0010', '東京都港区芝公園1-1-1'),
    (1, '取引先B', '取引先担当者B', '03-1234-0002', '100-0020', '東京都新宿区新宿2-2-2'),
    (2, '取引先C', '取引先担当者C', '03-8765-0003', '150-0030', '東京都目黒区目黒3-3-3');


-- 企業テーブルのテストデータ
INSERT INTO organization (name, representative_name, phone_number, postal_code, address)
VALUES
    ('株式会社サンプル', '山田 太郎', '03-1234-5678', '100-0001', '東京都千代田区丸の内1-1-1'),
    ('有限会社テスト', '鈴木 花子', '03-8765-4321', '150-0002', '東京都渋谷区渋谷2-2-2');
