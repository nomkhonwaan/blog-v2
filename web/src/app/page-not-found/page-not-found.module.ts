import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { PageNotFoundRoutingModule } from './page-not-found-routing.module';

import { SharedModule } from '../shared/share.module';

@NgModule({
  imports: [
    CommonModule,
    SharedModule,
    PageNotFoundRoutingModule,
  ],
})
export class PageNotFoundModule { }
