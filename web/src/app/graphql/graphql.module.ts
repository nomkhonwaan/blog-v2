import { HTTP_INTERCEPTORS } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { ApolloModule, APOLLO_OPTIONS } from 'apollo-angular';
import { HttpLink, HttpLinkModule } from 'apollo-angular-link-http';
import { InMemoryCache } from 'apollo-cache-inmemory';
import { ApolloLink } from 'apollo-link';
import { environment } from 'src/environments/environment';
import { AppHttpInterceptor } from '../app-http.interceptor';

const uri = environment.graphql.endpoint;

export function createApollo(httpLink: HttpLink) {
  return {
    link: ApolloLink.from([httpLink.create({ uri, withCredentials: true })]),
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
      deps: [HttpLink],
    },
    {
      provide: HTTP_INTERCEPTORS,
      useClass: AppHttpInterceptor,
      multi: true,
    },
  ],
})
export class GraphQLModule { }
