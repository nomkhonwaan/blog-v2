import { OnInit, Component, ChangeDetectionStrategy } from '@angular/core';

import { PostComponent } from './post.component';

@Component({
  selector: 'app-post-excerpt',
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <article [innerHTML]="summary"></article>
  `,
  styleUrls: ['./post-content.component.scss'],
})
export class PostExcerptComponent extends PostComponent implements OnInit {

  summary: string;

  ngOnInit(): void {
    const paragraphs: string[] = this.post.html.split('</p>');

    this.summary = paragraphs[0] + '</p>';
  }

}
