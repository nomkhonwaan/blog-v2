import { Component, OnInit, Directive, AfterViewInit, ElementRef } from '@angular/core';
import { Apollo } from 'apollo-angular';
import { ActivatedRoute } from '@angular/router';
import { faTimes, IconDefinition } from '@fortawesome/pro-light-svg-icons';
import gql from 'graphql-tag';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { ApolloQueryResult } from 'apollo-client';

@Directive({ selector: '[autoResize]' })
export class AutoResizeDirective implements AfterViewInit {

  constructor(private elementRef: ElementRef) { }

  ngAfterViewInit(): void {

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
  ) { }

  ngOnInit(): void {
    const slug: string = this.route.snapshot.paramMap.get('slug');

    (slug ? this.findPostBySlug(slug) : this.createNewPost()).subscribe((post: Post): void => {
      this.post = post;
    });
  }

  private createNewPost(): Observable<Post> {
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
    );
  }

}
