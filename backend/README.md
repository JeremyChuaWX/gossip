# Backend

## Technologies

- gofiber
- gorm

## Data

### User model

- `user_id`: uuid
- `username`: string
- `email`: string? (sending notifications)
- `password`: hashed string
- `comments`: comment[]
- `posts`: post[]
- `role`: role enum
- `isPublic`: boolean
- `subscribed`: post[]

### Post model

- `post_id`: uuid
- `author`: user
- `time_posted`: datetime
- `post_score`: integer
- `comments`: comment[]
- `tags`: tag[]
- `title`: string
- `body`: text

### Comment model

- `comment_id`: uuid
- `post`: post
- `parent_comment`: comment?
- `author`: user
- `time_commented`: datetime
- `comment_score`: integer
- `body`: text
- `replies`: comment[]

### Tag model

- `tag_id`: uuid
- `name`: string

## API endpoints

- api/
  - auth/
    - signin (POST): Sign in user with provided credentials
    - signup (POST): Sign up a new user with provided credentials
    - signout (GET): Sign out currently signed in user
    - refresh (GET): Refresh access tokens using provided refresh token
  - users/
    - get-me (GET): Get currently signed in user information
    - get-user/:id (GET): Get user by id
    - update-me (PUT): Update currently signed in user with provided information
    - toggle-visibility (GET): Toggle visibility of currently signed in user
    - delete-me (DELETE): Delete currently signed in user
  - posts/
    - /create-post (POST): Create new post with provided information
    - /get-posts (GET): Get all posts
    - /get-post/\[id\] (GET): Get post by id
    - /update-post/\[id\] (PUT): Update post by id with provided information
    - /update-postscore/\[id\] (PUT): Update post scores by id with provided information
    - /delete-post/\[id\] (DELETE): Delete post by id
  - comments/
    - /create-comment (POST): Create new comment with provided information
    - /get-comment/\[id\] (GET): Get comment by id
    - /update-comment/\[id\] (PUT): Update comment by id with provided information
    - /update-commentscore/\[id\] (PUT): Update comment scores by id with provided information
    - /delete-comment/\[id\] (DELETE): Delete comment by id
  - tags/
  - subscriptions/
  - taggable/
