import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { SingleComponent } from './single.component';
import { SingleRoutingModule } from './single-routing.module';

import { GraphQLModule } from '../graphql/graphql.module';
import { SharedModule } from '../shared/share.module';

@NgModule({
  imports: [
    CommonModule,
    GraphQLModule,
    SharedModule,
    SingleRoutingModule,
  ],
  declarations: [
    SingleComponent,
  ],
})
export class SingleModule { }
