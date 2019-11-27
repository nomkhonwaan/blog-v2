import { NgModule } from '@angular/core';

import { DateFormatPipe } from './date-format.pipe';
import { SafeHtmlPipe } from './safe-html.pipe';
import { StringSplitPipe } from './string-split.pipe';

@NgModule({
  declarations: [
    DateFormatPipe,
    SafeHtmlPipe,
    StringSplitPipe,
  ],
  exports: [
    DateFormatPipe,
    SafeHtmlPipe,
    StringSplitPipe
  ],
})
export class TemplateModule { }
