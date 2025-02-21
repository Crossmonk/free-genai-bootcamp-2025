# Backend Server Technical Specs

## Business Goal:

A language learning school wants to build a prototype of learning portal which will act as three things:

- Inventory of possible vocabulary that can be learned
- Act as a Learning record store (LRS), providing correct and wrong score on practice vocabulary
- A unified launchpad to launch different learning apps

## Technical Restrictions:

- Use SQLite3 as the database
- The backend server will be written in Go
- Does not require authentication/authorization, assume there is a single user
- API will be code using Gin
- Mage is a task runner for GO.
- The API will always return JSON
- The API will always receive JSON
- The API will be documented using Swagger and OpenAPI

### Directory Structure

```text
backend-go/
├── cmd/                      # Main applications of the project
│   └── api/                  # Our main API application
│       └── main.go          # Entry point
├── internal/                 # Private application and library code
│   ├── api/                 # API specific code
│   │   ├── handlers/        # HTTP request handlers
│   │   ├── middleware/      # HTTP middleware
│   │   └── router/         # Route definitions
│   ├── domain/             # Business/domain models
│   │   └── models/         # Data structures
│   ├── repository/         # Data access layer
│   │   └── sqlite/         # SQLite specific implementations
│   └── service/            # Business logic layer
├── pkg/                     # Library code that could be used by external applications
│   ├── config/             # Configuration handling
│   └── logger/             # Logging utilities
├── migrations/             # Database migrations
├── seeds/                  # Seed data files
├── api/                    # OpenAPI/Swagger specs
│   └── swagger.yaml
├── scripts/                # Scripts for development
├── test/                   # Additional test code
├── magefile.go            # Mage task definitions
├── go.mod
└── README.md
```

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

- word belongs to groups through word_groups
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

#### GET /api/words

Get paginated list of words with review statistics

**_ Response: _**

```json
{
  "data": {
    "words": [
      {
        "id": 1,
        "kanji": "食べる",
        "romaji": "taberu",
        "english": "to eat",
        "parts": {
          "verb_type": "ru-verb",
          "topic": "food"
        },
        "stats": {
          "correct_count": 10,
          "wrong_count": 2,
          "accuracy": 83.33
        }
      }
    ],
    "pagination": {
      "current_page": 1,
      "total_pages": 10,
      "total_items": 100,
      "items_per_page": 10
    }
  }
}
```

#### GET /api/words/:id

Get a single word by ID

**_ Response: _**

```json
{
  "data": {
    "id": 1,
    "kanji": "食べる",
    "romaji": "taberu",
    "english": "to eat",
    "parts": {
      "verb_type": "ru-verb",
      "topic": "food"
    },
    "stats": {
      "correct_count": 10,
      "wrong_count": 2,
      "accuracy": 83.33
    },
    "groups": [
      {
        "id": 1,
        "name": "Basic Verbs"
      }
    ]
  }
}
```

#### GET /api/groups

Get paginated list of word groups with word counts

**_ Response: _**

```json
{
  "data": {
    "groups": [
      {
        "id": 1,
        "name": "Basic Verbs",
        "words_count": 20,
        "last_studied_at": "2024-03-20T15:30:00Z"
      }
    ],
    "pagination": {
      "current_page": 1,
      "total_pages": 5,
      "total_items": 50,
      "items_per_page": 10
    }
  }
}
```

#### GET /api/group/:id

Get words from a specific group (This is intended to be used by target apps)

**_ Response: _**

```json
{
  "data": {
    "id": 1,
    "name": "Basic Verbs",
    "words_count": 20,
    "last_studied_at": "2024-03-20T15:30:00Z",
    "study_stats": {
      "total_reviews": 150,
      "correct_reviews": 120,
      "accuracy": 80.0
    }
  }
}
```

#### GET /api/group/:id/words

Get all words belonging to a specific group with their review statistics

Response:

```json
{
  "data": {
    "group": {
      "id": 1,
      "name": "Basic Verbs"
    },
    "words": [
      {
        "id": 1,
        "kanji": "食べる",
        "romaji": "taberu",
        "english": "to eat",
        "parts": {
          "verb_type": "ru-verb",
          "topic": "food"
        },
        "stats": {
          "correct_count": 10,
          "wrong_count": 2,
          "accuracy": 83.33
        }
      }
    ],
    "pagination": {
      "current_page": 1,
      "total_pages": 5,
      "total_items": 45,
      "items_per_page": 10
    }
  }
}
```

