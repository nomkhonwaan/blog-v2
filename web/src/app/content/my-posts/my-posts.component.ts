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
  posts: Array<Post> = [];

  /**
   * An original list of posts
   */
  actualPosts: Array<Post> = [];

  /**
   * List of FontAwesome icons
   */
  icons: { [name: string]: IconDefinition } = {
    faSpinnerThird,
  };

  /**
   * A current offset number
   */
  offset: number = 0;

  /**
   * A maximum items per page
   */
  itemsPerPage = 10;

  constructor(
    private apollo: Apollo,
    private changeDetectorRef: ChangeDetectorRef,
  ) { }

  ngOnInit(): void {
    this.renderMyPosts(this.offset, this.itemsPerPage + 1);
  }

  renderMyPosts(offset: number, limit: number): void {
    this.apollo.query({
      query: gql`
        {
          myPosts(offset: $offset, limit: $limit) {
            title slug
            status
            html
            publishedAt
            categories { name slug }
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
      finalize((): void => this.changeDetectorRef.markForCheck()),
    ).subscribe((posts: Array<Post>): void => {
      this.actualPosts = posts;
      this.posts = this.posts.concat(posts.slice(0, this.itemsPerPage));
    });
  }

  onScroll(): void {
    if (this.actualPosts.length > this.itemsPerPage) {
      this.offset += this.itemsPerPage;
      this.renderMyPosts(this.offset, this.itemsPerPage + 1);
    }
  }

}
