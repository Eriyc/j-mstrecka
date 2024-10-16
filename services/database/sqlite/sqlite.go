package sqlite

import (
	"database/sql"
	"fmt"
	"gostrecka/models"
	"gostrecka/services/database"
	"gostrecka/services/env"
	"log"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"github.com/sarulabs/di/v2"
	_ "github.com/tursodatabase/go-libsql"
)

type SqliteMiddleware struct {
	Db        *sql.DB
	Logger    *slog.Logger
	Container di.Container
}

var _ database.Database = (*SqliteMiddleware)(nil)

func New(container di.Container) *SqliteMiddleware {
	return &SqliteMiddleware{
		Container: container,
		Logger:    container.Get("logger").(*slog.Logger).With("service", "SQLITE"),
	}
}

func (m *SqliteMiddleware) setup() (err error) {
	err = SetupMigrations(m.Db)
	if err != nil {
		return err
	}

	err = m.Migrate(m.Db)
	if err != nil {
		return err
	}

	return nil
}

func (m *SqliteMiddleware) Connect() (err error) {
	cfg := m.Container.Get("config").(env.Config)

	_, err = os.Stat(cfg.DbUrl)
	if os.IsNotExist(err) {
		f, err := os.OpenFile(cfg.DbUrl, os.O_RDONLY|os.O_CREATE, 0666)
		if err != nil {
			m.Logger.Error("could not create database file", "error", err.Error())
			return err
		}
		m.Logger.Warn("could not find database file, creating it", "path", cfg.DbUrl)
		f.Close()
	}

	if m.Db, err = sql.Open("libsql", fmt.Sprintf("file:%s", cfg.DbUrl)); err != nil {
		m.Logger.Error("could not open database", "error", err.Error())
		return err
	}

	err = m.setup()
	return
}

func (m *SqliteMiddleware) Close() {
	if m.Db != nil {
		m.Db.Close()
	}
}

func (m *SqliteMiddleware) Status() error {
	return m.Db.Ping()
}

func (m *SqliteMiddleware) GetUser(id string) (user models.User, balance models.Balance, err error) {

	row := m.Db.QueryRow("SELECT id, name FROM users WHERE id = ?", id)
	err = row.Scan(&user.ID, &user.Name)
	if err != nil {
		return
	}

	row = m.Db.QueryRow(`
		SELECT 
			total_credits_earned,
			total_payments_made,
			total_debt_incurred, 
			remaining_credits, 
			debt_incurred 
		FROM user_credits 
		WHERE user_id = ?
	`, id)

	err = row.Scan(
		&balance.TotalCreditsEarned,
		&balance.TotalPaymentsMade,
		&balance.TotalDebtIncurred,
		&balance.RemainingCredits,
		&balance.DebtIncurred,
	)

	return
}

func (m *SqliteMiddleware) CreateUser(id string, name string) error {
	_, err := m.Db.Exec("INSERT INTO users (id, name) VALUES (?, ?)", id, name)
	if err != nil {
		return err
	}

	var upc = rand.Intn(90000000) + 10000000

	_, err = m.Db.Exec("INSERT INTO upcs (referable_id, referable_type, upc) VALUES (?, 'user', ?)", id, upc)
	return err

}

func (m *SqliteMiddleware) GetUpcType(upc string) (lookup models.UpcLookup, err error) {
	row := m.Db.QueryRow("SELECT referable_id, referable_type FROM upcs WHERE upc = ?", upc)
	err = row.Scan(&lookup.ReferableId, &lookup.Type)
	return
}

