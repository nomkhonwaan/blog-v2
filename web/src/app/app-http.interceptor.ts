import { Injectable } from '@angular/core';
import { HttpInterceptor, HttpEvent, HttpHandler, HttpRequest, HttpErrorResponse, HttpResponse } from '@angular/common/http';
import { Store } from '@ngrx/store';
import { Observable } from 'rxjs';
import { tap } from 'rxjs/operators';

import { isFetching, isNotFetching } from './app.actions';

@Injectable()
export class AppHttpInterceptor implements HttpInterceptor {

  constructor(private store: Store<AppState>) { }

  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    this.store.dispatch(isFetching());

    return next.handle(req)
      .pipe(
        tap(
          (event: any): void => {
            if (event instanceof HttpResponse) {
              this.store.dispatch(isNotFetching());
            }
          },
          (err: any): void => {
            if (err instanceof HttpErrorResponse) {
              this.store.dispatch(isNotFetching());
            }
          },
        ),
      );
  }

}
