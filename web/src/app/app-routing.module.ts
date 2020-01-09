import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AuthGuard } from './auth';
import { LoginComponent } from './login.component';
import { LogoutComponent } from './logout.component';

const routes: Routes = [
  {
    path: 'login',
    pathMatch: 'full',
    component: LoginComponent,
  },
  {
    path: 'logout',
    pathMatch: 'full',
    component: LogoutComponent,
  },
  {
    path: 'new-post',
    pathMatch: 'full',
    canActivate: [AuthGuard],
    loadChildren: () => import('./admin').then((m) => m.AdminModule),
  },
  {
    path: ':year/:month/:date/:slug/edit',
    pathMatch: 'full',
    canActivate: [AuthGuard],
    loadChildren: () => import('./admin').then((m) => m.AdminModule),
  },
  {
    path: '',
    loadChildren: () => import('./content').then((m) => m.ContentModule),
  },
];

@NgModule({
  imports: [
    RouterModule.forRoot(routes, {
      anchorScrolling: 'enabled',
      scrollPositionRestoration: 'enabled',
      useHash: false,
    })
  ],
  exports: [
    RouterModule,
  ],
})
export class AppRoutingModule { }
