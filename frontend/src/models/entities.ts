export interface User {
  id: string;
  username: string;
  email: string | null;
  posts: Post[];
  comments: Comment[];
  subscribed: Post[];
  is_public: boolean;
}

export interface Post {
  id: string;
  comments: Comment[];
  tags: Tag[];
  post_score: number;
  title: string;
  body: string;
}

export interface Comment {
  id: string;
  user_id: string;
  post_id: string;
  parent_id: string | null;
  replies: Comment[];
  comment_score: number;
  body: string;
}

export interface Tag {
  id: string;
  name: string;
  posts: Post[];
}

export interface Subscription {
  id: string;
  user_id: string;
  post_id: string;
}

export interface Taggable {
  id: string;
  post_id: string;
  tag_id: string;
}
