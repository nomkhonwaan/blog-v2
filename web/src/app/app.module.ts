import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { NgModule } from '@angular/core';
import { EffectsModule } from '@ngrx/effects';
import { StoreModule } from '@ngrx/store';
import { StoreDevtoolsModule } from '@ngrx/store-devtools';
import { APOLLO_OPTIONS } from 'apollo-angular';
import { HttpLink } from 'apollo-angular-link-http';

import { environment } from 'src/environments/environment';

import { AppRoutingModule } from './app-routing.module';
import { AuthModule } from './auth/auth.module';
import { GraphQLModule, createApollo } from './graphql.module';
import { CategoryModule } from './category/category.module';
import { RecentPostsModule } from './recent-posts/recent-posts.module';
import { SharedModule } from './shared/share.module';

import { AppComponent } from './app.component';

import * as appReducer from './app.reducer';
import { AppHttpInterceptor } from './app-http.interceptor';

@NgModule({
  declarations: [
    AppComponent,
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    AppRoutingModule,
    StoreModule.forRoot({ app: appReducer.reducer }, {
      runtimeChecks: {
        strictActionImmutability: true,
        strictStateImmutability: true,
      }
    }),
    EffectsModule.forRoot([]),
    StoreDevtoolsModule.instrument({
      maxAge: 25,
      logOnly: environment.production,
    }),
    BrowserAnimationsModule,
    AuthModule,
    SharedModule,
    CategoryModule,
    RecentPostsModule,
    GraphQLModule,
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
  bootstrap: [AppComponent]
})
export class AppModule { }
