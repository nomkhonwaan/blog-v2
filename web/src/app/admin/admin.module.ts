import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';

import { AdminRoutingModule } from './admin-routing.module';
import { AdminComponent } from './admin.component';
import { AutoResizeDirective } from './auto-resize.directive';
import { MyPostsComponent } from './my-posts.component';
import { PostEditorComponent } from './post-editor.component';
import { AttachmentsEditorComponent } from './post-editor/attachments-editor.component';
import { AttachmentViewerComponent } from './post-editor/attachment-viewer.component';
import { MarkdownEditorComponent } from './post-editor/markdown-editor.component';
import { StatusEditorComponent } from './post-editor/status-editor.component';
import { TitleEditorComponent } from './post-editor/title-editor.component';

import { ApiModule } from '../api/api.module';
import { GraphQLModule } from '../graphql/graphql.module';
import { SharedModule } from '../shared/share.module';

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
    AutoResizeDirective,
    MarkdownEditorComponent,
    MyPostsComponent,
    PostEditorComponent,
    StatusEditorComponent,
    TitleEditorComponent,
  ],
  bootstrap: [AdminComponent],
})
export class AdminModule { }
