import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { ArchiveRoutingModule } from './archive-routing.module';

import { GraphQLModule } from '../graphql/graphql.module';
import { SharedModule } from '../shared/share.module';

@NgModule({
  imports: [
    CommonModule,
    GraphQLModule,
    SharedModule,
    ArchiveRoutingModule,
  ],
})
export class ArchiveModule { }
