import { NgModule } from '@angular/core';

import { ButtonModule } from './button/button.module';
import { DialogModule } from './dialog/dialog.module';
import { FormModule } from './form/form.module';

@NgModule({
  imports: [],
  declarations: [],
  exports: [
    ButtonModule,
    DialogModule,
    FormModule,
  ],
  bootstrap: [],
})
export class SharedModule { }
