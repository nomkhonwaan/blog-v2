import { Component, OnInit } from '@angular/core';
import { Apollo } from 'apollo-angular';
import gql from 'graphql-tag';
import { ApolloQueryResult } from 'apollo-client';

@Component({
  selector: 'app-recent-posts',
  templateUrl: './recent-posts.component.html',
  styleUrls: ['./recent-posts.component.scss'],
})
export class RecentPostsComponent implements OnInit {

  latestPublishedPosts: Post[];

  constructor(private apollo: Apollo) { }

  ngOnInit(): void {
    this.apollo.watchQuery({
      query: gql`
        {
          latestPublishedPosts(offset: 0, limit: 5) {
            title
            slug
            html
            publishedAt
            categories {
              name slug
            }
            tags {
              name slug
            }
          }
        }
      `,
    }).valueChanges.subscribe((result: ApolloQueryResult<{ latestPublishedPosts: Post[] }>): void => {
      this.latestPublishedPosts = result.data.latestPublishedPosts
        .map(
          (p: Post): Post => Object.assign(p, {
            tags: [
              { name: 'Donec', slug: 'donec' },
              { name: 'Lorem', slug: 'lorem' },
              { name: 'Curabitur', slug: 'curabitur' },
            ],
          }));
    });
  }

}
