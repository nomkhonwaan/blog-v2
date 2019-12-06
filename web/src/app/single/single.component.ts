import { Component, OnInit, ChangeDetectionStrategy, ChangeDetectorRef } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { ActivatedRoute } from '@angular/router';
import { Apollo } from 'apollo-angular';
import { ApolloQueryResult } from 'apollo-client';
import gql from 'graphql-tag';
import { map, finalize } from 'rxjs/operators';

@Component({
  selector: 'app-single',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './single.component.html',
  styleUrls: ['./single.component.scss'],
})
export class SingleComponent implements OnInit {

  /**
   * A post object
   */
  post: Post;

  constructor(
    private apollo: Apollo,
    private route: ActivatedRoute,
    private title: Title,
    private changeDetectorRef: ChangeDetectorRef,
  ) { }

  ngOnInit(): void {
    const year: number = parseInt(this.route.snapshot.paramMap.get('year'), 10);

    // automatically redirect to version 1.0 website if the published year less than 2019
    if (year < 2019) {
      window.location.href = [
        'https://v1.nomkhonwaan.com',
        this.route.snapshot.paramMap.get('year'),
        this.route.snapshot.paramMap.get('month'),
        this.route.snapshot.paramMap.get('date'),
        this.route.snapshot.paramMap.get('slug'),
      ].join('/');

      return;
    }

    this.apollo.query({
      query: gql`
        {
          post(slug: $slug) {
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
            engagement {
              shareCount
            }
          }
        }
      `,
      variables: {
        slug: this.route.snapshot.paramMap.get('slug'),
      },
    }).pipe(
      map((result: ApolloQueryResult<{ post: Post }>): Post => result.data.post),
      finalize((): void => this.changeDetectorRef.markForCheck()),
    ).subscribe((post: Post): void => {
      this.title.setTitle(`${post.title} - Nomkhonwaan | Trust me I'm Petdo`);
      this.post = post;
    });
  }

}
