import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { ArchiveComponent } from './archive.component';
import { LatestPublishedPostsComponent } from './latest-published-posts.component';
import { PageNotFoundComponent } from './page-not-found.component';
import { RecentPostsComponent } from './recent-posts.component';
import { SingleComponent } from './single.component';

const routes: Routes = [
  {
    path: '', component: RecentPostsComponent,
    children: [
      { path: '', pathMatch: 'full', component: LatestPublishedPostsComponent },
      { path: 'category/:slug', component: ArchiveComponent, data: { type: 'category' } },
      { path: 'tag/:slug', component: ArchiveComponent, data: { type: 'tag' } },
      { path: ':year/:month/:date/:slug', component: SingleComponent },
      { path: '**', component: PageNotFoundComponent },
    ]
  },
];

@NgModule({
  imports: [
    RouterModule.forChild(routes),
  ],
  exports: [
    RouterModule,
  ],
})
export class RecentPostsRoutingModule { }
