import { animate, state, style, transition, trigger } from '@angular/animations';
import { Component, HostBinding, OnInit } from '@angular/core';
import { Store, select } from '@ngrx/store';
import { faBars, faSearch, faTimes, IconDefinition } from '@fortawesome/pro-light-svg-icons';
import { faGithubSquare, faMedium, IconDefinition as BrandIconDefinition } from '@fortawesome/free-brands-svg-icons';
import { Apollo } from 'apollo-angular';
import gql from 'graphql-tag';
import { Observable } from 'rxjs';

import { toggleSidebar } from './app.actions';
import { ApolloQueryResult } from 'apollo-client';
import { AuthService } from './auth/auth.service';
import { Router } from '@angular/router';

@Component({
  animations: [
    trigger('slideInOut', [
      state('true', style({ transform: 'translateX(0)' })),
      transition('* => true', [
        animate('.4s ease-in-out', style({ transform: 'translateX(0)' })),
      ]),
      transition('true => false', [
        style({ transform: 'translateX(0)' }),
        animate('.4s ease-in-out'),
      ]),
    ]),
  ],
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {

  app$: Observable<AppState>;

  faBars: IconDefinition = faBars;
  faSearch: IconDefinition = faSearch;
  faTimes: IconDefinition = faTimes;
  faGithubSquare: BrandIconDefinition = faGithubSquare;
  faMedium: BrandIconDefinition = faMedium;

  @HostBinding('@slideInOut')
  sidebarExpanded = false;

  /**
   * Use to display at sidebar as a sub-menu to the group of posts
   */
  categories: Category[];

  /**
   * Use to display at footer section as a current year of the copyright
   */
  fullYear: string;

  constructor(
    private apollo: Apollo,
    private store: Store<AppState>,
    private auth: AuthService,
    private router: Router,
  ) {
    this.app$ = store.pipe(select('app'));
    this.app$.subscribe(({ sidebar }: AppState): void => { this.sidebarExpanded = !sidebar.collapsed; });
  }

  ngOnInit(): void {
    if (this.auth.isAuthenticated()) {
      this.auth.renewTokens();
    }

    this.apollo.watchQuery({
      query: gql`
        {
          categories {
            name
            slug
          }
        }
      `
    }).valueChanges.subscribe((result: ApolloQueryResult<{ categories: Category[] }>): void => {
      this.categories = result.data.categories;
    });

    this.fullYear = new Date().getFullYear().toString();
  }

  toggleSidebar() {
    this.store.dispatch(toggleSidebar());
  }

}
