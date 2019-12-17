import { ChangeDetectionStrategy, ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { faSpinnerThird, IconDefinition } from '@fortawesome/pro-light-svg-icons';
import { Apollo } from 'apollo-angular';
import { ApolloQueryResult } from 'apollo-client';
import gql from 'graphql-tag';
import { finalize, map } from 'rxjs/operators';

@Component({
  selector: 'app-my-posts',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './my-posts.component.html',
  styleUrls: ['./my-posts.component.scss'],
})
export class MyPostsComponent implements OnInit {

  /**
   * List of posts
   */
  posts: Array<Post>;

  /**
   * Use to indicate loading status
   */
  isFetching = false;

  /**
   * List of FontAwesome icons
   */
  icons: { [name: string]: IconDefinition } = {
    faSpinnerThird,
  };

  constructor(
    private apollo: Apollo,
    private changeDetectorRef: ChangeDetectorRef,
  ) { }

  ngOnInit(): void {
    this.renderMyPosts(0, 11);
  }

  renderMyPosts(offset: number, limit: number): void {
    this.isFetching = true;

    this.apollo.query({
      query: gql`
        {
          myPosts(offset: $offset, limit: $limit) {
            title slug
            status
            publishedAt
            createdAt updatedAt
          }
        }
      `,
      variables: {
        offset,
        limit,
      },
      fetchPolicy: 'network-only',
    }).pipe(
      map((result: ApolloQueryResult<{ myPosts: Array<Post> }>): Array<Post> => result.data.myPosts),
      finalize((): void => {
        this.changeDetectorRef.markForCheck();
        this.isFetching = true;
      }),
    ).subscribe((posts: Array<Post>): void => {
      this.posts = posts;
    });
  }

}
