import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { ArchiveComponent } from './pages/archive.component';
import { PageNotFoundComponent } from './pages/page-not-found.component';

const routes: Routes = [
  { path: ':year/:month/:date/:slug', loadChildren: () => import('./single/single.module')
      .then((m) => m.SingleModule) },
  { path: 'category/:slug', component: ArchiveComponent, data: { type: 'category' } },
  { path: 'tag/:slug', component: ArchiveComponent, data: { type: 'tag' } },
  { path: '', pathMatch: 'full', loadChildren: () => import('./recent-posts/recent-posts.module')
      .then((m) => m.RecentPostsModule) },
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