func (m *SqliteMiddleware) GetProductIdent(id int64) (product models.Product, price models.ProductPrice, err error) {
	row := m.Db.QueryRow("SELECT product_id, name, total_stock FROM current_stock WHERE product_id = ?", id)
	err = row.Scan(&product.ID, &product.Name, &product.TotalStock)

	if err != nil {
		return
	}

	row = m.Db.QueryRow(`
		SELECT 
			id, 
			product_id, 
			purchase_price, 
			internal_price, 
			external_price,
			start_date,
			COALESCE(datetime(end_date), datetime(9999999999, 'unixepoch')) AS end_date
		FROM 
			product_price 
		WHERE 
			product_id = $1 
			AND datetime(start_date) <= datetime($2, 'unixepoch') 
			AND (
				end_date IS NULL 
				OR datetime(end_date) > datetime($2, 'unixepoch')
			)
		ORDER BY 
			start_date DESC 
		LIMIT 1;
	`, id, time.Now().Unix())

	err = row.Scan(
		&price.ID,
		&price.ProductID,
		&price.PurchasePrice,
		&price.InternalPrice,
		&price.ExternalPrice,
		&price.StartDate,
		&price.EndDate,
	)

	return
}

func (m *SqliteMiddleware) SearchProduct(name string) (products []models.ProductWithPrice, err error) {
	rows, err := m.Db.Query(`
		SELECT
			c.product_id,
			c.name,
			c.total_stock,
			p.purchase_price,
			p.internal_price,
			p.external_price,
			start_date,
			COALESCE(datetime(end_date), datetime(9999999999, 'unixepoch')) AS end_date
		FROM
			current_stock c
		LEFT JOIN 
			product_price p ON c.product_id = p.product_id
		WHERE
			LOWER(c.name) LIKE LOWER(?)
			AND datetime(start_date) <= datetime($2, 'unixepoch') 
			AND (
				end_date IS NULL 
				OR datetime(end_date) > datetime($2, 'unixepoch')
			)
		ORDER BY
			c.product_id DESC
	`, "%"+name+"%", time.Now().Unix(), time.Now().Unix())

	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var product models.Product
		var price models.ProductPrice

		err = rows.Scan(
			&product.ID,
			&product.Name,
			&product.TotalStock,
			&price.PurchasePrice,
			&price.InternalPrice,
			&price.ExternalPrice,
			&price.StartDate,
			&price.EndDate,
		)

		if err != nil {
			return
		}

		products = append(products, models.ProductWithPrice{Product: product, Price: price})

	}

	return
}

