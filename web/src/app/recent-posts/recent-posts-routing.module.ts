import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { ArchiveComponent } from './archive.component';
import { LatestPublishedPostsComponent } from './latest-published-posts.component';
import { RecentPostsComponent } from './recent-posts.component';

const routes: Routes = [
  {
    path: '', component: RecentPostsComponent,
    children: [
      { path: '', pathMatch: 'full', component: LatestPublishedPostsComponent },
      { path: 'category/:slug', component: ArchiveComponent, data: { type: 'category' } },
      { path: 'tag/:slug', component: ArchiveComponent, data: { type: 'tag' } },
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
