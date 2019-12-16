import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { AbstractPostComponent } from '../abstract-post.component';

@Component({
  selector: 'app-post-excerpt',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './excerpt.component.html',
  styleUrls: ['../content/content.component.scss'],
})
export class PostExcerptComponent extends AbstractPostComponent implements OnInit {

  /**
   * A first paragraph of the post for using as a summary
   */
  content: string;

  ngOnInit(): void {
    const paragraphs: string[] = this.post.html.split('</p>');

    this.content = paragraphs[0] + '</p>';
  }

}
