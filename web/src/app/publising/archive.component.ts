import { OnInit, Component, ChangeDetectionStrategy, ChangeDetectorRef } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { ActivatedRoute, ParamMap, Router } from '@angular/router';
import { Apollo } from 'apollo-angular';
import { ApolloQueryResult } from 'apollo-client';
import gql from 'graphql-tag';
import { map, finalize, first, switchMap } from 'rxjs/operators';

import { environment } from 'src/environments/environment';

@Component({
  selector: 'app-archive',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './archive.component.html',
  styleUrls: ['./recent-posts.component.scss'],
})
export class ArchiveComponent implements OnInit {

  /**
   * Type of the archive page
   */
  type: string;

  /**
   * An archive object
   */
  archive: Category | Tag;

  /**
   * List of published posts
   */
  posts: Array<Post>;

  /**
   * Use to redirect to the next page, if posts still have more
   */
  nextPageRouterLink: Array<string>;

  /**
   * Use to redirect to the previous page, if the current page is not 1
   */
  previousPageRouterLink: Array<string>;

  constructor(
    private apollo: Apollo,
    private route: ActivatedRoute,
    private title: Title,
    private changeDetectorRef: ChangeDetectorRef,
  ) { }

  ngOnInit(): void {
    this.type = (this.route.snapshot.data as { type: string }).type;

    this.route.paramMap.subscribe((paramMap: ParamMap): void => {
      const page = paramMap.has('page') ? parseInt(paramMap.get('page'), 10) : 1;
      const query: string = this.type === 'all' ? this.buildLatestPublishedPostsQuery() : this.buildArchiveQuery(this.type);
      const offset: number = (page - 1) * 5;

      this.apollo.query({
        query: gql`
          ${query}

          fragment PublishedPost on Post {
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
        `,
        variables: {
          slug: paramMap.get('slug'),
          offset,
        },
      }).pipe(
        finalize((): void => this.changeDetectorRef.markForCheck()),
      ).subscribe((result: ApolloQueryResult<{ category?: Category, tag?: Tag, latestPublishedPosts?: Array<Post> }>): void => {
        if (result.data.latestPublishedPosts) {
          this.posts = result.data.latestPublishedPosts;

          this.title.setTitle(`Recent Posts - ${environment.title}`);
        } else if (result.data.category || result.data.tag) {
          this.archive = result.data[this.type];
          this.posts = result.data[this.type].latestPublishedPosts;

          this.title.setTitle(`${this.archive.name} - ${environment.title}`);
        }

        if (page > 1) {
          this.title.setTitle(`Page ${page} · ${this.title.getTitle()}`);

          this.previousPageRouterLink = this.buildPreviousPageRouterLink(page);
        } else {
          this.previousPageRouterLink = null;
        }

        if (this.posts.length > 5) {
          this.nextPageRouterLink = this.buildNextPageRouterLink(page);
        } else {
          this.nextPageRouterLink = null;
        }
      });
    });
  }

  buildLatestPublishedPostsQuery(): string {
    return `
      {
        latestPublishedPosts(offset: $offset, limit: 6) {
          ...PublishedPost
        }
      }
    `;
  }

  buildArchiveQuery(type: string): string {
    return `
      {
        ${type}(slug: $slug) {
          name
          slug
          latestPublishedPosts(offset: $offset, limit: 6) {
            ...PublishedPost
          }
        }
      }
    `;
  }

  buildPreviousPageRouterLink(currentPage: number): Array<string> {
    // TODO: need to replace this silly function for building a dynamic URL on Angular
    if (this.type === 'all') {
      return ['/', (currentPage - 1).toString()];
    } else {
      return ['/', this.type, this.archive.slug, (currentPage - 1).toString()];
    }
  }

  buildNextPageRouterLink(currentPage: number): Array<string> {
    // TODO: need to replace this silly function for building a dynamic URL on Angular
    if (this.type === 'all') {
      return ['/', (currentPage + 1).toString()];
    } else {
      return ['/', this.type, this.archive.slug, (currentPage + 1).toString()];
    }
  }
}
