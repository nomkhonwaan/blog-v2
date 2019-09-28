import { Action, createReducer, on } from '@ngrx/store';
import update from 'immutability-helper';

import { toggleSidebar } from './app.actions';

const initialState: State = {
  sidebar: {
    collapsed: true,
  },
};

const appReducer = createReducer(
  initialState,
  on(toggleSidebar, (state) => update(state, { sidebar: { $toggle: ['collapsed'] } })),
);

export function reducer(state: State | undefined, action: Action) {
  return appReducer(state, action);
}
