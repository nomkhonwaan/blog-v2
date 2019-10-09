import { Component, OnInit } from '@angular/core';
import * as moment from 'moment';

import { PostComponent } from './post.component';

@Component({
  selector: 'app-post-title',
  template: `
    <a [routerLink]="link">
      <ng-content></ng-content>
    </a>
  `,
  styles: [
    `
      :host {
          color: #333;
          font-family: Athiti, sans-serif;
          font-weight: 500;
      }
    `,
    `
      ::ng-deep h2 {
          font-size: 3.6rem;
      }
    `,
  ],
})
export class PostTitleComponent extends PostComponent implements OnInit {

  link: string[];

  ngOnInit(): void {
    const publishedAt: moment.Moment = moment(this.post.publishedAt);

    this.link = [
      publishedAt.format('/YYYY/MM/DD'),
      this.post.slug,
    ];
  }

}
