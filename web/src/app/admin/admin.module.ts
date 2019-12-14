import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { ApiModule } from '../api/api.module';
import { GraphQLModule } from '../graphql/graphql.module';
import { SharedModule } from '../shared/share.module';
import { AdminRoutingModule } from './admin-routing.module';
import { AdminComponent } from './admin.component';
import { EditorComponent } from './editor';
import { HeightAsScrollDirective } from './editor/height-as-scroll.directive';
import { PostTitleEditorComponent } from './editor/title';
import { MyPostsComponent } from './my-posts';
import { AttachmentViewerComponent } from './post-editor/attachment-viewer.component';
import { AttachmentsEditorComponent } from './post-editor/attachments-editor.component';
import { MarkdownEditorComponent } from './post-editor/markdown-editor.component';
import { StatusEditorComponent } from './post-editor/status-editor.component';

@NgModule({
  imports: [
    AdminRoutingModule,
    ApiModule,
    CommonModule,
    GraphQLModule,
    FontAwesomeModule,
    FormsModule,
    SharedModule,
  ],
  declarations: [
    AdminComponent,
    AttachmentsEditorComponent,
    AttachmentViewerComponent,
    EditorComponent,
    HeightAsScrollDirective,
    MarkdownEditorComponent,
    MyPostsComponent,
    StatusEditorComponent,
    PostTitleEditorComponent,
  ],
  bootstrap: [AdminComponent],
})
export class AdminModule { }
