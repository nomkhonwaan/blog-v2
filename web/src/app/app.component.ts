import { Component } from '@angular/core';
import { Store, select } from '@ngrx/store';
import { Observable } from 'rxjs';

import { toggleSidebar } from './app.actions';
import { AppState } from './app.reducer';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {

  app$: Observable<AppState>;

  constructor(private store: Store<AppState>) {
    this.app$ = store.pipe(select('app'));
  }

  toggleSidebar() {
    this.store.dispatch(toggleSidebar());
  }

}
