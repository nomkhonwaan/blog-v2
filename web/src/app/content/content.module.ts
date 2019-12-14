import { HTTP_INTERCEPTORS } from '@angular/common/http';
import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';

import { ContentComponent } from './content.component';
import { ContentRoutingModule } from './content-routing.module';
import { RecentPostsComponent } from './recent-posts';

import { AppHttpInterceptor } from '../index';
import { AuthModule } from '../auth';
import { GraphQLModule } from '../graphql';
import { SharedModule } from '../shared';

@NgModule({
  imports: [
    AuthModule,
    CommonModule,
    ContentRoutingModule,
    FontAwesomeModule,
    GraphQLModule,
    SharedModule,
  ],
  declarations: [
    ContentComponent,
    RecentPostsComponent,
  ],
  providers: [
    {
      provide: HTTP_INTERCEPTORS,
      useClass: AppHttpInterceptor,
      multi: true,
    },
  ],
})
export class ContentModule { }
