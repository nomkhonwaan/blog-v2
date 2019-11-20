import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { PageNotFoundComponent } from './page-not-found.component';
import { PageNotFoundRoutingModule } from './page-not-found-routing.module';

import { SharedModule } from '../shared/share.module';

@NgModule({
  imports: [
    CommonModule,
    SharedModule,
    PageNotFoundRoutingModule,
  ],
  declarations: [
    PageNotFoundComponent,
  ],
})
export class PageNotFoundModule { }
