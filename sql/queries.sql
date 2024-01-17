-- name: upsert_recipe
INSERT INTO recipe (id, name, intro, description, instructions, length_total, length_hands_on, closing, version, author_name)
VALUES (:your_recipe_id, :your_recipe_name, :intro_text, :description_text, :instructions_text, :total_length_value, :hands_on_length_value, :closing_text, 1, :author_name)
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
