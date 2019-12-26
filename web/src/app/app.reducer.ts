import { Action, createReducer, on } from '@ngrx/store';
import update from 'immutability-helper';
import { isFetching, isNotFetching, setAuthentication, toggleEditorSidebar, toggleSidebar } from './app.actions';

const initialState: AppState = {
  isFetching: false,
  sidebar: {
    collapsed: true,
  },
  editor: {
    sidebar: {
      collapsed: true,
    },
  },
};

const appReducer = createReducer(
  initialState,
  on(isFetching, (state) => update<AppState>(state, { isFetching: { $set: true } })),
  on(isNotFetching, (state) => update<AppState>(state, { isFetching: { $set: false } })),
  on(setAuthentication, (state, { accessToken, idToken, userInfo }) => update<AppState>(state, {
    auth: {
      $set: {
        accessToken,
        idToken,
        userInfo,
      },
    },
  })),
  on(toggleSidebar, (state) => update<AppState>(state, { sidebar: { $toggle: ['collapsed'] } })),
  on(toggleEditorSidebar, (state) => update<AppState>(state, { editor: { sidebar: { $toggle: ['collapsed'] } } })),
);

export function reducer(state: AppState | undefined, action: Action) {
  return appReducer(state, action);
}
