import { Injectable } from '@angular/core';
import { HttpInterceptor, HttpEvent, HttpHandler, HttpRequest } from '@angular/common/http';
import { Store, select } from '@ngrx/store';
import { Observable, ObservableInput, of, merge } from 'rxjs';
import { finalize, first, mergeMap, debounceTime, tap } from 'rxjs/operators';

import { isFetching, isNotFetching } from './app.actions';
import { AuthService } from './auth/auth.service';

@Injectable()
export class AppHttpInterceptor implements HttpInterceptor {

  auth$: Observable<{ accessToken: string }>;

  constructor(private auth: AuthService, private store: Store<{ app: AppState }>) {
    this.auth$ = store.pipe(select('app', 'auth'));
  }

  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    // this.store.dispatch(isFetching());

    return this.withAuthorization(req).pipe(
      first(),
      mergeMap((req: HttpRequest<any>): ObservableInput<HttpEvent<any>> => {
        return next.handle(req);
      }),
      // finalize((): void => { this.store.dispatch(isNotFetching()); }),
    );
  }

  private withAuthorization(req: HttpRequest<any>): Observable<HttpRequest<any>> {
    return this.auth$.pipe(
      first(),
      mergeMap((auth?: { accessToken: string }): ObservableInput<any> => {
        if (this.auth.isAuthenticated()) {
          return of<HttpRequest<any>>(
            req.clone({
              setHeaders: {
                Authorization: `Bearer ${auth.accessToken}`,
              },
              withCredentials: true,
            }),
          );
        }

        return of<HttpRequest<any>>(req);
      }),
    );
  }

}
