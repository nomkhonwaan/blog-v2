import { DOCUMENT } from '@angular/common';
import { Component, OnInit, Directive, ElementRef, HostListener, Inject, AfterViewInit } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { ActivatedRoute } from '@angular/router';
import { faTimes, IconDefinition } from '@fortawesome/pro-light-svg-icons';
import { Apollo } from 'apollo-angular';
import { ApolloQueryResult } from 'apollo-client';
import gql from 'graphql-tag';
import { Observable } from 'rxjs';
import { map, tap } from 'rxjs/operators';

import { environment } from 'src/environments/environment';

@Directive({ selector: '[appAutoResize]' })
export class AutoResizeDirective implements AfterViewInit {

  constructor(@Inject(DOCUMENT) private document: Document, private elementRef: ElementRef) { }

  ngAfterViewInit(): void {
    const elem: HTMLElement = this.elementRef.nativeElement as HTMLElement;

    elem.style.height = 'auto';

    this.resize();
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
    body.style.height = (elem.scrollHeight + 256).toString() + 'px';
  }

}

@Component({
  selector: 'app-post-editor',
  templateUrl: './post-editor.component.html',
  styleUrls: ['./post-editor.component.scss'],
})
export class PostEditorComponent implements OnInit {

  /**
   * A post object
   */
  post: Post;

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
    private title: Title,
  ) { }

  ngOnInit(): void {
    const slug: string = this.route.snapshot.paramMap.get('slug');

    (slug ? this.findPostBySlug(slug) : this.createNewPost()).subscribe((post: Post): void => {
      this.post = post;
    });
  }

  onChangeTitle(): void {
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
    }).pipe(
      map((result: ApolloQueryResult<{ updatePostTitle: Post }>): Post => result.data.updatePostTitle),
    ).subscribe((post: Post): void => {
      this.title.setTitle(`Edit · ${post.title} - ${environment.title}`);
      this.post.slug = post.slug;
    })
  }

  onChangeMarkdown(): void {
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
      }
    }).pipe(
      map((result: ApolloQueryResult<{ updatePostContent: Post }>): Post => result.data.updatePostContent)
    ).subscribe((post: Post): void => {
      this.post.html = post.html;
    });
  }

  private createNewPost(): Observable<Post> {
    this.title.setTitle(`Draft a new post - ${environment.title}`);

    return this.apollo.mutate({
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
    );
  }

  private findPostBySlug(slug: string): Observable<Post> {
    return this.apollo.query({
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
    );
  }

}
