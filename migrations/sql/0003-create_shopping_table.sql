CREATE TABLE shopping (
    shopping_id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    is_purchased INTEGER NOT NULL DEFAULT 0 CHECK(is_purchased IN (0, 1))
);
