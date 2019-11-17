import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { RecentPostsComponent } from './recent-posts.component';

const routes: Routes = [
  { path: '', component: RecentPostsComponent },
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
