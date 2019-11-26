import { NgModule } from '@angular/core';
import { Store, select } from '@ngrx/store';
import { ApolloModule, APOLLO_OPTIONS } from 'apollo-angular';
import { HttpLinkModule, HttpLink } from 'apollo-angular-link-http';
import { InMemoryCache } from 'apollo-cache-inmemory';
import { GraphQLRequest, ApolloLink } from 'apollo-link';
import { setContext } from 'apollo-link-context';

import { environment } from 'src/environments/environment';

const uri = environment.graphql.endpoint;

export function createApollo(store: Store<AppState>, httpLink: HttpLink) {
  // let bearer: string;

  // store.pipe(select('app', 'auth', 'accessToken')).subscribe((accessToken: string): void => {
  //   if (accessToken !== '') {
  //     bearer = `Bearer ${accessToken}`;
  //   }
  // });

  // const auth: ApolloLink = setContext((operation: GraphQLRequest, prevContext: any): Promise<any> | any => {
  //   if (bearer !== '') {
  //     return {
  //       headers: {
  //         Authorization: bearer,
  //       },
  //     };
  //   } else {
  //     return {
  //       headers: {},
  //     };
  //   }
  // });

  return {
    link: ApolloLink.from([httpLink.create({ uri, withCredentials: true })]),
    cache: new InMemoryCache(),
  };
}

@NgModule({
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
