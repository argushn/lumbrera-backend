# Lumbrera - Gamified Bible Study App Architecture

**Version:** 2026.01.0
**Last Updated:** 2026-01-11
**Versioning:** CalVer (YYYY.MM.MICRO)

## Table of Contents

1. [Project Overview](#project-overview)
2. [Tech Stack](#tech-stack)
3. [Architecture Overview](#architecture-overview)
4. [Features & Gamification](#features--gamification)
5. [Project Structure](#project-structure)
6. [Data Models](#data-models)
7. [API Design](#api-design)
8. [Authentication & Authorization](#authentication--authorization)
9. [Testing Strategy](#testing-strategy)
10. [Development Workflow](#development-workflow)
11. [Deployment Strategy](#deployment-strategy)

---

## Project Overview

**Lumbrera** is a gamified mobile Bible study application inspired by Duolingo's learning approach. The app makes biblical studies engaging through:

- **Bite-sized lessons**: Daily structured learning paths
- **Gamification**: Points, streaks, achievements, and leaderboards
- **Social learning**: Community challenges and friend competitions
- **Progress tracking**: Comprehensive analytics and milestones

### Core Principles

- **Test-Driven Development (TDD)**: All features developed test-first
- **Monorepo**: Frontend and backend in single repository
- **CalVer**: Versioning scheme YYYY.MM.MICRO
- **API-First**: OpenAPI documented endpoints
- **Mobile-First**: Optimized for mobile experience

---

## Tech Stack

### Frontend
- **Framework**: Flutter 3.x
- **State Management**: Riverpod / Bloc
- **Navigation**: Go Router
- **HTTP Client**: Dio with interceptors
- **Local Storage**: Hive / SQLite
- **Testing**: Flutter Test, Integration Tests, E2E (Cucumber/Gherkin)

### Backend
- **Language**: Go 1.22+
- **Runtime**: AWS Lambda (Serverless)
- **Framework**: AWS Lambda Go, Serverless Framework v3
- **Database**: DynamoDB
- **Authentication**: Auth0 (Facebook, Google OAuth)
- **API Gateway**: AWS API Gateway with custom authorizers
- **Testing**: Go testing, Testify, Ginkgo/Gomega

### Infrastructure
- **Containerization**: Docker, Docker Compose
- **IaC**: Serverless Framework, AWS CloudFormation
- **CI/CD**: GitHub Actions
- **Monitoring**: AWS CloudWatch, X-Ray
- **Caching**: DynamoDB DAX (optional), CloudFront

### Development Tools
- **E2E Testing**: Cucumber/Gherkin with Godog (backend), Flutter Cucumber (frontend)
- **API Documentation**: OpenAPI 3.0, Swagger UI
- **Code Quality**: golangci-lint, flutter analyze
- **Git Hooks**: pre-commit, pre-push

---

## Architecture Overview

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Mobile App (Flutter)                     │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │  Auth    │  │  Lessons │  │  Social  │  │  Profile │   │
│  │  Module  │  │  Module  │  │  Module  │  │  Module  │   │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘   │
└─────────────────────────┬───────────────────────────────────┘
                          │ HTTPS/REST
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                    Auth0 Identity Platform                   │
│              (Facebook, Google OAuth 2.0)                    │
└─────────────────────────┬───────────────────────────────────┘
                          │ JWT Tokens
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                      AWS API Gateway                         │
│                  (Custom JWT Authorizer)                     │
└─────────────────────────┬───────────────────────────────────┘
                          │
           ┌──────────────┼──────────────┐
           ▼              ▼              ▼
    ┌──────────┐   ┌──────────┐   ┌──────────┐
    │  Lambda  │   │  Lambda  │   │  Lambda  │
    │  Users   │   │  Lessons │   │  Social  │
    └────┬─────┘   └────┬─────┘   └────┬─────┘
         │              │              │
         └──────────────┼──────────────┘
                        ▼
              ┌──────────────────┐
              │   DynamoDB       │
              │   Tables         │
              │  - Users         │
              │  - Lessons       │
              │  - Progress      │
              │  - Achievements  │
              │  - Leaderboards  │
              └──────────────────┘
```

### Service Architecture

#### Microservices (Lambda Functions)

1. **User Service**
   - User profile management
   - Preferences and settings
   - User statistics

2. **Lesson Service**
   - CRUD operations for lessons
   - Lesson content delivery
   - Lesson progression tracking

3. **Progress Service**
   - Track user progress
   - Calculate XP and levels
   - Manage streaks

4. **Achievement Service**
   - Award achievements
   - Track badges
   - Milestone notifications

5. **Social Service**
   - Leaderboards (global, friends, groups)
   - Friend management
   - Community challenges

6. **Quiz Service**
   - Quiz generation and validation
   - Answer checking
   - Performance analytics

---

## Features & Gamification

### Core Features

#### 1. Structured Learning Path
- **Daily Lessons**: Bite-sized 5-10 minute lessons
- **Skill Trees**: Organized by books, themes, or difficulty
- **Progressive Unlocking**: Complete prerequisites to unlock advanced content
- **Multiple Formats**: Reading, multiple choice, fill-in-the-blank, matching

#### 2. Gamification Elements

##### XP (Experience Points)
- **Lesson Completion**: 10-50 XP per lesson (based on difficulty)
- **Perfect Score**: +20% bonus XP
- **Daily Goal**: 50 XP minimum
- **Streak Bonus**: +5 XP per day of current streak

##### Levels
- **Level 1-100**: Exponential XP requirements
- **Level Benefits**: Unlock exclusive content, badges
- **Prestige System**: Reset at level 100 with special perks

##### Streaks
- **Daily Streaks**: Consecutive days of study
- **Streak Freeze**: 2 per month (DLC or premium)
- **Streak Milestones**: 7, 30, 100, 365 days
- **Recovery Window**: 24-hour window to maintain streak

##### Achievements/Badges
- **Scholar**: Complete 100 lessons
- **Devoted**: 30-day streak
- **Perfectionist**: 10 perfect scores in a row
- **Early Bird**: Study before 9 AM for 7 days
- **Night Owl**: Study after 9 PM for 7 days
- **Bookworm**: Complete entire book of the Bible
- **Social Butterfly**: Invite 5 friends

##### Leaderboards
- **Global**: Top 100 users worldwide
- **Friends**: Compete with connected friends
- **Local**: Top users in same country/region
- **Weekly/Monthly/All-Time**: Multiple timeframes
- **Division System**: Bronze, Silver, Gold, Diamond leagues

#### 3. Social Features
- **Friend System**: Add friends via Auth0 social connections
- **Friend Challenges**: Head-to-head competitions
- **Study Groups**: Join or create study groups
- **Share Progress**: Share achievements on social media
- **Community Challenges**: Weekly themed challenges

#### 4. Content Types

##### Lesson Formats
- **Reading Comprehension**: Read passage, answer questions
- **Memory Verses**: Memorize and recite verses
- **Multiple Choice**: Biblical trivia and facts
- **Fill in the Blank**: Complete missing words
- **Matching**: Match concepts, people, events
- **Timeline Ordering**: Arrange events chronologically
- **Map Activities**: Locate biblical places

##### Difficulty Levels
- **Beginner**: Basic stories and concepts
- **Intermediate**: Deeper theology and context
- **Advanced**: Original languages, historical context
- **Expert**: Scholarly interpretations, cross-references

#### 5. Progress Tracking
- **Daily Statistics**: Lessons completed, XP earned, time studied
- **Weekly Reports**: Email summaries of progress
- **Monthly Insights**: Strengths, weaknesses, recommendations
- **Yearly Review**: Annual accomplishments recap
- **Visual Charts**: Progress graphs and heatmaps

#### 6. Personalization
- **Study Reminders**: Customizable push notifications
- **Preferred Translation**: NIV, KJV, ESV, NLT, etc.
- **Difficulty Adjustment**: Adaptive learning based on performance
- **Learning Pace**: Relaxed, Regular, Intense modes
- **Theme Preferences**: Light/Dark mode, color schemes

---

## Project Structure

### Monorepo Layout

```
lumbrera/
├── backend/                          # Go backend services
│   ├── functions/                    # Lambda functions
│   │   ├── users/
│   │   │   ├── create/
│   │   │   │   ├── handler.go
│   │   │   │   └── handler_test.go
│   │   │   ├── get/
│   │   │   ├── update/
│   │   │   └── delete/
│   │   ├── lessons/
│   │   │   ├── create/
│   │   │   ├── get/
│   │   │   ├── list/
│   │   │   └── update/
│   │   ├── progress/
│   │   │   ├── update/
│   │   │   ├── get/
│   │   │   └── streak/
│   │   ├── achievements/
│   │   │   ├── award/
│   │   │   └── list/
│   │   ├── social/
│   │   │   ├── leaderboard/
│   │   │   ├── friends/
│   │   │   └── challenges/
│   │   └── auth/
│   │       ├── authorize/            # Auth0 JWT validation
│   │       └── webhook/              # Auth0 user sync
│   ├── internal/
│   │   ├── database/                 # Database layer
│   │   │   ├── dynamodb.go
│   │   │   ├── repository/
│   │   │   │   ├── user_repository.go
│   │   │   │   ├── lesson_repository.go
│   │   │   │   ├── progress_repository.go
│   │   │   │   └── achievement_repository.go
│   │   │   └── migrations/           # Table schemas
│   │   ├── models/                   # Domain models
│   │   │   ├── user.go
│   │   │   ├── lesson.go
│   │   │   ├── progress.go
│   │   │   ├── achievement.go
│   │   │   ├── leaderboard.go
│   │   │   └── quiz.go
│   │   ├── services/                 # Business logic
│   │   │   ├── user_service.go
│   │   │   ├── lesson_service.go
│   │   │   ├── gamification_service.go
│   │   │   └── leaderboard_service.go
│   │   ├── auth/                     # Auth0 integration
│   │   │   ├── jwt_validator.go
│   │   │   └── permissions.go
│   │   ├── utils/                    # Utilities
│   │   │   ├── response.go
│   │   │   ├── errors.go
│   │   │   └── validators.go
│   │   └── config/                   # Configuration
│   │       └── config.go
│   ├── tests/
│   │   ├── integration/              # Integration tests
│   │   │   └── api_test.go
│   │   └── e2e/                      # E2E Cucumber tests
│   │       ├── features/
│   │       │   ├── user.feature
│   │       │   ├── lessons.feature
│   │       │   ├── progress.feature
│   │       │   └── leaderboard.feature
│   │       └── steps/
│   │           ├── user_steps.go
│   │           └── lesson_steps.go
│   ├── go.mod
│   ├── go.sum
│   ├── serverless.yml                # Serverless config
│   ├── Makefile
│   └── README.md
│
├── mobile/                            # Flutter mobile app
│   ├── lib/
│   │   ├── main.dart
│   │   ├── app/
│   │   │   ├── app.dart
│   │   │   └── routes.dart
│   │   ├── core/
│   │   │   ├── config/
│   │   │   │   ├── env.dart
│   │   │   │   └── theme.dart
│   │   │   ├── network/
│   │   │   │   ├── api_client.dart
│   │   │   │   ├── auth_interceptor.dart
│   │   │   │   └── error_handler.dart
│   │   │   ├── storage/
│   │   │   │   └── local_storage.dart
│   │   │   └── utils/
│   │   │       ├── constants.dart
│   │   │       └── validators.dart
│   │   ├── features/
│   │   │   ├── auth/
│   │   │   │   ├── data/
│   │   │   │   │   ├── models/
│   │   │   │   │   ├── repositories/
│   │   │   │   │   └── datasources/
│   │   │   │   ├── domain/
│   │   │   │   │   ├── entities/
│   │   │   │   │   ├── repositories/
│   │   │   │   │   └── usecases/
│   │   │   │   └── presentation/
│   │   │   │       ├── providers/
│   │   │   │       ├── screens/
│   │   │   │       └── widgets/
│   │   │   ├── lessons/
│   │   │   │   ├── data/
│   │   │   │   ├── domain/
│   │   │   │   └── presentation/
│   │   │   ├── progress/
│   │   │   │   ├── data/
│   │   │   │   ├── domain/
│   │   │   │   └── presentation/
│   │   │   ├── social/
│   │   │   │   ├── data/
│   │   │   │   ├── domain/
│   │   │   │   └── presentation/
│   │   │   └── profile/
│   │   │       ├── data/
│   │   │       ├── domain/
│   │   │       └── presentation/
│   │   └── shared/
│   │       ├── widgets/
│   │       │   ├── buttons/
│   │       │   ├── cards/
│   │       │   └── dialogs/
│   │       └── extensions/
│   ├── test/
│   │   ├── unit/
│   │   ├── widget/
│   │   └── integration/
│   ├── integration_test/
│   │   └── app_test.dart
│   ├── test_driver/
│   │   └── e2e/
│   │       ├── features/
│   │       │   ├── user_journey.feature
│   │       │   └── lesson_flow.feature
│   │       └── steps/
│   ├── pubspec.yaml
│   ├── analysis_options.yaml
│   └── README.md
│
├── docs/                              # Documentation
│   ├── ARCHITECTURE.md               # This file
│   ├── API.md                        # API documentation
│   ├── DATABASE.md                   # Database schemas
│   ├── DEPLOYMENT.md                 # Deployment guide
│   ├── TESTING.md                    # Testing strategy
│   └── openapi.yaml                  # OpenAPI specification
│
├── infrastructure/                    # Infrastructure as Code
│   ├── docker/
│   │   ├── backend/
│   │   │   └── Dockerfile
│   │   └── mobile/
│   │       └── Dockerfile
│   ├── docker-compose.yml            # Local development
│   ├── docker-compose.test.yml       # Testing environment
│   └── scripts/
│       ├── init-dynamodb.sh
│       └── seed-data.sh
│
├── .github/
│   ├── workflows/
│   │   ├── backend-ci.yml
│   │   ├── mobile-ci.yml
│   │   ├── e2e-tests.yml
│   │   └── deploy.yml
│   └── PULL_REQUEST_TEMPLATE.md
│
├── scripts/                           # Utility scripts
│   ├── version.sh                    # CalVer versioning
│   ├── setup-dev.sh                  # Development setup
│   └── run-e2e.sh                    # Run E2E tests
│
├── .gitignore
├── .env.example
├── VERSION                            # CalVer version file
└── README.md                          # Main README
```

---

## Data Models

### DynamoDB Table Design

#### Single-Table Design Pattern

Using DynamoDB single-table design for optimal performance and cost efficiency.

**Table Name**: `lumbrera-{stage}`

**Indexes**:
- **Primary Key**: `PK` (Partition Key), `SK` (Sort Key)
- **GSI1**: `GSI1PK`, `GSI1SK` - For queries by user
- **GSI2**: `GSI2PK`, `GSI2SK` - For leaderboards and rankings

#### Entity Schemas

##### 1. User

```go
type User struct {
    // DynamoDB Keys
    PK        string `dynamodbav:"PK" json:"-"`              // USER#{UserID}
    SK        string `dynamodbav:"SK" json:"-"`              // PROFILE
    GSI1PK    string `dynamodbav:"GSI1PK" json:"-"`          // USER#{UserID}
    GSI1SK    string `dynamodbav:"GSI1SK" json:"-"`          // METADATA

    // Core Fields
    UserID    string `dynamodbav:"UserID" json:"user_id"`
    Auth0ID   string `dynamodbav:"Auth0ID" json:"auth0_id"`
    Email     string `dynamodbav:"Email" json:"email"`
    Name      string `dynamodbav:"Name" json:"name"`
    Avatar    string `dynamodbav:"Avatar" json:"avatar"`

    // Gamification
    XP        int    `dynamodbav:"XP" json:"xp"`
    Level     int    `dynamodbav:"Level" json:"level"`
    Streak    int    `dynamodbav:"Streak" json:"streak"`

    // Preferences
    Translation      string   `dynamodbav:"Translation" json:"translation"`
    DailyGoal        int      `dynamodbav:"DailyGoal" json:"daily_goal"`
    NotificationsOn  bool     `dynamodbav:"NotificationsOn" json:"notifications_on"`

    // Metadata
    CreatedAt string `dynamodbav:"CreatedAt" json:"created_at"`
    UpdatedAt string `dynamodbav:"UpdatedAt" json:"updated_at"`
    EntityType string `dynamodbav:"EntityType" json:"-"`     // USER
}
```

##### 2. Lesson

```go
type Lesson struct {
    // DynamoDB Keys
    PK        string `dynamodbav:"PK" json:"-"`              // LESSON#{LessonID}
    SK        string `dynamodbav:"SK" json:"-"`              // METADATA

    // Core Fields
    LessonID     string `dynamodbav:"LessonID" json:"lesson_id"`
    Title        string `dynamodbav:"Title" json:"title"`
    Description  string `dynamodbav:"Description" json:"description"`
    Book         string `dynamodbav:"Book" json:"book"`           // Genesis, Exodus, etc.
    Chapter      int    `dynamodbav:"Chapter" json:"chapter"`
    Verses       string `dynamodbav:"Verses" json:"verses"`       // e.g., "1-10"

    // Content
    ContentType  string `dynamodbav:"ContentType" json:"content_type"`  // reading, quiz, memory
    Content      string `dynamodbav:"Content" json:"content"`            // JSON content

    // Metadata
    Difficulty   string `dynamodbav:"Difficulty" json:"difficulty"`      // beginner, intermediate, advanced
    XPReward     int    `dynamodbav:"XPReward" json:"xp_reward"`
    EstDuration  int    `dynamodbav:"EstDuration" json:"est_duration"`   // minutes
    Order        int    `dynamodbav:"Order" json:"order"`                // Sequence in course

    // Unlocking
    Prerequisites []string `dynamodbav:"Prerequisites" json:"prerequisites"` // Lesson IDs

    // Timestamps
    CreatedAt string `dynamodbav:"CreatedAt" json:"created_at"`
    UpdatedAt string `dynamodbav:"UpdatedAt" json:"updated_at"`
    EntityType string `dynamodbav:"EntityType" json:"-"`     // LESSON
}
```

##### 3. UserProgress

```go
type UserProgress struct {
    // DynamoDB Keys
    PK        string `dynamodbav:"PK" json:"-"`              // USER#{UserID}
    SK        string `dynamodbav:"SK" json:"-"`              // PROGRESS#{LessonID}
    GSI1PK    string `dynamodbav:"GSI1PK" json:"-"`          // LESSON#{LessonID}
    GSI1SK    string `dynamodbav:"GSI1SK" json:"-"`          // USER#{UserID}

    // Core Fields
    UserID       string `dynamodbav:"UserID" json:"user_id"`
    LessonID     string `dynamodbav:"LessonID" json:"lesson_id"`

    // Progress Data
    Status       string `dynamodbav:"Status" json:"status"`           // not_started, in_progress, completed
    Score        int    `dynamodbav:"Score" json:"score"`             // 0-100
    Attempts     int    `dynamodbav:"Attempts" json:"attempts"`
    XPEarned     int    `dynamodbav:"XPEarned" json:"xp_earned"`

    // Timestamps
    StartedAt    string `dynamodbav:"StartedAt" json:"started_at"`
    CompletedAt  string `dynamodbav:"CompletedAt,omitempty" json:"completed_at,omitempty"`
    UpdatedAt    string `dynamodbav:"UpdatedAt" json:"updated_at"`
    EntityType   string `dynamodbav:"EntityType" json:"-"`           // PROGRESS
}
```

##### 4. Achievement

```go
type Achievement struct {
    // DynamoDB Keys
    PK        string `dynamodbav:"PK" json:"-"`              // ACHIEVEMENT#{AchievementID}
    SK        string `dynamodbav:"SK" json:"-"`              // METADATA

    // Core Fields
    AchievementID  string `dynamodbav:"AchievementID" json:"achievement_id"`
    Name           string `dynamodbav:"Name" json:"name"`
    Description    string `dynamodbav:"Description" json:"description"`
    Icon           string `dynamodbav:"Icon" json:"icon"`

    // Requirements
    Type           string `dynamodbav:"Type" json:"type"`              // streak, lessons, perfect_score
    Requirement    int    `dynamodbav:"Requirement" json:"requirement"` // e.g., 30 for 30-day streak

    // Metadata
    Rarity         string `dynamodbav:"Rarity" json:"rarity"`          // common, rare, epic, legendary
    XPReward       int    `dynamodbav:"XPReward" json:"xp_reward"`

    CreatedAt string `dynamodbav:"CreatedAt" json:"created_at"`
    EntityType string `dynamodbav:"EntityType" json:"-"`              // ACHIEVEMENT
}
```

##### 5. UserAchievement

```go
type UserAchievement struct {
    // DynamoDB Keys
    PK        string `dynamodbav:"PK" json:"-"`              // USER#{UserID}
    SK        string `dynamodbav:"SK" json:"-"`              // ACHIEVEMENT#{AchievementID}

    // Core Fields
    UserID         string `dynamodbav:"UserID" json:"user_id"`
    AchievementID  string `dynamodbav:"AchievementID" json:"achievement_id"`

    // Progress
    Progress       int    `dynamodbav:"Progress" json:"progress"`
    UnlockedAt     string `dynamodbav:"UnlockedAt,omitempty" json:"unlocked_at,omitempty"`
    EntityType     string `dynamodbav:"EntityType" json:"-"`          // USER_ACHIEVEMENT
}
```

##### 6. Leaderboard Entry

```go
type LeaderboardEntry struct {
    // DynamoDB Keys
    PK        string `dynamodbav:"PK" json:"-"`              // LEADERBOARD#{Type}#{Period}
    SK        string `dynamodbav:"SK" json:"-"`              // SCORE#{Score}#{UserID}
    GSI2PK    string `dynamodbav:"GSI2PK" json:"-"`          // USER#{UserID}
    GSI2SK    string `dynamodbav:"GSI2SK" json:"-"`          // LEADERBOARD#{Type}

    // Core Fields
    UserID       string `dynamodbav:"UserID" json:"user_id"`
    UserName     string `dynamodbav:"UserName" json:"user_name"`
    Avatar       string `dynamodbav:"Avatar" json:"avatar"`

    // Ranking Data
    Type         string `dynamodbav:"Type" json:"type"`           // global, friends, local
    Period       string `dynamodbav:"Period" json:"period"`       // weekly, monthly, all_time
    Score        int    `dynamodbav:"Score" json:"score"`         // XP or lessons completed
    Rank         int    `dynamodbav:"Rank" json:"rank"`

    // Metadata
    UpdatedAt    string `dynamodbav:"UpdatedAt" json:"updated_at"`
    EntityType   string `dynamodbav:"EntityType" json:"-"`       // LEADERBOARD_ENTRY
}
```

##### 7. DailyStreak

```go
type DailyStreak struct {
    // DynamoDB Keys
    PK        string `dynamodbav:"PK" json:"-"`              // USER#{UserID}
    SK        string `dynamodbav:"SK" json:"-"`              // STREAK#{Date}

    // Core Fields
    UserID       string `dynamodbav:"UserID" json:"user_id"`
    Date         string `dynamodbav:"Date" json:"date"`           // YYYY-MM-DD
    XPEarned     int    `dynamodbav:"XPEarned" json:"xp_earned"`
    LessonsCount int    `dynamodbav:"LessonsCount" json:"lessons_count"`
    GoalMet      bool   `dynamodbav:"GoalMet" json:"goal_met"`

    CreatedAt  string `dynamodbav:"CreatedAt" json:"created_at"`
    EntityType string `dynamodbav:"EntityType" json:"-"`         // DAILY_STREAK
}
```

---

## API Design

### REST API Conventions

- **Base URL**: `https://api.lumbrera.app/v1`
- **Authentication**: Bearer token (Auth0 JWT)
- **Content-Type**: `application/json`
- **Versioning**: URL path versioning (`/v1`, `/v2`)
- **Error Format**: RFC 7807 Problem Details

### Endpoints

#### Authentication

```
POST   /auth/token           # Exchange Auth0 token
GET    /auth/me              # Get current user profile
```

#### Users

```
POST   /users                # Create user (internal, Auth0 webhook)
GET    /users/{id}           # Get user by ID
PUT    /users/{id}           # Update user
DELETE /users/{id}           # Delete user
GET    /users/{id}/stats     # Get user statistics
PUT    /users/{id}/preferences # Update preferences
```

#### Lessons

```
GET    /lessons              # List lessons (paginated, filtered)
GET    /lessons/{id}         # Get lesson details
POST   /lessons              # Create lesson (admin only)
PUT    /lessons/{id}         # Update lesson (admin only)
DELETE /lessons/{id}         # Delete lesson (admin only)
GET    /lessons/recommended  # Get recommended lessons for user
```

#### Progress

```
GET    /progress             # Get user's overall progress
GET    /progress/lessons/{id} # Get progress for specific lesson
POST   /progress/lessons/{id}/start # Start a lesson
PUT    /progress/lessons/{id}/complete # Complete a lesson
POST   /progress/lessons/{id}/answer # Submit answer
GET    /progress/streak      # Get current streak info
```

#### Achievements

```
GET    /achievements         # List all achievements
GET    /achievements/{id}    # Get achievement details
GET    /users/{id}/achievements # Get user's achievements
```

#### Leaderboards

```
GET    /leaderboards/global  # Global leaderboard
GET    /leaderboards/friends # Friends leaderboard
GET    /leaderboards/local   # Local/regional leaderboard
```

#### Social

```
GET    /friends              # List friends
POST   /friends/{id}         # Add friend
DELETE /friends/{id}         # Remove friend
GET    /challenges           # List active challenges
POST   /challenges           # Create challenge
POST   /challenges/{id}/join # Join challenge
```

### Request/Response Examples

#### Complete a Lesson

**Request:**
```http
PUT /v1/progress/lessons/lesson_001/complete
Authorization: Bearer {jwt_token}
Content-Type: application/json

{
  "score": 95,
  "time_spent": 420,
  "answers": [
    {
      "question_id": "q1",
      "answer": "Moses",
      "correct": true
    }
  ]
}
```

**Response:**
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "user_progress": {
    "lesson_id": "lesson_001",
    "status": "completed",
    "score": 95,
    "xp_earned": 55,
    "completed_at": "2026-01-11T10:30:00Z"
  },
  "user_stats": {
    "total_xp": 1255,
    "level": 8,
    "streak": 15
  },
  "achievements_unlocked": [
    {
      "achievement_id": "ach_001",
      "name": "Perfect Score",
      "xp_reward": 10
    }
  ]
}
```

---

## Authentication & Authorization

### Auth0 Integration

#### Configuration

```yaml
Auth0:
  Domain: lumbrera.auth0.com
  ClientID: {FLUTTER_CLIENT_ID}
  Audience: https://api.lumbrera.app
  Scope: openid profile email
  Connections:
    - google-oauth2
    - facebook
```

#### Flow

1. **Mobile App** initiates OAuth flow with Auth0
2. **User** authenticates via Google/Facebook
3. **Auth0** returns JWT access token + ID token
4. **Mobile App** stores tokens securely
5. **API Requests** include `Authorization: Bearer {access_token}`
6. **API Gateway** validates JWT via custom authorizer
7. **Lambda** receives validated user context

#### Custom Authorizer (Lambda)

```go
// functions/auth/authorize/handler.go
func HandleAuthorize(request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
    token := extractToken(request.AuthorizationToken)

    // Validate JWT signature and claims
    claims, err := validateAuth0Token(token)
    if err != nil {
        return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
    }

    // Build IAM policy
    return generatePolicy(claims.Subject, "Allow", request.MethodArn, claims), nil
}
```

#### User Sync

Auth0 webhook triggers Lambda to create/update user in DynamoDB:

```go
// functions/auth/webhook/handler.go
func HandleAuth0Webhook(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    var webhookEvent Auth0Event
    json.Unmarshal([]byte(event.Body), &webhookEvent)

    switch webhookEvent.Type {
    case "s.signup":
        createUser(webhookEvent.Data.User)
    case "s.profile.update":
        updateUser(webhookEvent.Data.User)
    }

    return success(), nil
}
```

### Permissions

**Scopes:**
- `read:lessons` - View lesson content
- `write:progress` - Update user progress
- `read:leaderboards` - View leaderboards
- `write:social` - Manage friends, challenges
- `admin:lessons` - Create/edit lessons (admin only)

---

## Testing Strategy

### Test-Driven Development (TDD)

All features follow **Red-Green-Refactor** cycle:

1. **Red**: Write failing test
2. **Green**: Implement minimal code to pass
3. **Refactor**: Improve code while keeping tests green

### Backend Testing (Go)

#### Unit Tests
- **Coverage Target**: 80%+
- **Framework**: Go testing + Testify
- **Location**: `*_test.go` files alongside source

```go
// functions/lessons/create/handler_test.go
func TestCreateLesson(t *testing.T) {
    tests := []struct {
        name       string
        input      events.APIGatewayProxyRequest
        wantStatus int
        wantError  bool
    }{
        {
            name: "valid lesson creation",
            input: events.APIGatewayProxyRequest{
                Body: `{"title":"Genesis 1","book":"Genesis"}`,
            },
            wantStatus: 201,
            wantError: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            resp, err := HandleCreate(tt.input)
            assert.Equal(t, tt.wantStatus, resp.StatusCode)
            if tt.wantError {
                assert.Error(t, err)
            }
        })
    }
}
```

#### Integration Tests
- **Framework**: Go testing + DynamoDB Local
- **Location**: `backend/tests/integration/`

```go
// backend/tests/integration/lesson_test.go
func TestLessonLifecycle(t *testing.T) {
    // Setup
    db := setupTestDB()
    defer cleanupTestDB(db)

    // Create lesson
    lesson := createTestLesson(db)
    assert.NotEmpty(t, lesson.ID)

    // Retrieve lesson
    retrieved := getLesson(db, lesson.ID)
    assert.Equal(t, lesson.Title, retrieved.Title)
}
```

#### E2E Tests (Backend)
- **Framework**: Godog (Cucumber for Go)
- **Location**: `backend/tests/e2e/`

```gherkin
# backend/tests/e2e/features/lessons.feature
Feature: Lesson Management
  As a user
  I want to complete lessons
  So that I can learn Bible content

  Scenario: Complete a lesson successfully
    Given I am authenticated as user "john@example.com"
    And a lesson "Genesis 1:1" exists
    When I start the lesson "Genesis 1:1"
    And I submit answers with 100% accuracy
    Then the lesson should be marked as completed
    And I should receive 50 XP
    And my streak should increase by 1
```

### Frontend Testing (Flutter)

#### Unit Tests
- **Coverage Target**: 80%+
- **Framework**: Flutter Test
- **Location**: `mobile/test/unit/`

```dart
// mobile/test/unit/gamification_service_test.dart
void main() {
  group('GamificationService', () {
    test('calculates XP correctly for lesson completion', () {
      final service = GamificationService();
      final xp = service.calculateXP(
        score: 95,
        difficulty: Difficulty.intermediate,
        hasStreak: true,
      );

      expect(xp, equals(55)); // 50 base + 5 streak bonus
    });
  });
}
```

#### Widget Tests
- **Framework**: Flutter Test
- **Location**: `mobile/test/widget/`

```dart
// mobile/test/widget/lesson_card_test.dart
void main() {
  testWidgets('LessonCard displays correctly', (tester) async {
    await tester.pumpWidget(
      MaterialApp(
        home: LessonCard(
          lesson: Lesson(title: 'Genesis 1'),
        ),
      ),
    );

    expect(find.text('Genesis 1'), findsOneWidget);
    expect(find.byIcon(Icons.lock), findsNothing);
  });
}
```

#### Integration Tests
- **Framework**: Flutter Integration Test
- **Location**: `mobile/integration_test/`

```dart
// mobile/integration_test/lesson_flow_test.dart
void main() {
  IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  testWidgets('Complete lesson flow', (tester) async {
    app.main();
    await tester.pumpAndSettle();

    // Navigate to lessons
    await tester.tap(find.text('Lessons'));
    await tester.pumpAndSettle();

    // Start lesson
    await tester.tap(find.byKey(Key('lesson_1')));
    await tester.pumpAndSettle();

    // Complete lesson
    await tester.tap(find.text('Submit'));
    await tester.pumpAndSettle();

    // Verify XP increase
    expect(find.textContaining('XP'), findsOneWidget);
  });
}
```

### E2E Tests (Full Stack)

**Framework**: Cucumber/Gherkin with both backend (Godog) and frontend (flutter_gherkin)

```gherkin
# E2E feature spanning frontend and backend
Feature: User Journey - First Lesson
  As a new user
  I want to complete my first lesson
  So that I can start learning

  Scenario: New user completes first lesson
    Given I have installed the app
    When I sign up with Google
    Then I should see the onboarding screen
    When I complete the onboarding
    Then I should see the lesson dashboard
    When I select "Genesis 1:1" lesson
    And I complete all questions correctly
    Then I should see a celebration animation
    And my XP should be 50
    And I should unlock the "First Steps" achievement
```

### Test Environments

```yaml
Local:
  Backend: DynamoDB Local, Lambda containers
  Frontend: Flutter on emulator/simulator

CI:
  Backend: GitHub Actions + DynamoDB Local
  Frontend: GitHub Actions + Flutter Test

Staging:
  Backend: AWS Lambda (staging stage)
  Frontend: Firebase App Distribution

Production:
  Backend: AWS Lambda (prod stage)
  Frontend: App Store / Play Store
```

---

## Development Workflow

### Local Development Setup

```bash
# Clone repository
git clone https://github.com/your-org/lumbrera.git
cd lumbrera

# Run setup script
./scripts/setup-dev.sh

# Start local environment
docker-compose up -d

# Backend development
cd backend
make test
make run

# Frontend development
cd mobile
flutter pub get
flutter run
```

### Docker Compose Services

```yaml
version: '3.8'

services:
  # DynamoDB Local
  dynamodb:
    image: amazon/dynamodb-local
    ports:
      - "8000:8000"
    command: "-jar DynamoDBLocal.jar -sharedDb -dbPath ./data"
    volumes:
      - dynamodb-data:/home/dynamodblocal/data

  # Backend API (Local Lambda)
  api:
    build:
      context: ./backend
      dockerfile: ../infrastructure/docker/backend/Dockerfile
    ports:
      - "3000:3000"
    environment:
      - AWS_REGION=us-east-1
      - DYNAMODB_ENDPOINT=http://dynamodb:8000
      - AUTH0_DOMAIN=${AUTH0_DOMAIN}
      - AUTH0_AUDIENCE=${AUTH0_AUDIENCE}
    depends_on:
      - dynamodb
    volumes:
      - ./backend:/app

  # DynamoDB Admin UI
  dynamodb-admin:
    image: aaronshaf/dynamodb-admin
    ports:
      - "8001:8001"
    environment:
      - DYNAMO_ENDPOINT=http://dynamodb:8000
    depends_on:
      - dynamodb

volumes:
  dynamodb-data:
```

### Git Workflow

```bash
# Feature branch
git checkout -b feature/leaderboard-system

# TDD: Write test first
# Implement feature
# Run tests
make test

# Commit (CalVer auto-tagged)
git commit -m "feat: add global leaderboard endpoint"

# Push and create PR
git push origin feature/leaderboard-system
```

### CalVer Versioning

**Format**: `YYYY.MM.MICRO`

**Examples**:
- `2026.01.0` - January 2026, first release
- `2026.01.1` - January 2026, second release
- `2026.02.0` - February 2026, first release

**Version File**: `/VERSION`

```bash
# scripts/version.sh
#!/bin/bash
YEAR=$(date +%Y)
MONTH=$(date +%m)
MICRO=$(cat VERSION | cut -d. -f3)
NEXT_MICRO=$((MICRO + 1))

echo "${YEAR}.${MONTH}.${NEXT_MICRO}" > VERSION
```

### CI/CD Pipeline

**GitHub Actions Workflow**:

```yaml
# .github/workflows/backend-ci.yml
name: Backend CI

on:
  pull_request:
    paths:
      - 'backend/**'
  push:
    branches:
      - main

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.22'

      - name: Start DynamoDB Local
        run: docker-compose up -d dynamodb

      - name: Run Tests
        run: |
          cd backend
          make test

      - name: Upload Coverage
        uses: codecov/codecov-action@v3

  deploy:
    needs: test
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    steps:
      - name: Deploy to AWS
        run: |
          cd backend
          serverless deploy --stage prod
```

---

## Deployment Strategy

### Environments

| Environment | Backend Stage | Frontend | Purpose |
|-------------|---------------|----------|---------|
| Local | N/A | Emulator | Development |
| Dev | dev | TestFlight/Internal | Integration testing |
| Staging | staging | Beta | QA and user testing |
| Production | prod | App Stores | Live users |

### Backend Deployment

```bash
# Deploy to staging
cd backend
serverless deploy --stage staging

# Deploy to production
serverless deploy --stage prod

# Rollback if needed
serverless rollback --stage prod --timestamp {timestamp}
```

### Frontend Deployment

```bash
# Build for iOS
cd mobile
flutter build ios --release

# Build for Android
flutter build appbundle --release

# Deploy to TestFlight (iOS)
fastlane beta

# Deploy to Play Console (Android)
fastlane deploy_internal
```

### Monitoring & Observability

- **Logging**: CloudWatch Logs
- **Metrics**: CloudWatch Metrics, custom dashboards
- **Tracing**: AWS X-Ray for distributed tracing
- **Alerts**: CloudWatch Alarms + SNS notifications
- **Error Tracking**: Sentry (mobile), CloudWatch Insights (backend)

---

## Next Steps

1. **Phase 1 - Foundation** (Weeks 1-4)
   - ✅ Repository structure
   - Set up Docker Compose
   - Auth0 integration
   - Basic user CRUD

2. **Phase 2 - Core Features** (Weeks 5-8)
   - Lesson CRUD and delivery
   - Progress tracking
   - XP and leveling system
   - Streak tracking

3. **Phase 3 - Gamification** (Weeks 9-12)
   - Achievements system
   - Leaderboards
   - Social features

4. **Phase 4 - Polish** (Weeks 13-16)
   - UI/UX refinement
   - Performance optimization
   - Beta testing
   - App store submission

---

## References

- [DynamoDB Single-Table Design](https://aws.amazon.com/blogs/compute/creating-a-single-table-design-with-amazon-dynamodb/)
- [Auth0 Go SDK](https://github.com/auth0/go-auth0)
- [Flutter Clean Architecture](https://resocoder.com/flutter-clean-architecture-tdd/)
- [Serverless Framework Docs](https://www.serverless.com/framework/docs)
- [Godog (Cucumber for Go)](https://github.com/cucumber/godog)

---

**Document Version**: 2026.01.0
**Last Updated**: 2026-01-11
**Maintained By**: Lumbrera Development Team
