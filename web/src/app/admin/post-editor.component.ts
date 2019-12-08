import { trigger } from '@angular/animations';
import { DOCUMENT } from '@angular/common';
import { Component, OnInit, Directive, ElementRef, HostListener, Inject, AfterViewInit } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { ActivatedRoute, Router } from '@angular/router';
import { faImage, faSpinnerThird, faTimes, IconDefinition } from '@fortawesome/pro-light-svg-icons';
import { Apollo } from 'apollo-angular';
import { ApolloQueryResult } from 'apollo-client';
import { GraphQLError } from 'graphql';
import gql from 'graphql-tag';
import { map, tap } from 'rxjs/operators';

import { environment } from 'src/environments/environment';
import { BehaviorSubject } from 'rxjs';

@Directive({ selector: '[appAutoResize]' })
export class AutoResizeDirective implements AfterViewInit {

  constructor(@Inject(DOCUMENT) private document: Document, private elementRef: ElementRef) { }

  ngAfterViewInit(): void {
    setTimeout(() => this.resize());
  }

  @HostListener('change')
  onChange(): void {
    this.resize();
  }

  @HostListener('document:keypress', ['$event'])
  onKeyPress(event: KeyboardEvent): void {
    this.resize();
  }

  private resize(): void {
    const elem: HTMLElement = this.elementRef.nativeElement as HTMLElement;
    const body: HTMLElement = this.document.body;

    elem.style.height = elem.scrollHeight.toString() + 'px';
    body.style.height = body.scrollHeight.toString() + 'px';
  }

}

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
   * Use to display when GraphQL returns
   */
  isErrors$: BehaviorSubject<ReadonlyArray<GraphQLError>> = new BehaviorSubject(null);
  isFetching$: BehaviorSubject<boolean> = new BehaviorSubject(false);
  isUploadingAttachments = false;
  isUploadingAttachments$: BehaviorSubject<boolean> = new BehaviorSubject(false);

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
          slug
        }
        createdAt
        updatedAt
      }
    `,
  };

  constructor(
    private apollo: Apollo,
    private route: ActivatedRoute,
    private router: Router,
    private title: Title,
  ) { }

  ngOnInit(): void {
    const slug: string | null = this.route.snapshot.paramMap.get('slug');

    if (slug) {
      this.findPostBySlug(slug);
    } else {
      this.createNewPost();
    }
  }

  onChangeTitle(): void {
    this.isFetching$.next(true);

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
      tap((result: ApolloQueryResult<any>): void => { this.isErrors$.next(result.errors) }),
      map((result: ApolloQueryResult<{ updatePostTitle: Post }>): Post => result.data.updatePostTitle),
      tap((_: Post): void => { this.isFetching$.next(false); }),
    ).subscribe((post: Post): void => {
      this.title.setTitle(`Edit · ${post.title} - ${environment.title}`);
      this.post.slug = post.slug;
    })
  }

  onChangeMarkdown(): void {
    this.isFetching$.next(true);

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
      tap((result: ApolloQueryResult<any>): void => { this.isErrors$.next(result.errors) }),
      map((result: ApolloQueryResult<{ updatePostContent: Post }>): Post => result.data.updatePostContent),
      tap((_: Post): void => { this.isFetching$.next(false); }),
    ).subscribe((post: Post): void => {
      this.post.html = post.html;
    });
  }

  onChangeAttachments(files: FileList): void {
    console.log(files);
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

}
