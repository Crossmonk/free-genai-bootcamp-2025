INSERT INTO group_words (group_id, word_id)
SELECT g.id, w.id
FROM groups g, words w
WHERE g.name = 'Beginner Words' AND w.term IN ('Hello', 'Goodbye')
ON CONFLICT DO NOTHING; 