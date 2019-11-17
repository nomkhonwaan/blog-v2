import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { SingleComponent } from './single.component';
import { SingleRouterModule } from './single-router.module';

import { GraphQLModule } from '../graphql/graphql.module';
import { SharedModule } from '../shared/share.module';

@NgModule({
  imports: [
    CommonModule,
    GraphQLModule,
    SharedModule,
    SingleRouterModule,
  ],
  declarations: [
    SingleComponent,
  ],
})
export class SingleModule { }
