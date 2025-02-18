# Frontend Technical Specs

## Role

You are building a web-app for a language learning platform.

## Tech Stack

- React.js
- Tailwind CSS
- Shadcn UI
- Lucide Icons
- React Router
- React Query
- React Hook Form
- React Toastify
- React Icons 

## Pages

### Dashboard `/dashboard`

#### Purpose

The purpose of this page is to provide a summary of learning and act as the default page when a user visit the web-app.

#### Components

- Last Study Session
    - Shows last activity used
    - Shows when last activity was used
    - Summarizes wrong vs correct from last activity
    - Has a link to the group
- Study Progress
    - Total words study eg. 10/125
        - Across all study sessions show the total words studied of all possible words in our database
    - Dispplay a mastery percentage eg. 8%

- Quick stats
    - Success rate eg. 80%
    - Total study sessions eg. 10
    - Total active groups eg. 3
    - Study streak eg. 4 days
- Start Studying Button
    - Goes to study activities page


#### Endpoints

- GET /api/dashboard/last_study_session
- GET /api/dashboard/study_progress
- GET /api/dashboard/quick_stats

#### Groups

- List of groups
- Create new group
- Delete group

### Study Activities `/study_activities`

The purpose of this page is to show a collection of study activities with a tumbnail and its name, to either launch or view the study activity.

#### Components

- Study Activity Card
    - Thumbnail
    - Name
    - Description
    - Launch Button
    - View Button

#### Endpoints

- GET /api/study_activities


### Study Activity `/study_activity`

The purpose of this page is to show the study activity and its past study sessions.

#### Components

- Name of the study activity
- Thumbnail of study activity
- Description of study activity
- Launch button
- Study activities paginated by liust
    - id
    - activity name
    - group name
    - start time
    - end time (inferred by the last word_review_item submited)

#### Endpoints

- GET /api/study_activity/:id
- GET /api/study_activity/:id/study_sessions

### Study Activity Launch `/study_activity/:id/launch`

The purpose of this page is to launch the study activity.

#### Components

- Name of the Study Activity
- Launch form
    - Select field for group
    - Launch now button

### Behaviour

- After the form is submitted, a new tab opens with the study activity based on its URL provided by the backend.
- When a user clicks on the launch now button, the page will navigate to the study session show page.


#### Endpoints

- POST /api/study_activities
    - parameters:
        - group_id
        - study_activity_id


### Words
#### Purpose

The purpose of this page is to show a list of words in our database.

#### Components

- Paginated list of words
    - Columns:
        - Japanese
        - Romaji
        - English
        - correct count
        - wrong count
    - Pagination with 100 items per page
    - Clicking the Japanese word will take you to the word show page.

#### Endpoints

- GET /api/words

### Word Show `/word/:id`

The purpose of this page is to show a single word with its japanese, romaji, english, correct count, wrong count and its word parts.

#### Components

- Japanese
- Romaji
-English
- Study statistics
    - Correct count
    - Wrong Count
- Word groups
    - shown as a series of pills eg. tags
    - when a group name is clieked it will take us to the groups show page.

#### Endpoints

- GET /api/words/:id

### Word Group `/groups`

The purpose of this page is to show a list of word groupsin our databse.

#### Components

- Paginated group list
    - Columns:
        - Group Item
        - Word Count
    - Pagination with 100 items per page
    - Clicking the Group Item will take us to the group show page.

#### Endpoints

- GET /api/groups

### Group Show `/group/:id`

The purpose of this page is to show a single word group with its name, description, words and its study statistics.

#### Components

- Name of the group
- Description of the group
- Word Count
- Study statistics
    - Total word count
- Words in Group (paginated)
    - Should use the same components as the words index page.
- Study sessions (paginated)
    - Should use the same components as the study sessions index page.

#### Endpoints

- GET /api/group/:id (The group and study sessions)
- GET /api/group/:id/words
- GET /api/group/:id/study_sessions

### Study Sessions `/study_sessions`

The purpose of this page is to show a list of study sessions in our database.

#### Components 

- Paginated list of study sessions
    - Columns:
        - Id
        - Name
        - Group Name
        - Start Time
        - End Time
        - Number of Review Items
    - Pagination with 100 items per page
    - Clicking the Study Session Item will take us to the study session show page.

#### Endpoints

- GET /api/study_sessions

### Study Session Show `/study_session/:id`

The purpose of this page is to show a single study session with its name, group, start time, end time, and review items.

#### Components
- Study Session Details
    - Name of the study session
    - Group name
    - Start time
    - End time
    - Review items (paginated)
        - Should use the same components as the review items index page.
- Words Reviewed (paginated)
    - Should use the same components as the words index page.

#### Endpoints

- GET /api/study_session/:id
- GET /api/study_session/:id/words

### Settings `/settings`

The purpose of this page is to configure the study portal.

#### Components

- Theme selection eg. light, dark
- Language selection eg. English, Japanese
- Reset History
    - This will delete all the study sessions and words review items.
- Full Reset
    - This will truncate all tables and re-create with seed data.
- Load  Seed Data

#### Endpoints

- POST /api/settings/full_reset
- POST /api/settings/load_seed_data

