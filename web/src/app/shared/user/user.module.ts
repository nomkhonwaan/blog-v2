import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { RouterModule } from '@angular/router';

import { UserNavbarComponent } from './user-navbar.component';

@NgModule({
  imports: [
    BrowserAnimationsModule,
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
