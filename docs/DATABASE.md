# Database Schema Documentation

## Overview

Lumbrera uses **AWS DynamoDB** with a **single-table design** pattern for optimal performance, scalability, and cost efficiency.

**Table Name**: `lumbrera-{stage}`
- **Local**: `lumbrera-local`
- **Dev**: `lumbrera-dev`
- **Staging**: `lumbrera-staging`
- **Production**: `lumbrera-prod`

## Table Configuration

### Primary Index

| Attribute | Type | Key Type |
|-----------|------|----------|
| PK | String | Partition Key (HASH) |
| SK | String | Sort Key (RANGE) |

### Global Secondary Indexes

#### GSI1 - User-Centric Queries
| Attribute | Type | Key Type |
|-----------|------|----------|
| GSI1PK | String | Partition Key |
| GSI1SK | String | Sort Key |

**Use Cases**:
- Query all items belonging to a user
- Get user's progress across all lessons
- Retrieve user's achievements

#### GSI2 - Leaderboard and Rankings
| Attribute | Type | Key Type |
|-----------|------|----------|
| GSI2PK | String | Partition Key |
| GSI2SK | String | Sort Key |

**Use Cases**:
- Query leaderboard entries by type and period
- Get user's rank in different leaderboards
- Retrieve top performers

## Access Patterns

### User Management

| Pattern | Index | PK | SK |
|---------|-------|----|----|
| Get user profile | Primary | `USER#{UserID}` | `PROFILE` |
| Update user | Primary | `USER#{UserID}` | `PROFILE` |
| Get user by Auth0 ID | GSI1 | `AUTH0#{Auth0ID}` | `PROFILE` |

### Lessons

| Pattern | Index | PK | SK |
|---------|-------|----|----|
| Get lesson | Primary | `LESSON#{LessonID}` | `METADATA` |
| List lessons by book | Primary | `BOOK#{BookName}` | `LESSON#{Order}` |
| Get recommended lessons | GSI1 | `DIFFICULTY#{Level}` | `ORDER#{Order}` |

### Progress Tracking

| Pattern | Index | PK | SK |
|---------|-------|----|----|
| Get user's progress for lesson | Primary | `USER#{UserID}` | `PROGRESS#{LessonID}` |
| Get all progress for user | Primary | `USER#{UserID}` | begins_with(`PROGRESS#`) |
| Get all users who completed lesson | GSI1 | `LESSON#{LessonID}` | `USER#{UserID}` |
| Get user's daily streak | Primary | `USER#{UserID}` | `STREAK#{Date}` |

### Achievements

| Pattern | Index | PK | SK |
|---------|-------|----|----|
| Get achievement details | Primary | `ACHIEVEMENT#{AchievementID}` | `METADATA` |
| List all achievements | Primary | `ACHIEVEMENT#ALL` | `ACHIEVEMENT#{AchievementID}` |
| Get user's achievements | Primary | `USER#{UserID}` | begins_with(`ACHIEVEMENT#`) |
| Get achievement progress | Primary | `USER#{UserID}` | `ACHIEVEMENT#{AchievementID}` |

### Leaderboards

| Pattern | Index | PK | SK |
|---------|-------|----|----|
| Get global leaderboard | Primary | `LEADERBOARD#GLOBAL#{Period}` | `SCORE#{Score}#{UserID}` |
| Get friends leaderboard | Primary | `LEADERBOARD#FRIENDS#{UserID}#{Period}` | `SCORE#{Score}#{FriendID}` |
| Get local leaderboard | Primary | `LEADERBOARD#LOCAL#{Region}#{Period}` | `SCORE#{Score}#{UserID}` |
| Get user's rank | GSI2 | `USER#{UserID}` | `LEADERBOARD#{Type}#{Period}` |

### Social Features

| Pattern | Index | PK | SK |
|---------|-------|----|----|
| Get user's friends | Primary | `USER#{UserID}` | begins_with(`FRIEND#`) |
| Check friendship | Primary | `USER#{UserID}` | `FRIEND#{FriendID}` |
| Get friend requests | Primary | `USER#{UserID}` | begins_with(`FRIEND_REQUEST#`) |

---

## Entity Details

### 1. User Entity

**Primary Index**:
- **PK**: `USER#{UserID}`
- **SK**: `PROFILE`

**GSI1**:
- **GSI1PK**: `AUTH0#{Auth0ID}`
- **GSI1SK**: `PROFILE`

