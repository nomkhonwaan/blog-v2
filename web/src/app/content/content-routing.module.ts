import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { ContentComponent } from './content.component';
import { RecentPostsComponent } from './recent-posts';
import { ArchiveComponent } from './archive';

const routes: Routes = [
  {
    path: '', component: ContentComponent,
    children: [
      { path: '', pathMatch: 'full', component: RecentPostsComponent },
      { path: ':page', pathMatch: 'full', component: ArchiveComponent, data: { from: 'all' } },
      { path: 'category/:slug', pathMatch: 'full', component: ArchiveComponent, data: { from: 'category' } },
      { path: 'category/:slug/:page', pathMatch: 'full', component: ArchiveComponent, data: { from: 'category' } },
      { path: 'tag/:slug', pathMatch: 'full', component: ArchiveComponent, data: { from: 'tag' } },
      { path: 'tag/:slug/:page', pathMatch: 'full', component: ArchiveComponent, data: { from: 'tag' } },
    ],
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
export class ContentRoutingModule { }