func (m *SqliteMiddleware) CreateProduct(name string, purchasePrice float64, internalPrice float64, externalPrice float64) error {
	tx, err := m.Db.Begin()
	if err != nil {
		return err
	}

	row := tx.QueryRow("INSERT INTO products (name) VALUES (?) RETURNING id", name)
	var id string
	err = row.Scan(&id)
	if err != nil {
		log.Printf("Error creating product: %s", err)
		tx.Rollback()
		return err
	}

	var upc = rand.Intn(90000000) + 10000000
	_, err = tx.Exec("INSERT INTO upcs (referable_id, referable_type, upc) VALUES (?, 'product', ?)", id, upc)

	if err != nil {
		log.Printf("Error creating upc: %s", err)
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("INSERT INTO product_price (product_id, purchase_price, internal_price, external_price, start_date) VALUES (?, ?, ?, ?, datetime('now'))",
		id, purchasePrice, internalPrice, externalPrice)

	if err != nil {
		log.Printf("Error creating product_price: %s", err)
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (m *SqliteMiddleware) Strecka(user models.User, productId int64, amount int64) error {

	product, price, err := m.GetProductIdent(productId)
	if err != nil {
		return err
	}

	m.Db.Exec("INSERT INTO transactions (user_id, product_id, quantity, price_type, price_paid) VALUES ($1, $2, $3, 'internal', $4)",
		user.ID, product.ID, amount, price.InternalPrice)

	return nil
}

func (m *SqliteMiddleware) UpdatePrice(productId int64, purchasePrice float64, internalPrice float64, externalPrice float64) error {
	tx, err := m.Db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE product_price SET end_date = datetime('now') WHERE product_id = ? AND end_date IS NULL", productId)
	if err != nil {
		log.Printf("Error updating product_price: %s", err)
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("INSERT INTO product_price (product_id, purchase_price, internal_price, external_price, start_date) VALUES (?, ?, ?, ?, datetime('now'))",
		productId, purchasePrice, internalPrice, externalPrice)

	if err != nil {
		log.Printf("Error creating product_price: %s", err)
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (m *SqliteMiddleware) AddStock(productId int64, userId string, amount int64) error {
	tx, err := m.Db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO product_stock (product_id, added_by, added_date, quantity) VALUES (?, ?, datetime('now'), ?)", productId, userId, amount)
	if err != nil {
		log.Printf("Error adding stock: %s", err)
		tx.Rollback()
	}

	return tx.Commit()
}

func (m *SqliteMiddleware) GetLatestTransactions() (transactions []models.LatestTransaction, err error) {

	rows, err := m.Db.Query(`
		SELECT
            t.user_id,
            u.name,
            DATETIME(t.transaction_date),
            SUM(t.quantity) OVER (
                PARTITION BY t.user_id
                ORDER BY t.transaction_date
                ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
            ) AS cumulative_transaction_count
        FROM
            transactions t
        LEFT JOIN users u ON
            t.user_id = u.id
        WHERE
            t.transaction_date >= datetime('now', '-12 hour')  -- Adjust the period here
        ORDER BY
            t.transaction_date ASC;
		`)

	if err != nil {
		fmt.Printf("errro: %v\n", err)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var transaction models.LatestTransaction
		err = rows.Scan(
			&transaction.UserID,
			&transaction.UserName,
			&transaction.TransactionDate,
			&transaction.CumulativeTransactionCount,
		)

		if err != nil {
			fmt.Printf("errro: %v\n", err)
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return
}

func (m *SqliteMiddleware) GetTransactionLeaderboard() (leaderboard []models.TransactionLeaderboard, err error) {
	rows, err := m.Db.Query(`
	WITH total_quantities AS (
		-- Calculate cumulative quantity sums over the last 12 hours
		SELECT
			t.user_id,
			SUM(t.quantity) AS total_quantity_sum
		FROM
			transactions t
		WHERE
			t.transaction_date >= datetime('now', '-12 hour')  -- Last 12 hours
		GROUP BY
			t.user_id
	),
	recent_transactions AS (
		-- Calculate cumulative quantity sums and ranks over the last 15 minutes
		SELECT
			t.user_id,
			SUM(t.quantity) AS recent_quantity_sum,
			RANK() OVER (ORDER BY SUM(t.quantity) DESC) AS recent_rank
		FROM
			transactions t
		WHERE
			t.transaction_date >= datetime('now', '-15 minute')  -- Last 15 minutes
		GROUP BY
			t.user_id
	),
	previous_transactions AS (
		-- Calculate cumulative quantity sums and ranks for the period before the last 15 minutes but within the last 12 hours
		SELECT
			t.user_id,
			SUM(t.quantity) AS previous_quantity_sum,
			RANK() OVER (ORDER BY SUM(t.quantity) DESC) AS previous_rank
		FROM
			transactions t
		WHERE
			t.transaction_date >= datetime('now', '-12 hour')  -- Last 12 hours
			AND t.transaction_date < datetime('now', '-15 minute')  -- Before the last 15 minutes
		GROUP BY
			t.user_id
	),
	rank_changes AS (
		-- Determine rank changes by comparing ranks in the last 15 minutes with the previous period
		SELECT
			r.user_id,
			r.recent_rank,
			p.previous_rank,
			CASE
				WHEN p.previous_rank IS NULL THEN 'New'  -- User is new in the leaderboard
				WHEN r.recent_rank < p.previous_rank THEN '↑'  -- Rank increased
				WHEN r.recent_rank > p.previous_rank THEN '↓'  -- Rank decreased
				ELSE '='  -- Rank remained the same
			END AS rank_change_indicator
		FROM
			recent_transactions r
		LEFT JOIN
			previous_transactions p ON r.user_id = p.user_id
	)
	SELECT
		tq.user_id,
		u.name,
		tq.total_quantity_sum,
    	RANK() OVER (ORDER BY tq.total_quantity_sum DESC) AS current_rank,  -- Current rank based on last 12 hours
		COALESCE(rc.rank_change_indicator, '=') AS rank_change_indicator
	FROM
		total_quantities tq
	LEFT JOIN
		rank_changes rc ON tq.user_id = rc.user_id
	LEFT JOIN
		users u ON tq.user_id = u.id
	ORDER BY
		tq.total_quantity_sum DESC;
	`)

	if err != nil {
		fmt.Printf("errro: %v\n", err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var transaction models.TransactionLeaderboard
		err = rows.Scan(
			&transaction.UserID,
			&transaction.UserName,
			&transaction.TotalTransactionCount,
			&transaction.CurrentRank,
			&transaction.RankChangeIndicator,
		)

		if err != nil {
			fmt.Printf("errro: %v\n", err)
			return
		}

		leaderboard = append(leaderboard, transaction)
	}

	return
}

func (m *SqliteMiddleware) GetUserUpcs() (upcs []models.Upc, err error) {
	rows, err := m.Db.Query(`
		SELECT
			u.id,
			u.upc,
			u.referable_type,
			u.referable_id,
			users.name AS referable_name
		FROM
			upcs u
		LEFT JOIN
			users ON u.referable_id = users.id
		WHERE
			u.referable_type = 'user'
		ORDER BY
			u.id ASC;
	`)

	if err != nil {
		fmt.Printf("errro: %v\n", err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var upc models.Upc
		err = rows.Scan(
			&upc.ID,
			&upc.Upc,
			&upc.Referable,
			&upc.ReferableId,
			&upc.ReferableName,
		)

		if err != nil {
			fmt.Printf("errro: %v\n", err)
			return
		}

		upcs = append(upcs, upc)
	}

	return
}

func (m *SqliteMiddleware) GetProductUpcs() (upcs []models.Upc, err error) {
	rows, err := m.Db.Query(`
		SELECT
			u.id,
			u.upc,
			u.referable_type,
			u.referable_id,
			products.name AS referable_name
		FROM
			upcs u
		LEFT JOIN
			products ON u.referable_id = products.id
		WHERE
			referable_type = 'product'
		ORDER BY
			u.id ASC;
	`)

	if err != nil {
		fmt.Printf("errro: %v\n", err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var upc models.Upc
		err = rows.Scan(
			&upc.ID,
			&upc.Upc,
			&upc.Referable,
			&upc.ReferableId,
			&upc.ReferableName,
		)

		if err != nil {
			fmt.Printf("errro: %v\n", err)
			return
		}

		upcs = append(upcs, upc)
	}

	return
}

func (m *SqliteMiddleware) GetTransactionNumbers(UserId string) (items []models.TransactionNumber, err error) {
	rows, err := m.Db.Query(`
		SELECT
			t.product_id,
			p.name,
			SUM(t.quantity) AS quantity,
			SUM(t.price_paid) AS price_paid
		FROM
			transactions t
		LEFT JOIN
			products p ON t.product_id = p.id
		WHERE
			t.user_id = $1
		GROUP BY
			t.product_id
		ORDER BY
			quantity DESC
		LIMIT 10;
	`, UserId)

	if err != nil {
		fmt.Printf("errro: %v\n", err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var item models.TransactionNumber
		err = rows.Scan(
			&item.ProductID,
			&item.ProductName,
			&item.Quantity,
			&item.PricePaid,
		)

		if err != nil {
			fmt.Printf("errro: %v\n", err)
			return
		}

		items = append(items, item)
	}

	return
}
