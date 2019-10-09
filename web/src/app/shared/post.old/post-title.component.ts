import { Component, OnInit } from '@angular/core';
import * as moment from 'moment';

import { PostComponent } from './post.component';

@Component({
  selector: 'app-post-title',
  template: `
    <a [routerLink]="link" [ngSwitch]="type">
      <h1 *ngSwitchCase="single">{{post.title}}</h1>

      <h2 *ngSwitchCase="medium">{{post.title}}</h2>

      <h3 *ngSwitchCase="thumbnail">{{post.title}}</h3>
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
      h1 {
          /* font-size: rem; */
      }
    `,
    `
      h2 {
          font-size: 3.6rem;
      }
    `,
    `
      h3 {
          /* font-size: rem; */
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
