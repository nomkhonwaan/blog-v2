import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { MyPostsComponent } from './my-posts.component';
import { PostEditorComponent } from './post-editor.component';
import { AdminComponent } from './admin.component';

const routes: Routes = [
  {
    path: '',
    component: AdminComponent,
    children: [
      { path: '', pathMatch: 'full', component: PostEditorComponent },
      { path: 'new-post', component: PostEditorComponent },
    ]
  },

  // { path: '', component: PostEditorComponent },
  // { path: 'new-post', component: PostEditorComponent },
  // { path: 'posts', component: MyPostsComponent },
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
