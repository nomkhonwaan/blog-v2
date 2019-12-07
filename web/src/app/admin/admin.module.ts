import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';

import { AdminRoutingModule } from './admin-routing.module';
import { AdminComponent } from './admin.component';
import { MyPostsComponent } from './my-posts.component';
import { PostEditorComponent, AutoResizeDirective } from './post-editor.component';

import { GraphQLModule } from '../graphql/graphql.module';
import { SharedModule } from '../shared/share.module';

@NgModule({
  imports: [
    AdminRoutingModule,
    CommonModule,
    GraphQLModule,
    FontAwesomeModule,
    FormsModule,
    SharedModule,
  ],
  declarations: [
    AdminComponent,
    AutoResizeDirective,
    MyPostsComponent,
    PostEditorComponent,
  ],
  bootstrap: [AdminComponent],
})
export class AdminModule { }
