# Lumbrera - Gamified Bible Study App

[![Version](https://img.shields.io/badge/version-2026.01.0-blue.svg)](VERSION)
[![Go Version](https://img.shields.io/badge/go-1.22+-00ADD8.svg)](go.mod)
[![Flutter](https://img.shields.io/badge/flutter-3.x-02569B.svg)](mobile/)
[![License](https://img.shields.io/badge/license-Proprietary-red.svg)](LICENSE)

**Lumbrera** is a gamified mobile Bible study application that makes biblical learning engaging and fun, inspired by Duolingo's approach to language learning.

## 📖 Overview

Lumbrera combines structured Bible study with game mechanics to create an engaging learning experience:

- **📚 Bite-sized Lessons**: Daily 5-10 minute lessons organized by books and themes
- **🎮 Gamification**: XP, levels, streaks, achievements, and leaderboards
- **👥 Social Learning**: Friends, challenges, and community competitions
- **📊 Progress Tracking**: Comprehensive analytics and personalized insights
- **🌍 Multi-Translation**: Support for NIV, KJV, ESV, NLT, and more

## 🏗️ Architecture

### Tech Stack

**Frontend (Mobile)**
- Flutter 3.x
- Riverpod/Bloc for state management
- Auth0 for authentication (Google, Facebook)

**Backend**
- Go 1.22+
- AWS Lambda (Serverless)
- DynamoDB (Single-table design)
- API Gateway

**Infrastructure**
- Docker & Docker Compose
- Serverless Framework
- GitHub Actions (CI/CD)

**Testing**
- TDD Approach
- Cucumber/Gherkin (E2E)
- Go testing + Testify
- Flutter Test

### Project Structure

```
lumbrera/
├── backend/                 # Go backend (this repo)
│   ├── functions/          # Lambda functions
│   ├── internal/           # Shared code
│   ├── tests/             # Tests
│   └── docs/              # Documentation
├── mobile/                 # Flutter app (coming soon)
├── infrastructure/         # IaC and Docker
├── scripts/               # Utility scripts
└── docs/                  # Project documentation
```

See [ARCHITECTURE.md](docs/ARCHITECTURE.md) for detailed architecture documentation.

## 🚀 Quick Start

### Prerequisites

- **Go** 1.22 or higher
- **Docker** and Docker Compose
- **Make**
- **AWS CLI** (optional, for deployment)
- **Flutter** 3.x (for mobile development)

### Setup

```bash
# Clone the repository
git clone https://github.com/your-org/lumbrera.git
cd lumbrera

# Run the setup script
./scripts/setup-dev.sh

# The script will:
# - Install dependencies
# - Create necessary directories
# - Start Docker services
# - Initialize DynamoDB tables
# - Build the backend
# - Run tests
```

### Manual Setup

If you prefer manual setup:

```bash
# 1. Copy environment variables
cp .env.example .env
# Edit .env with your Auth0 credentials

# 2. Install Go dependencies
go mod download

# 3. Build backend
make build

# 4. Start Docker services
docker-compose up -d

# 5. Initialize database
docker-compose up init-db

# 6. Run tests
make test
```

## 💻 Development

### Available Commands

```bash
# Backend Development
make build          # Build Lambda functions
make test           # Run all tests
make test-unit      # Run unit tests only
make test-e2e       # Run E2E tests
make run            # Start local environment
make stop           # Stop local environment
make clean          # Clean build artifacts
make deploy-dev     # Deploy to dev environment
make deploy-prod    # Deploy to production

# Versioning (CalVer: YYYY.MM.MICRO)
./scripts/version.sh        # Show current version
./scripts/version.sh bump   # Bump micro version
./scripts/version.sh set 2026.01.5  # Set specific version

# Docker
docker-compose up -d                    # Start all services
docker-compose logs -f                  # View logs
docker-compose down                     # Stop all services
docker-compose -f docker-compose.test.yml up  # Run tests in Docker
```

### Local Services

After running `docker-compose up -d`, the following services are available:

| Service | URL | Description |
|---------|-----|-------------|
| DynamoDB | http://localhost:8000 | Local DynamoDB instance |
| DynamoDB Admin | http://localhost:8001 | Web UI for DynamoDB |
| Create Lesson API | http://localhost:8080 | Create lesson endpoint |
| Get Lesson API | http://localhost:8081 | Get lesson endpoint |

### Testing Locally

```bash
# Create a lesson
curl -X POST http://localhost:8080/2015-03-31/functions/function/invocations \
  -H "Content-Type: application/json" \
  -d '{
    "httpMethod": "POST",
    "body": "{\"title\":\"Genesis 1:1\",\"book\":\"Genesis\",\"chapter\":1}"
  }'

# Get a lesson
curl -X GET http://localhost:8081/2015-03-31/functions/function/invocations \
  -H "Content-Type: application/json" \
  -d '{
    "httpMethod": "GET",
    "pathParameters": {"id": "lesson-id-here"}
  }'
```

## 🧪 Testing

We follow **Test-Driven Development (TDD)** practices.

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
go test -v -cover ./...

# Run specific test
go test -v ./functions/lessons/create/...

# Run E2E tests
make test-e2e

# Run tests in Docker
docker-compose -f docker-compose.test.yml up --abort-on-container-exit
```

### Writing Tests

```go
// Unit Test Example
func TestCreateLesson(t *testing.T) {
    db := database.GetMockedClient()

    lesson := &models.Lesson{
        Title: "Genesis 1:1",
        Book:  "Genesis",
    }

    err := repository.SaveLesson(db, lesson)
    assert.NoError(t, err)
    assert.NotEmpty(t, lesson.ID)
}
```

### E2E Tests (Cucumber)

```gherkin
# tests/e2e/features/lessons.feature
Feature: Lesson Management
  Scenario: Create and retrieve a lesson
    Given I am authenticated
    When I create a lesson with title "Genesis 1:1"
    Then the lesson should be created successfully
    When I retrieve the lesson
    Then I should see the lesson details
```

## 📚 Documentation

- **[Architecture](docs/ARCHITECTURE.md)** - System architecture and design
- **[API Documentation](docs/openapi.yaml)** - OpenAPI 3.0 specification
- **[Database Schema](docs/DATABASE.md)** - DynamoDB schema and patterns
- **[Testing Strategy](docs/TESTING.md)** - Testing approach and guidelines
- **[Deployment Guide](docs/DEPLOYMENT.md)** - Deployment instructions

## 🔐 Authentication

Lumbrera uses **Auth0** for authentication with support for:

- **Google OAuth 2.0**
- **Facebook OAuth 2.0**

### Setup Auth0

1. Create an Auth0 account at https://auth0.com
2. Create a new application (Native for mobile)
3. Enable Google and Facebook connections
4. Update `.env` with your Auth0 credentials:

```bash
AUTH0_DOMAIN=your-tenant.auth0.com
AUTH0_CLIENT_ID=your-client-id
AUTH0_CLIENT_SECRET=your-client-secret
AUTH0_AUDIENCE=https://api.lumbrera.app
```

## 🎮 Features

### Gamification Elements

- **XP System**: Earn experience points for completing lessons
- **Levels**: Progress through 100 levels
- **Streaks**: Maintain daily study streaks
- **Achievements**: Unlock badges and rewards
- **Leaderboards**: Compete globally, locally, or with friends

### Lesson Types

- **Reading Comprehension**: Read and understand passages
- **Quizzes**: Multiple choice, true/false
- **Memory Verses**: Memorize scripture
- **Fill in the Blank**: Complete missing words
- **Matching**: Match concepts and people
- **Timeline**: Arrange events chronologically

### Social Features

- **Friends**: Connect with other learners
- **Challenges**: Compete in head-to-head challenges
- **Groups**: Join study groups
- **Share**: Share achievements on social media

## 📊 Database

Lumbrera uses **DynamoDB** with a single-table design for optimal performance.

### Key Entities

- **Users**: User profiles and preferences
- **Lessons**: Bible study content
- **Progress**: User progress tracking
- **Achievements**: Unlockable achievements
- **Leaderboards**: Rankings and scores
- **Streaks**: Daily activity tracking

See [DATABASE.md](docs/DATABASE.md) for detailed schema documentation.

## 🚢 Deployment

### Deploy to AWS

```bash
# Deploy to development
make deploy-dev

# Deploy to production
make deploy-prod

# Or use Serverless directly
cd backend
serverless deploy --stage prod
```

### Environment Stages

- **local**: Local development
- **dev**: Development environment
- **staging**: Staging environment
- **prod**: Production environment

### CI/CD

GitHub Actions automatically:
- Runs tests on pull requests
- Deploys to dev on merge to `develop`
- Deploys to staging on merge to `staging`
- Deploys to prod on merge to `main`

## 📝 Versioning

Lumbrera uses **Calendar Versioning (CalVer)** with the format `YYYY.MM.MICRO`.

Examples:
- `2026.01.0` - January 2026, first release
- `2026.01.1` - January 2026, second release
- `2026.02.0` - February 2026, first release

```bash
# View current version
./scripts/version.sh

# Bump version
./scripts/version.sh bump

# Set specific version
./scripts/version.sh set 2026.01.5
```

The version script automatically updates:
- `VERSION` file
- `package.json`
- `mobile/pubspec.yaml`
- `docs/ARCHITECTURE.md`
- `docs/openapi.yaml`

## 🤝 Contributing

We follow these practices:

1. **TDD**: Write tests before implementation
2. **Branch naming**: `feature/`, `bugfix/`, `hotfix/`
3. **Commit messages**: Conventional commits format
4. **Code review**: All PRs require review
5. **Documentation**: Update docs with code changes

### Git Workflow

```bash
# Create feature branch
git checkout -b feature/leaderboard-system

# Write tests (RED)
# Implement feature (GREEN)
# Refactor (REFACTOR)

# Commit changes
git commit -m "feat: add global leaderboard endpoint"

# Push and create PR
git push origin feature/leaderboard-system
```

## 📄 License

Proprietary - All rights reserved

## 🙋 Support

- **Documentation**: See `docs/` directory
- **Issues**: Create a GitHub issue
- **Questions**: Contact the development team

## 🗺️ Roadmap

### Phase 1 - Foundation (Weeks 1-4)
- [x] Repository structure
- [ ] Docker Compose setup
- [ ] Auth0 integration
- [ ] Basic user CRUD

### Phase 2 - Core Features (Weeks 5-8)
- [ ] Lesson CRUD and delivery
- [ ] Progress tracking
- [ ] XP and leveling system
- [ ] Streak tracking

### Phase 3 - Gamification (Weeks 9-12)
- [ ] Achievements system
- [ ] Leaderboards
- [ ] Social features

### Phase 4 - Polish (Weeks 13-16)
- [ ] UI/UX refinement
- [ ] Performance optimization
- [ ] Beta testing
- [ ] App store submission

## 🎯 Project Goals

1. Make Bible study engaging and accessible
2. Foster consistent daily study habits
3. Build a supportive learning community
4. Provide comprehensive biblical education
5. Track and celebrate spiritual growth

---

**Built with ❤️ for biblical learning**

**Version**: 2026.01.0 | **Last Updated**: 2026-01-11
