import { createReducer, on } from '@ngrx/store';
import update from 'immutability-helper';

import { fetchRecentPosts, fetchRecentPostsSuccess, fetchRecentPostsError } from './recent-posts.actions';

export class RecentPostsState {
  // List of latest published posts
  posts: Post[];

  constructor() {
    this.posts = [];
  }
}

export const recentPostsReducer = createReducer(new RecentPostsState(),
  on(fetchRecentPostsSuccess, (state, actions) => update(state, { $set: { posts: actions.posts } })),
  on(fetchRecentPostsError, console.error),
);
