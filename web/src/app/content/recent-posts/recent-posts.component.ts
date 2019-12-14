import { Component, ChangeDetectionStrategy, OnInit } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { Apollo } from 'apollo-angular';
import { environment } from 'src/environments/environment';
import gql from 'graphql-tag';

@Component({
  selector: 'app-recent-posts',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './recent-posts.component.html',
  styleUrls: ['./recent-posts.component.scss'],
})
export class RecentPostsComponent implements OnInit {

  /**
   * List of the latest published posts from all categories
   */
  latestPublishedPosts: Array<Post>;

  constructor(private apollo: Apollo, private title: Title) { }

  ngOnInit(): void {
    this.title.setTitle(environment.title);
  }

}
