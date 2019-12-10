import { OnInit, Component, ChangeDetectionStrategy, ChangeDetectorRef } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { ActivatedRoute, ParamMap } from '@angular/router';
import { Apollo } from 'apollo-angular';
import { ApolloQueryResult } from 'apollo-client';
import gql from 'graphql-tag';
import { map, finalize } from 'rxjs/operators';

import { environment } from 'src/environments/environment';

@Component({
  selector: 'app-archive',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './archive.component.html',
  styleUrls: ['./latest-published-posts.component.scss'],
})
export class ArchiveComponent implements OnInit {

  /**
   * An archive object
   */
  archive: Category | Tag;

  /**
   * Type of the archive page
   */
  type: string;

  constructor(
    private apollo: Apollo,
    private route: ActivatedRoute,
    private title: Title,
    private changeDetectorRef: ChangeDetectorRef,
  ) { }

  ngOnInit(): void {
    this.type = (this.route.snapshot.data as { type: string }).type;


    this.route.paramMap.subscribe((paramMap: ParamMap): void => {
      const query: string = this.type === 'all'
        ? this.buildLatestPublishedPostsQuery()
        : this.buildArchiveQuery(this.type);

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
          offset: 0,
        },
      }).pipe(
        map((result: ApolloQueryResult<{ archive: Category | Tag }>): Category | Tag => result.data[this.type]),
        finalize((): void => this.changeDetectorRef.markForCheck()),
      ).subscribe((archive: Category | Tag): void => {
        this.title.setTitle(`${archive.name} - ${environment.title}`);
        this.archive = archive;
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
}
