import { createAction, props } from '@ngrx/store';
import { Action } from 'rxjs/internal/scheduler/Action';

/**
 * Fetches latest published posts
 */
export const fetchRecentPosts = createAction('FETCH_RECENT_POSTS');

/**
 * Fetches latest published success
 */
export const fetchRecentPostsSuccess = createAction('FETCH_RECENT_POSTS_SUCCESS', props<{ posts: Post[] }>());

/**
 * Fetches latest published error
 */
export const fetchRecentPostsError = createAction('FETCH_RECENT_POSTS_ERROR');
