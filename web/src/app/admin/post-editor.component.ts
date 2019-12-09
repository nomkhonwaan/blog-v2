import { Component, OnInit } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { ActivatedRoute, Router } from '@angular/router';
import { faTimes, IconDefinition } from '@fortawesome/pro-light-svg-icons';
import { Apollo } from 'apollo-angular';
import { ApolloQueryResult } from 'apollo-client';
import { GraphQLError } from 'graphql';
import gql from 'graphql-tag';
import { map, tap } from 'rxjs/operators';

import { environment } from 'src/environments/environment';
import { Observable, forkJoin } from 'rxjs';
import { Store, select } from '@ngrx/store';

@Component({
  animations: [
  ],
  selector: 'app-post-editor',
  templateUrl: './post-editor.component.html',
  styleUrls: ['./post-editor.component.scss'],
})
export class PostEditorComponent implements OnInit {

  /**
   * A post object
   */
  post: Post;

  /**
   * Use to display when GraphQL returns errors
   */
  errors: ReadonlyArray<GraphQLError> = null;

  /**
   * Use to indicate loading status
   */
  isFetching = false;

  /**
   * An authenticated user information
   */
  userInfo$: Observable<UserInfo | null>;

  /**
   * To-be updated attachment
   */
  selectedAttachment: Attachment;

  faTimes: IconDefinition = faTimes;

  constructor(
    private apollo: Apollo,
    private route: ActivatedRoute,
    private router: Router,
    private store: Store<{ app: AppState }>,
    private title: Title,
  ) {
    this.userInfo$ = store.pipe(select('app', 'auth', 'userInfo'));
  }

  ngOnInit(): void {
    const slug: string | null = this.route.snapshot.paramMap.get('slug');

    if (slug) {
      this.findPostBySlug(slug);
    } else {
      this.createNewPost();
    }
  }

  onChaging(isFetching: boolean): void {
    this.isFetching = isFetching;
  }

  onChangingSuccess(post: Post): void {
    this.post = post;

    this.title.setTitle(`Edit · ${post.title} - ${environment.title}`);
  }

  onChagningErrors(errors: ReadonlyArray<GraphQLError>): void {
    this.errors = errors;
  }

  onSelectingAttachment(attachment: Attachment): void {
    this.selectedAttachment = attachment;
  }

  private createNewPost(): void {
    this.title.setTitle(`Draft a new post - ${environment.title}`);

    this.apollo.mutate({
      mutation: gql`
        mutation {
          createPost {
            slug createdAt
          }
        }
      `,
    }).pipe(
      map((result: ApolloQueryResult<{ createPost: Post }>): Post => result.data.createPost),
    ).subscribe((post: Post): void => {
      const createdAt: Date = new Date(
        new Date(post.createdAt)
          .toLocaleString('en-US', { timeZone: 'Asia/Bangkok' }),
      );

      this.router.navigate([
        createdAt.getFullYear().toString(),
        (createdAt.getMonth() + 1).toString(),
        createdAt.getDate().toString(),
        post.slug,
        'edit',
      ]);
    });
  }

  private findPostBySlug(slug: string): void {
    this.apollo.query({
      query: gql`
        {
          post(slug: $slug) {
            title
            slug
            markdown
            html
            authorId
            categories {
              name slug
            }
            tags {
              name slug
            }
            featuredImage {
              fileName slug
            }
            attachments {
              fileName slug
            }
            createdAt
            updatedAt
          }
        }
      `,
      variables: {
        slug,
      },
    }).pipe(
      map((result: ApolloQueryResult<{ post: Post }>): Post => result.data.post),
      tap((post: Post): void => {
        this.title.setTitle(`Edit · ${post.title || 'Untitled'} - ${environment.title}`);
      }),
    ).subscribe((post: Post): void => {
      this.post = post;
    });
  }
}
