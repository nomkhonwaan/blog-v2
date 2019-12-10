import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { ArchiveComponent } from './archive.component';
import { PageNotFoundComponent } from './page-not-found.component';
import { PublishingComponent } from './publishing.component';
import { RecentPostsComponent } from './recent-posts.component';
import { SingleComponent } from './single.component';

const routes: Routes = [
  {
    path: '', component: PublishingComponent,
    children: [
      { path: '', pathMatch: 'full', component: RecentPostsComponent },
      { path: ':page', pathMatch: 'full', component: ArchiveComponent, data: { type: 'all' } },
      { path: 'category/:slug', pathMatch: 'full', component: ArchiveComponent, data: { type: 'category' } },
      { path: 'category/:slug/:page', pathMatch: 'full', component: ArchiveComponent, data: { type: 'category' } },
      { path: 'tag/:slug', pathMatch: 'full', component: ArchiveComponent, data: { type: 'tag' } },
      { path: 'tag/:slug/:page', pathMatch: 'full', component: ArchiveComponent, data: { type: 'tag' } },
      { path: ':year/:month/:date/:slug', pathMatch: 'full', component: SingleComponent },
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
export class PublishingRoutingModule { }
