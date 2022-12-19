# Gossip 🤫

## Inspirations ✨

- Reddit

## Technology stack ⚙️

- React
- Ruby on rails API
- Postgres

## User requirements

### Types of users

- Visitor
- Poster
- Admin

### User stories

#### General user 👤

As a general user i can...

- View the entire post and see comments
- Search for posts by title
- Search for posts by tag
- Subscribe to posts and be notified of changes (e.g.: new comments)
- Unsubscribe to posts that I have lost interest in
- View the profile of other posters (given a public profile)
- Change my account details
  - Password
  - Email
  - Account visibility
- Thumbs up and down posts
- Thumbs up and down comments

#### Poster 📝

As a poster i can...

- _in addition to general user stories..._
- Start a forum post
- Tag my post to fit into a category
- Reply to comments on my post
- Edit comments I have made
- View a list of the posts I have started
- View a list of the posts I have commented in
- Be notified of changes (e.g.: new comments) in the posts I am involved in
- Delete my own comments and posts

#### Admin 🤖

As an admin i can...

- _in addition to poster stories..._
- Delete any comment or post
- Add and delete tags available to forum users
- Add and delete irrelevant tags in posts

## Definitions

### Comment

A comment is the basic form of interaction on the forum.
It has a title and a body.
A comment can be replied to by other comments.

### Post

A post is where comments will reside, and a main prompt is given to start the discussion.

### Tag

A tag is a tool for categorisation of posts and increasing discoverability.

## Features

- account management (login, creation, updating fields)
- CRUD operations for
  - comments
  - posts
  - tags
- thumbs up and down posts and comments
- add and remove tags from posts

## Data

### User model

- `user_id`: cuid
- `username`: string
- `email`: string? (sending notifications)
- `password`: hashed string
- `comments`: comment[]
- `isAdmin`: boolean
- `isPublic`: boolean
- `subscribed`: post[]

### Post model

- `post_id`: cuid
- `author`: user
- `time_posted`: datetime
- `post_score`: integer
- `comments`: comment[]
- `tags`: tag[]
- `title`: string
- `body`: text

### Comment model

- `comment_id`: cuid
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
  - user/
    - `[user_id]/`: returns account details of user with id `user_id`
  - post/
    - `[post_id]/`: returns post data object with id `post_id`
  - comment/
    - `[comment_id]/`: returns comment data object with id `comment_id`

## Testing strategy

### Frontend

- Jest

### Backend

- Postman
