import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';

import { RecentPostsComponent } from './recent-posts.component';
import { SingleComponent } from './single.component';

import { ArchiveComponent } from './archive.component';
import { GraphQLModule } from '../graphql/graphql.module';
import { SharedModule } from '../shared/share.module';
import { PageNotFoundComponent } from './page-not-found.component';

@NgModule({
  imports: [
    CommonModule,
    GraphQLModule,
    RouterModule,
    SharedModule,
  ],
  declarations: [
    ArchiveComponent,
    RecentPostsComponent,
    SingleComponent,
    PageNotFoundComponent,
  ],
})
export class PagesModule { }
