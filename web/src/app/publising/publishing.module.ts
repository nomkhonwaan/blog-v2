import { CommonModule } from '@angular/common';
import { HTTP_INTERCEPTORS } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';

import { ArchiveComponent } from './archive.component';
import { LottieAnimationDirective } from './lottie-animation.directive';
import { PageNotFoundComponent } from './page-not-found.component';
import { PublishingComponent } from './publishing.component';
import { PublishingRoutingModule } from './publishing-routing.module';
import { RecentPostsComponent } from './recent-posts.component';
import { SingleComponent } from './single.component';

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
    PublishingRoutingModule,
  ],
  declarations: [
    ArchiveComponent,
    LottieAnimationDirective,
    PageNotFoundComponent,
    PublishingComponent,
    RecentPostsComponent,
    SingleComponent,
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
export class PublishingModule { }
