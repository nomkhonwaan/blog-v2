import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { ButtonComponent } from './button';
import { DropdownButtonComponent } from './dropdown';
import { OutlineButtonComponent } from './outline';


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
