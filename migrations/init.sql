CREATE TABLE IF NOT EXISTS vehicles (
                                        id SERIAL PRIMARY KEY,
                                        type TEXT NOT NULL,
                                        brand TEXT NOT NULL,
                                        model TEXT NOT NULL,
                                        year INTEGER NOT NULL,
                                        motor TEXT NOT NULL,
                                        status TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS clients (
                                       id SERIAL PRIMARY KEY,
                                       name TEXT NOT NULL,
                                       email TEXT NOT NULL,
                                       phone BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS sales (
                                     id SERIAL PRIMARY KEY,
                                     client_id INTEGER NOT NULL REFERENCES clients(id),
    vehicle_id INTEGER NOT NULL REFERENCES vehicles(id) UNIQUE,
    price REAL NOT NULL,
    sale_date TIMESTAMP NOT NULL
    );

