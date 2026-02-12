-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Tables: Customers, Products, Transactions, Transaction Items, Point Redemptions
CREATE TABLE IF NOT EXISTS customers (
      id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
      name VARCHAR(100) NOT NULL,
      point INT DEFAULT 0,
      created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS products (
      id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
      name VARCHAR(100) NOT NULL,
      type VARCHAR(50) NOT NULL,
      flavor VARCHAR(50) NOT NULL,
      size VARCHAR(20) NOT NULL,
      price INT NOT NULL,
      stock INT NOT NULL,
      created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS transactions (
      id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
      customer_id uuid NOT NULL,
      transaction_date TIMESTAMP DEFAULT NOW(),
      total_price INT NOT NULL,
      CONSTRAINT fk_transactions_customer
            FOREIGN KEY(customer_id)
            REFERENCES customers(id)
            ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS transaction_items (
      id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
      transaction_id uuid NOT NULL,
      product_id uuid NOT NULL,
      quantity INT NOT NULL,
      price INT NOT NULL,
      CONSTRAINT fk_items_transaction
            FOREIGN KEY(transaction_id)
            REFERENCES transactions(id)
            ON DELETE CASCADE,
      CONSTRAINT fk_itmes_product
            FOREIGN KEY(product_id)
            REFERENCES products(id)
            ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS point_redemptions (
      id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
      customer_id uuid NOT NULL,
      product_id uuid NOT NULL,
      point_used INT NOT NULL,
      created_at TIMESTAMP DEFAULT NOW(),
      CONSTRAINT fk_redemption_customer
            FOREIGN KEY(customer_id)
            REFERENCES customers(id)
            ON DELETE CASCADE,
      CONSTRAINT fk_redemption_product
            FOREIGN KEY(product_id)
            REFERENCES products(id)
            ON DELETE CASCADE
);