**Attributes**:
```go
{
  "PK": "USER#user_123",
  "SK": "PROFILE",
  "GSI1PK": "AUTH0#auth0|abc123",
  "GSI1SK": "PROFILE",
  "EntityType": "USER",

  "UserID": "user_123",
  "Auth0ID": "auth0|abc123",
  "Email": "john@example.com",
  "Name": "John Doe",
  "Avatar": "https://cdn.lumbrera.app/avatars/user_123.jpg",

  "XP": 1250,
  "Level": 8,
  "Streak": 15,

  "Translation": "NIV",
  "DailyGoal": 50,
  "NotificationsOn": true,

  "CreatedAt": "2026-01-01T00:00:00Z",
  "UpdatedAt": "2026-01-11T10:30:00Z"
}
```

---

### 2. Lesson Entity

**Primary Index**:
- **PK**: `LESSON#{LessonID}`
- **SK**: `METADATA`

**Attributes**:
```go
{
  "PK": "LESSON#lesson_001",
  "SK": "METADATA",
  "EntityType": "LESSON",

  "LessonID": "lesson_001",
  "Title": "In the Beginning",
  "Description": "Explore the creation story in Genesis",
  "Book": "Genesis",
  "Chapter": 1,
  "Verses": "1-10",

  "ContentType": "quiz",
  "Content": {
    "questions": [
      {
        "id": "q1",
        "type": "multiple_choice",
        "question": "On which day did God create light?",
        "options": ["Day 1", "Day 2", "Day 3", "Day 4"],
        "correct_answer": "Day 1"
      }
    ]
  },

  "Difficulty": "beginner",
  "XPReward": 50,
  "EstDuration": 10,
  "Order": 1,
  "Prerequisites": [],

  "CreatedAt": "2026-01-01T00:00:00Z",
  "UpdatedAt": "2026-01-05T12:00:00Z"
}
```

**Secondary Pattern - By Book**:
- **PK**: `BOOK#Genesis`
- **SK**: `LESSON#001`

---

### 3. UserProgress Entity

**Primary Index**:
- **PK**: `USER#{UserID}`
- **SK**: `PROGRESS#{LessonID}`

**GSI1** (for querying by lesson):
- **GSI1PK**: `LESSON#{LessonID}`
- **GSI1SK**: `USER#{UserID}#SCORE#{Score}`

**Attributes**:
```go
{
  "PK": "USER#user_123",
  "SK": "PROGRESS#lesson_001",
  "GSI1PK": "LESSON#lesson_001",
  "GSI1SK": "USER#user_123#SCORE#095",
  "EntityType": "PROGRESS",

  "UserID": "user_123",
  "LessonID": "lesson_001",

  "Status": "completed",
  "Score": 95,
  "Attempts": 1,
  "XPEarned": 55,

  "StartedAt": "2026-01-11T10:00:00Z",
  "CompletedAt": "2026-01-11T10:10:00Z",
  "UpdatedAt": "2026-01-11T10:10:00Z"
}
```

---

### 4. DailyStreak Entity

**Primary Index**:
- **PK**: `USER#{UserID}`
- **SK**: `STREAK#{Date}`

**Attributes**:
```go
{
  "PK": "USER#user_123",
  "SK": "STREAK#2026-01-11",
  "EntityType": "DAILY_STREAK",

  "UserID": "user_123",
  "Date": "2026-01-11",
  "XPEarned": 65,
  "LessonsCount": 2,
  "GoalMet": true,

  "CreatedAt": "2026-01-11T10:10:00Z"
}
```

**Query Patterns**:
- Get last 30 days: Query PK=`USER#user_123` and SK between `STREAK#2025-12-12` and `STREAK#2026-01-11`
- Check today's activity: GetItem with PK=`USER#user_123` and SK=`STREAK#2026-01-11`

---

### 5. Achievement Entity

**Primary Index**:
- **PK**: `ACHIEVEMENT#{AchievementID}`
- **SK**: `METADATA`

**Attributes**:
```go
{
  "PK": "ACHIEVEMENT#ach_001",
  "SK": "METADATA",
  "EntityType": "ACHIEVEMENT",

  "AchievementID": "ach_001",
  "Name": "Scholar",
  "Description": "Complete 100 lessons",
  "Icon": "https://cdn.lumbrera.app/icons/scholar.png",

  "Type": "lessons_completed",
  "Requirement": 100,

  "Rarity": "epic",
  "XPReward": 500,

  "CreatedAt": "2026-01-01T00:00:00Z"
}
```

**Secondary Pattern - All Achievements List**:
- **PK**: `ACHIEVEMENT#ALL`
- **SK**: `ACHIEVEMENT#{AchievementID}`

---

### 6. UserAchievement Entity

