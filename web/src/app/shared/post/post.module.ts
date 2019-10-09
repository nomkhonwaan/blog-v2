import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';

import { PostAuthorComponent } from '../post.old/post-author.component';
import { PostCategoriesComponent } from './post-categories.component';
import { PostExcerptComponent } from './post-excerpt.component';
import { PostTitleComponent } from './post-title.component';

@NgModule({
  imports: [
    RouterModule,
  ],
  declarations: [
    PostAuthorComponent,
    PostCategoriesComponent,
    PostExcerptComponent,
    PostTitleComponent,
  ],
  exports: [
    PostAuthorComponent,
    PostCategoriesComponent,
    PostExcerptComponent,
    PostTitleComponent,
  ],
})
export class PostModule { }
