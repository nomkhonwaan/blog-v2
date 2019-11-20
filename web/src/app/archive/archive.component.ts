import { OnInit, Component, ChangeDetectionStrategy, ChangeDetectorRef } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { ActivatedRoute } from '@angular/router';
import { Apollo } from 'apollo-angular';
import { ApolloQueryResult } from 'apollo-client';
import gql from 'graphql-tag';
import { map, finalize } from 'rxjs/operators';

@Component({
  selector: 'app-archive',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './archive.component.html',
  styleUrls: ['./archive.component.scss'],
})
export class ArchiveComponent implements OnInit {

  /**
   * An archive object
   */
  archive: Category | Tag;

  constructor(
    private apollo: Apollo,
    private route: ActivatedRoute,
    private title: Title,
    private changeDetectorRef: ChangeDetectorRef,
  ) { }

  ngOnInit(): void {
    const type: string = (this.route.snapshot.data as { type: string }).type;
    this.apollo.query({
      query: gql`
        {
          ${type}(slug: $slug) {
            name
            slug
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
            }
          }
        }
      `,
      variables: {
        slug: this.route.snapshot.paramMap.get('slug'),
      },
    }).pipe(
      map((result: ApolloQueryResult<{ archive: Category | Tag }>): Category | Tag => result.data[type]),
      finalize((): void => this.changeDetectorRef.markForCheck()),
    ).subscribe((archive: Category | Tag): void => {
      this.title.setTitle(`${archive.name} - Nomkhonwaan | Trust me I'm Petdo`);
      this.archive = archive;
    });
  }

}
