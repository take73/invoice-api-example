INSERT INTO organization (name, representative_name, phone_number, postal_code, address)
VALUES
    ('株式会社TEST_A', '山田 太郎', '03-1234-5678', '100-0001', '東京都千代田区丸の内1-1-1'),
    ('株式会社TEST_B', '鈴木 花子', '03-8765-4321', '150-0002', '東京都渋谷区渋谷2-2-2');

INSERT INTO user (organization_id, name, email, password)
VALUES
    (1, 'testA', 'aaa@example.com', 'password123'),
    (2, 'testB', 'bbb@example.com', 'password123');