#### GET /api/group/:id/study_sessions

Get all study sessions for a specific group with performance metrics

Response:

```json
{
  "data": {
    "group": {
      "id": 1,
      "name": "Basic Verbs"
    },
    "study_sessions": [
      {
        "id": 1,
        "created_at": "2024-03-20T15:30:00Z",
        "study_activity": {
          "id": 1,
          "name": "Flashcards",
          "url": "https://example.com/flashcards"
        },
        "stats": {
          "total_reviews": 20,
          "correct_reviews": 15,
          "accuracy": 75.0,
          "duration_minutes": 15
        }
      }
    ],
    "pagination": {
      "current_page": 1,
      "total_pages": 3,
      "total_items": 25,
      "items_per_page": 10
    }
  }
}
```

#### GET /api/dashboard/last_study_session

Get details about the most recent study session

Response:

```json
{
  "data": {
    "session": {
      "id": 1,
      "created_at": "2024-03-20T15:30:00Z",
      "group": {
        "id": 1,
        "name": "Basic Verbs"
      },
      "study_activity": {
        "id": 1,
        "name": "Flashcards",
        "url": "https://example.com/flashcards"
      },
      "stats": {
        "total_reviews": 20,
        "correct_reviews": 15,
        "accuracy": 75.0,
        "duration_minutes": 15
      },
      "words_reviewed": [
        {
          "id": 1,
          "kanji": "食べる",
          "correct": true
        }
      ]
    }
  }
}
```

#### GET /api/dashboard/study_progress

Get study progress over time (default: last 7 days)

Response:

```json
{
  "data": {
    "time_range": {
      "start_date": "2024-03-14T00:00:00Z",
      "end_date": "2024-03-20T23:59:59Z"
    },
    "daily_stats": [
      {
        "date": "2024-03-20",
        "total_sessions": 3,
        "total_reviews": 60,
        "correct_reviews": 45,
        "accuracy": 75.0,
        "total_duration_minutes": 45,
        "groups_studied": [
          {
            "id": 1,
            "name": "Basic Verbs",
            "reviews": 30
          }
        ]
      }
    ],
    "overall_stats": {
      "total_sessions": 12,
      "total_reviews": 240,
      "average_daily_reviews": 34.3,
      "average_accuracy": 78.5,
      "total_duration_minutes": 180
    }
  }
}
```

#### GET /api/dashboard/quick_stats

Get summary statistics for quick overview

Response:

```json
{
  "data": {
    "total_words": 500,
    "total_groups": 25,
    "study_stats": {
      "today": {
        "sessions": 3,
        "reviews": 60,
        "accuracy": 75.0,
        "duration_minutes": 45
      },
      "this_week": {
        "sessions": 12,
        "reviews": 240,
        "accuracy": 78.5,
        "duration_minutes": 180
      },
      "all_time": {
        "total_sessions": 150,
        "total_reviews": 3000,
        "average_accuracy": 82.3,
        "total_duration_hours": 50
      }
    },
    "recent_activity": {
      "last_session_at": "2024-03-20T15:30:00Z",
      "last_group_studied": {
        "id": 1,
        "name": "Basic Verbs"
      },
      "streak_days": 5
    },
    "top_groups": [
      {
        "id": 1,
        "name": "Basic Verbs",
        "accuracy": 85.5,
        "total_reviews": 500
      }
    ]
  }
}
```

#### GET /api/study_activity/:id

Get details about a specific study activity

Response:

```json
{
  "data": {
    "activity": {
      "id": 1,
      "name": "Flashcards",
      "url": "https://example.com/flashcards",
      "stats": {
        "total_sessions": 50,
        "total_reviews": 1000,
        "average_accuracy": 78.5,
        "total_duration_minutes": 600
      },
      "last_used": "2024-03-20T15:30:00Z"
    }
  }
}
```

#### GET /api/study_activity/:id/study_sessions

Get all study sessions for a specific activity type

Response:

```json
{
  "data": {
    "activity": {
      "id": 1,
      "name": "Flashcards",
      "url": "https://example.com/flashcards"
    },
    "study_sessions": [
      {
        "id": 1,
        "created_at": "2024-03-20T15:30:00Z",
        "group": {
          "id": 1,
          "name": "Basic Verbs"
        },
        "stats": {
          "total_reviews": 20,
          "correct_reviews": 15,
          "accuracy": 75.0,
          "duration_minutes": 15
        }
      }
    ],
    "pagination": {
      "current_page": 1,
      "total_pages": 5,
      "total_items": 50,
      "items_per_page": 10
    }
  }
}
```

