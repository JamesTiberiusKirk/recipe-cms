-- name: schema_up
CREATE TABLE IF NOT EXISTS recipe (
    id              VARCHAR(255) PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    intro           VARCHAR(255),
    description     VARCHAR(255),
    instructions    VARCHAR(255),
    length_total    VARCHAR(255),
    length_hands_on VARCHAR(255),
    closing         VARCHAR(255),
    version         INTEGER,
    author_name     VARCHAR(255),

    FOREIGN KEY (author_name) REFERENCES "user" (name)
);
CREATE TABLE IF NOT EXISTS ingredient (
    recipe_id   VARCHAR(255) NOT NULL,
    array_index INTEGER      NOT NULL,
    field       VARCHAR(255) NOT NULL,
    name        VARCHAR(255) NOT NULL,
    amount      INTEGER,
    unit_name   VARCHAR(255),
    
    PRIMARY KEY (recipe_id, array_index, field) ON DELETE CASCADE,
    FOREIGN KEY (recipe_id) REFERENCES recipe (id),
    FOREIGN KEY (unit_name) REFERENCES unit (name)
);
CREATE TABLE IF NOT EXISTS tag (
    recipe_id   VARCHAR(255) REFERENCES recipe(id) ON DELETE CASCADE,
    tag_name    VARCHAR(255),

    PRIMARY KEY (recipe_id, tag_name),
    FOREIGN KEY (recipe_id) REFERENCES recipe (id)
);
CREATE TABLE IF NOT EXISTS "user" (
    name        VARCHAR(255) PRIMARY KEY,
    password    VARCHAR(255) NOT NULL
);
CREATE TABLE IF NOT EXISTS unit (
    name            VARCHAR(255) PRIMARY KEY,
    display_name    VARCHAR(255) NOT NULL
);
INSERT INTO unit (name, display_name) VALUES
    ('kg', 'kg'),
    ('g', 'g'),
    ('l', 'l'),
    ('ml', 'ml'),
    ('unit', 'unit');

-- name: schema_down
DROP TABLE IF EXISTS ingredient;
DROP TABLE IF EXISTS unit;
DROP TABLE IF EXISTS tag;
DROP TABLE IF EXISTS auth_user;
DROP TABLE IF EXISTS recipe;
DROP TABLE IF EXISTS migrations;

