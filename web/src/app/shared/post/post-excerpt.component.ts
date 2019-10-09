import { OnInit, Component } from '@angular/core';

import { PostComponent } from './post.component';

@Component({
  selector: 'app-post-excerpt',
  template: `
    <article [innerHTML]="summary"></article>
  `,
  styles: [
    `
      :host {
          font: normal 300 1.6rem Pridi, sans-serif;
      }
    `
  ],
})
export class PostExcerptComponent extends PostComponent implements OnInit {

  summary: string;

  ngOnInit(): void {
    const paragraphs: string[] = this.post.html.split('</p>');

    this.summary = paragraphs[0] + '</p>';
  }

}
