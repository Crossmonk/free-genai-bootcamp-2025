INSERT INTO groups (name, description)
VALUES 
    ('Beginner Words', 'Basic vocabulary for beginners'),
    ('Common Phrases', 'Everyday useful phrases')
ON CONFLICT DO NOTHING; 