import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';

import { PostTitleComponent } from './post-title.component';
import { PostExcerptComponent } from './post-excerpt.component';

@NgModule({
  imports: [
    RouterModule,
  ],
  declarations: [
    PostTitleComponent,
    PostExcerptComponent,
  ],
  exports: [
    PostTitleComponent,
    PostExcerptComponent,
  ],
})
export class PostModule { }
