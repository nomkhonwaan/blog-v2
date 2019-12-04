import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { AdminRoutingModule } from './admin-routing.module';

import { GraphQLModule } from '../graphql/graphql.module';
import { SharedModule } from '../shared/share.module';

@NgModule({
  imports: [
    AdminRoutingModule,
    CommonModule,
    GraphQLModule,
    SharedModule,
  ],
  declarations: [
  ],
})
export class AdminModule {}
