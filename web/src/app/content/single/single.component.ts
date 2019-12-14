import { ChangeDetectionStrategy, ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { ActivatedRoute } from '@angular/router';
import { Apollo } from 'apollo-angular';
import { ApolloQueryResult } from 'apollo-client';
import gql from 'graphql-tag';
import { finalize, map } from 'rxjs/operators';
import { environment } from 'src/environments/environment';

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
    const slug: string = this.route.snapshot.paramMap.get('slug');

    // When year less than 2019, redirect to the v1 sub-domain
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

    this.renderPost(slug);
  }

  renderPost(slug: string): void {
    this.apollo.query({
      query: gql`
        {
          post(slug: $slug) {
            title slug
            html
            publishedAt
            categories { name slug }
            tags { name slug }
            featuredImage { slug }
            engagement { shareCount }
          }
        }
      `,
      variables: {
        slug,
      },
    }).pipe(
      map((result: ApolloQueryResult<{ post: Post }>): Post => result.data.post),
      finalize((): void => this.changeDetectorRef.markForCheck()),
    ).subscribe((post: Post): void => {
      this.post = post;

      this.title.setTitle(`${post.title} - ${environment.title}`);
    });
  }
}
