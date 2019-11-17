import { OnInit, Component, ChangeDetectionStrategy, ChangeDetectorRef } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Apollo } from 'apollo-angular';
import { ApolloQueryResult } from 'apollo-client';
import gql from 'graphql-tag';
import { map, finalize } from 'rxjs/operators';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-archive',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './archive.component.html',
  styleUrls: ['./archive.component.scss'],
})
export class ArchiveComponent implements OnInit {

  archive$: Observable<Category | Tag>;

  constructor(private apollo: Apollo, private route: ActivatedRoute, private changeDetectorRef: ChangeDetectorRef) { }

  ngOnInit(): void {
    const type: string = (this.route.snapshot.data as { type: string }).type;
    this.archive$ = this.apollo.query({
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
      map((result: ApolloQueryResult<{ archive: Category | Tag }>): Category | Tag => result.data.archive),
      finalize((): void => this.changeDetectorRef.markForCheck()),
    );
  }

}
