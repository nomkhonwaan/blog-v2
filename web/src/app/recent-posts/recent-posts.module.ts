import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { RecentPostsComponent } from './recent-posts.component';
import { RecentPostsRoutingModule } from './recent-posts-routing.module';

import { GraphQLModule } from '../graphql/graphql.module';
import { SharedModule } from '../shared/share.module';

@NgModule({
  imports: [
    CommonModule,
    GraphQLModule,
    SharedModule,
    RecentPostsRoutingModule,
  ],
  declarations: [
    RecentPostsComponent,
  ],
})
export class RecentPostsModule {}
