import { NgModule } from '@angular/core';
import { ButtonModule } from './button';
import { DialogModule } from './dialog';
import { PostModule } from './post/post.module';
import { TemplateModule } from './template/template.module';
import { UserModule } from './user/user.module';

@NgModule({
  exports: [
    ButtonModule,
    DialogModule,
    PostModule,
    TemplateModule,
    UserModule,
  ],
})
export class SharedModule { }
