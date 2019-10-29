import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { ArchiveComponent } from './pages/archive.component';
import { RecentPostsComponent } from './pages/recent-posts.component';
import { SingleComponent } from './pages/single.component';
import { PageNotFoundComponent } from './pages/page-not-found.component';

const routes: Routes = [
  { path: ':year/:month/:date/:slug', component: SingleComponent },
  { path: 'category/:slug', component: ArchiveComponent, data: { type: 'category' } },
  { path: 'tag/:slug', component: ArchiveComponent, data: { type: 'tag' } },
  { path: '', component: RecentPostsComponent, pathMatch: 'full' },
  { path: '**', component: PageNotFoundComponent },
];

@NgModule({
  imports: [
    RouterModule.forRoot(routes, {
      scrollPositionRestoration: 'enabled',
    })
  ],
  exports: [
    RouterModule,
  ],
})
export class AppRoutingModule { }
