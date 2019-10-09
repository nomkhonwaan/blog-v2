import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';

import { MediumPostComponent } from './medium-post.component';
import { PostAuthorComponent } from './post-author.component';
import { PostCategoriesComponent } from './post-categories.component';
import { PostContentComponent } from './post-content.component';
import { PostPublishedAtComponent } from './post-published-at.component';
import { PostTagsComponent } from './post-tags.component';
import { PostTitleComponent } from './post-title.component';
import { SinglePostComponent } from './single-post.component';
import { ThumbnailPostComponent } from './thumbnail-post.component';

import { TemplateModule } from '../template/template.module';

@NgModule({
  imports: [
    CommonModule,
    RouterModule,
    TemplateModule,
  ],
  declarations: [
    MediumPostComponent,
    PostAuthorComponent,
    PostCategoriesComponent,
    PostContentComponent,
    PostPublishedAtComponent,
    PostTagsComponent,
    PostTitleComponent,
    SinglePostComponent,
    ThumbnailPostComponent,
  ],
  exports: [
    MediumPostComponent,
    SinglePostComponent,
    ThumbnailPostComponent,
  ],
})
export class PostModule { }
