import { NgModule } from '@angular/core';

import { DateFormatPipe } from './date-format.pipe';
import { StringSplitPipe } from './string-split.pipe';

@NgModule({
  declarations: [
    DateFormatPipe,
    StringSplitPipe,
  ],
  exports: [
    DateFormatPipe,
    StringSplitPipe
  ],
})
export class TemplateModule { }
