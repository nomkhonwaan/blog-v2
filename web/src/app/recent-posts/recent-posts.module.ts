import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { RouterModule } from '@angular/router';
import { EffectsModule } from '@ngrx/effects';

import { RecentPostsComponent } from './recent-posts.component';
import { RecentPostsEffects } from './recent-posts.effects';
import { StoreModule } from '@ngrx/store';
import { recentPostsReducer } from './recent-posts.reducer';

@NgModule({
  declarations: [
    RecentPostsComponent,
  ],
  imports: [
    BrowserModule,
    RouterModule,
    StoreModule.forFeature('recentPosts', recentPostsReducer),
    EffectsModule.forFeature([RecentPostsEffects]),
  ],
  providers: [],
  bootstrap: [],
})
export class RecentPostsModule { }
