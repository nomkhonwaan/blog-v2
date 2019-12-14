import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ArchiveComponent } from './archive';
import { ContentComponent } from './content.component';
import { RecentPostsComponent } from './recent-posts';
import { SingleComponent } from './single';
import { PageNotFoundComponent } from './page-not-found';

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
      { path: ':year/:month/:date/:slug', pathMatch: 'full', component: SingleComponent },
      { path: '**', component: PageNotFoundComponent },
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
