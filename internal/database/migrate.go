package database

import "log"

func AutoMigrate() error {
	queries := []string{

		// Enable UUID
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,

		// =========================
		// customers
		// =========================
		`
		CREATE TABLE IF NOT EXISTS customers (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name VARCHAR(100) NOT NULL,
			point INT DEFAULT 0,
			created_at TIMESTAMP DEFAULT NOW()
		);
		`,

		// =========================
		// products
		// =========================
		`
		CREATE TABLE IF NOT EXISTS products (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name VARCHAR(100) NOT NULL,
			type VARCHAR(50) NOT NULL,
			flavor VARCHAR(50) NOT NULL,
			size VARCHAR(20) NOT NULL,
			price INT NOT NULL,
			stock INT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		);
		`,

		// =========================
		// transactions
		// =========================
		`
		CREATE TABLE IF NOT EXISTS transactions (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			customer_id UUID NOT NULL,
			transaction_date TIMESTAMP DEFAULT NOW(),
			total_price INT NOT NULL,
			CONSTRAINT fk_transactions_customer
				FOREIGN KEY (customer_id)
				REFERENCES customers(id)
				ON DELETE CASCADE
		);
		`,

		// =========================
		// transaction_items
		// =========================
		`
		CREATE TABLE IF NOT EXISTS transaction_items (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			transaction_id UUID NOT NULL,
			product_id UUID NOT NULL,
			quantity INT NOT NULL,
			price INT NOT NULL,
			CONSTRAINT fk_items_transaction
				FOREIGN KEY (transaction_id)
				REFERENCES transactions(id)
				ON DELETE CASCADE,
			CONSTRAINT fk_items_product
				FOREIGN KEY (product_id)
				REFERENCES products(id)
				ON DELETE CASCADE
		);
		`,

		// =========================
		// point_redemptions
		// =========================
		`
		CREATE TABLE IF NOT EXISTS point_redemptions (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			customer_id UUID NOT NULL,
			product_id UUID NOT NULL,
			point_used INT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			CONSTRAINT fk_redemption_customer
				FOREIGN KEY (customer_id)
				REFERENCES customers(id)
				ON DELETE CASCADE,
			CONSTRAINT fk_redemption_product
				FOREIGN KEY (product_id)
				REFERENCES products(id)
				ON DELETE CASCADE
		);
		`,
	}

	for _, query := range queries {
		_, err := DB.Exec(query)
		if err != nil {
			log.Fatalf("❌ Failed to run migration: %v\nQuery:\n%s", err, query)
		}
	}

	log.Println("✅ Database migrated successfully (BSNACK)")
	return nil
}

func DropTables() error {
	queries := []string{
		`DROP TABLE IF EXISTS point_redemptions;`,
		`DROP TABLE IF EXISTS transaction_items;`,
		`DROP TABLE IF EXISTS transactions;`,
		`DROP TABLE IF EXISTS products;`,
		`DROP TABLE IF EXISTS customers;`,
	}

	for _, query := range queries {
		_, err := DB.Exec(query)
		if err != nil {
			log.Fatalf("❌ Failed to drop tables: %v\nQuery:\n%s", err, query)
		}
	}

	log.Println("✅ Tables dropped successfully (BSNACK)")
	return nil
}

