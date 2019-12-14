import { NgModule } from '@angular/core';
import { Routes } from '@angular/router';

import { ContentComponent } from './index';

const routes: Routes = [
  {
    path: '', component: ContentComponent,
  },
];

@NgModule({

})
export class ContentRoutingModule { }
