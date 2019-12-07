import { DOCUMENT } from '@angular/common';
import { Component, OnInit, Directive, ElementRef, HostListener, Inject } from '@angular/core';
import { Apollo } from 'apollo-angular';
import { ActivatedRoute } from '@angular/router';
import { faTimes, IconDefinition } from '@fortawesome/pro-light-svg-icons';
import gql from 'graphql-tag';
import { Observable } from 'rxjs';
import { map, tap } from 'rxjs/operators';
import { ApolloQueryResult } from 'apollo-client';
import { Title } from '@angular/platform-browser';

@Directive({ selector: '[autoResize]' })
export class AutoResizeDirective {

  constructor(@Inject(DOCUMENT) private document: Document, private elementRef: ElementRef) { }

  @HostListener('change')
  onChange(): void {
    this.resize();
  }

  @HostListener('document:keypress', ['$event'])
  onKeyPress(event: KeyboardEvent): void {
    this.resize();
  }

  private resize(): void {
    const elem: HTMLElement = <HTMLElement>this.elementRef.nativeElement;
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

  private createNewPost(): Observable<Post> {
    this.title.setTitle('Draft a new post - Nomkhonwaan | Trust me I\'m Petdo');

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
        this.title.setTitle(`Edit · ${post.title || 'Untitled'} - Nomkhonwaan | Trust me I'm Petdo`);
      }),
    );
  }

}
