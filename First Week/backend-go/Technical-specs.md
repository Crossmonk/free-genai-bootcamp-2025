# Backend Server Technical Specs

## Business Goal: 
A language learning school wants to build a prototype of learning portal which will act as three things:
- Inventory of possible vocabulary that can be learned
- Act as a  Learning record store (LRS), providing correct and wrong score on practice vocabulary
- A unified launchpad to launch different learning apps

## Technical Restrictions:
- Use SQLite3 as the database
- The backend server will be written in Go
- Does not require authentication/authorization, assume there is a single user
- API will be code using Gin
- The API will always return JSON
- The API will always receive JSON
- The API will be documented using Swagger and OpenAPI

## Database Schema
- words — Stores individual Japanese vocabulary words.
    - `id` (Primary Key): Unique identifier for each word
    - `kanji` (String, Required): The word written in Japanese kanji
    - `romaji` (String, Required): Romanized version of the word
    - `english` (String, Required): English translation of the word
    - `parts` (JSON, Required): Word components stored in JSON format

- groups — Manages collections of words.
    - `id` (Primary Key): Unique identifier for each group
    - `name` (String, Required): Name of the group
    - `words_count` (Integer, Default: 0): Counter cache for the number of words in the group

- word_groups — join-table enabling many-to-many relationship between words and groups.
    - `word_id` (Foreign Key): References words.id
    - `group_id` (Foreign Key): References groups.id

- study_activities — Defines different types of study activities available.
    - `id` (Primary Key): Unique identifier for each activity
    - `name` (String, Required): Name of the activity (e.g., "Flashcards", "Quiz")
    - `url` (String, Required): The full URL of the study activity

- study_sessions — Records individual study sessions.
    - `id` (Primary Key): Unique identifier for each session
    - `group_id` (Foreign Key): References groups.id
    - `study_activity_id` (Foreign Key): References study_activities.id
    - `created_at` (Timestamp, Default: Current Time): When the session was created

- word_review_items — Tracks individual word reviews within study sessions.
    - `id` (Primary Key): Unique identifier for each review
    - `word_id` (Foreign Key): References words.id
    - `study_session_id` (Foreign Key): References study_sessions.id
    - `correct` (Boolean, Required): Whether the answer was correct
    - `created_at` (Timestamp, Default: Current Time): When the review occurred

## Relationships

- word belongs to groups through  word_groups
- group belongs to words through word_groups
- session belongs to a group
- session belongs to a study_activity
- session has many word_review_items
- word_review_item belongs to a study_session
- word_review_item belongs to a word

## Design Notes
- All tables use auto-incrementing primary keys
- Timestamps are automatically set on creation where applicable
- Foreign key constraints maintain referential integrity
- JSON storage for word parts allows flexible component storage
- Counter cache on groups.words_count optimizes word counting queries

## API

### Routes

- GET /api/words - Get paginated list of words with review statistics
- GET /api/words/:id - Get a single word by ID
- GET /api/groups - Get paginated list of word groups with word counts
- GET /api/group/:id - Get words from a specific group (This is intended to be used by target apps)
- GET /api/group/:id/words - Get words from a specific group
- GET /api/group/:id/study_sessions - Get study sessions from a specific group
- GET /api/dashboard/last_study_session
- GET /api/dashboard/study_progress
- GET /api/dashboard/quick_stats
- GET /api/study_activity/:id
- GET /api/study_activity/:id/study_sessions
- GET /api/study_session/:id/words
- POST /api/study_activities
- POST /api/study_sessions - Create a new study session for a group
- POST /api/study_sessions/:id/review - Log a review attempt for a word during a study session
- POST /api/settings/full_reset
- POST /api/settings/load_seed_data


### Query Parameters

- GET /api/words
    - page: Page number (default: 1)
    - sort_by: Sort field ('kanji', 'romaji', 'english', 'correct_count', 'wrong_count') (default: 'kanji')
    - order: Sort order ('asc' or 'desc') (default: 'asc')

- GET /groups/:id
    - page: Page number (default: 1)
    - sort_by: Sort field ('name', 'words_count') (default: 'name')
    - order: Sort order ('asc' or 'desc') (default: 'asc')

- POST /api/study_sessions
    - group_id: ID of the group to study (required)
    - study_activity_id: ID of the study activity (required)


- POST /api/study_sessions/:id/review
    - word_id: ID of the word to review (required)
    - correct: Boolean indicating if the answer was correct (required)

- POST /api/study_activities
    - group_id: ID of the group to study (required)
    - study_activity_id: ID of the study activity (required)



