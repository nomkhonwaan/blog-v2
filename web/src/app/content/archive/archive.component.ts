import { ChangeDetectionStrategy, ChangeDetectorRef, Component, Input, OnInit } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { ActivatedRoute, ParamMap } from '@angular/router';
import { Apollo } from 'apollo-angular';
import { ApolloQueryResult } from 'apollo-client';
import gql from 'graphql-tag';
import { finalize, map } from 'rxjs/operators';
import { environment } from '../../../environments/environment';

@Component({
  selector: 'app-archive',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './archive.component.html',
  styleUrls: ['./archive.component.scss'],
})
export class ArchiveComponent implements OnInit {

  /**
   * When this property has been set to "all",
   * will retrieve the latest published posts from all categories
   * rather a single category or tag.
   */
  @Input()
  from: string;

  /**
   * List of posts
   */
  posts: Array<Post>;

  /**
   * An archive object
   */
  archive: Category | Tag;

  /**
   * Next and previous page URL will be displayed based on conditions
   */
  nextPage: Array<string>;
  previousPage: Array<string>;

  /**
   * A maximum items per page
   */
  itemsPerPage = 5;

  constructor(
    private apollo: Apollo,
    private route: ActivatedRoute,
    private title: Title,
    private changeDetectorRef: ChangeDetectorRef,
  ) { }

  ngOnInit(): void {
    if (!this.from) {
      this.from = (this.route.snapshot.data as { from: string }).from;
    }

    this.route.paramMap.pipe(
      map((paramMap: ParamMap): { page: number, slug?: string } => {
        const page: number = paramMap.has('page') ? parseInt(paramMap.get('page'), 10) : 1;
        const slug: string = paramMap.get('slug');

        return { page, slug };
      }),
    ).subscribe(({ page, slug }: { page: number, slug: string }): void => {
      const offset: number = (page - 1) * this.itemsPerPage;

      this.renderLatestPublishedPosts(page, slug, offset, this.itemsPerPage + 1);
    });
  }

  renderLatestPublishedPosts(page: number, slug: string, offset: number, limit: number): void {
    this.apollo.query({
      query: gql`
        ${this.buildQueryFrom(this.from)}

        fragment PublishedPost on Post {
          title slug
          html
          publishedAt
          categories { name slug }
          tags { name slug }
          featuredImage { slug }
        }
      `,
      variables: {
        slug,
        offset,
        limit,
      },
    }).pipe(
      finalize((): void => this.changeDetectorRef.markForCheck()),
    ).subscribe((result: ApolloQueryResult<{ archive?: Category | Tag, latestPublishedPosts?: Array<Post> }>): void => {
      if (result.data.latestPublishedPosts) {
        this.posts = result.data.latestPublishedPosts;

        this.title.setTitle(environment.title);
      } else if (result.data.archive) {
        this.archive = result.data.archive;
        this.posts = result.data.archive.latestPublishedPosts;

        this.title.setTitle(`${this.archive.name} - ${environment.title}`);
      }

      if (page > 1) {
        this.title.setTitle(`Page ${page} · ${this.title.getTitle()}`);

        this.previousPage = this.from === 'all'
          ? ['/', (page - 1).toString()]
          : ['/', this.from, this.archive.slug, (page - 1).toString()];
      } else {
        this.previousPage = null;
      }

      if (this.posts.length > 5) {
        this.nextPage = this.from === 'all'
          ? ['/', (page + 1).toString()]
          : ['/', this.from, this.archive.slug, (page + 1).toString()];
      } else {
        this.nextPage = null;
      }
    });
  }

  buildQueryFrom(from: string): string {
    if (from === 'all') {
      return `
        {
          latestPublishedPosts(offset: $offset, limit: $limit) {
            ...PublishedPost
          }
        }
      `;
    }

    return `
      {
        archive: ${from}(slug: $slug) {
          name slug
          latestPublishedPosts(offset: $offset, limit: $limit) {
            ...PublishedPost
          }
        }
      }
    `;
  }
}
