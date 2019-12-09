import { Input, Output, EventEmitter } from '@angular/core';
import { Apollo } from 'apollo-angular';
import gql from 'graphql-tag';
import { Observable } from 'rxjs';

import { ApiService } from 'src/app/api/api.service';
import { GraphQLError } from 'graphql';
import { ApolloQueryResult } from 'apollo-client';
import { tap, map, finalize } from 'rxjs/operators';

export abstract class EditorComponent {

  @Input()
  post: Post;

  @Output()
  change: EventEmitter<boolean> = new EventEmitter(false);

  @Output()
  changeErrors: EventEmitter<ReadonlyArray<GraphQLError>> = new EventEmitter(null);

  @Output()
  changeSuccess: EventEmitter<Post> = new EventEmitter(null);

  protected fragments: { [name: string]: any } = {
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

  constructor(protected apollo: Apollo, protected api: ApiService) { }

  protected mutate(query: string, variables: { [key: string]: any }): void {
    this.change.emit(true);

    this.apollo.mutate({
      mutation: gql`
        ${query}

        ${this.fragments.post}
      `,
      variables,
    }).pipe(
      tap((result: ApolloQueryResult<any>): void => { this.changeErrors.emit(result.errors); }),
      map((result: ApolloQueryResult<{ [key: string]: Post }>): Post => result.data[Object.keys(result.data)[0]]),
      finalize((): void => { this.change.emit(false); }),
    ).subscribe((post: Post): void => {
      this.changeSuccess.emit(post);
    });
  }

}
