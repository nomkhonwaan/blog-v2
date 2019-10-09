import { NgModule } from '@angular/core';

import { SinglePostComponent } from './single-post.component';

@NgModule({
  declarations: [
    SinglePostComponent,
  ],
  exports: [
    SinglePostComponent,
  ],
})
export class PostModule { }
