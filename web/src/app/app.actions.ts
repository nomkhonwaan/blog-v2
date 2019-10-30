import { createAction } from '@ngrx/store';

export const toggleSidebar = createAction('TOGGLE_SIDEBAR');
export const isFetching = createAction('IS_FETCHING');
export const isNotFetching = createAction('IS_NOT_FETCHING');
