import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';

import { UserNavbarComponent } from './user-navbar.component';

@NgModule({
  imports: [
    CommonModule,
    RouterModule,
  ],
  exports: [
    UserNavbarComponent,
  ],
  declarations: [
    UserNavbarComponent,
  ],
})
export class UserModule { }
