import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { ApiModule } from '../api/api.module';
import { GraphQLModule } from '../graphql/graphql.module';
import { SharedModule } from '../shared/share.module';
import { AdminRoutingModule } from './admin-routing.module';
import { AdminComponent } from './admin.component';
import { EditorComponent, HeightAsScrollDirective, PostAttachmentsEditorComponent, PostAttachmentViewerComponent, PostMarkdownEditorComponent, PostStatusEditorComponent, PostTitleEditorComponent } from './editor';
import { PostArchivesEditorComponent } from './editor/archives';

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
    EditorComponent,
    HeightAsScrollDirective,
    PostArchivesEditorComponent,
    PostAttachmentsEditorComponent,
    PostAttachmentViewerComponent,
    PostMarkdownEditorComponent,
    PostStatusEditorComponent,
    PostTitleEditorComponent,
  ],
  bootstrap: [AdminComponent],
})
export class AdminModule { }
