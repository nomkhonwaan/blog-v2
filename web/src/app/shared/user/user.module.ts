import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { UserNavbarComponent } from './navbar';

@NgModule({
  imports: [
    CommonModule,
    RouterModule,
  ],
  declarations: [
    UserNavbarComponent,
  ],
  exports: [
    UserNavbarComponent,
  ],
})
export class UserModule { }
