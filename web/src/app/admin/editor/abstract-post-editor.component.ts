import { EventEmitter, Input, Output } from '@angular/core';
import { Apollo } from 'apollo-angular';
import { ApolloQueryResult } from 'apollo-client';
import { GraphQLError } from 'graphql';
import gql from 'graphql-tag';
import { Observable } from 'rxjs';
import { finalize, map, tap } from 'rxjs/operators';
import { ApiService } from 'src/app/api/api.service';

export abstract class AbstractPostEditorComponent {

  /**
   * A post object
   */
  @Input()
  post: Post;

  /**
   * For emitting on updating event
   */
  @Output()
  changing: EventEmitter<boolean> = new EventEmitter(false);

  /**
   * For emitting a GraphQL error(s) response
   */
  @Output()
  changeErrors: EventEmitter<ReadonlyArray<GraphQLError>> = new EventEmitter(null);

  /**
   * For emitting an updated post object
   */
  @Output()
  changeSuccess: EventEmitter<Post> = new EventEmitter(null);

  constructor(protected apollo: Apollo, protected api: ApiService) { }

  protected mutate(query: string, variables: { [key: string]: any }): Observable<Post> {
    this.changing.emit(true);

    return this.apollo.mutate({
      mutation: gql`
        ${query}

        fragment EditablePost on Post {
          title slug
          status
          markdown html
          publishedAt
          authorId
          categories { name slug }
          tags { name slug }
          featuredImage { slug }
          attachments { fileName slug }
          createdAt updatedAt
        }
      `,
      variables,
    }).pipe(
      tap((result: ApolloQueryResult<any>): void => { this.changeErrors.emit(result.errors); }),
      map((result: ApolloQueryResult<{ [key: string]: Post }>): Post => result.data[Object.keys(result.data)[0]]),
      finalize((): void => { this.changing.emit(false); }),
    );
  }

}
