import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { RouterModule } from '@angular/router';

import { RecentPostsComponent } from './recent-posts.component';

@NgModule({
  declarations: [
    RecentPostsComponent,
  ],
  imports: [
    BrowserModule,
    RouterModule,
  ],
  providers: [],
  bootstrap: [],
})
export class RecentPostsModule { }
