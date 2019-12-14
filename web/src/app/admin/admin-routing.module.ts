import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AdminComponent } from './admin.component';
import { EditorComponent } from './editor';
import { MyPostsComponent } from './my-posts';

const routes: Routes = [
  {
    path: '',
    component: AdminComponent,
    children: [
      { path: '', component: EditorComponent },
      { path: 'new-post', pathMatch: 'full', component: EditorComponent },
      { path: 'my-posts', pathMatch: 'full', component: MyPostsComponent },
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
export class AdminRoutingModule { }
