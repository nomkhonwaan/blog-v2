import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { RecentPostsComponent } from './recent-posts/recent-posts.component';
import { CategoryComponent } from './category/category.component';
import { SingleComponent } from './single/single.component';

const routes: Routes = [
  { path: '', component: RecentPostsComponent, pathMatch: 'full' },
  { path: ':year/:month/:date/:slug', component: SingleComponent },
  { path: 'category/:slug', component: CategoryComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
