import { Action, createReducer, on } from '@ngrx/store';
import update from 'immutability-helper';

import {
  isFetching,
  isNotFetching,
  setAuthentication,
  toggleSidebar,
} from './app.actions';

const initialState: AppState = {
  isFetching: false,
  auth: {},
  sidebar: {
    collapsed: true,
  },
};

const appReducer = createReducer(
  initialState,
  on(isFetching, (state) => update<AppState>(state, { isFetching: { $set: true } })),
  on(isNotFetching, (state) => update<AppState>(state, { isFetching: { $set: false } })),
  on(setAuthentication, (state, { accessToken }) => update<AppState>(state, { auth: { accessToken: { $set: accessToken } } })),
  on(toggleSidebar, (state) => update<AppState>(state, { sidebar: { $toggle: ['collapsed'] } })),
);

export function reducer(state: AppState | undefined, action: Action) {
  return appReducer(state, action);
}
