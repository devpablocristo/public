-- Create customers table if it doesn't exist
CREATE TABLE IF NOT EXISTS customers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(20),
    age INTEGER,
    birth_date DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create index for common queries
CREATE INDEX IF NOT EXISTS idx_customers_email ON customers(email);
CREATE INDEX IF NOT EXISTS idx_customers_name_lastname ON customers(name, last_name);

-- Insert some sample data (optional, comment out if not needed)
INSERT OR IGNORE INTO customers (name, last_name, email, phone, age, birth_date) VALUES 
    ('John', 'Doe', 'john.doe@example.com', '+1234567890', 30, '1993-01-15 00:00:00'),
    ('Jane', 'Smith', 'jane.smith@example.com', '+1987654321', 25, '1998-05-20 00:00:00');