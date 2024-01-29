-- name: upsert_recipe
INSERT INTO recipe (id, name, intro, description, instructions, length_total, length_hands_on, closing, version, author_name)
VALUES (
    :id,
    :name,
    :intro,
    :description,
    :instructions,
    :total_length,
    :hands_on_length,
    :closing,
    :version,
    :author_name)
ON CONFLICT (id) DO UPDATE
    SET
COALESCE(NULLIF(EXCLUDED.command, ''), chat_journies.command),
        name =              COALESCE(NULLIF(EXCLUDED.name, ''), recipe.name,
        intro =             COALESCE(NULLIF(EXCLUDED.intro, ''), recipe.intro,
        description =       COALESCE(NULLIF(EXCLUDED.description, ''), recipe.description,
        instructions =      COALESCE(NULLIF(EXCLUDED.instructions, ''), recipe.instructions,
        length_total =      COALESCE(NULLIF(EXCLUDED.length_total, ''), recipe.length_total,
        length_hands_on =   COALESCE(NULLIF(EXCLUDED.length_hands_on, ''), recipe.length_hands_on,
        closing =           COALESCE(NULLIF(EXCLUDED.closing, ''), recipe.closing,
        version =           COALESCE(NULLIF(EXCLUDED.version, ''), recipe.version
        author_name =       COALESCE(NULLIF(EXCLUDED.author_name, ''), recipe.author_name;

-- name: insert_tags
INSERT INTO tag (recipe_id, tag_name)
SELECT :id, unnest(:tags::text[]);

-- name: get_all_recipes
SELECT
    r.id AS recipe_id,
    r.name AS recipe_name,
    r.intro AS recipe_intro,
    r.description AS recipe_description,
    r.instructions AS recipe_instructions,
    r.length_total AS recipe_length_total,
    r.length_hands_on AS recipe_length_hands_on,
    r.closing AS recipe_closing,
    r.version AS recipe_version,
    r.author_name AS recipe_author_name,
    u.password AS author_password,
    i.array_index AS ingredient_array_index,
    i.field AS ingredient_field,
    i.name AS ingredient_name,
    i.amount AS ingredient_amount,
    i.unit_name AS ingredient_unit_name,
    t.tag_name AS tag_name
FROM recipe r
JOIN "user" u ON r.author_name = u.name
LEFT JOIN ingredient i ON r.id = i.recipe_id
LEFT JOIN unit un ON i.unit_name = un.name
LEFT JOIN tag t ON r.id = t.recipe_id;
