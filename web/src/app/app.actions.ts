import { createAction, props } from '@ngrx/store';

export const isFetching = createAction('IS_FETCHING');
export const isNotFetching = createAction('IS_NOT_FETCHING');
export const setAuthorization = createAction('SET_AUTHORIZATION', props<{ accessToken: string }>());
export const toggleSidebar = createAction('TOGGLE_SIDEBAR');
