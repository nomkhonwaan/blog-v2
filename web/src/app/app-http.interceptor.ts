import { Injectable } from '@angular/core';
import { HttpInterceptor, HttpEvent, HttpHandler, HttpRequest } from '@angular/common/http';
import { Store, select } from '@ngrx/store';
import { Observable, ObservableInput, of } from 'rxjs';
import { first, mergeMap } from 'rxjs/operators';

import { AuthService } from './auth/auth.service';

@Injectable()
export class AppHttpInterceptor implements HttpInterceptor {

  auth$: Observable<{ accessToken: string }>;

  constructor(private auth: AuthService, private store: Store<{ app: AppState }>) {
    this.auth$ = store.pipe(select('app', 'auth'));
  }

  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    return this.withAuthorization(req).pipe(
      first(),
      mergeMap((authorizedRequest: HttpRequest<any>): ObservableInput<HttpEvent<any>> => {
        return next.handle(authorizedRequest);
      }),
    );
  }

  private withAuthorization(req: HttpRequest<any>): Observable<HttpRequest<any>> {
    return this.auth$.pipe(
      first(),
      mergeMap((auth?: { accessToken: string }): ObservableInput<any> => {
        if (this.auth.isLoggedIn()) {
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
