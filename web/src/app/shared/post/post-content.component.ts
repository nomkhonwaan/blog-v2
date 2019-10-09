import { PostComponent } from './post.component';
import { OnInit, Component } from '@angular/core';

@Component({
  selector: 'app-post-content',
  template: `
    <article [innerHTML]="content"></article>
  `,
  styles: [

  ],
})
export class PostContentComponent extends PostComponent implements OnInit {

  content: string;

  ngOnInit(): void {
    const paragraphs: string[] = this.post.html.split('</p>');

    this.content = paragraphs[0] + '</p>';
  }

}
