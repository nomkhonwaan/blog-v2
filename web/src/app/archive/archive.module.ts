import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { ArchiveComponent } from './archive.component';
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
  declarations: [
    ArchiveComponent,
  ],
})
export class ArchiveModule { }
