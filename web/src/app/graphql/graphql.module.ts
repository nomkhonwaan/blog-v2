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

export function createApollo(httpLink: HttpLink, store: Store<{ app: AppState }>) {
  const auth$: Observable<{ accessToken: string } > = store.pipe(select('app', 'auth'));

  const authContext: ApolloLink = setContext((operation: GraphQLRequest, prevContext: any): Promise<any> | any => {
    return auth$
      .pipe(take(1))
      .toPromise()
      .then((auth?: { accessToken?: string }) => {
        return auth && auth.accessToken ? {
          headers: {
            Authorization: `Bearer ${auth.accessToken}`,
          },
        } : {};
      });
  });

  return {
    link: ApolloLink.from([authContext, httpLink.create({ uri, withCredentials: true })]),
    cache: new InMemoryCache(),
  };
}

@NgModule({
  exports: [
    ApolloModule,
    HttpLinkModule,
  ],
  providers: [
    {
      provide: APOLLO_OPTIONS,
      useFactory: createApollo,
      deps: [HttpLink, Store],
    },
  ],
})
export class GraphQLModule { }
