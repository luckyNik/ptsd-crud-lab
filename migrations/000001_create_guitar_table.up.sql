CREATE TABLE IF NOT EXISTS guitar (
    id UUID PRIMARY KEY,
    manufacturer TEXT NOT NULL,
    string_count INT NOT NULL CHECK (string_count > 0),
    body_material TEXT NOT NULL,
    manufacture_date DATE NOT NULL
);
