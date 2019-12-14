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
import { LottieDirective } from './lottie.directive';
import { PageNotFoundComponent } from './page-not-found';
import { RecentPostsComponent } from './recent-posts';
import { SingleComponent } from './single';

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
    LottieDirective,
    ArchiveComponent,
    ContentComponent,
    RecentPostsComponent,
    SingleComponent,
    PageNotFoundComponent
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