**Primary Index**:
- **PK**: `USER#{UserID}`
- **SK**: `ACHIEVEMENT#{AchievementID}`

**Attributes**:
```go
{
  "PK": "USER#user_123",
  "SK": "ACHIEVEMENT#ach_001",
  "EntityType": "USER_ACHIEVEMENT",

  "UserID": "user_123",
  "AchievementID": "ach_001",

  "Progress": 45,
  "Requirement": 100,
  "UnlockedAt": null  // null if not unlocked
}
```

**After unlocking**:
```go
{
  "Progress": 100,
  "UnlockedAt": "2026-01-11T10:10:00Z"
}
```

---

### 7. LeaderboardEntry Entity

**Primary Index**:
- **PK**: `LEADERBOARD#{Type}#{Period}`
- **SK**: `SCORE#{PaddedScore}#{UserID}`

**GSI2** (for user rank lookup):
- **GSI2PK**: `USER#{UserID}`
- **GSI2SK**: `LEADERBOARD#{Type}#{Period}`

**Attributes**:
```go
{
  "PK": "LEADERBOARD#GLOBAL#WEEKLY",
  "SK": "SCORE#0000001250#user_123",
  "GSI2PK": "USER#user_123",
  "GSI2SK": "LEADERBOARD#GLOBAL#WEEKLY",
  "EntityType": "LEADERBOARD_ENTRY",

  "UserID": "user_123",
  "UserName": "John Doe",
  "Avatar": "https://cdn.lumbrera.app/avatars/user_123.jpg",

  "Type": "global",
  "Period": "weekly",
  "Score": 1250,
  "Rank": 42,

  "UpdatedAt": "2026-01-11T10:10:00Z"
}
```

**Score Padding**: Scores are zero-padded to 10 digits for proper sorting:
- 1250 → `0000001250`
- 9999999 → `0009999999`

**Leaderboard Types**:
- `LEADERBOARD#GLOBAL#WEEKLY`
- `LEADERBOARD#GLOBAL#MONTHLY`
- `LEADERBOARD#GLOBAL#ALL_TIME`
- `LEADERBOARD#FRIENDS#{UserID}#WEEKLY`
- `LEADERBOARD#LOCAL#{Region}#WEEKLY`

---

### 8. Friend Entity

**Primary Index**:
- **PK**: `USER#{UserID}`
- **SK**: `FRIEND#{FriendID}`

**GSI1** (bidirectional friendship):
- **GSI1PK**: `USER#{FriendID}`
- **GSI1SK**: `FRIEND#{UserID}`

**Attributes**:
```go
{
  "PK": "USER#user_123",
  "SK": "FRIEND#user_456",
  "GSI1PK": "USER#user_456",
  "GSI1SK": "FRIEND#user_123",
  "EntityType": "FRIEND",

  "UserID": "user_123",
  "FriendID": "user_456",
  "FriendName": "Jane Smith",
  "Avatar": "https://cdn.lumbrera.app/avatars/user_456.jpg",

  "Status": "accepted",  // pending, accepted, blocked
  "AddedAt": "2026-01-05T12:00:00Z"
}
```

---

## Table Creation Script

```bash
#!/bin/bash
# scripts/init-dynamodb.sh

STAGE=${1:-local}
TABLE_NAME="lumbrera-${STAGE}"

if [ "$STAGE" = "local" ]; then
  ENDPOINT="--endpoint-url http://localhost:8000"
else
  ENDPOINT=""
fi

# Create table
aws dynamodb create-table $ENDPOINT \
  --table-name $TABLE_NAME \
  --attribute-definitions \
    AttributeName=PK,AttributeType=S \
    AttributeName=SK,AttributeType=S \
    AttributeName=GSI1PK,AttributeType=S \
    AttributeName=GSI1SK,AttributeType=S \
    AttributeName=GSI2PK,AttributeType=S \
    AttributeName=GSI2SK,AttributeType=S \
  --key-schema \
    AttributeName=PK,KeyType=HASH \
    AttributeName=SK,KeyType=RANGE \
  --global-secondary-indexes \
    "[
      {
        \"IndexName\": \"GSI1\",
        \"KeySchema\": [
          {\"AttributeName\":\"GSI1PK\",\"KeyType\":\"HASH\"},
          {\"AttributeName\":\"GSI1SK\",\"KeyType\":\"RANGE\"}
        ],
        \"Projection\": {\"ProjectionType\":\"ALL\"},
        \"ProvisionedThroughput\": {
          \"ReadCapacityUnits\": 5,
          \"WriteCapacityUnits\": 5
        }
      },
      {
        \"IndexName\": \"GSI2\",
        \"KeySchema\": [
          {\"AttributeName\":\"GSI2PK\",\"KeyType\":\"HASH\"},
          {\"AttributeName\":\"GSI2SK\",\"KeyType\":\"RANGE\"}
        ],
        \"Projection\": {\"ProjectionType\":\"ALL\"},
        \"ProvisionedThroughput\": {
          \"ReadCapacityUnits\": 5,
          \"WriteCapacityUnits\": 5
        }
      }
    ]" \
  --provisioned-throughput \
    ReadCapacityUnits=5,WriteCapacityUnits=5

echo "Table $TABLE_NAME created successfully"
```

