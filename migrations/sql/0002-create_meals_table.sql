CREATE TABLE meals (
    meal_id INTEGER PRIMARY KEY,
    date TEXT NOT NULL,
    recipe_id INTEGER NOT NULL REFERENCES recipes(recipe_id) ON DELETE CASCADE
);

INSERT INTO meals (date, recipe_id)
VALUES
    ('2023-11-01', 1),
    ('2023-11-01', 5),

    ('2023-11-02', 2),
    ('2023-11-02', 3),

    ('2023-11-03', 4),
    ('2023-11-03', 7),

    ('2023-11-04', 8),
    ('2023-11-04', 9),

    ('2023-11-05', 6),
    ('2023-11-05', 10),

    ('2023-11-06', 1),
    ('2023-11-06', 3),

    ('2023-11-07', 4),
    ('2023-11-07', 8);
