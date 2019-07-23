import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { RecentPostsComponent } from './recent-posts/recent-posts.component';

const routes: Routes = [
  { path: '', component: RecentPostsComponent, pathMatch: 'full' },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
