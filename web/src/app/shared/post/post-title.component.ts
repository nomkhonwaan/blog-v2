import { Component, OnInit } from '@angular/core';
import * as moment from 'moment';

import { PostComponent } from './post.component';

@Component({
  selector: 'app-post-title',
  template: `test`,
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
