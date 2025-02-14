INSERT INTO study_activities (name, description) 
VALUES 
    ('Flashcards', 'Practice with flashcards'),
    ('Quiz', 'Test your knowledge')
ON CONFLICT DO NOTHING; 