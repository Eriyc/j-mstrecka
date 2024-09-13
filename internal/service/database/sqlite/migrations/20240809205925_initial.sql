-- Add migration script here
CREATE TABLE IF NOT EXISTS upcs (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    upc TEXT NOT NULL UNIQUE,
    referable_id INTEGER NOT NULL,
    referable_type TEXT NOT NULL CHECK(referable_type IN ('user', 'product'))
);

CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,

    upc_id INTEGER,
    FOREIGN KEY (upc_id) REFERENCES upcs(id)
);

CREATE TABLE IF NOT EXISTS user_payments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    payment_amount REAL NOT NULL,
    payment_date INTEGER NOT NULL DEFAULT (DATETIME('now')),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS products (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name TEXT NOT NULL,

    upc_id INTEGER,
    FOREIGN KEY (upc_id) REFERENCES upcs(id));

CREATE TABLE IF NOT EXISTS product_price (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    product_id INTEGER NOT NULL,
    purchase_price REAL NOT NULL,
    internal_price REAL NOT NULL,
    external_price REAL NOT NULL,
    start_date INTEGER NOT NULL,
    end_date INTEGER,
    FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE TABLE IF NOT EXISTS product_stock (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    product_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    added_date INTEGER NOT NULL DEFAULT (DATETIME('now')),
    added_by INTEGER NOT NULL,
    FOREIGN KEY (added_by) REFERENCES users(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    user_id TEXT,
    product_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    transaction_date INTEGER NOT NULL DEFAULT (DATETIME('now')),
    price_type TEXT NOT NULL CHECK(price_type IN ('purchase', 'internal', 'external')),
    price_paid REAL NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE VIEW IF NOT EXISTS user_credits AS
-- Calculate total credits earned from adding stock
WITH credit_calculations AS (
    SELECT
        u.id AS user_id,
        COALESCE(ROUND(SUM(ps.quantity * pp.purchase_price), 1), 0) AS total_credits_earned
    FROM
        users u
    LEFT JOIN
        product_stock ps ON u.id = ps.added_by
    LEFT JOIN
        product_price pp ON ps.product_id = pp.product_id
            AND ps.added_date BETWEEN pp.start_date AND IFNULL(pp.end_date, DATETIME('now'))
    GROUP BY
        u.id
),
-- Calculate total debt incurred from transactions
debt_calculations AS (
    SELECT
        u.id AS user_id,
        COALESCE(ROUND(SUM(t.quantity * pp.internal_price), 1), 0) AS total_debt_incurred
    FROM
        users u
    LEFT JOIN
        transactions t ON u.id = t.user_id
    LEFT JOIN
        product_price pp ON t.product_id = pp.product_id
            AND t.transaction_date BETWEEN pp.start_date AND IFNULL(pp.end_date, DATETIME('now'))
    GROUP BY
        u.id
),
-- Calculate total cash payments made by users
payment_calculations AS (
    SELECT
        u.id AS user_id,
        COALESCE(ROUND(SUM(up.payment_amount), 1), 0) AS total_payments_made
    FROM
        users u
    LEFT JOIN
        user_payments up ON u.id = up.user_id
    GROUP BY
        u.id
),
-- Combine credits earned, debt incurred, and payments made
total_calculations AS (
    SELECT
        cc.user_id,
        cc.total_credits_earned,
        dc.total_debt_incurred,
        pc.total_payments_made,
        COALESCE(cc.total_credits_earned + pc.total_payments_made - dc.total_debt_incurred, 0) AS net_balance
    FROM
        credit_calculations cc
    LEFT JOIN
        debt_calculations dc ON cc.user_id = dc.user_id
    LEFT JOIN
        payment_calculations pc ON cc.user_id = pc.user_id
)
SELECT
    tc.user_id,
    tc.total_credits_earned,
    tc.total_debt_incurred,
    tc.total_payments_made,
    CASE
        WHEN tc.net_balance >= 0 THEN tc.net_balance -- Positive balance or 0 indicates no debt
        ELSE 0 -- User has no credits left
    END AS remaining_credits,
    CASE
        WHEN tc.net_balance < 0 THEN ABS(tc.net_balance) -- Convert negative balance to positive debt value
        ELSE 0 -- No debt
    END AS debt_incurred
FROM
    total_calculations tc;

CREATE VIEW IF NOT EXISTS current_stock AS
WITH stock_summary AS (
    SELECT
        ps.product_id,
        SUM(ps.quantity) AS total_stock_added
    FROM
        product_stock ps
    GROUP BY
        ps.product_id
),
transaction_summary AS (
    SELECT
        t.product_id,
        SUM(t.quantity) AS total_stock_sold
    FROM
        transactions t
    GROUP BY
        t.product_id
)
SELECT
    p.id AS product_id,
    p.name,
    COALESCE(CAST(IFNULL(ss.total_stock_added, 0) - IFNULL(ts.total_stock_sold, 0) AS INTEGER), 0) AS total_stock
FROM
    products p
LEFT JOIN
    stock_summary ss ON p.id = ss.product_id
LEFT JOIN
    transaction_summary ts ON p.id = ts.product_id
GROUP BY
    p.id, p.name;

CREATE TRIGGER IF NOT EXISTS set_upc_product_update
AFTER INSERT ON upcs
WHEN NEW.referable_type = 'product'
BEGIN
    UPDATE products
    SET upc_id = NEW.id
    WHERE id = NEW.referable_id;
END;

CREATE TRIGGER IF NOT EXISTS set_upc_user_update
AFTER INSERT ON upcs
WHEN NEW.referable_type = 'user'
BEGIN
    UPDATE users
    SET upc_id = NEW.id
    WHERE id = NEW.referable_id;
END;