import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

const routes: Routes = [
  {
    path: ':year/:month/:date/:slug', loadChildren: () => import('./single/single.module')
      .then((m) => m.SingleModule),
  },
  {
    path: 'category/:slug', loadChildren: () => import('./archive/archive.module')
      .then((m) => m.ArchiveModule),
  },
  {
    path: 'tag/:slug', loadChildren: () => import('./archive/archive.module')
      .then((m) => m.ArchiveModule),
  },
  {
    path: '', pathMatch: 'full', loadChildren: () => import('./recent-posts/recent-posts.module')
      .then((m) => m.RecentPostsModule),
  },

  // {
  //   path: '**', loadChildren: () => import('./page-not-found/page-not-found.module')
  //     .then((m) => m.PageNotFoundModule),
  // },
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
