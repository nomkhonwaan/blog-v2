import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';

import { ButtonComponent } from './button.component';
import { DropdownButtonComponent } from './dropdown-button.component';
import { OutlineButtonComponent } from './outline-button.component';

@NgModule({
  imports: [
    CommonModule,
    FontAwesomeModule,
  ],
  declarations: [
    ButtonComponent,
    DropdownButtonComponent,
    OutlineButtonComponent,
  ],
  exports: [
    ButtonComponent,
    DropdownButtonComponent,
    OutlineButtonComponent,
  ],
})
export class ButtonModule { }
