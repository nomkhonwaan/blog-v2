import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Apollo } from 'apollo-angular';
import gql from 'graphql-tag';
import { ApolloQueryResult } from 'apollo-client';

@Component({
  selector: 'app-single',
  templateUrl: './single.component.html',
  styleUrls: ['./single.component.scss'],
})
export class SingleComponent implements OnInit {

  p: Post;

  constructor(private apollo: Apollo, private route: ActivatedRoute) { }

  ngOnInit(): void {
    this.apollo.watchQuery({
      query: gql`
        {
          post(idOrSlug: $idOrSlug) {
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
      `,
      variables: {
        idOrSlug: this.route.snapshot.paramMap.get('slug'),
      }
    }).valueChanges.subscribe((result: ApolloQueryResult<{ post: Post }>): void => {
      this.p = result.data.post;
    });
  }

}
