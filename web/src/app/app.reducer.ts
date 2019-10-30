import { Action, createReducer, on } from '@ngrx/store';
import update from 'immutability-helper';

import { toggleSidebar, isFetching, isNotFetching } from './app.actions';

const initialState: AppState = {
  isFetching: false,
  sidebar: {
    collapsed: true,
  },
};

const appReducer = createReducer(
  initialState,
  on(toggleSidebar, (state) => update<AppState>(state, { sidebar: { $toggle: ['collapsed'] } })),
  on(isFetching, (state) => update<AppState>(state, { isFetching: { $set: true } })),
  on(isNotFetching, (state) => update<AppState>(state, { isFetching: { $set: false } })),
);

export function reducer(state: AppState | undefined, action: Action) {
  return appReducer(state, action);
}
