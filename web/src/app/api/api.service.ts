import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Store } from '@ngrx/store';
import { Observable } from 'rxjs';

import { environment } from 'src/environments/environment';
import { RecentPostsState } from '../recent-posts/recent-posts.reducer';

@Injectable({
  providedIn: 'root',
})
export class ApiService {

  constructor(private http: HttpClient) { }

  /**
   * Fetches latest published posts from the back-end API
   */
  fetchRecentPosts(): Observable<Post[]> {
    return this.http.get<Post[]>(environment.myblog.url + '/v1/posts');
  }
}
