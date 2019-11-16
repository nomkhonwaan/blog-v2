import { Injectable } from '@angular/core';
import { HttpInterceptor, HttpEvent, HttpHandler, HttpRequest } from '@angular/common/http';
import { Store } from '@ngrx/store';
import { Observable } from 'rxjs';
import { finalize } from 'rxjs/operators';

import { isFetching, isNotFetching } from './app.actions';

@Injectable()
export class AppHttpInterceptor implements HttpInterceptor {

  constructor(private store: Store<AppState>) { }

  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    this.store.dispatch(isFetching());

    return next.handle(req).pipe(
      finalize((): void => {
        this.store.dispatch(isNotFetching());
      })
    );
  }

}
