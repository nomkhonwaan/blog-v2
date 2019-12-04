import { Action, createReducer, on } from '@ngrx/store';
import update from 'immutability-helper';

import { isLightboxOpened, isLightboxClosed } from './post.actions';

const initialState: PostState = {
  content: {
    fslightbox: {
      closed: true,
    },
  },
};

const postReducer = createReducer(
  initialState,
  on(isLightboxOpened, (state) => update<PostState>(state, { content: { fslightbox: { closed: { $set: false } } } })),
  on(isLightboxClosed, (state) => update<PostState>(state, { content: { fslightbox: { closed: { $set: true } } } })),
);

export function reducer(state: PostState | undefined, action: Action) {
  return postReducer(state, action);
}
