UPDATE unit SET display_name = "Unit" WHERE unit_name = "unit";
UPDATE unit SET display_name = "Part" WHERE unit_name = "part";
INSERT INTO unit (unit_name, display_name) VALUES
    ('can', 'Can'),
    ('clove', 'Clove'),
    ('small', 'Small'),
    ('medium', 'Medium'),
    ('large', 'Large'),
