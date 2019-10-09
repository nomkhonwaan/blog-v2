import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { CommonModule } from '@angular/common';

import { SingleComponent } from './single.component';

import { SharedModule } from '../shared/share.module';

@NgModule({
  imports: [
    CommonModule,
    RouterModule,
    SharedModule,
  ],
  declarations: [
    SingleComponent,
  ],
})
export class SingleModule { }