#### POST /api/study_activities

Create a new study activity

Request:

```json
{
  "name": "Matching Game",
  "url": "https://example.com/matching-game"
}
```

Response:

```json
{
  "data": {
    "id": 2,
    "name": "Matching Game",
    "url": "https://example.com/matching-game",
    "created_at": "2024-03-21T10:00:00Z"
  }
}
```

#### GET /api/study_session/:id/words

Get all words reviewed in a specific study session

Response:

```json
{
  "data": {
    "session": {
      "id": 1,
      "created_at": "2024-03-20T15:30:00Z",
      "study_activity": {
        "id": 1,
        "name": "Flashcards"
      },
      "group": {
        "id": 1,
        "name": "Basic Verbs"
      }
    },
    "word_reviews": [
      {
        "id": 1,
        "word": {
          "id": 1,
          "kanji": "食べる",
          "romaji": "taberu",
          "english": "to eat"
        },
        "correct": true,
        "reviewed_at": "2024-03-20T15:31:00Z"
      }
    ],
    "stats": {
      "total_reviews": 20,
      "correct_reviews": 15,
      "accuracy": 75.0,
      "duration_minutes": 15
    }
  }
}
```

#### POST /api/study_sessions
Create a new study session for a group

Request:
```json
{
  "group_id": 1,
  "study_activity_id": 1
}
```

Response:
```json
{
  "data": {
    "session": {
      "id": 1,
      "created_at": "2024-03-21T10:00:00Z",
      "group": {
        "id": 1,
        "name": "Basic Verbs"
      },
      "study_activity": {
        "id": 1,
        "name": "Flashcards",
        "url": "https://example.com/flashcards"
      }
    }
  }
}
```

#### POST /api/study_sessions/:id/review
Log a review attempt for a word during a study session

Request:
```json
{
  "word_id": 1,
  "correct": true
}
```

Response:
```json
{
  "data": {
    "review": {
      "id": 1,
      "word_id": 1,
      "study_session_id": 1,
      "correct": true,
      "created_at": "2024-03-21T10:01:00Z"
    },
    "updated_stats": {
      "total_reviews": 21,
      "correct_reviews": 16,
      "accuracy": 76.19,
      "duration_minutes": 16
    }
  }
}
```

#### POST /api/settings/full_reset
Reset all data in the system

Response:
```json
{
  "data": {
    "status": "success",
    "message": "All data has been reset",
    "timestamp": "2024-03-21T10:00:00Z"
  }
}
```

#### POST /api/settings/load_seed_data
Load initial seed data into the system

Response:
```json
{
  "data": {
    "status": "success",
    "message": "Seed data has been loaded",
    "summary": {
      "words_created": 100,
      "groups_created": 5,
      "study_activities_created": 3
    },
    "timestamp": "2024-03-21T10:00:00Z"
  }
}
```


### Query Parameters

- GET /api/words
  - page: Page number (default: 1)
  - sort_by: Sort field ('kanji', 'romaji', 'english', 'correct_count', 'wrong_count') (default: 'kanji')
  - order: Sort order ('asc' or 'desc') (default: 'asc')

- GET /groups/:id
  - page: Page number (default: 1)
  - sort_by: Sort field ('name', 'words_count') (default: 'name')
  - order: Sort order ('asc' or 'desc') (default: 'asc')

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

## Task Runner Tasks

Mage is a task runner for GO.
Lets list out possible tasks we need for our lang portal.

### Initialize Database
This task will initialize the SQLLite Database called `words.db`

### Migrate Database
This task will run a series of migrations sql files on the database.

Migration files are located in the `migrations` folder.
The migration files will be run in order of their file name.

The file names should look like:
```
000001_initial_migration.sql
000002_words.sql
000003_groups.sql
000004_study_activities.sql
000005_study_sessions.sql
000006_word_reviews.sql
```

### Seed Data
This task will import json files and transform them into target data for our database.

All seed files live in the `seeds` folder.
All seed files should be loaded.

In our task we should have DSL to specific each seed file and its expected group word name

```json
[
    {
          "kanji": "食べる",
          "romaji": "taberu",
          "english": "to eat"
    }

]
```