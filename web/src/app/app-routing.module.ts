import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { RecentPostsComponent } from './recent-posts/recent-posts.component';
import { CategoryComponent } from './category/category.component';

const routes: Routes = [
  { path: '', component: RecentPostsComponent, pathMatch: 'full' },
  { path: 'category/:slug', component: CategoryComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
