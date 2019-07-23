import { Injectable } from '@angular/core';
import { Actions, createEffect, ofType, Effect } from '@ngrx/effects';
import { Action } from '@ngrx/store';
import { EMPTY, Observable } from 'rxjs';
import { catchError, map, mergeMap } from 'rxjs/operators';

import { fetchRecentPosts, fetchRecentPostsSuccess, fetchRecentPostsError } from './recent-posts.actions';

import { ApiService } from '../api/api.service';

@Injectable()
export class RecentPostsEffects {

  fetchRecentPosts$ = createEffect(() => this.actions$.pipe(
    ofType(fetchRecentPosts),
    mergeMap(() => this.api.fetchRecentPosts()
      .pipe(
        map((posts: Post[]) => fetchRecentPostsSuccess({ posts })),
        catchError(() => EMPTY),
      )
    ),
  ));

  constructor(
    private actions$: Actions,
    private api: ApiService,
  ) { }
}
