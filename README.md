# Gossip 🤫

## Inspirations ✨

- Reddit

## Technology stack ⚙️

- React
- Go API
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

- _In addition to general user stories..._
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

- _In addition to poster stories..._
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

- Account management (login, creation, updating fields)
- CRUD operations for
  - Comments
  - Posts
  - Tags
- Thumbs up and down posts and comments
- Add and remove tags from posts

## Documentation

- [Frontend](./frontend/README.md)
- [Backend](./backend/README.md)
