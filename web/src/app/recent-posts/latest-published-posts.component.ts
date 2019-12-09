import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { Apollo } from 'apollo-angular';
import { ApolloQueryResult } from 'apollo-client';
import gql from 'graphql-tag';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

import { environment } from 'src/environments/environment';

@Component({
  selector: 'app-latest-published-posts',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './latest-published-posts.component.html',
  styleUrls: ['./latest-published-posts.component.scss'],
})
export class LatestPublishedPostsComponent implements OnInit {

  /**
   * Use to display as list of recent posts
   */
  latestPublishedPosts$: Observable<Post[]>;

  constructor(private apollo: Apollo, private title: Title) { }

  ngOnInit(): void {
    this.title.setTitle(environment.title);

    this.latestPublishedPosts$ = this.apollo.query({
      query: gql`
        {
          latestPublishedPosts(offset: 0, limit: 6) {
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
