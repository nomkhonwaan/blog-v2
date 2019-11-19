import { Component, ChangeDetectionStrategy } from '@angular/core';

import { PostComponent } from './post.component';

@Component({
  selector: 'app-post-tags',
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <ul class="_list-unstyled">
      <li *ngFor="let t of post.tags">
        <a [routerLink]="['/', 'tag', t.slug]">#{{t.name}}</a>
      </li>
    </ul>
  `,
  styles: [
    `
      li {
          display: inline-block;
          margin-left: 1.6rem;
      }
    `,
    `
      li:first-child {
          margin-left: 0;
      }
    `,
    `
      a {
          color: #0091ea;
          font: normal 300 1.5rem Lato, sans-serif;
      }
    `,
  ],
})
export class PostTagsComponent extends PostComponent { }
