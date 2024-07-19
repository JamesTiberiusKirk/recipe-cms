-- name: schema_up
CREATE TABLE IF NOT EXISTS unit (
    unit_name       TEXT PRIMARY KEY,
    display_name    TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS author (
    username    TEXT PRIMARY KEY,
    password    TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS recipe (
    id              TEXT PRIMARY KEY,
    recipe_name     TEXT NOT NULL,
    intro           TEXT,
    description     TEXT,
    instructions    TEXT,
    length_total    TEXT,
    length_hands_on TEXT,
    closing         TEXT,
    recipe_version  INTEGER,
    author_name     TEXT,
    images          TEXT[],

    FOREIGN KEY (author_name) REFERENCES author (username) ON UPDATE CASCADE
);
CREATE TABLE IF NOT EXISTS ingredient (
    recipe_id   TEXT         NOT NULL,
    array_index INTEGER      NOT NULL,
    field       TEXT         NOT NULL,

    ingredient_name     TEXT         NOT NULL,
    amount              INTEGER,
    unit_name           TEXT,
    
    PRIMARY KEY (recipe_id, array_index, field),
    FOREIGN KEY (recipe_id) REFERENCES recipe (id) ON DELETE CASCADE,
    FOREIGN KEY (unit_name) REFERENCES unit (unit_name)
);
CREATE TABLE IF NOT EXISTS tag (
    recipe_id   TEXT,
    tag_name    TEXT,

    PRIMARY KEY (recipe_id, tag_name),
    FOREIGN KEY (recipe_id) REFERENCES recipe (id) ON DELETE CASCADE
);
CREATE INDEX tag_index ON tag (tag_name);
INSERT INTO unit (unit_name, display_name) VALUES
    ('can', 'Can'),
    ('clove', 'Clove'),
    ('small', 'Small'),
    ('medium', 'Medium'),
    ('large', 'Large'),
    ('kg', 'kg'),
    ('g', 'g'),
    ('l', 'l'),
    ('ml', 'ml'),
    ('unit', 'Unit'),
    ('part', 'Part');

INSERT INTO author (username, password) VALUES
    ('TestUser', 'testPass');
INSERT INTO recipe (id, recipe_name, intro, description, instructions, length_total, length_hands_on, closing, recipe_version, author_name, images) VALUES
    (
        'testid-1234',
        'Chilly con carne',
        E'## title of the introduction\n This is a cool recipe',
        E'### Description\n - point 1 \n - point 2 ',
        E'### Instructions\n 1. step 1 \n 2. step 2 ',
        '8 hours',
        '20 minutes',
        E'This a cool recipe',
        1,
        'TestUser',
        ARRAY['https://thecozycook.com/wp-content/uploads/2022/11/Chili-Con-Carne-f2.jpg','https://www.ocado.com/cmscontent/recipe_image_large/36025764.jpg?brD4', 'https://www.foodleclub.com/wp-content/uploads/2020/09/chili-con-carne-1.jpg']
    ),
    (
        'testid-12345',
        'Chilly con carne v3',
        E'## title of the introduction\n This is a cool recipe',
        E'### Description\n - point 1 \n - point 2 ',
        E'### Instructions\n 1. step 1 \n 2. step 2 ',
        '8 hours',
        '20 minutes',
        E'This a cool recipe',
        2,
        'TestUser',
        ARRAY['https://thecozycook.com/wp-content/uploads/2022/11/Chili-Con-Carne-f2.jpg','https://www.ocado.com/cmscontent/recipe_image_large/36025764.jpg?brD4', 'https://www.foodleclub.com/wp-content/uploads/2020/09/chili-con-carne-1.jpg']
    ),
    (
        'testid-123456',
        'Chilly con carne v3',
        E'## title of the introduction\n This is a cool recipe',
        E'### Description\n - point 1 \n - point 2 ',
        E'### Instructions\n 1. step 1 \n 2. step 2 ',
        '8 hours',
        '20 minutes',
        E'This a cool recipe',
        3,
        'TestUser',
        ARRAY['https://thecozycook.com/wp-content/uploads/2022/11/Chili-Con-Carne-f2.jpg','https://www.ocado.com/cmscontent/recipe_image_large/36025764.jpg?brD4', 'https://www.foodleclub.com/wp-content/uploads/2020/09/chili-con-carne-1.jpg']
    ),
    (
        'testid-123457',
        'Chilly con carne v4',
        E'## title of the introduction\n This is a cool recipe',
        E'### Description\n - point 1 \n - point 2 ',
        E'### Instructions\n 1. step 1 \n 2. step 2 ',
        '8 hours',
        '20 minutes',
        E'This a cool recipe',
        4,
        'TestUser',
        ARRAY['https://thecozycook.com/wp-content/uploads/2022/11/Chili-Con-Carne-f2.jpg','https://www.ocado.com/cmscontent/recipe_image_large/36025764.jpg?brD4', 'https://www.foodleclub.com/wp-content/uploads/2020/09/chili-con-carne-1.jpg']
    ),
    (
        'testid-123458',
        'Chilly con carne v5',
        E'## title of the introduction\n This is a cool recipe',
        E'### Description\n - point 1 \n - point 2 ',
        E'### Instructions\n 1. step 1 \n 2. step 2 ',
        '8 hours',
        '20 minutes',
        E'This a cool recipe',
        5,
        'TestUser',
        ARRAY['https://thecozycook.com/wp-content/uploads/2022/11/Chili-Con-Carne-f2.jpg','https://www.ocado.com/cmscontent/recipe_image_large/36025764.jpg?brD4', 'https://www.foodleclub.com/wp-content/uploads/2020/09/chili-con-carne-1.jpg']
    );
INSERT INTO ingredient (recipe_id, array_index, field, ingredient_name, amount, unit_name) VALUES 
    ( 'testid-1234', 0, 'INGREDIENT', 'Beef mince 5%', 500, 'g' ),
    ( 'testid-1234', 1, 'INGREDIENT', 'Tomato can', 1, 'unit' ),
    ( 'testid-1234', 2, 'INGREDIENT', 'Red bell pepper', 100, 'g' ),
    ( 'testid-1234', 3, 'INGREDIENT', 'White onion', 100, 'g' ),
    ( 'testid-1234', 4, 'INGREDIENT', 'Kidney beans can', 1, 'unit' ),
    ( 'testid-1234', 0, 'SEASONING', 'Paprika', 5, 'part' ),
    ( 'testid-1234', 1, 'SEASONING', 'Cumin', 1, 'part' ),
    ( 'testid-1234', 2, 'SEASONING', 'Salt', 1, 'part' ),
    ( 'testid-1234', 3, 'SEASONING', 'Pepper', 1, 'part' ),
    ( 'testid-12345', 0, 'INGREDIENT', 'Beef mince 5%', 500, 'g' ),
    ( 'testid-12345', 1, 'INGREDIENT', 'Tomato can', 1, 'unit' ),
    ( 'testid-12345', 2, 'INGREDIENT', 'Red bell pepper', 100, 'g' ),
    ( 'testid-12345', 3, 'INGREDIENT', 'White onion', 100, 'g' ),
    ( 'testid-12345', 4, 'INGREDIENT', 'Kidney beans can', 1, 'unit' ),
    ( 'testid-12345', 0, 'SEASONING', 'Paprika', 5, 'part' ),
    ( 'testid-12345', 1, 'SEASONING', 'Cumin', 1, 'part' ),
    ( 'testid-12345', 2, 'SEASONING', 'Salt', 1, 'part' ),
    ( 'testid-12345', 3, 'SEASONING', 'Pepper', 1, 'part' ),
    ( 'testid-123456', 0, 'INGREDIENT', 'Beef mince 5%', 500, 'g' ),
    ( 'testid-123456', 1, 'INGREDIENT', 'Tomato can', 1, 'unit' ),
    ( 'testid-123456', 2, 'INGREDIENT', 'Red bell pepper', 100, 'g' ),
    ( 'testid-123456', 3, 'INGREDIENT', 'White onion', 100, 'g' ),
    ( 'testid-123456', 4, 'INGREDIENT', 'Kidney beans can', 1, 'unit' ),
    ( 'testid-123456', 0, 'SEASONING', 'Paprika', 5, 'part' ),
    ( 'testid-123456', 1, 'SEASONING', 'Cumin', 1, 'part' ),
    ( 'testid-123456', 2, 'SEASONING', 'Salt', 1, 'part' ),
    ( 'testid-123456', 3, 'SEASONING', 'Pepper', 1, 'part' ),
    ( 'testid-123457', 0, 'INGREDIENT', 'Beef mince 5%', 500, 'g' ),
    ( 'testid-123457', 1, 'INGREDIENT', 'Tomato can', 1, 'unit' ),
    ( 'testid-123457', 2, 'INGREDIENT', 'Red bell pepper', 100, 'g' ),
    ( 'testid-123457', 3, 'INGREDIENT', 'White onion', 100, 'g' ),
    ( 'testid-123457', 4, 'INGREDIENT', 'Kidney beans can', 1, 'unit' ),
    ( 'testid-123457', 0, 'SEASONING', 'Paprika', 5, 'part' ),
    ( 'testid-123457', 1, 'SEASONING', 'Cumin', 1, 'part' ),
    ( 'testid-123457', 2, 'SEASONING', 'Salt', 1, 'part' ),
    ( 'testid-123457', 3, 'SEASONING', 'Pepper', 1, 'part' ),
    ( 'testid-123458', 0, 'INGREDIENT', 'Beef mince 5%', 500, 'g' ),
    ( 'testid-123458', 1, 'INGREDIENT', 'Tomato can', 1, 'unit' ),
    ( 'testid-123458', 2, 'INGREDIENT', 'Red bell pepper', 100, 'g' ),
    ( 'testid-123458', 3, 'INGREDIENT', 'White onion', 100, 'g' ),
    ( 'testid-123458', 4, 'INGREDIENT', 'Kidney beans can', 1, 'unit' ),
    ( 'testid-123458', 0, 'SEASONING', 'Paprika', 5, 'part' ),
    ( 'testid-123458', 1, 'SEASONING', 'Cumin', 1, 'part' ),
    ( 'testid-123458', 2, 'SEASONING', 'Salt', 1, 'part' ),
    ( 'testid-123458', 3, 'SEASONING', 'Pepper', 1, 'part' );
INSERT INTO tag (recipe_id, tag_name) VALUES 
    ( 'testid-1234', 'slow cooker'),
    ( 'testid-1234', 'low kcal'),
    ( 'testid-1234', 'beef'),
    ( 'testid-1234', 'fiber'),
    ( 'testid-12345', 'slow cooker'),
    ( 'testid-12345', 'low kcal'),
    ( 'testid-12345', 'beef'),
    ( 'testid-12345', 'fiber'),
    ( 'testid-123456', 'slow cooker'),
    ( 'testid-123456', 'low kcal'),
    ( 'testid-123456', 'beef'),
    ( 'testid-123456', 'fiber'),
    ( 'testid-123457', 'slow cooker'),
    ( 'testid-123457', 'low kcal'),
    ( 'testid-123457', 'beef'),
    ( 'testid-123457', 'fiber'),
    ( 'testid-123458', 'slow cooker'),
    ( 'testid-123458', 'low kcal'),
    ( 'testid-123458', 'beef'),
    ( 'testid-123458', 'fiber');


-- name: schema_down
DROP TABLE IF EXISTS ingredient;
DROP TABLE IF EXISTS unit;
DROP TABLE IF EXISTS tag;
DROP TABLE IF EXISTS recipe;
DROP TABLE IF EXISTS author;
