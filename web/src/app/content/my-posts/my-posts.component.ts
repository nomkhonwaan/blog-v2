import { ChangeDetectionStrategy, ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { faSpinnerThird, IconDefinition } from '@nomkhonwaan/pro-light-svg-icons';
import { Apollo } from 'apollo-angular';
import { ApolloQueryResult } from 'apollo-client';
import gql from 'graphql-tag';
import { finalize, map } from 'rxjs/operators';
import { environment } from 'src/environments/environment';

@Component({
  selector: 'app-my-posts',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './my-posts.component.html',
  styleUrls: ['./my-posts.component.scss'],
})
export class MyPostsComponent implements OnInit {

  /**
   * All loaded posts
   */
  posts: Array<Post> = [];

  /**
   * List of posts that return from the service when reached the bottom of the page
   */
  newQueriedPosts: Array<Post> = [];

  /**
   * List of FontAwesome icons
   */
  icons: { [name: string]: IconDefinition } = {
    faSpinnerThird,
  };

  /**
   * A current offset number
   */
  offset = 0;

  /**
   * A maximum items per page
   */
  itemsPerPage = 10;

  constructor(
    private apollo: Apollo,
    private changeDetectorRef: ChangeDetectorRef,
    private title: Title,
  ) { }

  ngOnInit(): void {
    this.title.setTitle(`My posts - ${environment.title}`);

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
      this.newQueriedPosts = posts;
      this.posts = this.posts.concat(posts.slice(0, this.itemsPerPage));
    });
  }

  onScroll(): void {
    if (this.newQueriedPosts.length > this.itemsPerPage) {
      this.offset += this.itemsPerPage;
      this.renderMyPosts(this.offset, this.itemsPerPage + 1);
    }
  }

}
