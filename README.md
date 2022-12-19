# Gossip

## Inspirations

- Reddit

## Technology stack

- React
- Ruby on rails
- Postgres

## User requirements

### Types of users

- Visitor
- Poster
- Admin

### User stories

#### General user

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

#### Poster

As a poster i can...

- _in addition to general user stories..._
- Start a forum post
- Tag my post to fit into a category
- Reply to comments on my post
- Edit comments I have made
- View a list of the posts I have started
- View a list of the posts I have commented in
- Be notified of changes (e.g.: new comments) in the posts I am involved in
- Delete my own comments

#### Admin

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

## Features to be implemented

Basic features:

- Login using username and password

These are the features derived from the user stories:

- View all posts and comments
- Create new posts and comments
- Delete one's own comments
- Add tags to posts
- thumbs up and down posts and comments

## Design decisions

### Unable to delete your own posts

Posts are things that many people contribute on.
One person being able to delete the contributions of many does not seem right.
Only admin can delete posts.

## Data

### User model

- `user_id`: cuid
- `username`: string
- `email`: string (sending notifications)
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
- `comments`: comment[] (sorted by `comment_score`)
- `tags`: tag[]
- `title`: string
- `body`: string

### Comment model

- `comment_id`: cuid
- `post`: post
- `parent_comment`: comment
- `author`: user
- `time_commented`: datetime
- `comment_score`: integer
- `body`: string
- `replies`: comment[] (sorted by `comment_score`)

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
