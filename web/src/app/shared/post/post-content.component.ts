import { Component, OnInit } from '@angular/core';

import { PostComponent } from './post.component';

@Component({
  selector: 'app-post-content',
  template: `
    <article [innerHTML]="content"></article>
  `,
  styleUrls: ['./post-content.component.scss'],
})
export class PostContentComponent extends PostComponent implements OnInit {

  content: string;

  ngOnInit(): void {
    this.content = this.post.html.replace(new RegExp('/api/v1/attachments', 'g'), 'https://www.nomkhonwaan.com/api/v1/attachments');
  }
}