---

## Querying Examples

### 1. Get User's Completed Lessons

```go
result, err := client.Query(context.TODO(), &dynamodb.QueryInput{
    TableName:              aws.String("lumbrera-prod"),
    KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk_prefix)"),
    ExpressionAttributeValues: map[string]types.AttributeValue{
        ":pk":        &types.AttributeValueMemberS{Value: "USER#user_123"},
        ":sk_prefix": &types.AttributeValueMemberS{Value: "PROGRESS#"},
    },
    FilterExpression: aws.String("#status = :completed"),
    ExpressionAttributeNames: map[string]string{
        "#status": "Status",
    },
    ExpressionAttributeValues: map[string]types.AttributeValue{
        ":completed": &types.AttributeValueMemberS{Value: "completed"},
    },
})
```

### 2. Get Top 100 Global Leaderboard

```go
result, err := client.Query(context.TODO(), &dynamodb.QueryInput{
    TableName:              aws.String("lumbrera-prod"),
    KeyConditionExpression: aws.String("PK = :pk"),
    ExpressionAttributeValues: map[string]types.AttributeValue{
        ":pk": &types.AttributeValueMemberS{Value: "LEADERBOARD#GLOBAL#WEEKLY"},
    },
    ScanIndexForward: aws.Bool(false),  // Sort descending by score
    Limit:            aws.Int32(100),
})
```

### 3. Get User's Current Streak

```go
// Get last 30 days of activity
endDate := time.Now().Format("2006-01-02")
startDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")

result, err := client.Query(context.TODO(), &dynamodb.QueryInput{
    TableName:              aws.String("lumbrera-prod"),
    KeyConditionExpression: aws.String("PK = :pk AND SK BETWEEN :start AND :end"),
    ExpressionAttributeValues: map[string]types.AttributeValue{
        ":pk":    &types.AttributeValueMemberS{Value: "USER#user_123"},
        ":start": &types.AttributeValueMemberS{Value: fmt.Sprintf("STREAK#%s", startDate)},
        ":end":   &types.AttributeValueMemberS{Value: fmt.Sprintf("STREAK#%s", endDate)},
    },
    ScanIndexForward: aws.Bool(false),  // Most recent first
})

// Calculate streak from consecutive days
```

### 4. Get User's Achievement Progress

```go
result, err := client.Query(context.TODO(), &dynamodb.QueryInput{
    TableName:              aws.String("lumbrera-prod"),
    KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk_prefix)"),
    ExpressionAttributeValues: map[string]types.AttributeValue{
        ":pk":        &types.AttributeValueMemberS{Value: "USER#user_123"},
        ":sk_prefix": &types.AttributeValueMemberS{Value: "ACHIEVEMENT#"},
    },
})
```

---

## Migration Strategy

### Adding New Attributes

DynamoDB is schemaless, so adding new attributes is straightforward:

1. Update Go models
2. Update write operations to include new attributes
3. Old items will return null/empty for new attributes
4. Use `omitempty` in struct tags for backward compatibility

### Changing Access Patterns

1. Create new GSI if needed
2. Backfill existing data with new keys
3. Update application code to use new pattern
4. Deprecate old GSI (optional)

### Data Backfill Script

```go
// scripts/backfill/add_gsi_keys.go
func backfillGSI1Keys() error {
    // Scan all USER items
    result, err := client.Scan(context.TODO(), &dynamodb.ScanInput{
        TableName:        aws.String("lumbrera-prod"),
        FilterExpression: aws.String("EntityType = :type"),
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":type": &types.AttributeValueMemberS{Value: "USER"},
        },
    })

    // Update each item with GSI1 keys
    for _, item := range result.Items {
        userID := item["UserID"].(*types.AttributeValueMemberS).Value
        auth0ID := item["Auth0ID"].(*types.AttributeValueMemberS).Value

        _, err := client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
            TableName: aws.String("lumbrera-prod"),
            Key: map[string]types.AttributeValue{
                "PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID)},
                "SK": &types.AttributeValueMemberS{Value: "PROFILE"},
            },
            UpdateExpression: aws.String("SET GSI1PK = :gsi1pk, GSI1SK = :gsi1sk"),
            ExpressionAttributeValues: map[string]types.AttributeValue{
                ":gsi1pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("AUTH0#%s", auth0ID)},
                ":gsi1sk": &types.AttributeValueMemberS{Value: "PROFILE"},
            },
        })
    }

    return nil
}
```

