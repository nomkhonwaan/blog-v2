import { createAction, props } from '@ngrx/store';

export const isFetching = createAction('IS_FETCHING');
export const isNotFetching = createAction('IS_NOT_FETCHING');
export const setAuthentication = createAction('SET_AUTHENTICATION', props<{
  accessToken: string,
  idToken: string,
  userInfo: UserInfo | null,
}>());
export const toggleSidebar = createAction('TOGGLE_SIDEBAR');
