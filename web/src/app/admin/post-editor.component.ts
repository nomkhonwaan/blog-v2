import { Component, OnInit } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { ActivatedRoute, Router } from '@angular/router';
import { faImage, faSpinnerThird, faTimes, IconDefinition } from '@fortawesome/pro-light-svg-icons';
import { Apollo } from 'apollo-angular';
import { ApolloQueryResult } from 'apollo-client';
import { GraphQLError } from 'graphql';
import gql from 'graphql-tag';
import { map, tap, finalize } from 'rxjs/operators';

import { ApiService } from '../api/api.service';

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
   * Use to prevent a new upload request while performing
   */
  isUploadingAttachments = false;

  /**
   * An authenticated user information
   */
  userInfo$: Observable<UserInfo | null>;

  /**
   * To-be updated attachment
   */
  selectedAttachment: Attachment;

  faImage: IconDefinition = faImage;
  faSpinnerThird: IconDefinition = faSpinnerThird;
  faTimes: IconDefinition = faTimes;

  fragments: { [name: string]: any } = {
    post: gql`
      fragment EditablePost on Post {
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
    `,
  };

  constructor(
    private api: ApiService,
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

  onChangeTitle(): void {
    this.isFetching = true;

    this.apollo.mutate({
      mutation: gql`
        mutation {
          updatePostTitle(slug: $slug, title: $title) {
            ...EditablePost
          }
        }

        ${this.fragments.post}
      `,
      variables: {
        slug: this.post.slug,
        title: this.post.title,
      },
      errorPolicy: 'all',
    }).pipe(
      tap((result: ApolloQueryResult<any>): void => { this.errors = result.errors; }),
      map((result: ApolloQueryResult<{ updatePostTitle: Post }>): Post => result.data.updatePostTitle),
      finalize((): void => { this.isFetching = false; }),
    ).subscribe((post: Post): void => {
      this.title.setTitle(`Edit · ${post.title} - ${environment.title}`);
      this.post.slug = post.slug;
    });
  }

  onChangeMarkdown(): void {
    this.isFetching = true;

    this.apollo.mutate({
      mutation: gql`
        mutation {
          updatePostContent(slug: $slug, markdown: $markdown) {
            ...EditablePost
          }
        }

        ${this.fragments.post}
      `,
      variables: {
        slug: this.post.slug,
        markdown: this.post.markdown,
      },
      errorPolicy: 'all',
    }).pipe(
      tap((result: ApolloQueryResult<any>): void => { this.errors = result.errors; }),
      map((result: ApolloQueryResult<{ updatePostContent: Post }>): Post => result.data.updatePostContent),
      finalize((): void => { this.isFetching = false }),
    ).subscribe((post: Post): void => {
      this.post.html = post.html;
    });
  }

  onChangeAttachments(files: FileList): void {
    this.isUploadingAttachments = true;

    forkJoin(
      Array.
        from(files).
        map((file: File): Observable<Attachment> => this.api.uploadFile(file)),
    ).pipe(
      finalize((): void => { this.isUploadingAttachments = false; }),
    ).subscribe((attachments: Attachment[]): void => {
      this.updatePostAttachments(this.post.attachments.concat(attachments));
    });
  }

  private createNewPost(): void {
    this.title.setTitle(`Draft a new post - ${environment.title}`);

    this.apollo.mutate({
      mutation: gql`
        mutation {
          createPost {
            ...EditablePost
          }
        }

        ${this.fragments.post}
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
            ...EditablePost
          }
        }

        ${this.fragments.post}
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

  private updatePostAttachments(attachments: Attachment[]): void {
    this.isFetching = true;

    this.apollo.mutate({
      mutation: gql`
        mutation {
          updatePostAttachments(slug: $slug, attachmentSlugs: $attachmentSlugs) {
            ...EditablePost
          }
        }

        ${this.fragments.post}
      `,
      variables: {
        slug: this.post.slug,
        attachmentSlugs: attachments.map((attachment: Attachment): string => attachment.slug),
      },
    }).pipe(
      tap((result: ApolloQueryResult<any>): void => { this.errors = result.errors; }),
      map((result: ApolloQueryResult<{ updatePostAttachments: Post }>): Post => result.data.updatePostAttachments),
      finalize((): void => { this.isFetching = false; }),
    ).subscribe((post: Post): void => {
      this.post.attachments = post.attachments;
    });
  }

}
