import { createReducer, on } from '@ngrx/store';
import update from 'immutability-helper';

import { toggleSidebar } from './app.actions';

export class AppState {
  // Application sidebar specific state
  sidebar: {
    collapsed: boolean,
  };

  constructor() {
    this.sidebar = {
      collapsed: true,
    };
  }
}

export const appReducer = createReducer(new AppState(),
  on(toggleSidebar, (state) => update(state, { sidebar: { $toggle: ['collapsed'] } })),
);
