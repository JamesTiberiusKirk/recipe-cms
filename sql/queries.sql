-- name: get_all_recipes
SELECT
    r.id                AS id,
    r.recipe_name       AS recipename,
    r.intro             AS intro,
    r.description       AS description,
    r.instructions      AS instructions,
    r.length_total      AS lengthtotal,
    r.length_hands_on   AS lengthhandson,
    r.closing           AS closing,
    r.recipe_version    AS recipeversion,
    r.author_name       AS authorname,
    i.ingredient_name   AS ingredients_name,
    i.amount            AS ingredients_amount,
    iu.unit_name        AS ingredients_unit_name,
    iu.display_name     AS ingredients_unit_displayname,
    s.ingredient_name   AS seasonings_name,
    s.amount            AS seasonings_amount,
    su.unit_name        AS seasonings_unit_name,
    su.display_name     AS seasonings_unit_displayname,
    r.images            AS images,
    t.tag_name          AS tags
FROM recipe r 
LEFT JOIN ingredient AS i ON r.id = i.recipe_id AND i.field = 'INGREDIENT'
LEFT JOIN ingredient AS s ON r.id = s.recipe_id AND s.field = 'SEASONING'
LEFT JOIN unit AS iu ON i.unit_name = iu.unit_name
LEFT JOIN unit AS su ON s.unit_name = su.unit_name
LEFT JOIN tag        AS t ON r.id = t.recipe_id;

-- name: get_all_recipes_by_tag_name
SELECT
    r.id                AS id,
    r.recipe_name       AS recipename,
    r.intro             AS intro,
    r.description       AS description,
    r.instructions      AS instructions,
    r.length_total      AS lengthtotal,
    r.length_hands_on   AS lengthhandson,
    r.closing           AS closing,
    r.recipe_version    AS recipeversion,
    r.author_name       AS authorname,
    i.ingredient_name   AS ingredients_name,
    i.amount            AS ingredients_amount,
    iu.unit_name        AS ingredients_unit_name,
    iu.display_name     AS ingredients_unit_displayname,
    s.ingredient_name   AS seasonings_name,
    s.amount            AS seasonings_amount,
    su.unit_name        AS seasonings_unit_name,
    su.display_name     AS seasonings_unit_displayname,
    r.images            AS images,
    t.tag_name          AS tags
FROM recipe r 
LEFT JOIN ingredient AS i ON r.id = i.recipe_id AND i.field = 'INGREDIENT'
LEFT JOIN ingredient AS s ON r.id = s.recipe_id AND s.field = 'SEASONING'
LEFT JOIN unit AS iu ON i.unit_name = iu.unit_name
LEFT JOIN unit AS su ON s.unit_name = su.unit_name
LEFT JOIN tag        AS t ON r.id = t.recipe_id
WHERE t.tag_name = $1;

-- name: get_recipes_by_id
SELECT
    r.id                AS id,
    r.recipe_name       AS recipename,
    r.intro             AS intro,
    r.description       AS description,
    r.instructions      AS instructions,
    r.length_total      AS lengthtotal,
    r.length_hands_on   AS lengthhandson,
    r.closing           AS closing,
    r.recipe_version    AS recipeversion,
    r.author_name       AS authorname,
    i.ingredient_name   AS ingredients_name,
    i.amount            AS ingredients_amount,
    iu.unit_name        AS ingredients_unit_name,
    iu.display_name     AS ingredients_unit_displayname,
    s.ingredient_name   AS seasonings_name,
    s.amount            AS seasonings_amount,
    su.unit_name        AS seasonings_unit_name,
    su.display_name     AS seasonings_unit_displayname,
    r.images            AS images,
    t.tag_name          AS tags
FROM recipe r 
LEFT JOIN ingredient AS i ON r.id = i.recipe_id AND i.field = 'INGREDIENT'
LEFT JOIN ingredient AS s ON r.id = s.recipe_id AND s.field = 'SEASONING'
LEFT JOIN unit AS iu ON i.unit_name = iu.unit_name
LEFT JOIN unit AS su ON s.unit_name = su.unit_name
LEFT JOIN tag        AS t ON r.id = t.recipe_id
WHERE r.id = $1;

-- name: upsert_recipe
INSERT INTO recipe (id, recipe_name, intro, description, instructions, length_total, length_hands_on, closing, recipe_version, author_name, images)
VALUES (
    :id,
    :recipename,
    :intro,
    :description,
    :instructions,
    :lengthtotal,
    :lengthhandson,
    :closing,
    :recipeversion,
    :authorname,
    :images
)
ON CONFLICT (id) DO UPDATE SET
    recipe_name =       COALESCE(NULLIF(EXCLUDED.recipe_name, ''),     recipe.recipe_name),
    intro =             COALESCE(NULLIF(EXCLUDED.intro, ''),           recipe.intro),
    description =       COALESCE(NULLIF(EXCLUDED.description, ''),     recipe.description),
    instructions =      COALESCE(NULLIF(EXCLUDED.instructions, ''),    recipe.instructions),
    length_total =      COALESCE(NULLIF(EXCLUDED.length_total, ''),    recipe.length_total),
    length_hands_on =   COALESCE(NULLIF(EXCLUDED.length_hands_on, ''), recipe.length_hands_on),
    closing =           COALESCE(NULLIF(EXCLUDED.closing, ''),         recipe.closing),
    recipe_version =    COALESCE(NULLIF(EXCLUDED.recipe_version, 0),   recipe.recipe_version),
    author_name =       COALESCE(NULLIF(EXCLUDED.author_name, ''),     recipe.author_name),
    images =            COALESCE(EXCLUDED.images, recipe.images); -- doing the check for this in go

-- name: delete_all_tags_by_recipe_id
DELETE FROM tag WHERE recipe_id = :recipeid;

-- name: insert_tag
INSERT INTO tag (recipe_id, tag_name)
VALUES (:recipeid, :tagname)

-- name: delete_ingredient
DELETE FROM ingredient WHERE recipe_id = :recipeid AND field = :field;

-- name: upsert_ingredients
INSERT INTO ingredient (recipe_id, array_index, field, ingredient_name, amount, unit_name)
VALUES (:recipeid, :arrayindex, :field, :ingredientname, :amount, :unitname)
ON CONFLICT (recipe_id, array_index, field) DO UPDATE SET
    recipe_id =              COALESCE(NULLIF(EXCLUDED.recipe_id, ''),       ingredient.recipe_id),
    array_index =            COALESCE(NULLIF(EXCLUDED.array_index, 0),      ingredient.array_index),
    field =                  COALESCE(NULLIF(EXCLUDED.field, ''),           ingredient.field),
    ingredient_name =        COALESCE(NULLIF(EXCLUDED.ingredient_name, ''), ingredient.ingredient_name),
    amount =                 COALESCE(NULLIF(EXCLUDED.amount, 0),           ingredient.amount),
    unit_name =              COALESCE(NULLIF(EXCLUDED.unit_name, ''),       ingredient.unit_name);
