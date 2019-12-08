import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { Apollo } from 'apollo-angular';
import { ApolloQueryResult } from 'apollo-client';
import gql from 'graphql-tag';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { environment } from 'src/environments/environment';

@Component({
  selector: 'app-recent-posts',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './recent-posts.component.html',
  styleUrls: ['./recent-posts.component.scss'],
})
export class RecentPostsComponent implements OnInit {

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
