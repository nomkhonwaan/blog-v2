import { Component, ChangeDetectionStrategy } from '@angular/core';

import { PostComponent } from './post.component';

@Component({
  selector: 'app-post-categories',
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <ul class="_list-unstyled">
      <li *ngFor="let cat of post.categories">
        <a [routerLink]="['/', 'category', cat.slug]">{{cat.name}}</a>
      </li>
    </ul>
  `,
  styles: [
    `
      :host {
          color: #0091ea;
          font: normal 400 1.6rem Lato, sans-serif;
      }
    `,
    `
      li {
          display: inline-block;
      }
    `,
  ],
})
export class PostCategoriesComponent extends PostComponent { }
