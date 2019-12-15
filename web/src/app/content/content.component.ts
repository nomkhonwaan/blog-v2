import { animate, state, style, transition, trigger } from '@angular/animations';
import { Component, HostBinding, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { faHeart, IconDefinition as SolidIconDefinition } from '@fortawesome/fontawesome-free-solid';
import { faGithubSquare, faMedium, IconDefinition as BrandIconDefinition } from '@fortawesome/free-brands-svg-icons';
import { faBars, faSearch, faTimes, IconDefinition } from '@fortawesome/pro-light-svg-icons';
import { select, Store } from '@ngrx/store';
import { Apollo } from 'apollo-angular';
import { ApolloQueryResult } from 'apollo-client';
import gql from 'graphql-tag';
import { map } from 'rxjs/operators';
import { environment } from '../../environments/environment';
import { AuthService } from '../auth';
import { toggleSidebar } from '../index';

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
  selector: 'app-content',
  templateUrl: './content.component.html',
  styleUrls: ['./content.component.scss'],
})
export class ContentComponent implements OnInit {

  /**
   * Use to toggle application sidebar
   */
  @HostBinding('@slideInOut')
  hasSidebarExpanded = false;

  /**
   * Use to indicate whether loading animation should show or not
   */
  isFetching: boolean;

  /**
   * List of FontAwesome icons
   */
  icons: { [name: string]: IconDefinition | SolidIconDefinition | BrandIconDefinition } = {
    faBars,
    faHeart,
    faSearch,
    faTimes,
    faGithubSquare,
    faMedium,
  };

  /**
   * List of categories to-be rendered as sidebar menu-item(s)
   */
  categories: Array<Category>;

  /**
   * An authenticated user info object
   */
  userInfo: UserInfo | null;

  /**
   * A version number which will collect from Git tag
   */
  version: string = environment.version;

  /**
   * A revision number which will collect from Git commit ID
   */
  revision: string = environment.revision;

  constructor(
    private apollo: Apollo,
    private auth: AuthService,
    private router: Router,
    private store: Store<{ app: AppState }>,
  ) { }

  ngOnInit(): void {
    this.store.pipe(select('app')).subscribe((app: AppState): void => {
      this.hasSidebarExpanded = !app.sidebar.collapsed;
      this.userInfo = app.auth.userInfo;
    });

    if (this.isLoggedIn()) {
      this.auth.renewTokens();
    }

    this.renderSidebar();
  }

  isLoggedIn(): boolean {
    return this.auth.isLoggedIn();
  }

  login(): void {
    this.auth.login(this.router.url);
  }

  toggleSidebar(): void {
    this.store.dispatch(toggleSidebar());
  }

  renderSidebar(): void {
    this.apollo.query({
      query: gql`
        {
          categories { name slug }
        }
      `,
    }).pipe(
      map((result: ApolloQueryResult<{ categories: Array<Category> }>): Array<Category> => result.data.categories),
    ).subscribe((categories: Array<Category>): void => {
      this.categories = categories;
    });
  }

}
