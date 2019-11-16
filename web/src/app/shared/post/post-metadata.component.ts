import { Component, ChangeDetectionStrategy } from '@angular/core';

import { PostComponent } from './post.component';

@Component({
  selector: 'app-post-metadata',
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <app-post-author [post]="post"></app-post-author>

    <div class="published-at">{{ post.publishedAt | dateFormat }} on <app-post-categories [post]="post"></app-post-categories></div>
  `,
  styles: [
    `
      :host {
        display: grid;
        grid-template-columns: 6.4rem auto;
        grid-template-rows: 6.4rem auto;
      }
    `,
    `
      app-post-author {
        grid-column: 1/3;
      }
    `,
    `
      .published-at {
        color: #666;
        font: normal 400 1.6rem Lato, sans-serif;
        grid-column: 2/3;
        grid-row: 2/3;
        margin: -2.8rem 0 0 1.6rem;
      }
    `,
    `
      app-post-categories {
        display: inline-block;
      }
    `,
  ],
})
export class PostMetadataComponent extends PostComponent { }
