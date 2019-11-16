import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';
import { Store } from '@ngrx/store';
import { Apollo } from 'apollo-angular';
import { ApolloQueryResult } from 'apollo-client';
import gql from 'graphql-tag';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

@Component({
  selector: 'app-recent-posts',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './recent-posts.component.html',
  styleUrls: ['./recent-posts.component.scss'],
})
export class RecentPostsComponent implements OnInit {

  /**
   * Used to display as list of recent posts
   */
  latestPublishedPosts$: Observable<Post[]>;

  constructor(private apollo: Apollo, private store: Store<AppState>) { }

  ngOnInit(): void {
    this.latestPublishedPosts$ = this.apollo.query({
      query: gql`
        {
          latestPublishedPosts(offset: 0, limit: 5) {
            title
            slug
            html
            publishedAt
            categories {
              name slug
            }
            tags {
              name slug
            }
            featuredImage {
              slug
            }
          }
        }
      `,
    }).pipe(
      map((result: ApolloQueryResult<{ latestPublishedPosts: Post[] }>): Post[] => result.data.latestPublishedPosts),
    );
  }

}
