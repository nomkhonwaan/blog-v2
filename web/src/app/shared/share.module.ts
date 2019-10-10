import { NgModule } from '@angular/core';

import { ButtonModule } from './button/button.module';
import { DialogModule } from './dialog/dialog.module';
import { FormModule } from './form/form.module';
import { PostModule } from './post/post.module';
import { TemplateModule } from './template/template.module';

@NgModule({
  exports: [
    ButtonModule,
    DialogModule,
    FormModule,
    PostModule,
    TemplateModule,
  ],
})
export class SharedModule { }
