import { NgModule } from '@angular/core';

import { StringSplit } from './string-split.pipe';

@NgModule({
  declarations: [
    StringSplit,
  ],
  exports: [
    StringSplit,
  ],
})
export class TemplateModule { }
