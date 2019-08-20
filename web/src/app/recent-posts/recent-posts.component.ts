import { Component, OnInit } from '@angular/core';
import { Store, select } from '@ngrx/store';

import { RecentPostsState } from './recent-posts.reducer';
import { fetchRecentPosts } from './recent-posts.actions';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-recent-posts',
  templateUrl: './recent-posts.component.html',
  styleUrls: ['./recent-posts.component.scss'],
})
export class RecentPostsComponent implements OnInit {

  recentPosts$: Observable<RecentPostsState>;

  constructor(private store: Store<RecentPostsState>) {
    this.recentPosts$ = store.pipe(select('recentPosts'));
  }

  ngOnInit(): void {
    this.store.dispatch(fetchRecentPosts());
  }

}
