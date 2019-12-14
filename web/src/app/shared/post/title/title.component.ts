import { ChangeDetectionStrategy, Component, Input, OnInit } from '@angular/core';
import { PostComponent } from '../post.component';

@Component({
  selector: 'app-post-title',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './title.component.html',
  styleUrls: ['./title.component.scss'],
})
export class PostTitleComponent extends PostComponent implements OnInit {

  /**
   * Use to indiciate whether link should enable or not
   */
  @Input()
  link = true;

  /**
   * An array of string to be composed to router link
   */
  href: string[];

  ngOnInit(): void {
    const publishedAt: Date = new Date(
      new Date(this.post.publishedAt)
        .toLocaleString('en-US', { timeZone: 'Asia/Bangkok' }),
    );

    if (this.link) {
      this.href = [
        '/',
        publishedAt.getFullYear().toString(),
        (publishedAt.getMonth() + 1).toString(), // a month number start at 0 not 1
        publishedAt.getDate().toString(),
        this.post.slug,
      ];
    }
  }

  isDisabledLink(): boolean {
    return !!!this.link;
  }

}
