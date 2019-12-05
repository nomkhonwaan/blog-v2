import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { LoginComponent } from './login.component';

const routes: Routes = [
  {
    path: '', pathMatch: 'full',
    loadChildren: () => import('./recent-posts/recent-posts.module').then((m) => m.RecentPostsModule),
  },
  {
    path: 'category/:slug',
    data: { type: 'category' },
    loadChildren: () => import('./archive/archive.module').then((m) => m.ArchiveModule),
  },
  {
    path: 'tag/:slug',
    data: { type: 'tag' },
    loadChildren: () => import('./archive/archive.module').then((m) => m.ArchiveModule),
  },
  {
    path: 'login',
    component: LoginComponent,
  },
  {
    path: 'admin',
    loadChildren: () => import('./admin/admin.module').then((m) => m.AdminModule),
  },
  {
    path: ':year/:month/:date/:slug',
    loadChildren: () => import('./single/single.module').then((m) => m.SingleModule),
  },
  {
    path: '**',
    loadChildren: () => import('./page-not-found/page-not-found.module').then((m) => m.PageNotFoundModule),
  },
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
