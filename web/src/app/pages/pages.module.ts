import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';

import { RecentPostsComponent } from './recent-posts.component';
import { SingleComponent } from './single.component';

import { GraphQLModule } from '../graphql/graphql.module';
import { SharedModule } from '../shared/share.module';

@NgModule({
  imports: [
    CommonModule,
    GraphQLModule,
    RouterModule,
    SharedModule,
  ],
  declarations: [
    RecentPostsComponent,
    SingleComponent,
  ],
})
export class PagesModule { }
