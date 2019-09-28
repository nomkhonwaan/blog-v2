import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';

import { ButtonComponent } from './button.component';

@NgModule({
  imports: [
    CommonModule,
    FontAwesomeModule,
  ],
  declarations: [
    ButtonComponent,
  ],
  exports: [
    ButtonComponent,
  ],
  bootstrap: []
})
export class ButtonModule { }
