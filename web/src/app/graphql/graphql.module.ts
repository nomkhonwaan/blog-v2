import { NgModule } from '@angular/core';
import { Store, select } from '@ngrx/store';
import { ApolloModule, APOLLO_OPTIONS } from 'apollo-angular';
import { HttpLinkModule, HttpLink } from 'apollo-angular-link-http';
import { InMemoryCache } from 'apollo-cache-inmemory';
import { GraphQLRequest, ApolloLink } from 'apollo-link';
import { setContext } from 'apollo-link-context';
import { Observable } from 'rxjs';
import { take } from 'rxjs/operators';

import { environment } from 'src/environments/environment';

const uri = environment.graphql.endpoint;

export function createApollo(store: Store<AppState>, httpLink: HttpLink) {
  const accessToken$: Observable<string> = store.pipe(select('app', 'auth', 'accessToken'));

  const auth: ApolloLink = setContext((operation: GraphQLRequest, prevContext: any): Promise<any> | any => {
    return accessToken$
      .pipe(take(1))
      .toPromise()
      .then((accessToken: string): { headers?: { [Authorization:string]: string } } => {
        return accessToken ? {
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
        } : {};
      });
  });

  return {
    link: ApolloLink.from([auth, httpLink.create({ uri, withCredentials: true })]),
    cache: new InMemoryCache(),
  };
}

@NgModule({
  imports: [],
  exports: [ApolloModule, HttpLinkModule],
  providers: [
    {
      provide: APOLLO_OPTIONS,
      useFactory: createApollo,
      deps: [Store, HttpLink],
    },
  ],
})
export class GraphQLModule { }
