import { animate, state, style, transition, trigger } from '@angular/animations';
import { Component, HostBinding, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { Store, select } from '@ngrx/store';
import { faBars, faSearch, faTimes, IconDefinition } from '@fortawesome/pro-light-svg-icons';
import { faHeart, IconDefinition as SolidIconDefinition } from '@fortawesome/fontawesome-free-solid';
import { faGithubSquare, faMedium, IconDefinition as BrandIconDefinition } from '@fortawesome/free-brands-svg-icons';
import { Apollo } from 'apollo-angular';
import { ApolloQueryResult } from 'apollo-client';
import gql from 'graphql-tag';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

import { toggleSidebar } from '../app.actions';
import { AuthService } from '../auth/auth.service';

import { environment } from '../../environments/environment';

const coffeeCup = require('../../assets/lottie-files/lf30_editor_pohhBA.json');

@Component({
  animations: [
    trigger('slideInOut', [
      state('true', style({ transform: 'translateX(0)' })),
      state('false', style({ transform: 'translateX(-25.6rem)' })),
      transition('* => true', [
        animate('.4s ease-in-out', style({ transform: 'translateX(0)' })),
      ]),
      transition('true => false', [
        animate('.4s ease-in-out', style({ transform: 'translateX(-25.6rem)' })),
      ]),
    ]),
  ],
  selector: 'app-publishing',
  templateUrl: './publishing.component.html',
  styleUrls: ['./publishing.component.scss'],
})
export class PublishingComponent implements OnInit {

  faBars: IconDefinition = faBars;
  faHeart: SolidIconDefinition = faHeart;
  faSearch: IconDefinition = faSearch;
  faTimes: IconDefinition = faTimes;
  faGithubSquare: BrandIconDefinition = faGithubSquare;
  faMedium: BrandIconDefinition = faMedium;

  /**
   * Use to toggle sidebar pane for showing or hiding
   */
  @HostBinding('@slideInOut')
  hasSidebarExpanded = false;

  /**
   * Use to display loading animation while fetching resources
   */
  isFetching = false;

  /**
   * Use to display at sidebar as a sub-menu to the group of posts
   */
  categories$: Observable<Category[]>;

  /**
   * Use to render with animation directive
   */
  loadingAnimationData: any;

  /**
   * An authenticated user information
   */
  userInfo$: Observable<UserInfo | null>;

  /**
   * Use to display at footer section as a build version number
   */
  version: string = environment.version;

  /**
   * Use to display at footer section as a commit ID
   */
  revision: string = environment.revision;

  constructor(
    private apollo: Apollo,
    private auth: AuthService,
    private router: Router,
    private store: Store<{ app: AppState }>,
  ) {
    this.loadingAnimationData = coffeeCup;

    this.userInfo$ = store.pipe(select('app', 'auth', 'userInfo'));

    store.pipe(select('app')).subscribe((app: AppState): void => {
      this.isFetching = app.isFetching;
      this.hasSidebarExpanded = !app.sidebar.collapsed;
    });
  }

  ngOnInit(): void {
    // Try to check and renew an authentication token if possible
    this.renewTokenIfAuthenticated();

    // Perform a query to the GraphQL server for retrieving a list of categories for displaying on sidebar menu
    this.queryAllCategories();
  }

  isAuthenticated(): boolean {
    return this.auth.isAuthenticated();
  }

  login(): void {
    this.auth.login(this.router.url);
  }

  renewTokenIfAuthenticated(): void {
    if (this.auth.isAuthenticated()) {
      this.auth.renewTokens();
    }
  }

  queryAllCategories(): void {
    this.categories$ = this.apollo.watchQuery({
      query: gql`
          {
            categories {
              name
              slug
            }
          }
        `
    }).valueChanges.pipe(
      map((result: ApolloQueryResult<{ categories: Category[] }>): Category[] => result.data.categories),
    );
  }

  toggleSidebar() {
    this.store.dispatch(toggleSidebar());
  }

}
