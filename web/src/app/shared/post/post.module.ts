import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';

import { PostAuthorComponent } from './post-author.component';
import { PostCategoriesComponent } from './post-categories.component';
import { PostExcerptComponent } from './post-excerpt.component';
import { PostMetadataComponent } from './post-metadata.component';
import { PostTitleComponent } from './post-title.component';

@NgModule({
  imports: [
    RouterModule,
  ],
  declarations: [
    PostAuthorComponent,
    PostCategoriesComponent,
    PostExcerptComponent,
    PostMetadataComponent,
    PostTitleComponent,
  ],
  exports: [
    PostAuthorComponent,
    PostCategoriesComponent,
    PostExcerptComponent,
    PostMetadataComponent,
    PostTitleComponent,
  ],
})
export class PostModule { }
