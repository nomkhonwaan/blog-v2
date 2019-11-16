import { animate, state, style, transition, trigger } from '@angular/animations';
import { Component, HostBinding, OnInit, Directive, ElementRef, Input, NgZone } from '@angular/core';
import { Store, select } from '@ngrx/store';
import { faBars, faSearch, faTimes, IconDefinition } from '@fortawesome/pro-light-svg-icons';
import { faGithubSquare, faMedium, IconDefinition as BrandIconDefinition } from '@fortawesome/free-brands-svg-icons';
import { Apollo } from 'apollo-angular';
import gql from 'graphql-tag';
import Lottie from 'lottie-web';
import { Observable } from 'rxjs';
import { map, debounceTime } from 'rxjs/operators';

import { toggleSidebar } from './app.actions';
import { ApolloQueryResult } from 'apollo-client';
import { AuthService } from './auth/auth.service';
import { Router } from '@angular/router';

import { environment } from '../environments/environment';

const coffeeCup = require('../assets/lottie-files/lf30_editor_pohhBA.json');

@Directive({ selector: '[appAnimation]' })
export class AnimationDirective implements OnInit {

  @Input()
  data: any;

  constructor(private el: ElementRef, private ngZone: NgZone) { }

  ngOnInit(): void {
    this.ngZone.runOutsideAngular((): void => {
      Lottie.loadAnimation({
        container: this.el.nativeElement,
        renderer: 'svg',
        loop: true,
        autoplay: true,
        animationData: this.data,
      });
    });
  }

}

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
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit {

  faBars: IconDefinition = faBars;
  faSearch: IconDefinition = faSearch;
  faTimes: IconDefinition = faTimes;
  faGithubSquare: BrandIconDefinition = faGithubSquare;
  faMedium: BrandIconDefinition = faMedium;

  /**
   * Used to toggle sidebar pane for showing or hiding
   */
  @HostBinding('@slideInOut')
  sidebarExpanded = false;

  /**
   * Used to display loading animation while fetching resources
   */
  isFetching = false;

  /**
   * Used to display at sidebar as a sub-menu to the group of posts
   */
  categories$: Observable<Category[]>;

  /**
   * Used to render with animation directive
   */
  loadingAnimationData: any;

  /**
   * Used to display at footer section as a build version number
   */
  version: string = environment.version;

  /**
   * Used to display at footer section as a commit ID
   */
  revision: string = environment.revision;

  /**
   * Used to display at footer section as a current year of the copyright
   */
  fullYear: string;

  constructor(
    private apollo: Apollo,
    private store: Store<AppState>,
    private auth: AuthService,
    private router: Router,
  ) {
    this.loadingAnimationData = coffeeCup;

    store
      .pipe(select('app', 'isFetching'))
      .pipe(debounceTime(0))
      .subscribe((isFetching: boolean): void => {
        this.isFetching = isFetching;
      });

    store
      .pipe(select('app'))
      .subscribe(({ sidebar }: AppState): void => {
        this.sidebarExpanded = !sidebar.collapsed;
      });
  }

  ngOnInit(): void {
    // Try to check and renew an authentication token if possible
    this.renewTokenIfAuthenticated();

    // Perform a query to the GraphQL server for retrieving a list of categories for displaying on sidebar menu
    this.queryAllCategories();

    // Get the current year on user's browser for displaying as a copyright year,
    // I know you can fool it but who care ¯\_(ツ)_/¯
    this.fullYear = new Date().getFullYear().toString();
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