---

## Best Practices

### 1. Use Consistent Key Patterns

Always use the same delimiter and naming:
- ✅ `USER#user_123`
- ❌ `user:user_123`
- ❌ `user_123`

### 2. Pad Numeric Values for Sorting

When using numbers in sort keys:
```go
// Bad: Won't sort correctly
SK: fmt.Sprintf("SCORE#%d", score)  // "SCORE#10" comes before "SCORE#2"

// Good: Pads to fixed width
SK: fmt.Sprintf("SCORE#%010d", score)  // "SCORE#0000000010"
```

### 3. Use Sparse Indexes

Not all items need GSI keys. Only add GSI keys when needed:
```go
// Only leaderboard entries need GSI2
if entityType == "LEADERBOARD_ENTRY" {
    item["GSI2PK"] = fmt.Sprintf("USER#%s", userID)
    item["GSI2SK"] = fmt.Sprintf("LEADERBOARD#%s#%s", leaderboardType, period)
}
```

### 4. Batch Operations

Use BatchGetItem and BatchWriteItem for efficiency:
```go
// Get multiple lessons at once
keys := []map[string]types.AttributeValue{
    {
        "PK": &types.AttributeValueMemberS{Value: "LESSON#lesson_001"},
        "SK": &types.AttributeValueMemberS{Value: "METADATA"},
    },
    {
        "PK": &types.AttributeValueMemberS{Value: "LESSON#lesson_002"},
        "SK": &types.AttributeValueMemberS{Value: "METADATA"},
    },
}

result, err := client.BatchGetItem(context.TODO(), &dynamodb.BatchGetItemInput{
    RequestItems: map[string]types.KeysAndAttributes{
        "lumbrera-prod": {Keys: keys},
    },
})
```

### 5. Conditional Writes

Use conditions to prevent race conditions:
```go
// Only update if version matches (optimistic locking)
_, err := client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
    TableName: aws.String("lumbrera-prod"),
    Key:       key,
    UpdateExpression: aws.String("SET XP = XP + :xp, #version = #version + :one"),
    ConditionExpression: aws.String("#version = :expected_version"),
    ExpressionAttributeNames: map[string]string{
        "#version": "Version",
    },
    ExpressionAttributeValues: map[string]types.AttributeValue{
        ":xp":               &types.AttributeValueMemberN{Value: "50"},
        ":one":              &types.AttributeValueMemberN{Value: "1"},
        ":expected_version": &types.AttributeValueMemberN{Value: "42"},
    },
})
```

---

## Capacity Planning

### Read/Write Capacity Units

**Development**: 5 RCU / 5 WCU (provisioned)

**Production**: On-Demand pricing recommended
- No capacity planning needed
- Automatically scales
- Pay per request

### Cost Estimation (Monthly)

**Assumptions**:
- 10,000 active users
- 50 requests per user per day
- Average item size: 1KB

**On-Demand Pricing**:
- Read requests: 10,000 × 50 × 30 = 15M requests/month
- Write requests: 10,000 × 10 × 30 = 3M requests/month
- Storage: ~50GB

**Estimated Cost**: ~$100-200/month

---

## Testing

### Local DynamoDB

```bash
# Start local DynamoDB
docker run -p 8000:8000 amazon/dynamodb-local

# Create tables
./scripts/init-dynamodb.sh local

# Run tests
AWS_ENDPOINT=http://localhost:8000 go test ./...
```

### Test Data

```go
// tests/fixtures/users.go
func CreateTestUser(db *dynamodb.Client, userID string) *models.User {
    user := &models.User{
        UserID:  userID,
        Auth0ID: fmt.Sprintf("auth0|%s", userID),
        Email:   fmt.Sprintf("%s@test.com", userID),
        Name:    "Test User",
        XP:      100,
        Level:   1,
        Streak:  0,
    }

    // Save to DynamoDB
    repository.SaveUser(db, user)

    return user
}
```

---

**Document Version**: 2026.01.0
**Last Updated**: 2026-01-11
