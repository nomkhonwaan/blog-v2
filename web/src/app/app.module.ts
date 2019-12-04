import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { NgModule } from '@angular/core';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { EffectsModule } from '@ngrx/effects';
import { StoreModule } from '@ngrx/store';
import { StoreDevtoolsModule } from '@ngrx/store-devtools';

import { environment } from 'src/environments/environment';

import { AppRoutingModule } from './app-routing.module';
import { AuthModule } from './auth/auth.module';
import { GraphQLModule } from './graphql/graphql.module';
import { SharedModule } from './shared/share.module';

import { AppComponent, AnimationDirective } from './app.component';
import { AppLoginComponent } from './app-login.component';

import * as appReducer from './app.reducer';
import { AppHttpInterceptor } from './app-http.interceptor';

@NgModule({
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    FontAwesomeModule,
    HttpClientModule,
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
    AppRoutingModule,
    AuthModule,
    GraphQLModule,
    SharedModule,
  ],
  declarations: [
    AppComponent,
    AppLoginComponent,
    AnimationDirective,
  ],
  providers: [
    {
      provide: HTTP_INTERCEPTORS,
      useClass: AppHttpInterceptor,
      multi: true,
    },
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
