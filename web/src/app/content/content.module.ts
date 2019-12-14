import { CommonModule } from '@angular/common';
import { HTTP_INTERCEPTORS } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { AuthModule } from '../auth';
import { GraphQLModule } from '../graphql';
import { AppHttpInterceptor } from '../index';
import { SharedModule } from '../shared';
import { ArchiveComponent } from './archive';
import { ContentRoutingModule } from './content-routing.module';
import { ContentComponent } from './content.component';
import { RecentPostsComponent } from './recent-posts';

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
    ArchiveComponent,
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
