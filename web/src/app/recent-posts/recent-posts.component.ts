import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';
import { Title } from '@angular/platform-browser';
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

  constructor(private apollo: Apollo, private title: Title) { }

  ngOnInit(): void {
    this.title.setTitle(`Nomkhonwaan | Trust me I'm Petdo`);
    this.latestPublishedPosts$ = this.apollo.query({
      query: gql`
        {
          latestPublishedPosts(offset: 0, limit: 3) {
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
