SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO invoice (
    organization_id, client_id, issue_date, payment_amount, fee, fee_rate, tax, tax_rate, total_amount, due_date, status
) VALUES
    (1, 1, '2024-01-01', 10000.00, 400.00, 0.04, 40.00, 0.1, 10440.00, '2024-01-10', 'pending'),
    (1, 2, '2024-01-05', 20000.00, 800.00, 0.04, 80.00, 0.1, 20880.00, '2024-01-15', 'processing'),
    (2, 3, '2024-01-10', 30000.00, 1200.00, 0.04, 120.00, 0.1, 31320.00, '2024-01-20', 'paid'),
    (2, 4, '2024-01-15', 40000.00, 1600.00, 0.04, 160.00, 0.1, 41760.00, '2024-02-01', 'error'),
    (3, 5, '2024-02-01', 50000.00, 2000.00, 0.04, 200.00, 0.1, 52200.00, '2024-02-10', 'pending');

SET FOREIGN_KEY_CHECKS = 1;
