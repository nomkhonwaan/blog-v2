import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';

import { PostAuthorComponent } from './post-author.component';
import { PostCategoriesComponent } from './post-categories.component';
import { PostExcerptComponent } from './post-excerpt.component';
import { PostMetadataComponent } from './post-metadata.component';
import { PostTagsComponent } from './post-tags.component';
import { PostTitleComponent } from './post-title.component';

@NgModule({
  imports: [
    CommonModule,
    RouterModule,
  ],
  declarations: [
    PostAuthorComponent,
    PostCategoriesComponent,
    PostExcerptComponent,
    PostMetadataComponent,
    PostTagsComponent,
    PostTitleComponent,
  ],
  exports: [
    PostAuthorComponent,
    PostCategoriesComponent,
    PostExcerptComponent,
    PostMetadataComponent,
    PostTagsComponent,
    PostTitleComponent,
  ],
})
export class PostModule { }
