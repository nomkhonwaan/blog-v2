import { Component, OnInit, ChangeDetectionStrategy, ChangeDetectorRef } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { ActivatedRoute } from '@angular/router';
import { Apollo } from 'apollo-angular';
import { ApolloQueryResult } from 'apollo-client';
import gql from 'graphql-tag';
import { map, finalize } from 'rxjs/operators';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-single',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './single.component.html',
  styleUrls: ['./single.component.scss'],
})
export class SingleComponent implements OnInit {

  /**
   * A single post object
   */
  post: Post;

  constructor(
    private apollo: Apollo,
    private route: ActivatedRoute,
    private title: Title,
    private changeDetectorRef: ChangeDetectorRef,
  ) { }

  ngOnInit(): void {
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
      this.title.setTitle(post.title + ' - ' + this.title.getTitle());
      this.post = post;
    });
  }

}
