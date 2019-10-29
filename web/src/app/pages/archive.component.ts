import { OnInit, Component } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Apollo } from 'apollo-angular';
import { ApolloQueryResult } from 'apollo-client';
import gql from 'graphql-tag';

@Component({
  selector: 'app-archive',
  templateUrl: './archive.component.html',
  styleUrls: ['./archive.component.scss'],
})
export class ArchiveComponent implements OnInit {

  archive: Category | Tag;

  constructor(private apollo: Apollo, private route: ActivatedRoute) { }

  ngOnInit(): void {
    const type: string = (this.route.snapshot.data as { type: string }).type;
    this.apollo.watchQuery({
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
    }).valueChanges.subscribe((result: ApolloQueryResult<{ archive: Category | Tag }>): void => {
      this.archive = result.data.archive;
    });
  }

}
