import { CommonModule } from '@angular/common';
import { HTTP_INTERCEPTORS } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';

import { ArchiveComponent } from './archive.component';
import { LatestPublishedPostsComponent } from './latest-published-posts.component';
import { LottieAnimationDirective } from './lottie-animation.directive';
import { RecentPostsComponent } from './recent-posts.component';
import { RecentPostsRoutingModule } from './recent-posts-routing.module';

import { AppHttpInterceptor } from '../app-http.interceptor';
import { AuthModule } from '../auth/auth.module';
import { GraphQLModule } from '../graphql/graphql.module';
import { SharedModule } from '../shared/share.module';

@NgModule({
  imports: [
    AuthModule,
    CommonModule,
    FontAwesomeModule,
    GraphQLModule,
    SharedModule,
    RecentPostsRoutingModule,
  ],
  declarations: [
    ArchiveComponent,
    LatestPublishedPostsComponent,
    LottieAnimationDirective,
    RecentPostsComponent,
  ],
  providers: [
    {
      provide: HTTP_INTERCEPTORS,
      useClass: AppHttpInterceptor,
      multi: true,
    },
  ],
  bootstrap: [RecentPostsComponent]
})
export class RecentPostsModule { }
